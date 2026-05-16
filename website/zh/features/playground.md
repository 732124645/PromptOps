# Playground

## 是什么

**Playground** 是一个交互式控制台,用来拿一个 Prompt 去试真实(或 mock)
模型 —— 不用写任何代码。填好变量、选好提供方、运行,就能同时看到渲染后的
Prompt 和模型输出。

## 为什么重要

它打通了「编辑 Prompt」和「知道它好不好用」之间的闭环。你可以反复打磨文案,
即时看到效果,然后再发布。

## 模型提供方

| 提供方 | 需要什么 |
|---|---|
| `mock` | 什么都不要 —— 离线运行,返回确定性的回显。适合测试。 |
| `openai` | 一个 API Key |
| `claude` | 一个 API Key |
| `gemini` | 一个 API Key |
| `ollama` | 一个 Base URL(例如 `http://localhost:11434`)指向本地 Ollama |

API Key 在每次运行时填写,只用于那一次调用 —— 不会被存储。

## 怎么用

1. 打开 **Playground**。
2. 可选:挑一个已有 Prompt 载入其内容,或直接输入内容。
3. Playground 会识别每个 `{{变量}}` 并为其生成输入框 —— 逐个填好。
4. 选择 **模型提供方** 和模型。如果只想看渲染结果、不想用 API Key,先用 `mock`。
5. 点击 **运行**。

结果面板显示:

- **渲染后的 Prompt** —— 你的内容,所有 `{{变量}}` 已替换。
- **模型输出** —— 提供方返回的内容,以及所用的提供方和模型。

## 幕后

Playground 调用 `POST /api/playground/run`。每次运行也会被记录,供
[可观测性](./observability)使用 —— 之后能在运行日志里看到。

## 下一步

- 要把「Prompt + 提供方 + 模型」存成一个可复用、可一键运行的单元,见
  [Agent](./agents)。
- 要把多个步骤串起来,见 [Workflow](./workflows)。
