---
layout: home

hero:
  name: PromptOps
  text: AI Prompt Runtime Platform
  tagline: Decouple prompts, agents and workflows from code — with hot-reload, version control, gray releases and runtime observability.
  actions:
    - theme: brand
      text: Quickstart
      link: /guide/quickstart
    - theme: alt
      text: Introduction
      link: /guide/introduction
    - theme: alt
      text: GitHub
      link: https://github.com/732124645/PromptOps

features:
  - title: Prompt Runtime
    details: Prompts as engineering-managed resources — CRUD, search, categories, multiple environments (dev / test / prod) and version control.
  - title: Hot Reload
    details: Over WebSocket, prompt changes are pushed to connected SDKs in real time — no AI-service restart required.
  - title: Three SDKs
    details: Node, Python and Java SDKs (Python & Java zero-dependency) — fetch prompts by key, render template variables, subscribe to hot-reload.
  - title: Playground
    details: Fill in variables, call a model (OpenAI / Claude / Ollama / Gemini / mock) and instantly see the rendered prompt and output.
  - title: Agents & Workflows
    details: Config-driven agents (prompt + provider + model); a workflow engine chaining render → model → transform steps.
  - title: Gray Release
    details: Weighted A/B traffic split per key + environment — the SDK returns one of two versions by weight.
  - title: Observability & Audit
    details: Every mutation and model call is logged; run counts, token estimates and latency are aggregated.
  - title: Roles & Workspaces
    details: Three-tier RBAC (admin / editor / viewer); resources grouped into team workspaces.
---
