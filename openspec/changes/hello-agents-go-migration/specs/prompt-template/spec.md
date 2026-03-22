## ADDED Requirements

### Requirement: 模板定义

Prompt 模板 SHALL 支持变量占位符，使用 `{variable}` 格式。

#### Scenario: 创建模板
- **WHEN** 创建包含 `{name}` 占位符的模板
- **THEN** 模板正确存储并可被渲染

### Requirement: 变量替换

Prompt 模板 SHALL 支持变量替换生成最终 Prompt。

#### Scenario: 单变量替换
- **WHEN** 提供变量 name="World" 渲染模板 "Hello, {name}!"
- **THEN** 系统返回 "Hello, World!"

#### Scenario: 多变量替换
- **WHEN** 提供多个变量渲染模板
- **THEN** 系统替换所有匹配的变量

#### Scenario: 变量缺失
- **WHEN** 模板中的变量未提供值
- **THEN** 系统保留原始占位符或返回错误

### Requirement: 预定义模板

系统 SHALL 提供常用的预定义 Prompt 模板。

#### Scenario: ReAct 模板
- **WHEN** 使用预定义的 ReAct 模板
- **THEN** 模板包含 tools、question、history 等占位符

#### Scenario: Simple Agent 模板
- **WHEN** 使用预定义的简单对话模板
- **THEN** 模板包含 system_prompt 和 user_input 占位符

#### Scenario: Plan-and-Solve Planner 模板
- **WHEN** 使用预定义的 Planner 模板
- **THEN** 模板包含 question 占位符用于生成计划

#### Scenario: Plan-and-Solve Executor 模板
- **WHEN** 使用预定义的 Executor 模板
- **THEN** 模板包含 question、plan、history、current_step 占位符

#### Scenario: Reflection Act 模板
- **WHEN** 使用预定义的 Act 模板
- **THEN** 模板包含 task 占位符

#### Scenario: Reflection Reflect 模板
- **WHEN** 使用预定义的 Reflect 模板
- **THEN** 模板包含 task、code 占位符

#### Scenario: Reflection Improve 模板
- **WHEN** 使用预定义的 Improve 模板
- **THEN** 模板包含 task、last_code_attempt、feedback 占位符

### Requirement: 模板管理

系统 SHALL 支持模板的注册和获取。

#### Scenario: 注册模板
- **WHEN** 注册自定义模板
- **THEN** 可通过名称获取该模板

#### Scenario: 获取模板
- **WHEN** 通过名称获取已注册模板
- **THEN** 系统返回对应的模板实例

#### Scenario: 模板不存在
- **WHEN** 获取未注册的模板名称
- **THEN** 系统返回错误