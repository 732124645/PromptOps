# 权限与工作区

有两件相互独立的事让团队对 PromptOps 的使用井然有序:**角色**决定*一个人
能做什么*,**工作区**决定*资源如何归组*。

## 角色(RBAC)

每个用户拥有三种角色之一:

| 角色 | 能做什么 |
|---|---|
| `viewer` | 读取一切;运行 Prompt、Agent 与 Workflow |
| `editor` | viewer 的全部,外加创建 / 编辑 / 删除 / 发布 |
| `admin` | editor 的全部,外加用户管理和删除工作区 |

一句话:**viewer** 读和运行,**editor** 还能写,**admin** 还能管理。

### 登录

你用两种方式之一进行身份认证:

- **用户名 + 密码** —— 普通用户账号。服务端签发一个会话 Token;SDK 和 UI
  以 `Authorization: Bearer <token>` 形式携带它。
- **静态 Token** —— `PROMPTOPS_TOKEN` 的值(默认 `promptops-dev-token`)
  始终作为 `admin` Token 生效。它对 SDK 和脚本很方便;任何正式部署都应修改它。

默认管理员账号是 `admin` / `admin`,首次启动时创建。

### 管理用户

`admin` 打开 **用户** 页面来创建账号、修改角色或删除用户。角色在该用户的
下一次请求时立即生效。

## 工作区

**工作区**把 Prompt、Agent、Workflow 按团队或项目归组。这些资源每一个都
恰好属于一个工作区。

- UI 顶栏有工作区切换器 —— 选一个,列表就限定到它。
- 始终有一个内置的 `default` 工作区。
- 任何已登录用户都能创建工作区;删除工作区仅限 `admin`。
- 删除工作区**不会**删除其中的资源 —— 它们只是不再归属于它。

工作区关乎*组织方式*(团队 / 项目隔离)。Prompt 上的 `env` 字段关乎
*部署阶段*(`dev` / `test` / `prod`)—— 两者相互独立。

## 对应的 API

| 操作 | 接口 | 所需角色 |
|---|---|---|
| 登录 | `POST /api/login` | — |
| 当前用户 | `GET /api/me` | viewer |
| 用户管理 | `/api/users` `…/:id` | admin |
| 列出 / 创建工作区 | `GET` / `POST /api/workspaces` | viewer / editor |
| 删除工作区 | `DELETE /api/workspaces/:id` | admin |

除 `/api/login` 外所有 `/api/*` 路由都需要 `Authorization` 头;写操作需要
`editor`,用户管理需要 `admin`。

## 下一步

- [可观测性](./observability) —— 审计日志记录了谁改了什么。
