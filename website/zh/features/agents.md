# Agent

## 是什么

**Agent** 是一个保存下来、可复用的组合 —— Prompt 加上运行它的模型配置:
`Prompt + 提供方 + 模型`。[Playground](./playground) 是草稿纸,Agent 则是
一份有名字、可反复运行的配置。

## 为什么重要

实际工作中你拥有的不只是 Prompt 文本 —— 而是「这个 Prompt,用这个提供方,
配这个模型」。Agent 把这一整个单元固化一次,团队里任何人都能一致地运行它,
不用重新填配置。

## 一个 Agent 由什么组成

| 字段 | 说明 |
|---|---|
| `key` | 稳定标识,例如 `support.agent` |
| `name` | 人类可读的名称 |
| `description` | 这个 Agent 是做什么的 |
| `prompt` | Prompt 文本,可含 `{{变量}}` 占位符 |
| `provider` | `mock` / `openai` / `claude` / `ollama` / `gemini` |
| `model` | 模型名;留空则用提供方默认模型 |

## 怎么用

1. 进入 **Agents** 页面,点击 **新建 Agent**。
2. 填写 key,选择提供方和模型,写好 Prompt。
3. 点击 **保存**。
4. 在该 Agent 页面填好识别出的 `{{变量}}`,若提供方需要则填 API Key,点击
   **运行 Agent**。

结果会展示渲染后的 Prompt 和模型输出,和 Playground 一样 —— 但配置被保存了
下来,下次可直接用。

## 幕后

运行 Agent 会调用 `POST /api/agents/:id/run`。每次运行都会被记录,供
[可观测性](./observability)使用。

## Agent 与 Workflow 的区别

- **Agent** 是*单个* Prompt-模型步骤。
- [**Workflow**](./workflows) 把*多个*步骤串起来,把每一步的输出喂给下一步。

## 下一步

- [Workflow](./workflows) —— 编排多个步骤。
- [可观测性](./observability) —— 查看 Agent 运行记录。
