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
  NModal,
  NTable,
  NEmpty,
  useMessage,
} from 'naive-ui'
import { api, type Prompt, type PromptVersion } from '../api/client'
import { lineDiff, type DiffLine } from '../utils/diff'

const route = useRoute()
const router = useRouter()
const message = useMessage()

const isNew = computed(() => !route.params.id)
const id = ref<string>((route.params.id as string) || '')
const saving = ref(false)

const form = ref<Partial<Prompt>>({
  key: '',
  name: '',
  content: '',
  version: 'v1',
  env: 'dev',
  category: '',
  tags: '',
  model: 'gpt-4o',
})

const envOptions = [
  { label: 'dev', value: 'dev' },
  { label: 'test', value: 'test' },
  { label: 'prod', value: 'prod' },
]

const versions = ref<PromptVersion[]>([])
const showVersions = ref(false)

const showDiff = ref(false)
const diffVersion = ref('')
const diffLines = ref<DiffLine[]>([])

function openDiff(v: PromptVersion) {
  diffVersion.value = v.version
  diffLines.value = lineDiff(v.content, form.value.content || '')
  showDiff.value = true
}

const contentPlaceholder = '你是一个 {{language}} 专家。\n请分析下面代码:\n{{code}}'

// Detected {{variable}} placeholders in the prompt content.
const variables = computed(() => {
  const matches = (form.value.content || '').match(/{{\s*[\w.]+\s*}}/g) || []
  return [...new Set(matches.map((m) => m.replace(/[{}]/g, '').trim()))]
})

// Built in script so the template never contains a literal "}}" inside an
// interpolation, which the Vue template parser would close prematurely.
function varLabel(name: string): string {
  return '{{' + name + '}}'
}

async function load() {
  if (isNew.value) return
  try {
    const { data } = await api.get(id.value)
    form.value = data.data
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
      const { data } = await api.create(form.value)
      id.value = data.data.id
      message.success('已创建')
      router.replace({ name: 'prompt-edit', params: { id: id.value } })
    } else {
      await api.update(id.value, form.value)
      message.success('已保存')
    }
  } catch {
    message.error('保存失败')
  } finally {
    saving.value = false
  }
}

async function publish() {
  if (isNew.value) {
    message.warning('请先保存 Prompt')
    return
  }
  await api.publish(id.value)
  message.success(`已发布版本 ${form.value.version}`)
}

async function openVersions() {
  if (isNew.value) {
    message.warning('请先保存 Prompt')
    return
  }
  const { data } = await api.versions(id.value)
  versions.value = data.data
  showVersions.value = true
}

async function rollback(v: PromptVersion) {
  await api.rollback(id.value, v.version)
  message.success(`已回滚到 ${v.version}`)
  showVersions.value = false
  load()
}

onMounted(load)
</script>

<template>
  <div class="page">
    <div class="toolbar">
      <h2>{{ isNew ? '新建 Prompt' : `编辑: ${form.key}` }}</h2>
      <n-space>
        <n-button @click="router.push('/')">返回</n-button>
        <n-button v-if="!isNew" @click="openVersions">版本历史</n-button>
        <n-button v-if="!isNew" @click="publish">发布版本</n-button>
        <n-button type="primary" :loading="saving" @click="save">保存</n-button>
      </n-space>
    </div>

    <n-card :bordered="true">
      <n-form>
        <n-grid :cols="2" :x-gap="16">
          <n-grid-item>
            <n-form-item label="Key (例: code.review)">
              <n-input v-model:value="form.key" :disabled="!isNew" placeholder="code.review" />
            </n-form-item>
          </n-grid-item>
          <n-grid-item>
            <n-form-item label="名称">
              <n-input v-model:value="form.name" placeholder="Prompt 名称" />
            </n-form-item>
          </n-grid-item>
          <n-grid-item>
            <n-form-item label="环境">
              <n-select v-model:value="form.env" :options="envOptions" />
            </n-form-item>
          </n-grid-item>
          <n-grid-item>
            <n-form-item label="版本">
              <n-input v-model:value="form.version" placeholder="v1" />
            </n-form-item>
          </n-grid-item>
          <n-grid-item>
            <n-form-item label="分类">
              <n-input v-model:value="form.category" placeholder="例: code" />
            </n-form-item>
          </n-grid-item>
          <n-grid-item>
            <n-form-item label="模型">
              <n-input v-model:value="form.model" placeholder="gpt-4o" />
            </n-form-item>
          </n-grid-item>
          <n-grid-item :span="2">
            <n-form-item label="标签 (逗号分隔)">
              <n-input v-model:value="form.tags" placeholder="review, sql" />
            </n-form-item>
          </n-grid-item>
        </n-grid>

        <n-form-item label="Prompt 内容">
          <n-input
            v-model:value="form.content"
            type="textarea"
            class="mono"
            :autosize="{ minRows: 12, maxRows: 28 }"
            :placeholder="contentPlaceholder"
          />
        </n-form-item>

        <div class="vars">
          <span class="vars-label">检测到的变量:</span>
          <n-space :size="6">
            <n-tag v-for="v in variables" :key="v" size="small" type="warning">{{ varLabel(v) }}</n-tag>
            <span v-if="!variables.length" class="muted">无</span>
          </n-space>
        </div>
      </n-form>
    </n-card>

    <n-modal
      v-model:show="showVersions"
      preset="card"
      title="版本历史"
      style="width: 720px"
    >
      <n-empty v-if="!versions.length" description="暂无已发布版本" />
      <n-table v-else :bordered="false" :single-line="false">
        <thead>
          <tr>
            <th>版本</th>
            <th>环境</th>
            <th>发布时间</th>
            <th>内容预览</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="v in versions" :key="v.id">
            <td><n-tag size="small">{{ v.version }}</n-tag></td>
            <td>{{ v.env }}</td>
            <td>{{ new Date(v.created_at).toLocaleString() }}</td>
            <td class="preview">{{ v.content.slice(0, 60) }}{{ v.content.length > 60 ? '…' : '' }}</td>
            <td>
              <n-space :size="4">
                <n-button size="tiny" @click="openDiff(v)">对比当前</n-button>
                <n-button size="tiny" type="primary" ghost @click="rollback(v)">回滚</n-button>
              </n-space>
            </td>
          </tr>
        </tbody>
      </n-table>
    </n-modal>

    <n-modal
      v-model:show="showDiff"
      preset="card"
      :title="`Diff · ${diffVersion} → 当前`"
      style="width: 760px"
    >
      <div class="diff-legend">
        <span class="del">− 版本 {{ diffVersion }}</span>
        <span class="add">+ 当前内容</span>
      </div>
      <n-empty v-if="!diffLines.length" description="无内容" />
      <pre v-else class="diff"><span v-for="(line, i) in diffLines" :key="i" :class="['diff-line', line.type]">{{ line.type === 'add' ? '+ ' : line.type === 'del' ? '− ' : '  ' }}{{ line.text }}</span></pre>
    </n-modal>
  </div>
</template>

<style scoped>
.toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}
.vars {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-top: 4px;
}
.vars-label {
  font-size: 13px;
  color: #888;
}
.muted {
  color: #888;
  font-size: 13px;
}
.preview {
  font-family: 'SF Mono', Menlo, Consolas, monospace;
  font-size: 12px;
  color: #aaa;
}
.diff-legend {
  display: flex;
  gap: 16px;
  font-size: 12px;
  margin-bottom: 8px;
}
.diff-legend .del {
  color: #e88080;
}
.diff-legend .add {
  color: #63e2b7;
}
.diff {
  margin: 0;
  max-height: 60vh;
  overflow: auto;
  border-radius: 6px;
  background: rgba(255, 255, 255, 0.04);
  font-family: 'SF Mono', Menlo, Consolas, monospace;
  font-size: 12.5px;
}
.diff-line {
  display: block;
  padding: 0 10px;
  white-space: pre-wrap;
  word-break: break-word;
}
.diff-line.add {
  background: rgba(99, 226, 183, 0.14);
  color: #9ff0d3;
}
.diff-line.del {
  background: rgba(232, 128, 128, 0.14);
  color: #f0a8a8;
}
.diff-line.same {
  color: #999;
}
</style>
