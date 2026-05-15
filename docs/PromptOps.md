# PromptOps

## Open-source runtime platform for AI prompts, agents, and workflows.

---

# 1. 项目介绍

PromptOps 是一个：

```txt
面向 AI 应用的 Prompt Runtime 平台。
```

它不是普通的 Prompt 管理工具。

PromptOps 的目标是：

# 将 Prompt、Agent、Workflow 从代码中解耦出来。

让 AI 应用：

- 支持 Prompt 热更新
- 支持 Prompt 版本管理
- 支持 Prompt 发布
- 支持 Prompt 灰度
- 支持 Prompt Runtime
- 支持 Agent 配置化
- 支持 AI Workflow 配置化

---

# 2. 为什么做这个项目

目前 AI 项目存在很多问题：

| 问题 | 描述 |
|---|---|
| Prompt 写死代码 | 修改 Prompt 需要重新发布 |
| Prompt 无法版本管理 | 无法回滚 |
| Prompt 无法热更新 | AI 服务必须重启 |
| Prompt 混乱 | 多团队难协作 |
| 无法灰度发布 | Prompt 风险高 |
| Prompt 无法观测 | 无法知道效果 |

---

# 3. PromptOps 要解决什么

PromptOps 的核心目标：

```txt
让 Prompt 成为真正可工程化管理的资源。
```

包括：

- Build Prompt
- Version Prompt
- Deploy Prompt
- Runtime Prompt
- Observe Prompt

---

# 4. 项目定位

PromptOps 并不是：

- AI Chat
- AI IDE
- AI 助手

而是：

# AI Runtime Infrastructure

它更偏：

- Infra
- Runtime
- Platform
- PromptOps

---

# 5. 核心概念

# 5.1 Prompt

Prompt 是系统核心资源。

例如：

```yaml
key: code.review
version: v3
model: gpt-4o
```

---

# 5.2 Agent

Agent：

- Prompt
- Tools
- Memory
- Model

组合而成。

---

# 5.3 Workflow

Workflow：

```txt
Prompt → Tool → Prompt → Model
```

形成完整 AI 流程。

---

# 6. 第一阶段目标（MVP）

第一阶段不要做太复杂。

重点：

# 先做 Prompt Runtime。

---

# MVP 功能

| 功能 | 是否必须 |
|---|---|
| Prompt CRUD | ✅ |
| Prompt 分类 | ✅ |
| Prompt 搜索 | ✅ |
| Prompt 版本管理 | ✅ |
| 多环境(dev/test/prod) | ✅ |
| Prompt SDK 获取 | ✅ |
| Prompt 热更新 | ✅ |
| Web UI | ✅ |
| SQLite 存储 | ✅ |

---

# 7. MVP 页面设计

# 7.1 登录页

功能：

- 登录
- Token 登录

---

# 7.2 Prompt 列表页

功能：

- Prompt 搜索
- 标签
- 分类
- 环境切换

---

# 7.3 Prompt 编辑页

功能：

- Monaco Editor
- Prompt 编辑
- 变量高亮
- JSON 配置

---

# 7.4 Prompt Playground

功能：

- 输入变量
- 调用模型
- 查看结果

支持：

- OpenAI
- Ollama
- Claude
- Gemini

---

# 7.5 Prompt 版本页

功能：

- 历史版本
- Diff
- 回滚

---

# 8. Prompt Runtime 设计

# 8.1 Prompt Key

类似：

```txt
code.review
sql.generator
agent.customer.service
workflow.report.summary
```

---

# 8.2 Prompt 版本

支持：

```txt
v1
v2
v3
```

---

# 8.3 Prompt 环境

支持：

| 环境 |
|---|
| dev |
| test |
| prod |

---

# 8.4 Prompt 热更新

这是核心能力。

AI 服务：

```txt
无需重启
```

即可更新 Prompt。

---

# 9. 技术架构

# 前端

推荐：

| 技术 | 用途 |
|---|---|
| Vue3 | 前端框架 |
| Vite | 构建工具 |
| TypeScript | 类型系统 |
| Naive UI | UI组件 |
| Pinia | 状态管理 |
| Monaco Editor | Prompt 编辑器 |

---

# 后端

推荐：

| 技术 | 用途 |
|---|---|
| Go | 核心服务 |
| Gin | Web框架 |
| GORM | ORM |
| SQLite | 数据库存储 |
| WebSocket | Prompt 热更新 |

---

# 为什么选择 Go

因为：

PromptOps 更偏：

- Infra
- Runtime
- Platform

Go 更适合：

- 高并发
- 单文件部署
- SDK 开发
- WebSocket
- 后期 Agent Runtime

---

# 10. 数据库设计

# prompt

| 字段 | 类型 |
|---|---|
| id | varchar |
| key | varchar |
| name | varchar |
| content | text |
| version | varchar |
| env | varchar |
| tags | text |
| model | varchar |
| created_at | datetime |

---

# prompt_version

| 字段 | 类型 |
|---|---|
| id | varchar |
| prompt_id | varchar |
| version | varchar |
| content | text |
| created_at | datetime |

---

# 11. API 设计

# 获取 Prompt

```http
GET /api/prompts/{key}
```

返回：

```json
{
  "key": "sql.generator",
  "version": "v3",
  "content": "你是一个SQL专家..."
}
```

---

# 创建 Prompt

```http
POST /api/prompts
```

---

# 发布 Prompt

```http
POST /api/prompts/publish
```

---

# 回滚 Prompt

```http
POST /api/prompts/rollback
```

---

# 12. SDK 设计

# Node SDK

```ts
const client = new PromptOpsClient({
  server: "http://localhost:8080",
  namespace: "prod"
})

const prompt = await client.getPrompt("sql.generator")
```

---

# Java SDK

```java
PromptOpsClient client = new PromptOpsClient();

Prompt prompt = client.getPrompt("sql.generator");
```

---

# 13. Prompt 模板变量

支持：

```txt
{{name}}
{{code}}
{{question}}
```

例如：

```txt
你是一个 {{language}} 专家。

请分析下面代码：

{{code}}
```

---

# 14. Prompt 热更新方案

推荐：

# WebSocket

客户端：

- 启动连接
- 监听 Prompt 更新
- 自动刷新本地缓存

---

# 15. 项目目录结构

```txt
promptops/
├── web/
├── server/
├── sdk/
│   ├── java/
│   ├── node/
│   └── python/
├── docs/
├── docker/
└── scripts/
```

---

# 16. Docker 部署

# docker-compose.yml

```yaml
version: '3'

services:
  promptops:
    build: .
    ports:
      - "8080:8080"
    volumes:
      - ./data:/app/data
```

---

# 17. 第二阶段规划

第二阶段增加：

| 功能 | 描述 |
|---|---|
| Prompt 灰度发布 | AB Test |
| Prompt 审计日志 | 修改记录 |
| Team Workspace | 团队空间 |
| Prompt 权限 | RBAC |
| Prompt Marketplace | Prompt市场 |
| Workflow Runtime | 工作流 |
| Agent Runtime | Agent |
| Prompt Observability | 可观测性 |

---

# 18. 第三阶段规划

第三阶段：

# AI Runtime Platform

增加：

- Model Gateway
- Agent Runtime
- Workflow Engine
- Tool Registry
- MCP Support
- Token Metrics
- AI Gateway

---

# 19. UI 风格建议

推荐风格：

```txt
Infra + Modern + Minimal
```

参考：

- Vercel
- Supabase
- Langfuse
- Linear
- Raycast

---

# 20. GitHub 信息

# 项目名

```txt
PromptOps
```

---

# GitHub 仓库

```txt
promptops
```

---

# GitHub 简介

```txt
Open-source runtime platform for AI prompts, agents, and workflows.
```

---

# GitHub Topics

```txt
ai
prompt
llm
agent
runtime
workflow
promptops
golang
vue
```

---

# 21. Claude Code 开发建议

# 第一阶段

Claude Code 先完成：

- Go 后端
- SQLite
- Prompt CRUD
- Vue3 UI
- Prompt 列表页
- Prompt 编辑页

---

# 第二阶段

完成：

- SDK
- WebSocket
- Prompt 热更新

---

# 第三阶段

完成：

- Playground
- AI 测试
- Prompt 版本管理

---

# 22. 最终目标

PromptOps 最终目标：

```txt
Build and operate AI behavior in production.
```
