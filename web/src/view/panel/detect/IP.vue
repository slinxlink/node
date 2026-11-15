<!-- view/panel/detect/IP.vue -->
<template>
    <Box title="IP 检测">
        <template #action>
            <span v-if="current?.UpdatedAt">{{ formatTime(current.UpdatedAt) }}</span>
            <Select v-model:value="source" :items="sources" />
            <Button class="run" :disabled="loading" @click="detect">
                <Icon :name="loading ? 'hourglass_empty' : 'travel_explore'" />
                <span>{{ loading ? '检测中' : '开始检测' }}</span>
            </Button>
        </template>

        <div class="grid">
            <div v-for="item in items" :key="item.label" class="item">
                <div class="label">{{ item.label }}</div>
                <div class="value">
                    <Stamp v-if="item.color" :color="item.color">{{ item.text }}</Stamp>
                    <span v-else>{{ item.text }}</span>
                </div>
            </div>
        </div>
    </Box>
</template>

<script setup lang="ts">
import Box from '@/component/layout/Box.vue'
import Button from '@/component/widget/Button.vue'
import Icon from '@/component/widget/Icon.vue'
import Select from '@/component/widget/Select.vue'
import Stamp from '@/component/widget/Stamp.vue'
import { getIP, fetchIP } from '@/api/detect'
import type { IPInfo } from '@/api/detect'
import { toast } from '@/composable/toast'
import { countryFlag } from '@/util/country'
import { formatTime } from '@/util/format'
import { translateType, typeColor } from '@/util/iptype'

// 各个数据源查到的结果分开存着，切换数据源直接显示对应的那条；只看 IPv4
const all = ref<IPInfo[]>([])
const source = ref('ip.net.coffee')
const loading = ref(false)

const sources = [
    { label: 'ip.net.coffee', value: 'ip.net.coffee' },
    { label: 'ippure.com', value: 'ippure.com' },
    { label: 'ip-api.com', value: 'ip-api.com' },
    { label: 'ip2location.io', value: 'ip2location.io' },
]

const current = computed(() =>
    all.value.find((item) => item.Source === source.value && item.IPVersion === 'v4'),
)

// 国家前面带旗帜；类型带颜色，没有归类的就是普通文字；原生性只有两个数据源能判断，其余显示短横线
const items = computed(() => {
    const ip = current.value
    const native = nativeText(ip?.Native)

    return [
        { label: 'IP 地址', text: ip?.IP || '-' },
        { label: '经纬度', text: coordinate(ip) },
        { label: 'IP 原生性', text: native.text, color: native.color },
        { label: 'IP 类型', text: translateType(ip?.IPType), color: typeColor(ip?.IPType) },
        { label: '国家', text: `${countryFlag(ip?.Country)} ${ip?.Country || '-'}`.trim() },
        { label: '城市', text: ip?.City || '-' },
        { label: 'ASN', text: ip?.ASN || '-' },
        { label: '运营商', text: ip?.ASNOrg || '-' },
    ]
})

// 纬度、经度，保留 4 位；都是 0 就是没查到
function coordinate(ip?: IPInfo): string {
    if (!ip || (ip.Latitude === 0 && ip.Longitude === 0)) return '-'

    return `${ip.Latitude.toFixed(4)}, ${ip.Longitude.toFixed(4)}`
}

function nativeText(value?: string) {
    if (value === 'native') return { text: '原生', color: 'green' }
    if (value === 'broadcast') return { text: '广播', color: 'yellow' }

    return { text: '-', color: '' }
}

async function detect() {
    loading.value = true

    try {
        const data = await fetchIP(source.value)
        const index = all.value.findIndex(
            (item) => item.Source === source.value && item.IPVersion === 'v4',
        )

        if (index >= 0) {
            all.value[index] = data
        } else {
            all.value.push(data)
        }
    } catch (err) {
        toast.failure((err as { error?: string })?.error || '检测失败')
    } finally {
        loading.value = false
    }
}

onMounted(async () => {
    all.value = await getIP()
})
</script>

<style scoped>
.grid {
    display: grid;
    grid-template-columns: repeat(4, minmax(0, 1fr));
    gap: 10px;

    .item {
        display: flex;
        flex-direction: column;
        gap: 5px;
        min-width: 0;
    }

    .label {
        color: var(--text-300);
        font-size: 12px;
    }

    .value {
        display: flex;
        align-items: center;
        height: 22px; /* 和标签(Stamp)一样高，文字的行就能垂直居中对齐 */
        min-width: 0;
        color: var(--text-100);
        font-weight: bold;

        span {
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
        }
    }

    @media (max-width: 1024px) {
        grid-template-columns: repeat(2, minmax(0, 1fr));
    }

    @media (max-width: 768px) {
        grid-template-columns: minmax(0, 1fr);
    }
}
</style>
