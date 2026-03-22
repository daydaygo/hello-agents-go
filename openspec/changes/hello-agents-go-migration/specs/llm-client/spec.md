## ADDED Requirements

### Requirement: LLM 客户端初始化

系统 SHALL 支持通过 Option 模式配置 LLM 客户端，包括 API Key、Base URL、Model、Timeout 等参数。

#### Scenario: 使用默认配置初始化
- **WHEN** 仅提供 API Key
- **THEN** 系统使用默认 Base URL (https://api.openai.com/v1) 和默认模型

#### Scenario: 使用自定义配置初始化
- **WHEN** 提供完整的配置参数
- **THEN** 系统使用自定义的配置创建客户端

### Requirement: 同步调用 LLM

系统 SHALL 支持同步调用 LLM API，返回完整响应文本。

#### Scenario: 基础对话
- **WHEN** 发送用户消息
- **THEN** 系统返回 LLM 的完整响应文本

#### Scenario: 带系统提示词
- **WHEN** 发送系统提示词和用户消息
- **THEN** 系统返回遵循系统提示词的响应

### Requirement: 流式调用 LLM

系统 SHALL 支持流式调用 LLM API，逐块返回响应内容。

#### Scenario: 流式输出
- **WHEN** 请求流式响应
- **THEN** 系统通过 channel 或回调逐块返回内容

### Requirement: Tool Calling 支持

系统 SHALL 支持 OpenAI Tool Calling 功能。

#### Scenario: 工具调用请求
- **WHEN** LLM 返回工具调用请求
- **THEN** 系统解析工具名称和参数

#### Scenario: 工具结果提交
- **WHEN** 提交工具执行结果
- **THEN** 系统继续对话并返回最终响应

### Requirement: 错误处理

系统 SHALL 正确处理 API 错误并返回详细错误信息。

#### Scenario: API 错误
- **WHEN** API 返回错误响应
- **THEN** 系统返回包含状态码和错误信息的 error

#### Scenario: 网络错误
- **WHEN** 网络连接失败
- **THEN** 系统返回包装后的网络错误