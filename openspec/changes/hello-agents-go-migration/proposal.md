## Why

Python 项目 hello-agents 是一个 AI Agent 学习教程，包含 16 章节代码。现有项目缺少 Go 语言版本的 Agent 框架实现，Go 开发者缺乏本地化的学习资源。本项目旨在将 hello-agents 核心功能迁移到 Go 语言，提供一个极简版、生产可用的 Agent 框架。

## What Changes

- **新增**: Go 语言版本的 Agent 框架核心实现
- **新增**: LLM 客户端封装 (基于 OpenAI 官方 Go SDK v3)
- **新增**: SimpleAgent 和 ReActAgent 实现
- **新增**: 工具注册与执行机制
- **新增**: Prompt 模板管理
- **新增**: 示例代码和测试用例

### 迁移范围 (对应 Python 章节)

| 章节 | 内容 | 迁移状态 |
|-----|------|---------|
| Chapter 1-2 | Agent 入门/ELIZA | 不迁移 (概念性代码) |
| Chapter 3 | NLP 基础 (BPE/Transformer) | 不迁移 (非 Agent 核心) |
| Chapter 4 | ReAct/Reflection/Plan-and-Solve | ✅ 迁移核心模式 |
| Chapter 5 | n8n 工作流配置 | 不迁移 (配置文件) |
| Chapter 6 | 第三方框架 (LangGraph/AutoGen) | 不迁移 |
| Chapter 7 | 自定义 Agent 实现 | ✅ 迁移核心实现 |
| Chapter 8 | Memory/RAG | 迁移简化版 (working memory) |
| Chapter 9-12 | 高级特性 | 不迁移 (超出极简范围) |

## Capabilities

### New Capabilities

- `llm-client`: LLM 客户端封装，支持 OpenAI 兼容 API、流式响应、Tool Calling
- `simple-agent`: 基础对话 Agent，支持历史消息、系统提示词
- `react-agent`: 推理-行动 Agent，支持 Thought-Action-Observation 循环
- `tool-registry`: 工具注册表，支持工具注册、执行、描述生成
- `prompt-template`: Prompt 模板管理，支持变量替换

### Modified Capabilities

(无 - 这是新项目)

## Impact

### 技术栈

| 组件 | Python | Go |
|-----|--------|-----|
| LLM SDK | openai Python | github.com/openai/openai-go/v3 |
| 环境变量 | python-dotenv | os.Getenv() |
| JSON | json | encoding/json/v2 |
| 随机数 | random | math/rand/v2 |
| 日志 | logging | log/slog |

### 代码结构

```
hello-agents-go/
├── go.mod
├── pkg/
│   ├── llm/          # LLM 客户端
│   ├── agent/        # Agent 实现
│   ├── tool/         # 工具系统
│   └── prompt/       # Prompt 模板
├── examples/         # 示例代码
└── internal/         # 内部实现
```

### 依赖

- Go >= 1.25
- github.com/openai/openai-go/v3 (官方 SDK)