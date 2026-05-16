# API 参考

所有接口以 `/api` 为前缀。除 `/api/login`、`/health` 外,均需在请求头携带
`Authorization: Bearer <token>`。`/ws` WebSocket 也需要 token —— 由于浏览器无法
在 WebSocket 连接上设置请求头,token 以握手时的 `?token=` 查询参数传入。

权限:写操作需 `editor` 及以上角色,用户管理需 `admin`;`viewer` 仅可读与运行。

## 认证

| 方法 | 路径 | 说明 |
|---|---|---|
| POST | `/api/login` | 登录(用户名 / 密码,或静态 Token) |
| POST | `/api/logout` | 注销当前会话 |
| GET | `/api/me` | 当前用户与角色 |

## Prompt

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/api/prompts` | 列表 / 搜索(`q`、`env`、`category`、`tag`、`workspace`) |
| POST | `/api/prompts` | 创建 |
| GET / PUT / DELETE | `/api/prompts/:id` | 获取 / 更新 / 删除 |
| GET | `/api/prompts/:id/versions` | 版本历史 |
| POST | `/api/prompts/publish` | 发布(快照当前内容为版本) |
| POST | `/api/prompts/rollback` | 回滚到指定版本 |
| GET / PUT / DELETE | `/api/prompts/:id/rollout` | 灰度发布(AB)配置 |
| GET | `/api/sdk/prompts/:key?env=` | SDK 运行时按 key 获取 Prompt |

## Playground / Agent / Workflow

| 方法 | 路径 | 说明 |
|---|---|---|
| POST | `/api/playground/run` | 渲染 Prompt 并调用模型提供方 |
| GET | `/api/playground/providers` | 列出可用模型提供方 |
| GET / POST / PUT / DELETE | `/api/agents` `/api/agents/:id` | Agent 增删改查 |
| POST | `/api/agents/:id/run` | 运行 Agent |
| GET / POST / PUT / DELETE | `/api/workflows` `/api/workflows/:id` | Workflow 增删改查 |
| POST | `/api/workflows/:id/run` | 运行 Workflow,返回逐步轨迹 |

## 治理与观测

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/api/audit` | 审计日志(分页:`?limit=` `?offset=`) |
| GET | `/api/runs` | 运行日志(分页:`?limit=` `?offset=`) |
| GET | `/api/runs/stats` | 运行指标聚合 |
| GET | `/api/clients` | 在线热更新连接(SDK 与浏览器) |
| GET / POST / PUT / DELETE | `/api/users` `/api/users/:id` | 用户管理(仅 admin) |
| GET / POST / DELETE | `/api/workspaces` `/api/workspaces/:id` | 团队工作区 |

## 实时

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/ws` | WebSocket 热更新事件流(通过 `?token=` 鉴权) |
