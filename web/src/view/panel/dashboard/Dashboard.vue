<template>
    <main>
        <div class="page">
            <div class="body">
                <div class="card fluid">
                    <div class="gauge-card">
                        <Gauge name="CPU" :pct="system.cpu?.percent ?? 0" color="#378ADD" :detail="`${system.cpu?.cores ?? 0} Core`">
                            <template #tip>
                                <div>逻辑处理器: {{ system.cpu?.cores ?? 0 }}</div>
                                <div>频率: {{ system.cpu?.mhz ?? 0 }} MHz</div>
                            </template>
                        </Gauge>
                        <Gauge name="RAM"  :pct="system.ram?.percent ?? 0"  color="#1D9E75" :detail="`${formatBytes(system.ram?.used)} / ${formatBytes(system.ram?.total)}`" />
                        <Gauge name="Swap" :pct="system.swap?.percent ?? 0" color="#888780" :detail="`${formatBytes(system.swap?.used)} / ${formatBytes(system.swap?.total)}`" />
                        <Gauge name="存储"  :pct="system.disk?.percent ?? 0" color="#D85A30" :detail="`${formatBytes(system.disk?.used)} / ${formatBytes(system.disk?.total)}`" />
                    </div>
                </div>
                <CoreCard />
                <PanelCard />
                <div class="card md">
                    <div class="header">
                        <div class="title">运行时间</div>
                    </div>
                    <div class="body">
                        <div class="info">
                            <div class="title">核心</div>
                            <div class="detail">
                                <i class="icon">bolt</i>
                                <span>{{ formatUptime(core.StartedAt) }}</span>
                            </div>
                        </div>
                        <div class="info">
                            <div class="title">面板</div>
                            <div class="detail">
                                <i class="icon">monitor</i>
                                <span>{{ formatUptime(Config.StartedAt) }}</span>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="card md">
                    <div class="header">
                        <div class="title">整体速率</div>
                    </div>
                    <div class="body">
                        <div class="info">
                            <div class="title">上传</div>
                            <div class="detail">
                                <i class="icon">upload</i>
                                <span>{{ formatBytes(system.network?.upload ?? 0) }}/s</span>
                            </div>
                        </div>
                        <div class="info">
                            <div class="title">下载</div>
                            <div class="detail">
                                <i class="icon">download</i>
                                <span>{{ formatBytes(system.network?.download ?? 0) }}/s</span>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="card md">
                    <div class="header">
                        <div class="title">进程开销</div>
                    </div>
                    <div class="body">
                        <div class="info">
                            <div class="title">RAM</div>
                            <div class="detail">
                                <i class="icon">dns</i>
                                <span>{{ formatBytes(process.memory ?? 0) }}</span>
                            </div>
                        </div>
                        <div class="info">
                            <div class="title">线程</div>
                            <div class="detail">
                                <i class="icon">account_tree</i>
                                <span>{{ process.threads ?? 0 }}</span>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="card md">
                    <div class="header">
                        <div class="title">总数据</div>
                    </div>
                    <div class="body">
                        <div class="info">
                            <div class="title">发送</div>
                            <div class="detail">
                                <i class="icon">cloud_upload</i>
                                <span>{{ formatBytes(stats.Upload ?? 0) }}</span>
                            </div>
                        </div>
                        <div class="info">
                            <div class="title">接收</div>
                            <div class="detail">
                                <i class="icon">cloud_download</i>
                                <span>{{ formatBytes(stats.Download ?? 0) }}</span>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="card md">
                    <div class="header">
                        <div class="title">连接数</div>
                    </div>
                    <div class="body">
                        <div class="info">
                            <div class="title">TCP</div>
                            <div class="detail">
                                <i class="icon">swap_horiz</i>
                                <span>{{ system.connections?.tcp ?? 0 }}</span>
                            </div>
                        </div>
                        <div class="info">
                            <div class="title">UDP</div>
                            <div class="detail">
                                <i class="icon">swap_horiz</i>
                                <span>{{ system.connections?.udp ?? 0 }}</span>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="card md">
                    <div class="header">
                        <div class="title">IP 地址</div>
                    </div>
                    <div class="body">
                        <div class="info">
                            <div class="title">IPv4</div>
                            <div class="detail">
                                <i class="icon">public</i>
                                <span>{{ Config?.IPv4 }}</span>
                            </div>
                        </div>
                        <div class="info">
                            <div class="title">IPv6</div>
                            <div class="detail">
                                <i class="icon">public</i>
                                <span>{{ Config?.IPv6 || '无' }}</span>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="card fluid">
                    <Chart :logs="logs" />
                </div>
                <Quick />
                <Ad />
            </div>
        </div>
    </main>
</template>

<script setup lang="ts">PanelCard
import Gauge from '@/component/widget/Gauge.vue'
import CoreCard from '@/component/widget/CoreCard.vue'
import PanelCard from '@/component/widget/PanelCard.vue'
import Chart from '@/component/widget/Chart.vue'
import Quick from '@/view/panel/dashboard/widget/Quick.vue'
import Ad from '@/view/panel/dashboard/widget/Ad.vue'
import { getConfig } from '@/api/config'
import { getCore, getCoreProcess } from '@/api/core'
import { getSystem, getStats, getSystemLog } from '@/api/system'
import { formatBytes, formatUptime } from '@/util/format'

const Config = ref<any>({})
const core = ref<any>({})
const process = ref<any>({})
const system = ref<any>({})
const stats = ref<any>({})
const logs = ref<any[]>([])

async function fetchSystem() { system.value = await getSystem() }
async function fetchCore() { core.value = await getCore() }
async function fetchProcess() { process.value = await getCoreProcess() }
async function fetchStats() { stats.value = await getStats() }
async function fetchLog() { logs.value = await getSystemLog() }

const timers: ReturnType<typeof setInterval>[] = []

onMounted(async () => {
    const [config] = await Promise.all([getConfig()])
    Config.value = config

    fetchCore()
    fetchProcess()
    fetchSystem()
    fetchLog()
    fetchStats()

    timers.push(
        setInterval(fetchSystem, 5000),
        setInterval(fetchCore, 60 * 1000),
        setInterval(fetchProcess, 5000),
        setInterval(fetchStats, 15 * 60 * 1000),
        setInterval(fetchLog, 15 * 60 * 1000),
        setInterval(async () => { Config.value = await getConfig() }, 60 * 1000),
    )
})

onUnmounted(() => timers.forEach(clearInterval))
</script>

<style scoped>
.gauge-card {
    display: flex;
    flex-direction: row;
    align-items: center;
    gap: 10px;
    justify-content: space-between;

    > * {
        flex: 1;
    }

    @media (max-width: 768px) {
        flex-wrap: wrap;

        > * {
            flex: 0 0 calc(50% - 5px);
        }
    }
}
</style>