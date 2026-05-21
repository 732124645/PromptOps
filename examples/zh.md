# PromptOps 示例（中文）

> The English index is at [`README.md`](README.md).

`examples/` 目录下的工程都是端到端可运行的最小示例，专门用来展示
PromptOps 在真实业务里是怎么用的。每个示例都假设你已经在本地启动了
PromptOps（见仓库根 `README.md` 的 Quick Start 或 Docker 部署章节）。

## 当前示例

| 目录 | 语言 | 说明 |
|---|---|---|
| [`node-hot-reload`](node-hot-reload) | Node.js | 最小消费者：拉取一条 prompt 并通过 WebSocket 订阅热更新；在 Web UI 修改并发布之后，终端立即反映新内容，无需重启 |
| [`sql-agent`](sql-agent) | Node.js | 自然语言 → SQL 智能体：从 PromptOps 拉 `sql.generator` 模板，调用 Playground 的 `mock` 提供方生成 SQL；通过修改 prompt 就能在运行时改变 SQL 风格规则 |

## 欢迎贡献的示例

如果你想贡献新的示例，可以从下面这些方向选：

- `python-hot-reload` —— Python 版本的最小热更新示例
- `java-springboot-hot-reload` —— Spring Boot 集成示例
- `customer-service-agent` —— 多版本话术 + 灰度发布的客服场景
- `workflow-demo` —— 多步 Workflow 编排（render → model → transform）的演示
- `code-review-agent` —— Prompt 驱动的代码评审 agent

贡献流程见 [`../CONTRIBUTING.md`](../CONTRIBUTING.md)。
