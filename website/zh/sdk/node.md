# Node SDK

零依赖,使用 Node 22+ 内置的 `fetch` 与 `WebSocket`。源码位于 `sdk/node`。

## 用法

```js
import { PromptOpsClient } from '@promptops/client'

const client = new PromptOpsClient({
  server: 'http://localhost:8080',
  namespace: 'prod',
  token: 'promptops-dev-token',
})

// 按 key 获取 Prompt(带本地缓存)
const prompt = await client.getPrompt('sql.generator')

// 渲染模板变量
const text = await client.render('sql.generator', { question: '查询所有用户' })

// 订阅热更新:Prompt 在后端变更后自动刷新缓存
client.on('update', (e) => console.log('已热更新:', e.key))
client.watch()
```

## API

| 方法 | 说明 |
|---|---|
| `new PromptOpsClient({ server, namespace?, token? })` | 创建客户端 |
| `getPrompt(key, { refresh? })` | 获取 Prompt,默认走缓存 |
| `render(key, vars)` | 获取并渲染 `{{变量}}` |
| `on(event, cb)` | 监听 `update` / `connect` / `disconnect` / `error` |
| `watch()` | 打开 WebSocket,自动刷新缓存 |
| `close()` | 关闭连接 |

`renderTemplate(content, vars)` 也可单独导出使用。
