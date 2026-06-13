<!-- components/RouteDetect.vue -->
<template>
    <div class="detect">
        <div class="header">
            <div class="title-wrap">
                <div class="title">回程检测</div>
                <span v-if="data.length" class="updated">
                    {{ formatTime(data[0]?.UpdatedAt) }}
                </span>
            </div>
            <div class="action">
                <button :disabled="loading" @click="handleFetch">
                    <i class="icon">{{ loading ? 'hourglass_empty' : 'travel_explore' }}</i>
                    <span>{{ loading ? '检测中' : '开始检测' }}</span>
                </button>
            </div>
        </div>

        <div v-for="city in cities" :key="city.key" class="route-group">
            <div class="route-group-title">{{ city.label }}</div>
            <div class="route-row">
                <div class="route-isp">
                    <span class="route-isp-name">电信</span>
                    <span :class="lineClass('telecom', getLine(city.key, 'telecom'))">
                        {{ getLine(city.key, 'telecom') || '-' }}
                    </span>
                </div>
                <div class="route-isp">
                    <span class="route-isp-name">联通</span>
                    <span :class="lineClass('unicom', getLine(city.key, 'unicom'))">
                        {{ getLine(city.key, 'unicom') || '-' }}
                    </span>
                </div>
                <div class="route-isp">
                    <span class="route-isp-name">移动</span>
                    <span :class="lineClass('mobile', getLine(city.key, 'mobile'))">
                        {{ getLine(city.key, 'mobile') || '-' }}
                    </span>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { getRoute, fetchRoute } from '@/api/detect'
import { formatTime } from '@/util/format'

const loading = ref(false)
const data = ref<any[]>([])

const cities = [
    { key: 'shanghai',  label: '上海' },
    { key: 'beijing',   label: '北京' },
    { key: 'guangzhou', label: '广州' },
]

function getLine(city: string, isp: string): string {
    const record = data.value.find((r: any) => r.City === city)
    if (!record) return ''
    if (isp === 'telecom') return record.Telecom ?? ''
    if (isp === 'unicom')  return record.Unicom ?? ''
    if (isp === 'mobile')  return record.Mobile ?? ''
    return ''
}

function lineClass(isp: string, line: string): string {
    if (!line || line === '-' || line === '未知' || line === '检测失败') return 'c-empty'

    const green: Record<string, string[]> = {
        telecom: ['CN2 GIA'],
        unicom:  ['9929'],
        mobile:  ['CMIN2'],
    }
    const yellow: Record<string, string[]> = {
        telecom: ['CN2 GT'],
        unicom:  ['4837'],
        mobile:  ['CMI'],
    }

    if (green[isp]?.includes(line))  return 'c-yes'
    if (yellow[isp]?.includes(line)) return 'c-reject'
    return 'c-no'
}

async function handleFetch() {
    loading.value = true
    try {
        const res = await fetchRoute()
        data.value = res.data ?? []
    } finally {
        loading.value = false
    }
}

onMounted(async () => {
    const res = await getRoute()
    data.value = res.data ?? []
})

watch(loading, (val) => {
    document.body.style.cursor = val ? 'wait' : ''
})
</script>

<style scoped>
.route-group {
    display: flex;
    flex-direction: column;
    gap: 10px;
}

.route-group-title {
    font-size: 12px;
    color: var(--color-text-dark);
}

.route-row {
    display: flex;
    gap: 32px;
}

.route-isp {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 13px;
}

.route-isp-name {
    color: var(--color-text-light);
    font-weight: bold;
}

@media (max-width: 768px) {
    .route-row {
        flex-direction: column;
        gap: 8px;
    }
}
</style>