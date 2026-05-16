---
layout: home

hero:
  name: PromptOps
  text: AI Prompt Runtime 平台
  tagline: 把 Prompt、Agent、Workflow 从代码中解耦 —— 支持热更新、版本管理、灰度发布与运行时观测。
  actions:
    - theme: brand
      text: 快速上手
      link: /guide/quickstart
    - theme: alt
      text: 项目介绍
      link: /guide/introduction
    - theme: alt
      text: GitHub
      link: https://github.com/732124645/PromptOps

features:
  - title: Prompt Runtime
    details: Prompt 作为可工程化管理的资源 —— CRUD、搜索、分类、多环境(dev / test / prod)与版本管理。
  - title: 热更新
    details: 基于 WebSocket,Prompt 变更实时推送到已连接的 SDK,AI 服务无需重启即可生效。
  - title: 三语言 SDK
    details: Node、Python、Java 三个零依赖 SDK,按 key 获取 Prompt、渲染模板变量、订阅热更新。
  - title: Playground
    details: 填入变量、调用模型(OpenAI / Claude / Ollama / Gemini / mock)、即时查看渲染结果与输出。
  - title: Agent 与 Workflow
    details: 配置化 Agent(Prompt + 提供方 + 模型);Workflow 步骤引擎串联 render → model → transform。
  - title: 灰度发布
    details: 按 key + 环境的 AB 流量切分,SDK 获取 Prompt 时按权重返回两个版本之一。
  - title: 可观测性与审计
    details: 记录所有变更与模型调用,聚合运行次数、Token 估算与延迟指标。
  - title: 权限与团队空间
    details: admin / editor / viewer 三级 RBAC;资源按 Workspace 团队空间归组。
---
