<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
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
import {
  api,
  type AuditEntry,
  type RunEntry,
  type RunStats,
  type ClientEntry,
} from '../api/client'

const message = useMessage()
const { t } = useI18n()
const loading = ref(false)
const stats = ref<RunStats | null>(null)
const runs = ref<RunEntry[]>([])
const audit = ref<AuditEntry[]>([])
const clients = ref<ClientEntry[]>([])

async function load() {
  loading.value = true
  try {
    const [s, r, a, cl] = await Promise.all([
      api.runStats(),
      api.listRuns(),
      api.listAudit(),
      api.listClients(),
    ])
    stats.value = s.data
    runs.value = r.data.data
    audit.value = a.data.data
    clients.value = cl.data.data
  } catch {
    message.error(t('common.loadFailed'))
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
      <h2>{{ t('observability.title') }}</h2>
      <n-button size="small" @click="load">{{ t('common.refresh') }}</n-button>
    </div>

    <n-spin :show="loading">
      <n-grid v-if="stats" :cols="4" :x-gap="12" :y-gap="12">
        <n-gi>
          <n-card>
            <div class="stat-num">{{ stats.total }}</div>
            <div class="stat-label">{{ t('observability.totalRuns') }}</div>
          </n-card>
        </n-gi>
        <n-gi>
          <n-card>
            <div class="stat-num">{{ stats.ok }} / {{ stats.error }}</div>
            <div class="stat-label">{{ t('observability.okError') }}</div>
          </n-card>
        </n-gi>
        <n-gi>
          <n-card>
            <div class="stat-num">{{ stats.prompt_tokens + stats.output_tokens }}</div>
            <div class="stat-label">{{ t('observability.estTokens') }}</div>
          </n-card>
        </n-gi>
        <n-gi>
          <n-card>
            <div class="stat-num">{{ stats.avg_latency_ms }} ms</div>
            <div class="stat-label">{{ t('observability.avgLatency') }}</div>
          </n-card>
        </n-gi>
      </n-grid>

      <div v-if="stats && stats.by_provider.length" class="providers">
        <span class="muted">{{ t('observability.byProvider') }}</span>
        <n-space :size="6">
          <n-tag v-for="p in stats.by_provider" :key="p.provider" size="small" type="info">
            {{ p.provider }} · {{ p.count }}
          </n-tag>
        </n-space>
      </div>

      <n-tabs type="line" style="margin-top: 16px">
        <n-tab-pane name="runs" :tab="t('observability.tabRuns')">
          <n-empty
            v-if="!runs.length"
            :description="t('observability.noRuns')"
            style="margin: 32px 0"
          />
          <n-table v-else :bordered="false" :single-line="false">
            <thead>
              <tr>
                <th>{{ t('observability.colTime') }}</th>
                <th>{{ t('observability.colSource') }}</th>
                <th>{{ t('observability.colProvider') }}</th>
                <th>{{ t('common.model') }}</th>
                <th>{{ t('observability.colTokens') }}</th>
                <th>{{ t('observability.colLatency') }}</th>
                <th>{{ t('observability.colStatus') }}</th>
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

        <n-tab-pane name="audit" :tab="t('observability.tabAudit')">
          <n-empty
            v-if="!audit.length"
            :description="t('observability.noAudit')"
            style="margin: 32px 0"
          />
          <n-table v-else :bordered="false" :single-line="false">
            <thead>
              <tr>
                <th>{{ t('observability.colTime') }}</th>
                <th>{{ t('observability.colAction') }}</th>
                <th>{{ t('observability.colResource') }}</th>
                <th>Key</th>
                <th>{{ t('observability.colSummary') }}</th>
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

        <n-tab-pane name="clients" :tab="`${t('observability.tabClients')} (${clients.length})`">
          <n-empty
            v-if="!clients.length"
            :description="t('observability.noClients')"
            style="margin: 32px 0"
          />
          <n-table v-else :bordered="false" :single-line="false">
            <thead>
              <tr>
                <th>{{ t('observability.colClientType') }}</th>
                <th>{{ t('observability.colApp') }}</th>
                <th>{{ t('observability.colNamespace') }}</th>
                <th>{{ t('observability.colRemoteAddr') }}</th>
                <th>{{ t('observability.colUserAgent') }}</th>
                <th>{{ t('observability.colConnectedAt') }}</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="cl in clients" :key="cl.id">
                <td>
                  <n-tag size="small" :type="cl.client_type === 'browser' ? 'default' : 'info'">
                    {{ cl.client_type }}
                  </n-tag>
                </td>
                <td>{{ cl.app || '-' }}</td>
                <td>{{ cl.namespace || '-' }}</td>
                <td><code>{{ cl.remote_addr }}</code></td>
                <td class="ua">{{ cl.user_agent || '-' }}</td>
                <td>{{ fmtTime(cl.connected_at) }}</td>
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
.ua {
  max-width: 320px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-size: 12px;
  color: #aaa;
}
</style>
