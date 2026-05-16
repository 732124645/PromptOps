# 快速上手

## 环境要求

- Go 1.24+
- Node.js 18+

## 启动后端

```bash
cd server
go mod tidy
go run .
```

后端默认监听 `:8080`,SQLite 数据库位于 `server/data/promptops.db`。

可配置的环境变量:

| 变量 | 默认值 | 说明 |
|---|---|---|
| `PROMPTOPS_ADDR` | `:8080` | 监听地址 |
| `PROMPTOPS_DB` | `data/promptops.db` | SQLite 文件路径 |
| `PROMPTOPS_TOKEN` | `promptops-dev-token` | 静态管理员 Token |

## 启动前端

```bash
cd web
npm install
npm run dev
```

打开 `http://localhost:5173`,默认管理员账号 `admin` / `admin`。

## 创建第一个 Prompt

1. 登录后进入 **Prompts** 页面,点击「新建 Prompt」。
2. 填写 `key`(例如 `sql.generator`)、环境、内容,内容里可用 `{{变量}}` 占位符。
3. 保存后点击「发布版本」生成一个不可变快照。

## 在 Playground 试运行

进入 **Playground**,选择刚创建的 Prompt,填入变量值,选择模型提供方
(`mock` 无需 API Key,可离线试跑),点击「运行」即可看到渲染后的 Prompt
与模型输出。

## 用 SDK 取用

```js
import { PromptOpsClient } from '@promptops/client'

const client = new PromptOpsClient({
  server: 'http://localhost:8080',
  namespace: 'prod',
})

const text = await client.render('sql.generator', { question: '查询所有用户' })

client.on('update', (e) => console.log('prompt 已热更新:', e.key))
client.watch()
```

更多见 [SDK 文档](/zh/sdk/node) 与 [API 参考](/zh/api)。
