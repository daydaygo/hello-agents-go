## 1. 项目初始化

- [ ] 1.1 创建 go.mod 文件，指定 Go 1.25 和依赖, 包名 github.com/daydaygo/hello-agents-go
- [ ] 1.2 创建项目目录结构 (pkg/llm, pkg/agent, pkg/tool, pkg/prompt)
- [ ] 1.3 创建 .env.example 和 AGENTS.md 配置文件
- [ ] 1.4 创建 README.md 项目说明文档

## 2. LLM 客户端实现

- [ ] 2.1 创建 pkg/llm/client.go，定义 Client 结构体和 Option 接口
- [ ] 2.2 实现 WithAPIKey、WithBaseURL、WithModel、WithTimeout 等 Option 函数
- [ ] 2.3 实现 NewClient 构造函数
- [ ] 2.4 实现 Invoke 方法 (同步调用)
- [ ] 2.5 实现 StreamInvoke 方法 (流式调用)
- [ ] 2.6 实现 Message 结构体和消息构建辅助函数
- [ ] 2.7 编写 llm/client_test.go 单元测试

## 3. Prompt 模板实现

- [ ] 3.1 创建 pkg/prompt/template.go，定义 Template 结构体
- [ ] 3.2 实现 Render 方法 (变量替换)
- [ ] 3.3 创建预定义模板 (ReAct 模板、Simple Agent 模板)
- [ ] 3.4 创建预定义模板 (Plan-and-Solve Planner/Executor 模板)
- [ ] 3.5 创建预定义模板 (Reflection Act/Reflect/Improve 模板)
- [ ] 3.6 实现 TemplateRegistry 模板管理器
- [ ] 3.7 编写 prompt/template_test.go 单元测试

## 4. 工具系统实现

- [ ] 4.1 创建 pkg/tool/tool.go，定义 Tool 接口
- [ ] 4.2 创建 pkg/tool/registry.go，定义 ToolRegistry 结构体
- [ ] 4.3 实现 Register、Unregister、Get、Execute 方法
- [ ] 4.4 实现 GetToolsDescription 方法 (生成工具描述)
- [ ] 4.5 实现内置 Calculator 工具
- [ ] 4.6 编写 tool/registry_test.go 单元测试

## 5. SimpleAgent 实现

- [ ] 5.1 创建 pkg/agent/agent.go，定义 Agent 接口
- [ ] 5.2 创建 pkg/agent/simple.go，定义 SimpleAgent 结构体
- [ ] 5.3 实现 SimpleAgent 构造函数和 Option 配置
- [ ] 5.4 实现 Run 方法 (基础对话)
- [ ] 5.5 实现历史消息管理 (AddMessage、GetHistory、ClearHistory)
- [ ] 5.6 实现可选的工具调用支持
- [ ] 5.7 编写 agent/simple_test.go 单元测试

## 6. ReActAgent 实现

- [ ] 6.1 创建 pkg/agent/react.go，定义 ReActAgent 结构体
- [ ] 6.2 实现 ReActAgent 构造函数和 Option 配置
- [ ] 6.3 实现 Run 方法 (Thought-Action-Observation 循环)
- [ ] 6.4 实现 Thought/Action 解析函数
- [ ] 6.5 实现 Finish 检测和最终答案返回
- [ ] 6.6 实现执行历史记录管理
- [ ] 6.7 实现自定义 Prompt 模板支持
- [ ] 6.8 编写 agent/react_test.go 单元测试

## 7. PlanAndSolveAgent 实现

- [ ] 7.1 创建 pkg/agent/plan_solve.go，定义 PlanAndSolveAgent 结构体
- [ ] 7.2 实现 PlanAndSolveAgent 构造函数和 Option 配置
- [ ] 7.3 实现 planPhase 方法 (生成任务计划)
- [ ] 7.4 实现 executePhase 方法 (逐个执行子任务)
- [ ] 7.5 实现 Task 结构体 (描述、状态、结果)
- [ ] 7.6 实现计划验证和调整机制
- [ ] 7.7 实现自定义 Prompt 模板支持
- [ ] 7.8 编写 agent/plan_solve_test.go 单元测试

## 8. ReflectionAgent 实现

- [ ] 8.1 创建 pkg/agent/reflection.go，定义 ReflectionAgent 结构体
- [ ] 8.2 实现 ReflectionAgent 构造函数和 Option 配置
- [ ] 8.3 实现 Run 方法 (执行-反思-改进循环)
- [ ] 8.4 实现 actPhase 方法 (初始执行)
- [ ] 8.5 实现 reflectPhase 方法 (自我反思评估)
- [ ] 8.6 实现 improvePhase 方法 (基于反思改进)
- [ ] 8.7 实现最大反思次数限制 (max_reflections)
- [ ] 8.8 实现自定义 Prompt 模板支持
- [ ] 8.9 编写 agent/reflection_test.go 单元测试

## 9. 示例代码

- [ ] 9.1 创建 examples/simple-agent/main.go (基础对话示例)
- [ ] 9.2 创建 examples/react-agent/main.go (ReAct 示例)
- [ ] 9.3 创建 examples/plan-solve-agent/main.go (Plan-and-Solve 示例)
- [ ] 9.4 创建 examples/reflection-agent/main.go (Reflection 示例)
- [ ] 9.5 创建 examples/tool-calling/main.go (工具调用示例)

## 10. 完善与验证

- [ ] 10.1 运行 gofmt 格式化所有代码
- [ ] 10.2 运行 golangci-lint 代码检查
- [ ] 10.3 运行 modernize 工具检查现代 Go 写法
- [ ] 10.4 确保所有测试通过 (go test ./...)
- [ ] 10.5 更新 README.md 使用说明