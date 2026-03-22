# AGENTS.md - Python to Go Migration Rules

This file defines the rules and patterns for migrating Python code to Go in this project.

## 1. 项目结构映射

| Python | Go |
|--------|-----|
| `code/` | `pkg/` |
| `code/llm/` | `pkg/llm/` |
| `code/agent/` | `pkg/agent/` |
| `code/tool/` | `pkg/tool/` |
| `code/prompt/` | `pkg/prompt/` |

## 2. 类型映射

| Python | Go |
|--------|-----|
| `str` | `string` |
| `int` | `int` |
| `float` | `float64` |
| `bool` | `bool` |
| `list` | `[]T` (slice) |
| `dict` | `map[string]T` |
| `Optional[T]` | `*T` (pointer) |
| `Any` | `interface{}` 或 `any` |

## 3. 错误处理

- Python: 使用 `try/except` 和返回 `None`
- Go: 返回 `error` 作为最后一个返回值，不吞异常

```python
# Python
def do_something():
    try:
        return result, None
    except Exception as e:
        return None, str(e)
```

```go
// Go
func DoSomething() (Result, error) {
    result, err := internal()
    if err != nil {
        return Result{}, fmt.Errorf("do something: %w", err)
    }
    return result, nil
}
```

## 4. 配置模式 (Option Pattern)

使用 Option 模式替代 Python 的 kwargs:

```python
# Python
class Client:
    def __init__(self, api_key, base_url=None, model=None, timeout=30):
        ...
```

```go
// Go
type Option func(*Client)

func WithBaseURL(url string) Option {
    return func(c *Client) { c.baseURL = url }
}

func NewClient(apiKey string, opts ...Option) *Client {
    c := &Client{apiKey: apiKey}
    for _, opt := range opts {
        opt(c)
    }
    return c
}
```

## 5. 接口设计

优先使用接口，而非继承:

```python
# Python
class Agent(ABC):
    @abstractmethod
    def run(self, input: str) -> str:
        pass
```

```go
// Go
type Agent interface {
    Run(ctx context.Context, input string) (string, error)
}
```

## 6. 日志

使用 `log/slog` 替代 Python 的 `logging`:

```python
# Python
import logging
logger = logging.getLogger(__name__)
logger.info("message", extra={"key": "value"})
```

```go
// Go
import "log/slog"
slog.Info("message", "key", "value")
```

## 7. JSON 处理

使用 `encoding/json` (Go 1.25+ 使用 v2):

```python
# Python
import json
data = json.loads(json_str)
json_str = json.dumps(data)
```

```go
// Go
import "encoding/json"
var data Data
json.Unmarshal([]byte(jsonStr), &data)
jsonBytes, _ := json.Marshal(data)
```

## 8. 依赖管理

| Python | Go |
|--------|-----|
| `openai` | `github.com/openai/openai-go/v3` |
| `python-dotenv` | `os.Getenv()` |
| `logging` | `log/slog` |
| `random` | `math/rand/v2` |

## 9. 测试

使用标准库 `testing` 包:

```python
# Python
import unittest
class TestClient(unittest.TestCase):
    def test_invoke(self):
        ...
```

```go
// Go
func TestClient_Invoke(t *testing.T) {
    ...
}
```

## 10. 包命名规范

- 包名使用小写，无下划线
- 包名应简洁且有意义
- 避免使用 `util`、`common` 等通用名称

## 11. 代码风格

- 使用 `gofmt` 格式化代码
- 使用 `golangci-lint` 检查代码
- 导出函数必须有文档注释
- 错误信息不以大写开头，不以句号结尾

## 12. 并发处理

使用 `context.Context` 进行超时和取消控制:

```go
func (c *Client) Invoke(ctx context.Context, input string) (string, error) {
    select {
    case <-ctx.Done():
        return "", ctx.Err()
    default:
        // 执行操作
    }
}
```

## 13. 环境变量

使用 `os.Getenv()` 读取环境变量:

```python
# Python
import os
api_key = os.getenv("OPENAI_API_KEY")
```

```go
// Go
import "os"
apiKey := os.Getenv("OPENAI_API_KEY")
```

## 14. 文档注释

导出类型、函数、变量必须有文档注释:

```go
// Client represents an LLM client.
type Client struct {
    // apiKey is the API key for authentication.
    apiKey string
}

// NewClient creates a new LLM client with the given API key.
func NewClient(apiKey string, opts ...Option) *Client {
    ...
}
```

## 15. 不迁移内容

以下 Python 内容不迁移:
- NLP 基础代码 (BPE/Transformer)
- n8n 配置文件
- 第三方框架代码 (LangGraph/AutoGen)
- 高级特性 (MCP/A2A/RAG)