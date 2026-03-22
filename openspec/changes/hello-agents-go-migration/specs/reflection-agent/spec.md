## ADDED Requirements

### Requirement: ReflectionAgent 基础属性

ReflectionAgent SHALL 具有名称(name)、LLM 客户端(llm)、Memory、最大反思次数(max_reflections)等基础属性。

#### Scenario: 创建 Agent
- **WHEN** 使用名称和 LLM 客户端创建 ReflectionAgent
- **THEN** Agent 具有指定的名称、LLM 客户端和 Memory 实例

#### Scenario: 配置最大反思次数
- **WHEN** 创建时指定 max_reflections
- **THEN** Agent 最多进行指定次数的反思-改进循环

### Requirement: 初始执行 (Act Phase)

ReflectionAgent SHALL 首先执行初始方案。

#### Scenario: 初始执行
- **WHEN** 用户提出任务
- **THEN** Agent 执行初始方案并将结果存入 Memory

### Requirement: 反思评估 (Reflect Phase)

ReflectionAgent SHALL 对执行结果进行自我反思评估。

#### Scenario: 反思执行结果
- **WHEN** 进入反思阶段
- **THEN** Agent 分析上次执行结果，找出改进点

#### Scenario: 判断无需改进
- **WHEN** 反思结果显示 "无需改进" 或类似表述
- **THEN** Agent 提前结束迭代并返回当前结果

### Requirement: 改进执行 (Improve Phase)

ReflectionAgent SHALL 基于反思结果改进执行。

#### Scenario: 改进执行
- **WHEN** 反思发现改进点
- **THEN** Agent 基于反馈生成改进后的新执行结果

### Requirement: Memory 管理

ReflectionAgent SHALL 管理执行和反思的轨迹记录。

#### Scenario: 添加执行记录
- **WHEN** 执行一次方案
- **THEN** 系统将执行结果存入 Memory

#### Scenario: 添加反思记录
- **WHEN** 进行一次反思
- **THEN** 系统将反思反馈存入 Memory

#### Scenario: 获取轨迹
- **WHEN** 需要构建上下文
- **THEN** 系统返回格式化的执行-反思轨迹文本

#### Scenario: 获取最近执行
- **WHEN** 需要反思上次执行
- **THEN** 系统返回最近一次的执行结果

### Requirement: Run 方法

ReflectionAgent SHALL 实现 Agent 接口的 Run 方法。

#### Scenario: 完整流程
- **WHEN** 调用 Run 方法处理任务
- **THEN** 系统依次执行 Act -> [Reflect -> Improve]循环 -> 返回最终结果

#### Scenario: 达到最大反思次数
- **WHEN** 反思次数达到 max_reflections
- **THEN** Agent 返回当前最佳结果

### Requirement: 自定义 Prompt 模板

ReflectionAgent SHALL 支持自定义 Act、Reflect、Improve 阶段的 Prompt 模板。

#### Scenario: 使用自定义模板
- **WHEN** 创建时指定 act_prompt、reflect_prompt、improve_prompt
- **THEN** Agent 使用自定义模板