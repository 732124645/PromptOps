# Python SDK

纯标准库实现,内含一个最小 WebSocket 客户端。源码位于 `sdk/python`。

## 用法

```python
from promptops import PromptOpsClient

client = PromptOpsClient(
    "http://localhost:8080",
    namespace="prod",
    token="promptops-dev-token",
)

# 按 key 获取 Prompt(带本地缓存)
prompt = client.get_prompt("sql.generator")

# 渲染模板变量
text = client.render("sql.generator", {"question": "查询所有用户"})

# 订阅热更新
client.watch(on_update=lambda e: print("已热更新:", e["key"]))
```

## API

| 方法 | 说明 |
|---|---|
| `PromptOpsClient(server, namespace="prod", token=...)` | 创建客户端 |
| `get_prompt(key, refresh=False)` | 获取 Prompt,默认走缓存 |
| `render(key, variables)` | 获取并渲染 `{{变量}}` |
| `watch(on_update=...)` | 后台线程连接 WebSocket,自动刷新缓存 |
| `close()` | 关闭连接 |

`render_template(content, variables)` 也可单独导入使用。
