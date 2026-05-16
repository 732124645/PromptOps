# Java SDK

Pure JDK — uses `java.net.http` (HttpClient + WebSocket) with a small built-in
JSON parser. Source lives in `sdk/java`; requires Java 21+.

## Usage

```java
import io.promptops.PromptOpsClient;
import java.util.Map;

PromptOpsClient client = new PromptOpsClient(
    "http://localhost:8080", "prod", "promptops-dev-token");

// Fetch a prompt by key (locally cached)
var prompt = client.getPrompt("sql.generator");

// Render template variables
String text = client.render("sql.generator", Map.of("question", "list all users"));

// Subscribe to hot-reload
client.watch(key -> System.out.println("hot-reloaded: " + key));
```

## API

| Method | Description |
|---|---|
| `new PromptOpsClient(server[, namespace, token])` | Create a client |
| `getPrompt(key[, refresh])` | Fetch a prompt; cached by default |
| `render(key, variables)` | Fetch and render `{{variables}}` |
| `watch(onUpdate)` | Open the WebSocket and auto-refresh the cache |
| `close()` | Close the connection |

`PromptOpsClient.renderTemplate(content, variables)` is a static method you can
call directly.
