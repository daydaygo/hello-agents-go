## ADDED Requirements

### Requirement: 工具注册

ToolRegistry SHALL 支持工具的注册和注销。

#### Scenario: 注册工具
- **WHEN** 调用 Register 方法注册工具
- **THEN** 工具被添加到注册表，可通过名称访问

#### Scenario: 注销工具
- **WHEN** 调用 Unregister 方法注销工具
- **THEN** 工具从注册表移除

#### Scenario: 重复注册
- **WHEN** 注册同名工具
- **THEN** 后注册的工具覆盖先前的工具

### Requirement: 工具执行

ToolRegistry SHALL 支持通过名称执行工具。

#### Scenario: 执行已注册工具
- **WHEN** 调用 Execute 方法执行已注册的工具
- **THEN** 系统调用工具的 Execute 方法并返回结果

#### Scenario: 执行未注册工具
- **WHEN** 调用 Execute 方法执行未注册的工具
- **THEN** 系统返回错误提示工具不存在

### Requirement: 工具描述生成

ToolRegistry SHALL 生成所有工具的描述文本，用于 Prompt 中告知 LLM 可用工具。

#### Scenario: 获取工具描述
- **WHEN** 调用 GetToolsDescription 方法
- **THEN** 系统返回所有工具的名称、描述、参数格式的文本

#### Scenario: 空注册表
- **WHEN** 注册表为空时获取描述
- **THEN** 系统返回提示文本 "暂无可用工具"

### Requirement: 工具列表

ToolRegistry SHALL 支持列出所有已注册工具。

#### Scenario: 获取工具列表
- **WHEN** 调用 ListTools 方法
- **THEN** 系统返回所有已注册工具名称列表

### Requirement: Tool 接口

系统 SHALL 定义 Tool 接口，包含 Name、Description、Execute 方法。

#### Scenario: 实现自定义工具
- **WHEN** 用户实现 Tool 接口
- **THEN** 该工具可被注册到 ToolRegistry

### Requirement: 内置工具

系统 SHALL 提供内置的 Calculator 工具。

#### Scenario: 计算器工具
- **WHEN** 执行 Calculator 工具并传入数学表达式
- **THEN** 系统返回计算结果