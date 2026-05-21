# `sql.generator` — starter prompt

Copy the **Content** block below into a new prompt in the PromptOps Web UI:

- **Key**: `sql.generator`
- **Env**: `prod`
- **Model**: `mock-1` (or any model — Playground only uses this as metadata)

## Content

```
You are a senior data engineer. Generate a single {{dialect}} SQL query that
answers the user's question. Follow these rules:

1. Return only the SQL. Do not explain it.
2. Use snake_case for all aliases.
3. Always include a LIMIT clause when the result could be unbounded.
4. Prefer explicit JOINs over subqueries when both are valid.

User question:
{{question}}
```

## Variables expected by the example

| Variable | Source | Example |
|---|---|---|
| `dialect` | `DIALECT` env var (default `mysql`) | `mysql` |
| `question` | CLI argument | `"Top 10 customers by revenue this month"` |

## Try a runtime change

After the example is running, edit this prompt and click **Publish** to see
the SQL style change without restarting the agent. For instance, swap rule
2 for *"Use CamelCase for all aliases"* and re-run the agent — the next
generation will follow the new rule.
