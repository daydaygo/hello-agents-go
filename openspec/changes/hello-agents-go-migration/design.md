## Context

本项目是将 Python 版 hello-agents 框架迁移到 Go 语言。源项目是一个 AI Agent 学习教程，包含 16 个章节的代码示例。目标是一个极简版 Go Agent 框架，核心功能可用、代码简洁、易于学习。

**技术约束**:
- Go >= 1.25
- 使用 `math/rand/v2` 和 `encoding/json/v2`
- 使用官方 OpenAI Go SDK v3
- 不使用 langchaingo (不支持 langchain 1.0)
- 测试仅使用标准库 `testing` 包

## Goals / Non-Goals

**Goals:**
1. 实现 LLM 客户端封装，支持 OpenAI 兼容 API
2. 实现 SimpleAgent，支持基础对话和历史消息管理
3. 实现 ReActAgent，支持 Thought-Action-Observation 循环
4. 实现 PlanAndSolveAgent，支持计划制定与逐步执行
5. 实现 ReflectionAgent，支持执行后自我反思优化
6. 实现工具注册与执行机制
7. 实现 Prompt 模板管理
8. 提供可运行的示例代码

**Non-Goals:**
1. 不实现 MCP/A2A 协议支持
2. 不实现 RAG/向量数据库集成
3. 不实现多 Agent 协作
4. 不实现模型训练 (SFT/GRPO)
5. 不迁移 NLP 基础代码 (BPE/Transformer)

## Decisions

### D1: 使用 OpenAI 官方 Go SDK v3

**决策**: 使用 `github.com/openai/openai-go/v3`

**原因**:
- 官方维护，API 完整
- 支持 Responses API 和 Chat Completions API
- 支持 Tool Calling 和流式响应
- 类型安全

**替代方案**:
- langchaingo: 不支持 langchain 1.0，不采用
- 自行封装 HTTP: 工作量大，易出错

### D2: Option 设计模式配置

**决策**: 使用 Option 设计模式进行配置

**原因**:
- Go 惯用模式
- 配置灵活，扩展性好
- 无需测试配置代码

### D3: Agent 接口设计

**决策**: 定义最小 Agent 接口

```go
type Agent interface {
    Run(ctx context.Context, input string) (string, error)
}
```

**原因**:
- 简单明了，易于实现
- 支持不同 Agent 类型扩展
- 符合 Go 接口设计原则

### D4: 工具接口设计

**决策**: 定义 Tool 接口

```go
type Tool interface {
    Name() string
    Description() string
    Execute(input string) (string, error)
}
```

**原因**:
- 简单通用
- 易于实现自定义工具
- 支持工具描述生成

### D5: 错误处理策略

**决策**: 
- 所有错误返回 error，不吞异常
- 使用 log/slog 记录完整错误堆栈
- 用户友好错误信息在前台处理

**原因**:
- 符合 Go 最佳实践
- 便于调试和问题定位
- 安全性考虑

### D6: 三种 Agent 范式设计

**决策**: 实现 ReAct、Plan-and-Solve、Reflection 三种核心 Agent 范式

| 范式 | 核心机制 | 适用场景 |
|-----|---------|---------|
| ReAct | Thought → Action → Observation 循环 | 需要即时决策的任务 |
| Plan-and-Solve | Plan(计划) → Execute(执行) → Verify(验证) | 复杂多步骤任务 |
| Reflection | Act(执行) → Reflect(反思) → Improve(改进) | 需要迭代优化的任务 |

**原因**:
- 覆盖主流 Agent 设计模式
- 对应 Python 源项目 Chapter 4 核心内容
- 学习价值高，展示不同推理策略

### D7: ReActAgent 设计

**决策**: 实现 Thought-Action-Observation 循环

```
用户输入 → [Thought → Action → Observation]循环 → 最终答案
```

**关键特性**:
- 最大步数限制 (max_steps)
- 支持自定义 Prompt 模板
- 执行历史记录

### D8: PlanAndSolveAgent 设计

**决策**: 两阶段执行：Planning → Execution

```
用户输入 → Planning(生成任务列表) → Execution(逐个执行) → 最终答案
```

**关键特性**:
- 计划分解为可执行的子任务
- 逐步执行并验证
- 支持计划调整

### D9: ReflectionAgent 设计

**决策**: 执行后反思改进

```
用户输入 → 执行初始方案 → 反思评估 → 改进执行 → 最终答案
```

**关键特性**:
- 自我评估执行结果
- 识别改进点
- 迭代优化直到满意

## Risks / Trade-offs

| Risk | Mitigation |
|------|------------|
| OpenAI SDK API 变更 | 锁定版本，定期更新 |
| Go 1.25 特性依赖 | 明确版本要求，CI 检查 |
| 功能精简过度 | 保留扩展接口，后续可增强 |
| 测试覆盖不足 | 核心路径必须有测试 |

## Migration Plan

**Phase 1: 基础设施**
1. 创建项目结构
2. 初始化 go.mod
3. 实现 LLM 客户端
4. 实现 Prompt 模板
5. 实现工具系统

**Phase 2: 基础 Agent**
6. 实现 SimpleAgent

**Phase 3: 核心范式 Agent**
7. 实现 ReActAgent
8. 实现 PlanAndSolveAgent
9. 实现 ReflectionAgent

**Phase 4: 完善**
10. 添加示例代码
11. 完善测试和文档

## Open Questions

1. 是否需要支持多 LLM Provider (Azure, Anthropic)?
   - 当前决策: 仅支持 OpenAI 兼容 API，通过 base_url 配置

2. 是否需要持久化对话历史?
   - 当前决策: 不持久化，仅内存存储，保持极简