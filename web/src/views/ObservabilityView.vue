<script setup lang="ts">
import { ref, onMounted } from 'vue'
import {
  NCard,
  NGrid,
  NGi,
  NTabs,
  NTabPane,
  NTable,
  NTag,
  NEmpty,
  NSpin,
  NButton,
  NSpace,
  useMessage,
} from 'naive-ui'
import { api, type AuditEntry, type RunEntry, type RunStats } from '../api/client'

const message = useMessage()
const loading = ref(false)
const stats = ref<RunStats | null>(null)
const runs = ref<RunEntry[]>([])
const audit = ref<AuditEntry[]>([])

async function load() {
  loading.value = true
  try {
    const [s, r, a] = await Promise.all([api.runStats(), api.listRuns(), api.listAudit()])
    stats.value = s.data
    runs.value = r.data.data
    audit.value = a.data.data
  } catch {
    message.error('加载失败')
  } finally {
    loading.value = false
  }
}

function fmtTime(s: string) {
  return new Date(s).toLocaleString()
}

onMounted(load)
</script>

<template>
  <div class="page">
    <div class="toolbar">
      <h2>观测</h2>
      <n-button size="small" @click="load">刷新</n-button>
    </div>

    <n-spin :show="loading">
      <n-grid v-if="stats" :cols="4" :x-gap="12" :y-gap="12">
        <n-gi>
          <n-card>
            <div class="stat-num">{{ stats.total }}</div>
            <div class="stat-label">总运行次数</div>
          </n-card>
        </n-gi>
        <n-gi>
          <n-card>
            <div class="stat-num">{{ stats.ok }} / {{ stats.error }}</div>
            <div class="stat-label">成功 / 失败</div>
          </n-card>
        </n-gi>
        <n-gi>
          <n-card>
            <div class="stat-num">{{ stats.prompt_tokens + stats.output_tokens }}</div>
            <div class="stat-label">估算 Token 总量</div>
          </n-card>
        </n-gi>
        <n-gi>
          <n-card>
            <div class="stat-num">{{ stats.avg_latency_ms }} ms</div>
            <div class="stat-label">平均延迟</div>
          </n-card>
        </n-gi>
      </n-grid>

      <div v-if="stats && stats.by_provider.length" class="providers">
        <span class="muted">按提供方:</span>
        <n-space :size="6">
          <n-tag v-for="p in stats.by_provider" :key="p.provider" size="small" type="info">
            {{ p.provider }} · {{ p.count }}
          </n-tag>
        </n-space>
      </div>

      <n-tabs type="line" style="margin-top: 16px">
        <n-tab-pane name="runs" tab="运行日志">
          <n-empty v-if="!runs.length" description="暂无运行记录" style="margin: 32px 0" />
          <n-table v-else :bordered="false" :single-line="false">
            <thead>
              <tr>
                <th>时间</th>
                <th>来源</th>
                <th>提供方</th>
                <th>模型</th>
                <th>Token (in / out)</th>
                <th>延迟</th>
                <th>状态</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="r in runs" :key="r.id">
                <td>{{ fmtTime(r.created_at) }}</td>
                <td><n-tag size="small">{{ r.source }}</n-tag></td>
                <td>{{ r.provider }}</td>
                <td>{{ r.model || '-' }}</td>
                <td>{{ r.prompt_tokens }} / {{ r.output_tokens }}</td>
                <td>{{ r.latency_ms }} ms</td>
                <td>
                  <n-tag size="small" :type="r.status === 'ok' ? 'success' : 'error'">
                    {{ r.status }}
                  </n-tag>
                </td>
              </tr>
            </tbody>
          </n-table>
        </n-tab-pane>

        <n-tab-pane name="audit" tab="审计日志">
          <n-empty v-if="!audit.length" description="暂无审计记录" style="margin: 32px 0" />
          <n-table v-else :bordered="false" :single-line="false">
            <thead>
              <tr>
                <th>时间</th>
                <th>操作</th>
                <th>资源</th>
                <th>Key</th>
                <th>说明</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="a in audit" :key="a.id">
                <td>{{ fmtTime(a.created_at) }}</td>
                <td><n-tag size="small">{{ a.action }}</n-tag></td>
                <td>{{ a.resource }}</td>
                <td><code>{{ a.key || '-' }}</code></td>
                <td>{{ a.summary || '-' }}</td>
              </tr>
            </tbody>
          </n-table>
        </n-tab-pane>
      </n-tabs>
    </n-spin>
  </div>
</template>

<style scoped>
.toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}
.stat-num {
  font-size: 24px;
  font-weight: 700;
}
.stat-label {
  font-size: 12px;
  color: #888;
  margin-top: 4px;
}
.providers {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-top: 14px;
}
.muted {
  color: #888;
  font-size: 13px;
}
code {
  font-family: 'SF Mono', Menlo, Consolas, monospace;
}
</style>
