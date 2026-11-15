<template>
    <Page class="dashboard" title="仪表盘" :motion="themeMotion">
        <div class="grid">
            <Box class="full">
                <div class="gauges">
                    <Gauge
                        name="CPU"
                        :pct="system?.cpu.percent ?? 0"
                        color="#378ADD"
                        :detail="`${system?.cpu.cores ?? 0} Core`"
                    >
                        <template #tip>
                            <div>逻辑处理器: {{ system?.cpu.cores ?? 0 }}</div>
                            <div>频率: {{ system?.cpu.mhz ?? 0 }} MHz</div>
                        </template>
                    </Gauge>
                    <Gauge
                        name="RAM"
                        :pct="system?.ram.percent ?? 0"
                        color="#1D9E75"
                        :detail="usage(system?.ram)"
                    />
                    <Gauge
                        name="Swap"
                        :pct="system?.swap.percent ?? 0"
                        color="#888780"
                        :detail="usage(system?.swap)"
                    />
                    <Gauge
                        name="存储"
                        :pct="system?.disk.percent ?? 0"
                        color="#D85A30"
                        :detail="usage(system?.disk)"
                    />
                </div>
            </Box>
            <CoreCard />
            <PanelCard />
            <Metric
                title="运行时间"
                :items="[
                    { label: '核心', icon: 'bolt', value: formatUptime(core?.StartedAt ?? '') },
                    { label: '面板', icon: 'monitor', value: formatUptime(config?.StartedAt ?? '') },
                ]"
            />
            <Metric
                title="整体速率"
                :items="[
                    { label: '上传', icon: 'upload', value: `${formatBytes(system?.network.upload)}/s` },
                    { label: '下载', icon: 'download', value: `${formatBytes(system?.network.download)}/s` },
                ]"
            />
            <Metric
                title="进程开销"
                :items="[
                    { label: 'RAM', icon: 'dns', value: formatBytes(process?.memory) },
                    { label: '线程', icon: 'account_tree', value: process?.threads ?? 0 },
                ]"
            />
            <Metric
                title="总数据"
                :items="[
                    { label: '发送', icon: 'cloud_upload', value: formatBytes(stats?.Upload) },
                    { label: '接收', icon: 'cloud_download', value: formatBytes(stats?.Download) },
                ]"
            />
            <Metric
                title="连接数"
                :items="[
                    { label: 'TCP', icon: 'swap_horiz', value: system?.connections.tcp ?? 0 },
                    { label: 'UDP', icon: 'swap_horiz', value: system?.connections.udp ?? 0 },
                ]"
            />
            <Metric
                title="IP 地址"
                :items="[
                    { label: 'IPv4', icon: 'public', value: config?.IPv4 ?? '' },
                    { label: 'IPv6', icon: 'public', value: config?.IPv6 || '无' },
                ]"
            />
            <Chart class="full" :logs="logs" />
            <Quick />
            <Ad />
        </div>
    </Page>
</template>

<script setup lang="ts">
import Page from '@/component/layout/Page.vue'
import Box from '@/component/layout/Box.vue'
import Gauge from '@/component/widget/Gauge.vue'
import Metric from '@/component/module/Metric.vue'
import Chart from '@/component/module/Chart.vue'
import CoreCard from '@/view/panel/dashboard/CoreCard.vue'
import PanelCard from '@/view/panel/dashboard/PanelCard.vue'
import Quick from '@/view/panel/dashboard/Quick.vue'
import Ad from '@/view/panel/dashboard/Ad.vue'
import { getConfig } from '@/api/config'
import type { Config } from '@/api/config'
import { getCore, getCoreProcess } from '@/api/core'
import type { Core, CoreProcess } from '@/api/core'
import { getSystem, getStats, getSystemLog } from '@/api/system'
import type { SystemStatus, Stats, SystemLog } from '@/api/system'
import { formatBytes, formatUptime } from '@/util/format'
import { themeMotion } from '@/composable/theme'

const config = ref<Config>()
const core = ref<Core>()
const process = ref<CoreProcess>()
const system = ref<SystemStatus>()
const stats = ref<Stats>()
const logs = ref<SystemLog[]>([])

const usage = (item?: { used: number; total: number }) =>
    `${formatBytes(item?.used)} / ${formatBytes(item?.total)}`

async function loadConfig() {
    config.value = await getConfig()
}

async function loadCore() {
    core.value = await getCore()
}

// 核心没运行时拿不到进程信息，保持上一次的值
async function loadProcess() {
    try {
        process.value = await getCoreProcess()
    } catch {
        process.value = undefined
    }
}

async function loadSystem() {
    system.value = await getSystem()
}

async function loadStats() {
    stats.value = await getStats()
}

async function loadLog() {
    logs.value = await getSystemLog()
}

const timers: ReturnType<typeof setInterval>[] = []

onMounted(() => {
    loadConfig()
    loadCore()
    loadProcess()
    loadSystem()
    loadLog()
    loadStats()

    timers.push(
        setInterval(loadSystem, 5000),
        setInterval(loadCore, 60 * 1000),
        setInterval(loadProcess, 5000),
        setInterval(loadStats, 15 * 60 * 1000),
        setInterval(loadLog, 15 * 60 * 1000),
        setInterval(loadConfig, 60 * 1000),
    )
})

onUnmounted(() => timers.forEach(clearInterval))
</script>

<style scoped>
.dashboard {
    &:deep(.motion button) {
        border-radius: 999px;
    }
}

.grid {
    display: grid;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    gap: 20px;
    padding-bottom: 20px;

    .full {
        grid-column: 1 / -1;
    }

    @media (max-width: 768px) {
        grid-template-columns: minmax(0, 1fr);
    }
}

.gauges {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 10px;

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
