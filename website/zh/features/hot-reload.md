# 热更新

## 是什么

**热更新**把 Prompt 变更实时推送给运行中的应用。编辑者一保存 Prompt,所有
已连接的 SDK 就自动刷新缓存 —— 不重启、不发版。

## 为什么重要

没有它,改一个 Prompt 就要重启 AI 服务。有了它,Prompt 成为运行时资源:
改文案、点保存,生产环境毫秒级生效。

## 工作原理

1. SDK 调用 `watch()`,与服务端 `/ws` 建立一条 WebSocket 连接。
2. 当某个 Prompt 被创建、更新或删除时,服务端向所有已连接客户端**广播**
   一个小事件 —— `{ type, key, env, version }`。
3. 每个 SDK 检查该事件:如果它缓存了这个 `key`,就重新拉取最新内容,并触发
   一个 `update` 事件供你的代码响应。

Web 界面也用同样方式连接 —— Prompts 页面在 WebSocket 打开时会显示
**「热更新已连接」** 标记。

## 在 SDK 里使用

```js
import { PromptOpsClient } from '@promptops/client'

const client = new PromptOpsClient({
  server: 'http://localhost:8080',
  namespace: 'prod',
  appName: 'checkout-service', // 会显示在在线客户端面板里
})

// 预热缓存
await client.render('sql.generator', { question: '查询所有用户' })

// 响应实时变更
client.on('update', (e) => console.log('Prompt 已热更新:', e.key))
client.on('connect', () => console.log('已开始监听'))
client.watch()
```

调用 `watch()` 后,界面里对 `sql.generator` 的任何编辑都会触发 `update`
回调,并自动刷新缓存内容。

## 查看谁连接了

每一条 `watch()` 连接都会被记录。可观测性页面在 **在线客户端** 里列出它们
—— 见[可观测性](./observability)。

## 下一步

- [Prompt 管理](./prompts) —— 被热更新的资源本身。
- [SDK](/zh/sdk/node) —— 完整客户端参考。
