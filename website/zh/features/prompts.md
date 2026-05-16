# Prompt 管理

## 是什么

**Prompt** 是 PromptOps 里最核心的资源 —— 一段有名字、可版本化的文本,你的
AI 应用在运行时按需获取它,而不是把它写死在代码里。

不再把 Prompt 文本写进代码、每次改动都重新发版;而是存在这里,通过 `key`
获取。编辑者改文案,你的应用无需发布即可生效。

## 一个 Prompt 由什么组成

| 字段 | 说明 |
|---|---|
| `key` | SDK 获取时用的稳定标识,例如 `sql.generator` |
| `name` | 人类可读的名称 |
| `content` | Prompt 正文,可含 `{{变量}}` 占位符 |
| `env` | 环境 —— `dev`、`test` 或 `prod` |
| `version` | 当前工作版本标签,例如 `v1` |
| `category` | 自由文本分类,例如 `code` |
| `tags` | 逗号分隔的标签,例如 `review, sql` |
| `model` | 该 Prompt 建议使用的模型,例如 `gpt-4o` |

同一个 `key` 在每个环境下各存在一份 —— `dev` 里的 `sql.generator` 和 `prod`
里的 `sql.generator` 是两个独立的 Prompt。

## 模板变量

Prompt 正文里可以写 `{{变量}}` 占位符,在渲染时被填充 —— 由 SDK 的
`render()`、Playground 或 Agent 完成:

```
你是一个 {{language}} 专家。
请分析下面的代码:
{{code}}
```

编辑器会自动识别正文里所有 `{{变量}}` 并列出来,你随时知道一个 Prompt
需要哪些输入。

## 怎么用

1. 进入 **Prompts** 页面,点击 **新建 Prompt**。
2. 填写 `key`、环境、内容,以及可选的元数据。
3. 点击 **保存**。
4. 之后用搜索框(匹配 key / 名称 / 内容)和环境筛选器查找 Prompt。

## 在应用里获取 Prompt

```js
import { PromptOpsClient } from '@promptops/client'

const client = new PromptOpsClient({ server: 'http://localhost:8080', namespace: 'prod' })

// 用你传入的值填充 {{变量}} 占位符
const text = await client.render('sql.generator', { question: '查询所有用户' })
```

## 下一步

- 保存的 Prompt 是可编辑、可变的。要冻结它,见[版本与发布](./versioning)。
- 要把改动即时推送给运行中的应用,见[热更新](./hot-reload)。
