# 介绍

PromptOps 是一个面向 AI 应用的 **Prompt Runtime 平台**。它把 Prompt、Agent、
Workflow 从代码中解耦出来,让 AI 应用支持 Prompt 的热更新、版本管理、发布、
灰度与运行时观测。

核心理念:

> 让 Prompt 成为真正可工程化管理的资源 —— Build → Version → Deploy → Runtime → Observe。

## 为什么需要它

| 问题 | 描述 |
|---|---|
| Prompt 写死在代码里 | 改一次 Prompt 就要重新发版 |
| 无法版本管理 | 出问题难以回滚 |
| 无法热更新 | 更新 Prompt 必须重启 AI 服务 |
| 多团队协作混乱 | 缺少归属与权限 |
| 无法灰度 | 新 Prompt 直接全量,风险高 |
| 无法观测 | 不知道调用了多少、效果如何 |

## 能力总览

| 模块 | 说明 |
|---|---|
| Prompt Runtime | CRUD、搜索、分类、多环境、版本发布与回滚 |
| 热更新 | WebSocket 实时推送 Prompt 变更 |
| SDK | Node / Python / Java,零第三方依赖 |
| Playground | 填变量、调模型、看结果;版本 Diff 对比 |
| Agent | 配置化的 Prompt + 提供方 + 模型 |
| Workflow | render → model → transform 步骤编排引擎 |
| 可观测性 | 审计日志、运行日志、Token 与延迟统计 |
| RBAC | admin / editor / viewer 三级角色 |
| 灰度发布 | 按权重的 AB 版本流量切分 |
| 团队空间 | 资源按 Workspace 归组 |

## 技术架构

- **后端**:Go + Gin + GORM + SQLite,WebSocket 热更新。
- **前端**:Vue3 + Vite + Naive UI + Pinia。
- **SDK**:Node(内置 fetch / WebSocket)、Python(标准库)、Java(java.net.http)。
- **部署**:单二进制 / Docker。

下一步:[快速上手](./quickstart) · [部署](./deployment)。
