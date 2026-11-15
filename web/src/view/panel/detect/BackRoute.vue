<!-- view/panel/detect/BackRoute.vue -->
<template>
    <Box title="回程检测">
        <template #action>
            <span v-if="data.length">{{ formatTime(data[0]?.UpdatedAt ?? '') }}</span>
            <Button class="run" :disabled="loading" @click="detect">
                <Icon :name="loading ? 'hourglass_empty' : 'travel_explore'" />
                <span>{{ loading ? '检测中' : '开始检测' }}</span>
            </Button>
        </template>

        <div v-for="city in cities" :key="city.key" class="group">
            <div class="title">{{ city.label }}</div>
            <div class="row">
                <div v-for="isp in isps" :key="isp.key" class="isp">
                    <span class="name">{{ isp.label }}</span>
                    <span :style="{ color: lineColor(isp.key, line(city.key, isp.key)) }">
                        {{ line(city.key, isp.key) || '-' }}
                    </span>
                </div>
            </div>
        </div>
    </Box>
</template>

<script setup lang="ts">
import Box from '@/component/layout/Box.vue'
import Button from '@/component/widget/Button.vue'
import Icon from '@/component/widget/Icon.vue'
import { getBackRoute, fetchBackRoute } from '@/api/detect'
import type { BackRoute } from '@/api/detect'
import { toast } from '@/composable/toast'
import { formatTime } from '@/util/format'

type Isp = 'telecom' | 'unicom' | 'mobile'

const data = ref<BackRoute[]>([])
const loading = ref(false)

const cities = [
    { key: 'shanghai', label: '上海' },
    { key: 'beijing', label: '北京' },
    { key: 'guangzhou', label: '广州' },
]

const isps: { key: Isp; label: string }[] = [
    { key: 'telecom', label: '电信' },
    { key: 'unicom', label: '联通' },
    { key: 'mobile', label: '移动' },
]

// 顶级线路绿色，次一级黄色，其余普通线路红色；没数据或检测失败是灰色
const best: Record<Isp, string> = { telecom: 'CN2 GIA', unicom: '9929', mobile: 'CMIN2' }
const good: Record<Isp, string> = { telecom: 'CN2 GT', unicom: '4837', mobile: 'CMI' }

function line(city: string, isp: Isp): string {
    const record = data.value.find((item) => item.City === city)

    if (!record) return ''

    return { telecom: record.Telecom, unicom: record.Unicom, mobile: record.Mobile }[isp] ?? ''
}

function lineColor(isp: Isp, value: string): string {
    if (!value || ['-', '未知', '检测失败'].includes(value)) return 'var(--text-300)'
    if (value === best[isp]) return 'var(--green)'
    if (value === good[isp]) return 'var(--yellow)'

    return 'var(--red)'
}

async function detect() {
    loading.value = true

    try {
        data.value = await fetchBackRoute()
    } catch (err) {
        toast.failure((err as { error?: string })?.error || '检测失败')
    } finally {
        loading.value = false
    }
}

onMounted(async () => {
    data.value = await getBackRoute()
})
</script>

<style scoped>
.group {
    display: flex;
    flex-direction: column;
    gap: 10px;

    .title {
        color: var(--text-300);
        font-size: 12px;
    }
}

.row {
    display: flex;
    gap: 32px;

    @media (max-width: 768px) {
        flex-direction: column;
        gap: 8px;
    }
}

.isp {
    display: flex;
    align-items: center;
    gap: 8px;

    .name {
        color: var(--text-100);
        font-weight: bold;
    }
}
</style>
