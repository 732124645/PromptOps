# PromptOps

[English](README.md) | 简体中文

面向 AI 应用的开源 Prompt、Agent、Workflow 运行时平台。

PromptOps 是一个面向 AI 应用的 **Prompt Runtime 平台**。它将 Prompt、Agent、
Workflow 从代码中解耦出来,让 AI 应用支持 Prompt 的热更新、版本管理、发布、
灰度与运行时观测。

核心能力:**Build → Version → Deploy → Runtime → Observe** Prompt。

## 当前进度

MVP、SDK / 热更新、Playground、Agent / Workflow 运行时、可观测性与审计均已完成。

| 模块 | 技术栈 | 说明 |
|---|---|---|
| 后端 | Go + Gin + GORM + SQLite | Prompt CRUD / 搜索 / 版本管理 / 多环境 / SDK 接口 |
| 热更新 | WebSocket (`/ws`) | Prompt 变更实时推送给已连接客户端 |
| 前端 | Vue3 + Vite + Naive UI + Pinia | Prompt / Agent / Workflow / Playground / 观测 页面 |
| SDK | Node(依赖 `ws`)/ Python / Java(零依赖) | 按 key 获取 Prompt、模板渲染、WebSocket 热更新 |
| Playground | 模型网关(mock / OpenAI / Claude / Ollama / Gemini) | 填变量、调用模型、查看结果;版本 Diff 对比 |
| Agent | 配置化 Agent(Prompt + 提供方 + 模型) | 保存可复用配置并一键运行 |
| Workflow | 步骤引擎(render → model → transform) | 编排多步流程,串联输出,查看逐步轨迹 |
| 可观测性 | 审计日志 + 运行日志 + Token 统计 + 在线客户端 | 记录所有变更与模型调用,聚合运行指标 |
| 权限 (RBAC) | 用户 / 角色 / 会话(PBKDF2 口令) | admin / editor / viewer 三级角色,按角色控制接口 |
| 灰度发布 | Rollout(按 key + 环境的 AB 流量切分) | SDK 获取 Prompt 时按权重返回两个版本之一 |
| 团队空间 | Workspace | Prompt / Agent / Workflow 按工作区归组,前端可切换 |

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

默认登录账号:`admin` / `admin`(首次启动自动创建)。静态 Token
`promptops-dev-token` 仍可作为管理员凭据,供 SDK 与脚本使用。

## Docker 部署

```bash
docker compose up --build   # 构建前端 + 后端,访问 http://localhost:8080
```

## API 概览

| 方法 | 路径 | 说明 |
|---|---|---|
| POST | `/api/login` | 登录(用户名 / 密码,或静态 Token) |
| GET | `/api/me` | 当前用户与角色 |
| GET/POST/PUT/DELETE | `/api/users` `/api/users/:id` | 用户管理(仅 admin) |
| GET/POST/DELETE | `/api/workspaces` `/api/workspaces/:id` | 团队工作区管理 |
| GET | `/api/prompts` | 列表 / 搜索(`q`、`env`、`category`、`tag`) |
| POST | `/api/prompts` | 创建 |
| GET/PUT/DELETE | `/api/prompts/:id` | 获取 / 更新 / 删除 |
| GET | `/api/prompts/:id/versions` | 版本历史 |
| POST | `/api/prompts/publish` | 发布(快照当前内容为版本) |
| POST | `/api/prompts/rollback` | 回滚到指定版本 |
| GET/PUT/DELETE | `/api/prompts/:id/rollout` | 灰度发布(AB)配置 |
| GET | `/api/sdk/prompts/:key?env=` | SDK 运行时按 key 获取 Prompt |
| POST | `/api/playground/run` | 渲染 Prompt 并调用模型提供方 |
| GET | `/api/playground/providers` | 列出可用模型提供方 |
| GET/POST/PUT/DELETE | `/api/agents` `/api/agents/:id` | Agent 增删改查 |
| POST | `/api/agents/:id/run` | 运行 Agent |
| GET/POST/PUT/DELETE | `/api/workflows` `/api/workflows/:id` | Workflow 增删改查 |
| POST | `/api/workflows/:id/run` | 运行 Workflow,返回逐步轨迹 |
| GET | `/api/audit` | 审计日志(变更记录) |
| GET | `/api/runs` | 运行日志(模型调用记录) |
| GET | `/api/runs/stats` | 运行指标聚合(次数 / Token / 延迟) |
| GET | `/api/clients` | 在线热更新连接(SDK 与浏览器) |
| GET | `/ws` | WebSocket 热更新事件流 |

除 `/api/login`、`/ws`、`/health` 外,所有 `/api/*` 需要
`Authorization: Bearer <token>`。写操作需 `editor` 及以上角色,用户管理需
`admin` 角色;`viewer` 仅可读与运行。

## SDK

PromptOps 提供轻量的运行时 SDK,支持按 key 获取 Prompt、`{{变量}}` 模板
渲染,以及通过 WebSocket 的热更新(`watch()` —— Prompt 在后端变更后自动刷新
本地缓存)。

- `sdk/node` —— Node.js(Node 18+,仅依赖 `ws`;使用内置 `fetch`)
- `sdk/python` —— Python(零依赖,纯标准库,含最小 WebSocket 客户端)
- `sdk/java` —— Java(零依赖,纯 JDK,使用 `java.net.http`)

Node 示例:

```js
import { PromptOpsClient } from '@promptops/client'

const client = new PromptOpsClient({ server: 'http://localhost:8080', namespace: 'prod' })

const text = await client.render('sql.generator', { question: '查询所有用户' })

client.on('update', (e) => console.log('prompt 已热更新:', e.key))
client.watch()
```

## CI

`.github/workflows/ci.yml` 在 GitHub runner 上验证:后端(`go build`、`go test`、
冒烟测试)、前端(`npm run build`)、三个 SDK 的单元测试,以及一个集成任务 ——
启动真实服务端并跑通各 SDK 的「获取 → 渲染 → WebSocket 热更新」全链路。

## 文档

在线文档站(VitePress,由 GitHub Pages 托管):<https://732124645.github.io/PromptOps/>
—— 站点源码位于 `website/`,经 `.github/workflows/docs.yml` 自动构建发布。

完整的项目设计文档见 [docs/PromptOps.md](docs/PromptOps.md),涵盖项目定位、
核心概念、MVP 范围、技术架构、数据库与 API 设计、SDK、热更新方案、部署以及
路线图规划。
