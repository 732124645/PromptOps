<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  NButton,
  NInput,
  NSelect,
  NSpace,
  NCard,
  NForm,
  NFormItem,
  NGrid,
  NGridItem,
  NTag,
  NAlert,
  NEmpty,
  useMessage,
} from 'naive-ui'
import { api, type WorkflowStep, type WorkflowRunResult } from '../api/client'

const route = useRoute()
const router = useRouter()
const message = useMessage()

const isNew = computed(() => !route.params.id)
const id = ref<string>((route.params.id as string) || '')
const saving = ref(false)

const form = ref<{ key: string; name: string; description: string; steps: WorkflowStep[] }>({
  key: '',
  name: '',
  description: '',
  steps: [],
})

const stepTypeOptions = [
  { label: 'render — 渲染模板', value: 'render' },
  { label: 'model — 调用模型', value: 'model' },
  { label: 'transform — 文本变换', value: 'transform' },
]
const providerOptions = [
  { label: 'mock', value: 'mock' },
  { label: 'openai', value: 'openai' },
  { label: 'claude', value: 'claude' },
  { label: 'ollama', value: 'ollama' },
  { label: 'gemini', value: 'gemini' },
]
const opOptions = [
  { label: 'upper — 转大写', value: 'upper' },
  { label: 'lower — 转小写', value: 'lower' },
  { label: 'trim — 去首尾空白', value: 'trim' },
]
const templatePlaceholder = '模板内容,可使用 变量占位符,以及 input 变量代表上一步的输出'

function addStep() {
  form.value.steps.push({
    name: `步骤 ${form.value.steps.length + 1}`,
    type: 'render',
    template: '',
    provider: 'mock',
    model: '',
    op: 'upper',
  })
}

function removeStep(i: number) {
  form.value.steps.splice(i, 1)
}

function moveStep(i: number, delta: number) {
  const j = i + delta
  if (j < 0 || j >= form.value.steps.length) return
  const steps = form.value.steps
  const tmp = steps[i]
  steps[i] = steps[j]
  steps[j] = tmp
}

// Variables referenced by render steps, excluding the reserved {{input}}.
const detectedVars = computed(() => {
  const set = new Set<string>()
  for (const step of form.value.steps) {
    if (step.type !== 'render') continue
    const matches = (step.template || '').match(/{{\s*([\w.]+)\s*}}/g) || []
    for (const m of matches) {
      const name = m.replace(/[{}]/g, '').trim()
      if (name !== 'input') set.add(name)
    }
  }
  return [...set]
})

const variables = ref<Record<string, string>>({})
const apiKey = ref('')
const running = ref(false)
const result = ref<WorkflowRunResult | null>(null)
const error = ref('')

const needsKey = computed(() =>
  form.value.steps.some(
    (s) => s.type === 'model' && ['openai', 'claude', 'gemini'].includes(s.provider || ''),
  ),
)

async function load() {
  if (isNew.value) return
  try {
    const { data } = await api.getWorkflow(id.value)
    form.value = {
      key: data.data.key,
      name: data.data.name,
      description: data.data.description,
      steps: data.data.steps,
    }
  } catch {
    message.error('加载失败')
  }
}

async function save() {
  if (!form.value.key) {
    message.warning('Key 不能为空')
    return
  }
  saving.value = true
  try {
    if (isNew.value) {
      const { data } = await api.createWorkflow(form.value)
      id.value = data.data.id
      message.success('已创建')
      router.replace({ name: 'workflow-edit', params: { id: id.value } })
    } else {
      await api.updateWorkflow(id.value, form.value)
      message.success('已保存')
    }
  } catch {
    message.error('保存失败')
  } finally {
    saving.value = false
  }
}

async function run() {
  if (isNew.value) {
    message.warning('请先保存 Workflow')
    return
  }
  running.value = true
  error.value = ''
  result.value = null
  try {
    const { data } = await api.runWorkflow(id.value, {
      variables: variables.value,
      api_key: apiKey.value || undefined,
    })
    result.value = data
  } catch (e) {
    const resp = (e as { response?: { data?: { error?: string } } }).response
    error.value = resp?.data?.error || '运行失败'
  } finally {
    running.value = false
  }
}

onMounted(load)
</script>

<template>
  <div class="page">
    <div class="toolbar">
      <h2>{{ isNew ? '新建 Workflow' : `Workflow: ${form.key}` }}</h2>
      <n-space>
        <n-button @click="router.push('/workflows')">返回</n-button>
        <n-button type="primary" :loading="saving" @click="save">保存</n-button>
      </n-space>
    </div>

    <n-card title="配置" :bordered="true">
      <n-form>
        <n-grid :cols="2" :x-gap="16">
          <n-grid-item>
            <n-form-item label="Key">
              <n-input v-model:value="form.key" :disabled="!isNew" placeholder="summarize.flow" />
            </n-form-item>
          </n-grid-item>
          <n-grid-item>
            <n-form-item label="名称">
              <n-input v-model:value="form.name" placeholder="Workflow 名称" />
            </n-form-item>
          </n-grid-item>
          <n-grid-item :span="2">
            <n-form-item label="描述">
              <n-input v-model:value="form.description" placeholder="这个 Workflow 做什么" />
            </n-form-item>
          </n-grid-item>
        </n-grid>
      </n-form>
    </n-card>

    <n-card title="步骤" :bordered="true" style="margin-top: 16px">
      <template #header-extra>
        <n-button size="small" @click="addStep">+ 添加步骤</n-button>
      </template>

      <n-empty v-if="!form.steps.length" description="还没有步骤,点击右上角添加" />
      <n-card
        v-for="(step, i) in form.steps"
        :key="i"
        size="small"
        class="step-card"
        :bordered="true"
      >
        <div class="step-head">
          <n-tag type="info" size="small">#{{ i + 1 }}</n-tag>
          <n-input
            v-model:value="step.name"
            size="small"
            placeholder="步骤名"
            style="max-width: 180px"
          />
          <n-select
            v-model:value="step.type"
            :options="stepTypeOptions"
            size="small"
            style="width: 190px"
          />
          <div class="grow" />
          <n-button size="tiny" :disabled="i === 0" @click="moveStep(i, -1)">↑</n-button>
          <n-button size="tiny" :disabled="i === form.steps.length - 1" @click="moveStep(i, 1)">
            ↓
          </n-button>
          <n-button size="tiny" type="error" ghost @click="removeStep(i)">删除</n-button>
        </div>
        <div class="step-body">
          <n-input
            v-if="step.type === 'render'"
            v-model:value="step.template"
            type="textarea"
            class="mono"
            :autosize="{ minRows: 3, maxRows: 12 }"
            :placeholder="templatePlaceholder"
          />
          <n-space v-else-if="step.type === 'model'">
            <n-select
              v-model:value="step.provider"
              :options="providerOptions"
              style="width: 150px"
            />
            <n-input v-model:value="step.model" placeholder="模型(可留空)" style="width: 220px" />
          </n-space>
          <n-select v-else v-model:value="step.op" :options="opOptions" style="width: 220px" />
        </div>
      </n-card>
    </n-card>

    <n-card title="运行" :bordered="true" style="margin-top: 16px">
      <p v-if="isNew" class="muted">请先保存 Workflow 再运行。</p>
      <template v-else>
        <n-form>
          <n-form-item v-if="detectedVars.length" label="变量">
            <n-space vertical style="width: 100%">
              <div v-for="v in detectedVars" :key="v" class="var-row">
                <n-tag type="warning" size="small">{{ v }}</n-tag>
                <n-input v-model:value="variables[v]" :placeholder="`${v} 的值`" />
              </div>
            </n-space>
          </n-form-item>
          <n-form-item v-if="needsKey" label="API Key">
            <n-input
              v-model:value="apiKey"
              type="password"
              show-password-on="click"
              placeholder="含真实模型步骤,需要 API Key"
            />
          </n-form-item>
          <n-button type="primary" :loading="running" @click="run">运行 Workflow</n-button>
        </n-form>

        <n-alert v-if="error" type="error" title="运行失败" style="margin-top: 14px">
          {{ error }}
        </n-alert>
        <template v-if="result">
          <div class="result-label">步骤轨迹</div>
          <div v-for="(s, i) in result.steps" :key="i" class="trace">
            <div class="trace-head">
              <n-tag size="small">#{{ i + 1 }}</n-tag>
              <span class="trace-name">{{ s.name }}</span>
              <n-tag size="small" type="info">{{ s.type }}</n-tag>
            </div>
            <pre class="block">{{ s.output }}</pre>
          </div>
          <div class="result-label">最终输出</div>
          <pre class="block final">{{ result.output }}</pre>
        </template>
      </template>
    </n-card>
  </div>
</template>

<style scoped>
.toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}
.muted {
  color: #888;
  font-size: 13px;
}
.step-card {
  margin-bottom: 10px;
}
.step-head {
  display: flex;
  align-items: center;
  gap: 8px;
}
.step-body {
  margin-top: 10px;
}
.grow {
  flex: 1;
}
.var-row {
  display: flex;
  align-items: center;
  gap: 8px;
}
.result-label {
  font-size: 13px;
  color: #888;
  margin: 14px 0 6px;
}
.trace {
  margin-bottom: 8px;
}
.trace-head {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 4px;
}
.trace-name {
  font-size: 13px;
}
.block {
  margin: 0;
  padding: 12px;
  border-radius: 6px;
  background: rgba(255, 255, 255, 0.05);
  font-family: 'SF Mono', Menlo, Consolas, monospace;
  font-size: 13px;
  white-space: pre-wrap;
  word-break: break-word;
}
.block.final {
  background: rgba(99, 226, 183, 0.12);
}
</style>
