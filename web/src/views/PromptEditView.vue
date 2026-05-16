<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import {
  NButton,
  NInput,
  NInputNumber,
  NSelect,
  NSpace,
  NCard,
  NForm,
  NFormItem,
  NGrid,
  NGridItem,
  NTag,
  NModal,
  NSwitch,
  NTable,
  NEmpty,
  useMessage,
} from 'naive-ui'
import { api, type Prompt, type PromptVersion } from '../api/client'
import { lineDiff, type DiffLine } from '../utils/diff'
import { useWorkspaceStore } from '../stores/workspace'

const route = useRoute()
const router = useRouter()
const message = useMessage()
const workspace = useWorkspaceStore()
const { t, locale } = useI18n()

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

// Example text kept out of the i18n catalog: vue-i18n would try to compile
// the {{variable}} braces as nested placeholders and throw.
const contentPlaceholder = computed(() =>
  locale.value === 'en'
    ? 'You are a {{language}} expert.\nPlease review the code below:\n{{code}}'
    : '你是一个 {{language}} 专家。\n请分析下面代码:\n{{code}}',
)

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

const rollout = ref({ enabled: false, variant_a: '', variant_b: '', weight_a: 50 })

async function load() {
  if (isNew.value) return
  try {
    const { data } = await api.get(id.value)
    form.value = data.data
    const r = await api.getRollout(id.value)
    if (r.data.data) {
      rollout.value = {
        enabled: r.data.data.enabled,
        variant_a: r.data.data.variant_a,
        variant_b: r.data.data.variant_b,
        weight_a: r.data.data.weight_a,
      }
    }
  } catch {
    message.error(t('common.loadFailed'))
  }
}

async function saveRollout() {
  if (isNew.value) {
    message.warning(t('common.saveFirst', { name: t('entity.prompt') }))
    return
  }
  try {
    await api.setRollout(id.value, rollout.value)
    message.success(t('prompt.rolloutSaved'))
  } catch {
    message.error(t('prompt.rolloutSaveFailed'))
  }
}

async function save() {
  if (!form.value.key) {
    message.warning(t('common.keyRequired'))
    return
  }
  saving.value = true
  try {
    if (isNew.value) {
      const { data } = await api.create({
        ...form.value,
        workspace_id: workspace.currentId,
      })
      id.value = data.data.id
      message.success(t('common.created'))
      router.replace({ name: 'prompt-edit', params: { id: id.value } })
    } else {
      await api.update(id.value, form.value)
      message.success(t('common.saved'))
    }
  } catch {
    message.error(t('common.saveFailed'))
  } finally {
    saving.value = false
  }
}

async function publish() {
  if (isNew.value) {
    message.warning(t('common.saveFirst', { name: t('entity.prompt') }))
    return
  }
  await api.publish(id.value)
  message.success(t('prompt.published', { version: form.value.version }))
}

async function openVersions() {
  if (isNew.value) {
    message.warning(t('common.saveFirst', { name: t('entity.prompt') }))
    return
  }
  const { data } = await api.versions(id.value)
  versions.value = data.data
  showVersions.value = true
}

async function rollback(v: PromptVersion) {
  await api.rollback(id.value, v.version)
  message.success(t('prompt.rolledBack', { version: v.version }))
  showVersions.value = false
  load()
}

onMounted(load)
</script>

<template>
  <div class="page">
    <div class="toolbar">
      <h2>{{ isNew ? t('prompt.newTitle') : t('prompt.editTitle', { key: form.key }) }}</h2>
      <n-space>
        <n-button @click="router.push('/')">{{ t('common.back') }}</n-button>
        <n-button v-if="!isNew" @click="openVersions">{{ t('prompt.versionHistory') }}</n-button>
        <n-button v-if="!isNew" @click="publish">{{ t('prompt.publish') }}</n-button>
        <n-button type="primary" :loading="saving" @click="save">{{ t('common.save') }}</n-button>
      </n-space>
    </div>

    <n-card :bordered="true">
      <n-form>
        <n-grid :cols="2" :x-gap="16">
          <n-grid-item>
            <n-form-item :label="t('prompt.keyLabel')">
              <n-input v-model:value="form.key" :disabled="!isNew" placeholder="code.review" />
            </n-form-item>
          </n-grid-item>
          <n-grid-item>
            <n-form-item :label="t('common.name')">
              <n-input v-model:value="form.name" :placeholder="t('prompt.namePlaceholder')" />
            </n-form-item>
          </n-grid-item>
          <n-grid-item>
            <n-form-item :label="t('common.env')">
              <n-select v-model:value="form.env" :options="envOptions" />
            </n-form-item>
          </n-grid-item>
          <n-grid-item>
            <n-form-item :label="t('common.version')">
              <n-input v-model:value="form.version" placeholder="v1" />
            </n-form-item>
          </n-grid-item>
          <n-grid-item>
            <n-form-item :label="t('prompt.colCategory')">
              <n-input
                v-model:value="form.category"
                :placeholder="t('prompt.categoryPlaceholder')"
              />
            </n-form-item>
          </n-grid-item>
          <n-grid-item>
            <n-form-item :label="t('common.model')">
              <n-input v-model:value="form.model" placeholder="gpt-4o" />
            </n-form-item>
          </n-grid-item>
          <n-grid-item :span="2">
            <n-form-item :label="t('prompt.tagsLabel')">
              <n-input v-model:value="form.tags" placeholder="review, sql" />
            </n-form-item>
          </n-grid-item>
        </n-grid>

        <n-form-item :label="t('prompt.contentLabel')">
          <n-input
            v-model:value="form.content"
            type="textarea"
            class="mono"
            :autosize="{ minRows: 12, maxRows: 28 }"
            :placeholder="contentPlaceholder"
          />
        </n-form-item>

        <div class="vars">
          <span class="vars-label">{{ t('prompt.detectedVars') }}</span>
          <n-space :size="6">
            <n-tag v-for="v in variables" :key="v" size="small" type="warning">
              {{ varLabel(v) }}
            </n-tag>
            <span v-if="!variables.length" class="muted">{{ t('common.none') }}</span>
          </n-space>
        </div>
      </n-form>
    </n-card>

    <n-card :title="t('prompt.rolloutTitle')" :bordered="true" style="margin-top: 16px">
      <p v-if="isNew" class="muted">{{ t('prompt.rolloutSaveFirst') }}</p>
      <template v-else>
        <n-form>
          <n-form-item :label="t('prompt.rolloutEnable')">
            <n-switch v-model:value="rollout.enabled" />
          </n-form-item>
          <n-grid :cols="3" :x-gap="16">
            <n-grid-item>
              <n-form-item :label="t('prompt.variantAVersion')">
                <n-input v-model:value="rollout.variant_a" placeholder="v1" />
              </n-form-item>
            </n-grid-item>
            <n-grid-item>
              <n-form-item :label="t('prompt.variantBVersion')">
                <n-input v-model:value="rollout.variant_b" placeholder="v2" />
              </n-form-item>
            </n-grid-item>
            <n-grid-item>
              <n-form-item :label="t('prompt.weightA')">
                <n-input-number v-model:value="rollout.weight_a" :min="0" :max="100" />
              </n-form-item>
            </n-grid-item>
          </n-grid>
          <p class="muted">{{ t('prompt.rolloutHint') }}</p>
          <n-button type="primary" @click="saveRollout">{{ t('prompt.saveRollout') }}</n-button>
        </n-form>
      </template>
    </n-card>

    <n-modal
      v-model:show="showVersions"
      preset="card"
      :title="t('prompt.versionHistory')"
      style="width: 720px"
    >
      <n-empty v-if="!versions.length" :description="t('prompt.noVersions')" />
      <n-table v-else :bordered="false" :single-line="false">
        <thead>
          <tr>
            <th>{{ t('common.version') }}</th>
            <th>{{ t('common.env') }}</th>
            <th>{{ t('prompt.colPublishedAt') }}</th>
            <th>{{ t('prompt.colPreview') }}</th>
            <th>{{ t('common.actions') }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="v in versions" :key="v.id">
            <td><n-tag size="small">{{ v.version }}</n-tag></td>
            <td>{{ v.env }}</td>
            <td>{{ new Date(v.created_at).toLocaleString() }}</td>
            <td class="preview">
              {{ v.content.slice(0, 60) }}{{ v.content.length > 60 ? '…' : '' }}
            </td>
            <td>
              <n-space :size="4">
                <n-button size="tiny" @click="openDiff(v)">{{ t('prompt.diffCurrent') }}</n-button>
                <n-button size="tiny" type="primary" ghost @click="rollback(v)">
                  {{ t('prompt.rollback') }}
                </n-button>
              </n-space>
            </td>
          </tr>
        </tbody>
      </n-table>
    </n-modal>

    <n-modal
      v-model:show="showDiff"
      preset="card"
      :title="t('prompt.diffTitle', { version: diffVersion })"
      style="width: 760px"
    >
      <div class="diff-legend">
        <span class="del">{{ t('prompt.diffOld', { version: diffVersion }) }}</span>
        <span class="add">{{ t('prompt.diffNew') }}</span>
      </div>
      <n-empty v-if="!diffLines.length" :description="t('prompt.diffEmpty')" />
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
