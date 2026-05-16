# Java SDK

纯 JDK 实现,使用 `java.net.http`(HttpClient + WebSocket),自带一个最小
JSON 解析器。源码位于 `sdk/java`,需要 Java 21+。

## 用法

```java
import io.promptops.PromptOpsClient;
import java.util.Map;

PromptOpsClient client = new PromptOpsClient(
    "http://localhost:8080", "prod", "promptops-dev-token");

// 按 key 获取 Prompt(带本地缓存)
var prompt = client.getPrompt("sql.generator");

// 渲染模板变量
String text = client.render("sql.generator", Map.of("question", "查询所有用户"));

// 订阅热更新
client.watch(key -> System.out.println("已热更新: " + key));
```

## API

| 方法 | 说明 |
|---|---|
| `new PromptOpsClient(server[, namespace, token])` | 创建客户端 |
| `getPrompt(key[, refresh])` | 获取 Prompt,默认走缓存 |
| `render(key, variables)` | 获取并渲染 `{{变量}}` |
| `watch(onUpdate)` | 打开 WebSocket,自动刷新缓存 |
| `close()` | 关闭连接 |

`PromptOpsClient.renderTemplate(content, variables)` 是可直接调用的静态方法。
