## ADDED Requirements

### Requirement: Agent 基础属性

SimpleAgent SHALL 具有名称(name)、LLM 客户端(llm)、系统提示词(system_prompt)等基础属性。

#### Scenario: 创建 Agent
- **WHEN** 使用名称和 LLM 客户端创建 SimpleAgent
- **THEN** Agent 具有指定的名称和 LLM 客户端实例

#### Scenario: 设置系统提示词
- **WHEN** 创建时指定系统提示词
- **THEN** Agent 在对话中使用该系统提示词

### Requirement: 基础对话

SimpleAgent SHALL 支持基础对话功能，返回 LLM 响应。

#### Scenario: 单轮对话
- **WHEN** 用户发送消息
- **THEN** Agent 返回 LLM 的响应文本

#### Scenario: 多轮对话
- **WHEN** 用户连续发送多条消息
- **THEN** Agent 保持对话上下文，每轮响应考虑历史消息

### Requirement: 历史消息管理

SimpleAgent SHALL 管理对话历史消息。

#### Scenario: 获取历史消息
- **WHEN** 查询历史消息
- **THEN** 系统返回所有用户和助手的消息列表

#### Scenario: 清除历史消息
- **WHEN** 清除历史消息
- **THEN** 系统重置消息历史为空

### Requirement: 工具调用支持

SimpleAgent SHALL 支持可选的工具调用功能。

#### Scenario: 无工具对话
- **WHEN** 未配置工具注册表
- **THEN** Agent 进行纯文本对话

#### Scenario: 有工具对话
- **WHEN** 配置工具注册表且 LLM 请求工具调用
- **THEN** Agent 执行工具并返回工具结果给 LLM

### Requirement: 配置选项

SimpleAgent SHALL 支持通过 Option 模式配置。

#### Scenario: 配置最大工具迭代次数
- **WHEN** 设置 max_tool_iterations 选项
- **THEN** Agent 最多执行指定次数的工具迭代