<!-- component/module/Chart.vue -->
<template>
    <Box title="系统日志">
        <template #action>
            <div class="range">
                <Select v-model:value="hours" :items="ranges" />
            </div>
        </template>
        <Tab class="tabs" v-model:value="active" :items="tabs" />
        <canvas ref="canvas"></canvas>
    </Box>
</template>

<script setup lang="ts">
import {
    Chart as ChartJS,
    LineController,
    LineElement,
    PointElement,
    LinearScale,
    Filler,
    CategoryScale,
    Tooltip,
} from 'chart.js'
import type { Plugin } from 'chart.js'
import Box from '@/component/layout/Box.vue'
import Tab from '@/component/widget/Tab.vue'
import type { TabItem } from '@/component/widget/Tab.vue'
import Select from '@/component/widget/Select.vue'
import type { SystemLog } from '@/api/system'
import { resolved } from '@/composable/theme'
import { formatBytes } from '@/util/format'

ChartJS.register(
    LineController,
    LineElement,
    PointElement,
    LinearScale,
    Filler,
    CategoryScale,
    Tooltip,
)

const props = defineProps<{
    logs: SystemLog[]
}>()

type Key = 'CPU' | 'RAM' | 'Load' | 'Upload' | 'Download'

const colors: Record<Key, string> = {
    CPU: '#378ADD',
    RAM: '#1D9E75',
    Load: '#888780',
    Upload: '#D85A30',
    Download: '#9B59B6',
}

const tabs: [TabItem, ...TabItem[]] = [
    { label: 'CPU', value: 'CPU' },
    { label: '内存', value: 'RAM' },
    { label: '负载', value: 'Load' },
    { label: '上传', value: 'Upload' },
    { label: '下载', value: 'Download' },
]

const ranges = [
    { label: '3小时', value: 3 },
    { label: '6小时', value: 6 },
    { label: '12小时', value: 12 },
    { label: '24小时', value: 24 },
]

const active = ref<string>('CPU')
const hours = ref<number | null>(3)
const canvas = ref<HTMLCanvasElement>()
let chart: ChartJS | undefined

const key = computed(() => active.value as Key)

const filtered = computed(() => {
    const cutoff = Date.now() - (hours.value ?? 3) * 60 * 60 * 1000
    return props.logs.filter((l) => new Date(l.CreatedAt).getTime() >= cutoff)
})

const labels = computed(() =>
    filtered.value.map((l) => {
        const d = new Date(l.CreatedAt)
        return `${d.getHours().toString().padStart(2, '0')}:${d.getMinutes().toString().padStart(2, '0')}`
    }),
)

const data = computed(() => filtered.value.map((l) => l[key.value]))

function formatValue(value: number) {
    if (key.value === 'CPU' || key.value === 'RAM') return value.toFixed(1) + '%'
    if (key.value === 'Upload' || key.value === 'Download') return formatBytes(value / 900) + '/s'
    return value.toFixed(2)
}

// 画布读不了 CSS 变量，取当前主题下的实际颜色
function theme() {
    const css = getComputedStyle(document.body)
    const read = (name: string) => css.getPropertyValue(name).trim()
    return { text: read('--text-300'), grid: read('--background-100'), tip: read('--background-300'), tipText: read('--text-100') }
}

// 鼠标悬停时画竖线
const crosshair: Plugin = {
    id: 'crosshair',
    afterDraw(chart) {
        const active = chart.tooltip?.getActiveElements() ?? []
        const first = active[0]
        if (!first) return
        const { ctx, scales } = chart
        const y = scales.y
        if (!y) return
        ctx.save()
        ctx.beginPath()
        ctx.moveTo(first.element.x, y.top)
        ctx.lineTo(first.element.x, y.bottom)
        ctx.lineWidth = 1
        ctx.strokeStyle = theme().grid
        ctx.stroke()
        ctx.restore()
    },
}

function paint() {
    if (!chart) return
    const color = colors[key.value]
    const c = theme()

    chart.data.labels = labels.value

    const dataset = chart.data.datasets[0]
    if (dataset) {
        dataset.data = data.value
        dataset.borderColor = color
        dataset.backgroundColor = color + '20'
    }

    const { x, y } = chart.options.scales ?? {}
    if (x?.ticks) x.ticks.color = c.text
    if (y?.ticks) y.ticks.color = c.text
    if (y?.grid) y.grid.color = c.grid

    const tip = chart.options.plugins?.tooltip
    if (tip) {
        tip.backgroundColor = c.tip
        tip.titleColor = c.text
        tip.bodyColor = c.tipText
    }

    chart.update()
}

onMounted(() => {
    if (!canvas.value) return

    chart = new ChartJS(canvas.value, {
        type: 'line',
        plugins: [crosshair],
        data: {
            labels: [],
            datasets: [{ data: [], borderWidth: 2, fill: true, pointRadius: 0, tension: 0.4 }],
        },
        options: {
            responsive: true,
            interaction: { mode: 'index', intersect: false },
            plugins: {
                legend: { display: false },
                tooltip: {
                    displayColors: false,
                    borderWidth: 0,
                    callbacks: {
                        label: (ctx) => formatValue(ctx.parsed.y ?? 0),
                    },
                },
            },
            scales: {
                x: {
                    ticks: { font: { size: 11 }, maxTicksLimit: 8 },
                    grid: { display: false },
                    border: { display: false },
                },
                y: {
                    ticks: {
                        font: { size: 11 },
                        maxTicksLimit: 5,
                        callback: (value) => formatValue(Number(value)),
                    },
                    border: { display: false },
                },
            },
        },
    })

    paint()
})

onUnmounted(() => chart?.destroy())

watch([active, hours, () => props.logs], paint)

// 主题切换后等类名生效再重新取色
watch(resolved, () => nextTick(paint))
</script>

<style scoped>
.range {
    width: 120px;
    margin-right: -7px; /* 下拉 34px 比标题行 20px 高，上下各凸出 7px，离卡片上沿 13px；右边也往外挪 7px，同样是 13px */
}

/* Tab 默认上下留 20px 是给大标题行用的，这里紧凑一点 */
.tab.tabs {
    padding: 10px 0;
}

canvas {
    max-height: 200px;
}
</style>
