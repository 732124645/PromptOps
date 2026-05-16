# Python SDK

Pure standard library, including a minimal WebSocket client. Source lives in
`sdk/python`.

## Usage

```python
from promptops import PromptOpsClient

client = PromptOpsClient(
    "http://localhost:8080",
    namespace="prod",
    token="promptops-dev-token",
)

# Fetch a prompt by key (locally cached)
prompt = client.get_prompt("sql.generator")

# Render template variables
text = client.render("sql.generator", {"question": "list all users"})

# Subscribe to hot-reload
client.watch(on_update=lambda e: print("hot-reloaded:", e["key"]))
```

## API

| Method | Description |
|---|---|
| `PromptOpsClient(server, namespace="prod", token=...)` | Create a client |
| `get_prompt(key, refresh=False)` | Fetch a prompt; cached by default |
| `render(key, variables)` | Fetch and render `{{variables}}` |
| `watch(on_update=...)` | Connect the WebSocket in a background thread |
| `close()` | Close the connection |

`render_template(content, variables)` can also be imported standalone.
