<!-- components/Chart.vue -->
<template>
    <div class="chart-wrap">
        <div class="chart-header">
            <span class="chart-title">系统日志</span>
            <Select v-model="hours" :options="hourOptions" small/>
        </div>
        <div class="chart-tabs">
            <button 
                v-for="tab in tabs" 
                :key="tab.key"
                :class="{ active: activeTab === tab.key }"
                @click="activeTab = tab.key">
                {{ tab.label }}
            </button>
        </div>
        <canvas ref="canvas"></canvas>
    </div>
</template>

<script setup lang="ts">
import { Chart, LineController, LineElement, PointElement, LinearScale, Filler, CategoryScale, Tooltip } from 'chart.js'
import Select from '@/component/ui/Select.vue'
import { formatBytes } from '@/util/format'

Chart.register(LineController, LineElement, PointElement, LinearScale, Filler, CategoryScale, Tooltip)

const props = defineProps<{
    logs: Array<{
        CPU: number
        RAM: number
        Load: number
        Upload: number
        Download: number
        CreatedAt: string
    }>
}>()

const tabs = [
    { key: 'CPU',      label: 'CPU',  color: '#378ADD' },
    { key: 'RAM',      label: '内存',  color: '#1D9E75' },
    { key: 'Load',     label: '负载', color: '#888780' },
    { key: 'Upload',   label: '上传', color: '#D85A30' },
    { key: 'Download', label: '下载', color: '#9B59B6' },
]

const hourOptions = [
    { label: '3小时',  value: 3  },
    { label: '6小时',  value: 6  },
    { label: '12小时', value: 12 },
    { label: '24小时', value: 24 },
]

const crosshairPlugin = {
    id: 'crosshair',
    afterDraw(chart: any) {
        if (!chart.tooltip || !chart.tooltip._active || !chart.tooltip._active.length) return
        const ctx = chart.ctx
        const x = chart.tooltip._active[0].element.x
        const topY = chart.scales.y.top
        const bottomY = chart.scales.y.bottom
        ctx.save()
        ctx.beginPath()
        ctx.moveTo(x, topY)
        ctx.lineTo(x, bottomY)
        ctx.lineWidth = 1
        ctx.strokeStyle = 'rgba(255,255,255,0.2)'
        ctx.stroke()
        ctx.restore()
    }
}

const activeTab = ref('CPU')
const hours = ref(3)
const canvas = ref<HTMLCanvasElement>()
let chart: Chart

const currentTab = computed(() => tabs.find(t => t.key === activeTab.value)!)

const filteredLogs = computed(() => {
    const cutoff = new Date(Date.now() - hours.value * 60 * 60 * 1000)
    return props.logs.filter(l => new Date(l.CreatedAt) >= cutoff)
})

const labels = computed(() => filteredLogs.value.map(l => {
    const d = new Date(l.CreatedAt)
    return `${d.getHours().toString().padStart(2, '0')}:${d.getMinutes().toString().padStart(2, '0')}`
}))

const data = computed(() => filteredLogs.value.map(l => l[activeTab.value as keyof typeof l] as number))

function formatValue(value: number) {
    if (activeTab.value === 'CPU' || activeTab.value === 'RAM') {
        return value.toFixed(1) + '%'
    }
    if (activeTab.value === 'Upload' || activeTab.value === 'Download') {
        return formatBytes(value / 900) + '/s'
    }
    return value.toFixed(2)
}

onMounted(() => {
    chart = new Chart(canvas.value!, {
        type: 'line',
        plugins: [crosshairPlugin],
        data: {
            labels: labels.value,
            datasets: [{
                data: data.value,
                borderColor: currentTab.value.color,
                borderWidth: 2,
                fill: true,
                backgroundColor: currentTab.value.color + '20',
                pointRadius: 0,
                tension: 0.4,
            }]
        },
        options: {
            responsive: true,
            interaction: {
                mode: 'index',
                intersect: false,
            },
            plugins: {
                legend: { display: false },
                tooltip: {
                    displayColors: false,
                    backgroundColor: 'rgba(30,30,30,0.9)',
                    titleColor: '#aaa',
                    bodyColor: '#fff',
                    borderColor: 'rgba(255,255,255,0.1)',
                    borderWidth: 1,
                    callbacks: {
                        label: (ctx) => formatValue((ctx.parsed.y ?? 0) as number)
                    }
                }
            },
            scales: {
                x: {
                    display: true,
                    ticks: {
                        color: '#888780',
                        font: { size: 11 },
                        maxTicksLimit: 8,
                    },
                    grid: { display: false },
                    border: { display: false },
                },
                y: {
                    display: true,
                    ticks: {
                        color: '#888780',
                        font: { size: 11 },
                        maxTicksLimit: 5,
                        callback: (value) => formatValue(Number(value)),
                    },
                    grid: {
                        color: 'rgba(255,255,255,0.05)',
                    },
                    border: { display: false },
                },
            }
        }
    })
})

watch([activeTab, hours, () => props.logs], () => {
    if (!chart) return
    chart.data.labels = labels.value
    const dataset = chart.data.datasets?.[0]
    if (dataset) {
        dataset.data = data.value
        dataset.borderColor = currentTab.value.color
        dataset.backgroundColor = currentTab.value.color + '20'
        chart.update()
    }
})
</script>

<style scoped>
.chart-wrap {
    display: flex;
    flex-direction: column;
    gap: 10px;
}
.chart-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
}
.chart-title {
    font-size: 14px;
    font-weight: bold;
    color: var(--color-text-light);
}
.chart-tabs {
    display: flex;
    flex-direction: row;
    gap: 5px;
    flex-wrap: wrap;
}
.chart-tabs button {
    padding: 5px 12px;
    border-radius: 999px;
    font-size: 12px;
    background-color: transparent;
    color: var(--color-text-dark);
    
    &.active {
        background-color: var(--color-bg-light);
        color: var(--color-text-light);
    }
    &:hover {
        background-color: var(--color-bg-light);
    }
}
canvas {
    max-height: 200px;
}
</style>