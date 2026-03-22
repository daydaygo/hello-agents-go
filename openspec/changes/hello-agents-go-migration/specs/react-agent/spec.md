## ADDED Requirements

### Requirement: ReAct 循环

ReActAgent SHALL 实现 Thought-Action-Observation 推理循环。

#### Scenario: 完整推理循环
- **WHEN** 用户提出需要多步骤解决的问题
- **THEN** Agent 执行 Thought -> Action -> Observation 循环直到得出答案

#### Scenario: 达到最大步数
- **WHEN** 推理步数达到 max_steps 限制
- **THEN** Agent 返回当前最佳答案或错误提示

### Requirement: Thought 解析

ReActAgent SHALL 正确解析 LLM 输出中的 Thought 部分。

#### Scenario: 解析 Thought
- **WHEN** LLM 输出包含 "Thought: ..." 格式
- **THEN** 系统提取 Thought 内容用于日志或调试

### Requirement: Action 解析与执行

ReActAgent SHALL 正确解析 Action 并调用相应工具。

#### Scenario: 解析工具调用
- **WHEN** LLM 输出包含 "Action: tool_name[args]" 格式
- **THEN** 系统提取工具名称和参数并执行工具

#### Scenario: 解析 Finish
- **WHEN** LLM 输出包含 "Finish: answer" 格式
- **THEN** 系统返回最终答案

#### Scenario: 无效 Action 格式
- **WHEN** Action 格式不正确
- **THEN** 系统返回错误提示并请求 LLM 重新生成

### Requirement: 工具注册

ReActAgent SHALL 持有工具注册表以执行工具调用。

#### Scenario: 注册工具
- **WHEN** 创建 Agent 时提供工具注册表
- **THEN** Agent 可执行注册表中的所有工具

#### Scenario: 工具不存在
- **WHEN** LLM 请求执行不存在的工具
- **THEN** Agent 返回错误 Observation

### Requirement: 历史记录管理

ReActAgent SHALL 管理 Thought-Action-Observation 历史记录。

#### Scenario: 记录执行历史
- **WHEN** 执行推理循环
- **THEN** 系统记录每一步的 Thought、Action、Observation

### Requirement: 自定义 Prompt 模板

ReActAgent SHALL 支持自定义 Prompt 模板。

#### Scenario: 使用自定义模板
- **WHEN** 创建时指定 custom_prompt
- **THEN** Agent 使用自定义模板格式化输入