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
  NPagination,
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

const PAGE_SIZE = 50

const loading = ref(false)
const stats = ref<RunStats | null>(null)
const runs = ref<RunEntry[]>([])
const audit = ref<AuditEntry[]>([])
const clients = ref<ClientEntry[]>([])

const runsPage = ref(1)
const runsTotal = ref(0)
const auditPage = ref(1)
const auditTotal = ref(0)

async function fetchStats() {
  stats.value = (await api.runStats()).data
}

async function fetchRuns() {
  const { data } = await api.listRuns({
    limit: PAGE_SIZE,
    offset: (runsPage.value - 1) * PAGE_SIZE,
  })
  runs.value = data.data
  runsTotal.value = data.total
}

async function fetchAudit() {
  const { data } = await api.listAudit({
    limit: PAGE_SIZE,
    offset: (auditPage.value - 1) * PAGE_SIZE,
  })
  audit.value = data.data
  auditTotal.value = data.total
}

async function fetchClients() {
  clients.value = (await api.listClients()).data.data
}

async function load() {
  loading.value = true
  try {
    await Promise.all([fetchStats(), fetchRuns(), fetchAudit(), fetchClients()])
  } catch {
    message.error(t('common.loadFailed'))
  } finally {
    loading.value = false
  }
}

function onRunsPage(page: number) {
  runsPage.value = page
  fetchRuns().catch(() => message.error(t('common.loadFailed')))
}

function onAuditPage(page: number) {
  auditPage.value = page
  fetchAudit().catch(() => message.error(t('common.loadFailed')))
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
          <template v-else>
            <n-table :bordered="false" :single-line="false">
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
            <div v-if="runsTotal > PAGE_SIZE" class="pager">
              <n-pagination
                :page="runsPage"
                :page-size="PAGE_SIZE"
                :item-count="runsTotal"
                @update:page="onRunsPage"
              />
            </div>
          </template>
        </n-tab-pane>

        <n-tab-pane name="audit" :tab="t('observability.tabAudit')">
          <n-empty
            v-if="!audit.length"
            :description="t('observability.noAudit')"
            style="margin: 32px 0"
          />
          <template v-else>
            <n-table :bordered="false" :single-line="false">
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
            <div v-if="auditTotal > PAGE_SIZE" class="pager">
              <n-pagination
                :page="auditPage"
                :page-size="PAGE_SIZE"
                :item-count="auditTotal"
                @update:page="onAuditPage"
              />
            </div>
          </template>
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
.pager {
  display: flex;
  justify-content: flex-end;
  margin-top: 14px;
}
</style>
