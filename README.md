# PromptOps

Open-source runtime platform for AI prompts, agents, and workflows.

PromptOps 是一个面向 AI 应用的 **Prompt Runtime 平台**。它将 Prompt、Agent、
Workflow 从代码中解耦出来,让 AI 应用支持 Prompt 的热更新、版本管理、发布、
灰度与运行时观测。

核心能力:**Build → Version → Deploy → Runtime → Observe** Prompt。

## 当前进度

第一阶段 MVP + 第二阶段 SDK / 热更新已完成。

| 模块 | 技术栈 | 说明 |
|---|---|---|
| 后端 | Go + Gin + GORM + SQLite | Prompt CRUD / 搜索 / 版本管理 / 多环境 / SDK 接口 |
| 热更新 | WebSocket (`/ws`) | Prompt 变更实时推送给已连接客户端 |
| 前端 | Vue3 + Vite + Naive UI + Pinia | 登录页 / Prompt 列表页 / Prompt 编辑页 |
| SDK | Node / Python / Java(均零依赖) | 按 key 获取 Prompt、模板渲染、WebSocket 热更新 |

## 项目结构

```txt
PromptOps/
├── server/      # Go 后端 (Gin + GORM + SQLite)
├── web/         # Vue3 前端 (Vite + Naive UI)
├── sdk/         # 运行时 SDK (node / python / java)
├── docs/        # 设计文档
├── scripts/     # 冒烟测试等脚本
├── Dockerfile
└── docker-compose.yml
```

## 本地开发

### 后端

```bash
cd server
go mod tidy
go run .            # 默认监听 :8080,SQLite 文件位于 server/data/promptops.db
```

环境变量:`PROMPTOPS_ADDR`(默认 `:8080`)、`PROMPTOPS_DB`、`PROMPTOPS_TOKEN`
(默认 `promptops-dev-token`)。

### 前端

```bash
cd web
npm install
npm run dev         # http://localhost:5173,已配置代理到后端 :8080
```

默认登录 Token:`promptops-dev-token`。

## Docker 部署

```bash
docker compose up --build   # 构建前端 + 后端,访问 http://localhost:8080
```

## API 概览

| 方法 | 路径 | 说明 |
|---|---|---|
| POST | `/api/login` | Token 登录 |
| GET | `/api/prompts` | 列表 / 搜索(`q`、`env`、`category`、`tag`) |
| POST | `/api/prompts` | 创建 |
| GET/PUT/DELETE | `/api/prompts/:id` | 获取 / 更新 / 删除 |
| GET | `/api/prompts/:id/versions` | 版本历史 |
| POST | `/api/prompts/publish` | 发布(快照当前内容为版本) |
| POST | `/api/prompts/rollback` | 回滚到指定版本 |
| GET | `/api/sdk/prompts/:key?env=` | SDK 运行时按 key 获取 Prompt |
| GET | `/ws` | WebSocket 热更新事件流 |

除 `/api/login`、`/ws`、`/health` 外,所有 `/api/*` 需要
`Authorization: Bearer <token>`。

## SDK

PromptOps 提供**零依赖**的运行时 SDK,支持按 key 获取 Prompt、`{{变量}}` 模板
渲染,以及通过 WebSocket 的热更新(`watch()` —— Prompt 在后端变更后自动刷新
本地缓存)。

- `sdk/node` —— Node.js(Node 22+,使用内置 `fetch` / `WebSocket`)
- `sdk/python` —— Python(纯标准库,含最小 WebSocket 客户端)
- `sdk/java` —— Java(纯 JDK,使用 `java.net.http`)

Node 示例:

```js
import { PromptOpsClient } from '@promptops/client'

const client = new PromptOpsClient({ server: 'http://localhost:8080', namespace: 'prod' })

const text = await client.render('sql.generator', { question: '查询所有用户' })

client.on('update', (e) => console.log('prompt 已热更新:', e.key))
client.watch()
```

## CI

`.github/workflows/ci.yml` 在 GitHub runner 上验证:后端(`go build` + `go test` +
冒烟测试)、前端(`npm run build`)、三个 SDK 的单元测试,以及一个集成任务 ——
启动真实服务端并跑通各 SDK 的「获取 → 渲染 → WebSocket 热更新」全链路。

## Documentation

完整的项目设计文档见 [docs/PromptOps.md](docs/PromptOps.md),涵盖项目定位、
核心概念、MVP 范围、技术架构、数据库与 API 设计、SDK、热更新方案、部署以及
路线图规划。
