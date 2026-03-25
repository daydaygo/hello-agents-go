# Eino Agent Visualization

This document visualizes the internal architecture and execution flow of the three Eino agent patterns implemented in this project.

## Overview

```mermaid
graph LR
    subgraph Eino Agent Patterns
        A[ChatModel Agent] --- B[ReAct Agent] --- C[Supervisor Multi-Agent]
    end

    A --> A1[Simple tool-calling loop]
    B --> B1[Reasoning + Action loop]
    C --> C1[Multi-agent coordination]
```

## 1. ChatModel Agent

The ChatModel Agent (`adk.NewChatModelAgent`) is a basic agent that loops between the ChatModel and tools until no more tool calls are produced.

### Internal Graph

```mermaid
graph TD
    START((START)) --> GenModelInput[GenModelInput<br><i>Combine instruction + messages</i>]
    GenModelInput --> ChatModel[ChatModel<br><i>LLM Generate/Stream</i>]
    ChatModel --> Branch{Has tool calls?}
    Branch -- Yes --> ToolsNode[ToolsNode<br><i>Execute tool calls</i>]
    ToolsNode --> EventSend[Send tool result events]
    EventSend --> ChatModel
    Branch -- No --> Output[Agent Output<br><i>Return message</i>]
    Output --> END((END))

    style START fill:#2d9,stroke:#333,color:#000
    style END fill:#f66,stroke:#333,color:#fff
    style ChatModel fill:#69f,stroke:#333,color:#fff
    style ToolsNode fill:#f90,stroke:#333,color:#fff
    style Branch fill:#ff0,stroke:#333,color:#000
```

### Execution Flow

```mermaid
sequenceDiagram
    participant User
    participant Runner
    participant Agent as ChatModelAgent
    participant Model as ChatModel (LLM)
    participant Tools as ToolsNode

    User->>Runner: runner.Query(ctx, input)
    Runner->>Agent: Run(ctx, AgentInput)

    loop MaxIterations (default: 20)
        Agent->>Model: Generate(messages + instruction)
        Model-->>Agent: Response message

        alt Has tool calls
            Agent->>Tools: Execute tool calls
            Tools-->>Agent: Tool results
            Agent->>Agent: Append tool results to messages
        else No tool calls
            Agent-->>Runner: Final message
        end
    end

    Runner-->>User: Stream events via AsyncIterator
```

### Code Mapping

| Component | Source | Description |
|-----------|--------|-------------|
| `adk.NewChatModelAgent` | `github.com/cloudwego/eino/adk` | Agent constructor |
| `adk.NewRunner` | `github.com/cloudwego/eino/adk` | Event-driven runner |
| `ToolsConfig` | `compose.ToolsNodeConfig` | Tool registration |
| `utils.InferTool` | `eino/components/tool/utils` | Auto-infer tool schema from function |

### Example: Greeting Agent

```
examples/eino/chatmodel-agent/main.go

OpenAI ChatModel
    ├── Agent: "assistant"
    │   ├── Instruction: "You are a helpful assistant..."
    │   └── Tools:
    │       └── greeting(name) → "Hello, {name}!"
    └── Runner (streaming enabled)
```

---

## 2. ReAct Agent

The ReAct Agent (`react.NewAgent`) implements the Reasoning-Action pattern using a `compose.Graph` with state management. It loops between the ChatModel and tools, with a branch condition that checks for tool calls.

### Internal Graph (compose.Graph)

```mermaid
graph TD
    START((START)) --> chat[chat<br><b>ChatModel</b><br><i>Model with tool calling</i>]

    chat --> branch1{StreamToolCallChecker<br><i>Has tool calls?</i>}

    branch1 -- Yes --> tools[tools<br><b>ToolsNode</b><br><i>Execute all tool calls</i>]
    branch1 -- No --> END((END))

    tools --> branch2{ReturnDirectly?}
    branch2 -- No --> chat
    branch2 -- Yes --> direct_return[direct_return<br><i>Extract direct result</i>]
    direct_return --> END

    style START fill:#2d9,stroke:#333,color:#000
    style END fill:#f66,stroke:#333,color:#fff
    style chat fill:#69f,stroke:#333,color:#fff
    style tools fill:#f90,stroke:#333,color:#fff
    style branch1 fill:#ff0,stroke:#333,color:#000
    style branch2 fill:#ff0,stroke:#333,color:#000
    style direct_return fill:#c6f,stroke:#333,color:#fff
```

### State Management

```mermaid
graph LR
    subgraph "ReAct State"
        Messages["Messages []Message<br><i>Accumulated conversation</i>"]
        RDTCID["ReturnDirectlyToolCallID<br><i>Short-circuit flag</i>"]
    end
```

The ReAct agent uses `compose.Graph` local state to track:
- **Messages**: Accumulated message history across all iterations
- **ReturnDirectlyToolCallID**: When set, causes the agent to return the tool result directly

### Node Details

| Node Key | Node Name | Type | Description |
|----------|-----------|------|-------------|
| `chat` | `ChatModel` | `ChatModelNode` | Calls LLM with tool info bound; `modelPreHandle` appends input to state and applies `MessageModifier` |
| `tools` | `Tools` | `ToolsNode` | Executes tool calls from model response; `toolsNodePreHandle` checks for `ReturnDirectly` |
| `direct_return` | - | `LambdaNode` | Extracts the tool result matching `ReturnDirectlyToolCallID` |

### Branch Logic

| From Node | Condition | Target |
|-----------|-----------|--------|
| `chat` | `StreamToolCallChecker` returns true (has tool calls) | `tools` |
| `chat` | No tool calls detected | `END` |
| `tools` | `ReturnDirectlyToolCallID` is set | `direct_return` |
| `tools` | No return-directly flag | `chat` (loop) |

### Execution Flow

```mermaid
sequenceDiagram
    participant User
    participant Agent as ReAct Agent
    participant State as Graph State
    participant Model as ChatModel
    participant Tools as ToolsNode

    User->>Agent: Stream(ctx, messages)

    loop MaxStep (default: 20)
        Agent->>State: Append input to state.Messages
        Agent->>Model: Generate(state.Messages + modifier)
        Model-->>Agent: Response (message stream)

        alt Has tool calls (StreamToolCallChecker)
            Agent->>State: Append assistant message
            Agent->>Tools: Execute tool calls
            Tools-->>Agent: Tool results []*Message
            Agent->>State: Append tool results

            alt ReturnDirectly tool called
                Agent-->>User: Return tool result directly
            else Normal flow
                Note over Agent: Continue loop
            end
        else No tool calls
            Agent-->>User: Return final message
        end
    end
```

### Example: Calculator Agent

```
examples/eino/react-agent/main.go

compose.Graph[[]Message, Message]
    ├── Nodes:
    │   ├── chat (ChatModel) ← OpenAI gpt-4o with tool calling
    │   ├── tools (ToolsNode) ← calculator tool
    │   └── direct_return (Lambda) ← short-circuit return
    ├── Edges:
    │   ├── START → chat
    │   └── direct_return → END
    ├── Branches:
    │   ├── chat → {tools, END}
    │   └── tools → {chat, direct_return}
    └── Options:
        ├── MaxRunSteps: 20
        └── NodeTriggerMode: AnyPredecessor
```

---

## 3. Supervisor Multi-Agent

The Supervisor pattern (`supervisor.New`) creates a hierarchical multi-agent system where a supervisor agent coordinates sub-agents via `transfer_to_agent` tool calls.

### Architecture

```mermaid
graph TD
    subgraph Supervisor Container
        SUP[supervisor<br><b>ChatModelAgent</b><br><i>Routes tasks to specialists</i>]
    end

    subgraph Sub-Agents
        MATH[math_agent<br><b>ChatModelAgent</b><br><i>Math specialist</i>]
        GEN[general_agent<br><b>ChatModelAgent</b><br><i>General knowledge</i>]
    end

    User((User)) --> SUP
    SUP -- "transfer_to_agent(math_agent)" --> MATH
    SUP -- "transfer_to_agent(general_agent)" --> GEN
    MATH -- "transfer_to_agent(supervisor)" --> SUP
    GEN -- "transfer_to_agent(supervisor)" --> SUP
    SUP --> Response((Response))

    style SUP fill:#69f,stroke:#333,color:#fff
    style MATH fill:#f90,stroke:#333,color:#fff
    style GEN fill:#9c6,stroke:#333,color:#000
    style User fill:#2d9,stroke:#333,color:#000
    style Response fill:#f66,stroke:#333,color:#fff
```

### Transfer Mechanism

```mermaid
graph LR
    subgraph "Built-in Tools"
        Transfer["transfer_to_agent<br><i>name: string</i>"]
    end

    subgraph "Supervisor sees"
        T1["transfer_to_agent → math_agent"]
        T2["transfer_to_agent → general_agent"]
    end

    subgraph "Sub-agent sees (DeterministicTransfer)"
        T3["transfer_to_agent → supervisor<br><i>(auto-configured, only option)</i>"]
    end
```

The supervisor pattern works through the ADK's agent transfer mechanism:

1. **Supervisor** gets `transfer_to_agent` tool with all sub-agent names as options
2. **Sub-agents** get `transfer_to_agent` tool via `DeterministicTransferTo`, pre-configured to only transfer back to the supervisor
3. When a transfer tool call is made, the ADK runtime switches execution context to the target agent

### Event Flow

```mermaid
sequenceDiagram
    participant User
    participant Runner
    participant Supervisor
    participant MathAgent as math_agent
    participant GenAgent as general_agent

    User->>Runner: Query("What is 15 + 27?")
    Runner->>Supervisor: Run(input)
    Supervisor->>Supervisor: LLM decides: math question
    Supervisor-->>Runner: Event: TransferToAgent(math_agent)
    Runner->>MathAgent: Run(delegated input)

    MathAgent->>MathAgent: LLM + tools (add, multiply, evaluate)
    MathAgent-->>Runner: Event: message content
    MathAgent-->>Runner: Event: TransferToAgent(supervisor)
    Runner->>Supervisor: Resume with result

    Supervisor-->>Runner: Event: final message
    Runner-->>User: Stream all events
```

### Agent Hierarchy

```
examples/eino/supervisor/main.go

supervisorContainer (unified tracing)
└── supervisor (ChatModelAgent)
    ├── Instruction: "Route math questions to math_agent..."
    ├── Model: OpenAI gpt-4o
    ├── Sub-agents:
    │   ├── math_agent (ChatModelAgent)
    │   │   ├── Instruction: "You are a math specialist..."
    │   │   ├── Tools:
    │   │   │   ├── add(a, b) → a + b
    │   │   │   ├── multiply(a, b) → a * b
    │   │   │   └── evaluate_expression(expr) → result
    │   │   └── DeterministicTransferTo: [supervisor]
    │   │
    │   └── general_agent (ChatModelAgent)
    │       ├── Instruction: "You are a general knowledge specialist..."
    │       ├── Tools: (none)
    │       └── DeterministicTransferTo: [supervisor]
    │
    └── Runner (streaming enabled)
```

---

## Comparison: Three Agent Patterns

```mermaid
graph TB
    subgraph "ChatModel Agent"
        direction TB
        C1[Model] --> C2{Tools?}
        C2 -- Yes --> C3[Execute] --> C1
        C2 -- No --> C4[Done]
    end

    subgraph "ReAct Agent"
        direction TB
        R1[Model] --> R2{Tools?}
        R2 -- Yes --> R3[Execute] --> R4{Direct?}
        R4 -- No --> R1
        R4 -- Yes --> R5[Done]
        R2 -- No --> R5
    end

    subgraph "Supervisor"
        direction TB
        S1[Supervisor] --> S2{Delegate}
        S2 --> S3[Sub-Agent A]
        S2 --> S4[Sub-Agent B]
        S3 --> S1
        S4 --> S1
        S1 --> S5[Done]
    end
```

| Feature | ChatModel Agent | ReAct Agent | Supervisor |
|---------|----------------|-------------|------------|
| **Package** | `eino/adk` | `eino/flow/agent/react` | `eino/adk/prebuilt/supervisor` |
| **Graph Type** | Internal loop | `compose.Graph` (explicit) | Agent tree (implicit) |
| **State** | `ChatModelAgentState` | Custom `state` struct | Per-agent state |
| **Max Iterations** | 20 (default) | 12 (default) | Per-agent limits |
| **Tool Support** | Yes | Yes | Per sub-agent |
| **Multi-Agent** | No | No | Yes |
| **Return Directly** | `ToolsConfig.ReturnDirectly` | `ToolReturnDirectly` | N/A |
| **Streaming** | Via Runner | Native `Stream()` | Via Runner |
| **Tracing** | Per-agent callbacks | Graph-level callbacks | Unified container |

---

## Eino Compose Graph Primitives

The Eino agents are built on `compose.Graph`, a stateful directed graph execution engine:

```mermaid
graph LR
    subgraph "compose.Graph Primitives"
        Node[Node<br><i>ChatModel / Tools / Lambda</i>]
        Edge[Edge<br><i>Direct connection</i>]
        Branch[Branch<br><i>Conditional routing</i>]
        State[State<br><i>Shared mutable state</i>]
    end

    subgraph "Execution"
        Compile[graph.Compile] --> Runnable
        Runnable --> Invoke[Invoke<br><i>Sync</i>]
        Runnable --> Stream[Stream<br><i>Async</i>]
    end
```

### Key Concepts

- **Node**: A processing unit (ChatModel, ToolsNode, Lambda function)
- **Edge**: Direct connection from one node to another (`AddEdge`)
- **Branch**: Conditional routing based on output inspection (`AddBranch`)
- **State**: Per-run mutable state accessible by all nodes via `compose.ProcessState`
- **Compile**: Converts the graph definition into an executable `Runnable`
- **MaxRunSteps**: Safety limit to prevent infinite loops
- **NodeTriggerMode**: `AnyPredecessor` allows a node to fire when any incoming edge delivers data
