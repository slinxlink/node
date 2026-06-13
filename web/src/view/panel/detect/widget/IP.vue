<!-- components/IPDetect.vue -->
<template>
    <div class="detect">
        <div class="header">
            <div class="title-wrap">
                <div class="title">IP 检测</div>
                <span v-if="ipData?.UpdatedAt" class="updated">
                    {{ formatTime(ipData.UpdatedAt) }}
                </span>
            </div>
            <div class="action">
                <Select v-model="source" :options="sourceOptions" small />
                <button :disabled="loading" @click="handleFetch">
                    <i class="icon">{{ loading ? 'hourglass_empty' : 'travel_explore' }}</i>
                    <span>{{ loading ? '检测中' : '开始检测' }}</span>
                </button>
            </div>
        </div>
        <div class="ip-detect-body">
            <div class="ip-detect-item">
                <div class="ip-detect-label">IP 地址</div>
                <div class="ip-detect-value">
                    <i class="icon">public</i>
                    <span>{{ ipData?.IP || '-' }}</span>
                </div>
            </div>
            <div class="ip-detect-item">
                <div class="ip-detect-label">ASN</div>
                <div class="ip-detect-value">
                    <span>{{ ipData?.ASN || '-' }}</span>
                </div>
            </div>
            <div class="ip-detect-item">
                <div class="ip-detect-label">运营商</div>
                <div class="ip-detect-value">
                    <span>{{ ipData?.ASNOrg || '-' }}</span>
                </div>
            </div>
            <div class="ip-detect-item">
                <div class="ip-detect-label">注册地</div>
                <div class="ip-detect-value">
                    <span>{{ countryFlag(ipData?.RegisterCountry) }} {{ ipData?.RegisterCountry || '-' }}</span>
                </div>
            </div>
            <div class="ip-detect-item">
                <div class="ip-detect-label">国家</div>
                <div class="ip-detect-value">
                    <span>{{ countryFlag(ipData?.Country) }} {{ ipData?.Country || '-' }}</span>
                </div>
            </div>
            <div class="ip-detect-item">
                <div class="ip-detect-label">城市</div>
                <div class="ip-detect-value">
                    <span>{{ ipData?.City || '-' }}</span>
                </div>
            </div>
            <div class="ip-detect-item">
                <div class="ip-detect-label">IP 类型</div>
                <div class="ip-detect-value">
                    <span :class="typeClass(ipData?.IPType)">{{ translateType(ipData?.IPType) }}</span>
                </div>
            </div>
            <div class="ip-detect-item">
                <div class="ip-detect-label">运营商类型</div>
                <div class="ip-detect-value">
                    <span :class="typeClass(ipData?.OrgType)">{{ translateType(ipData?.OrgType) }}</span>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import Select from '@/component/ui/Select.vue'
import { getIP, fetchIP } from '@/api/detect'
import { translateType, typeClass } from '@/util/iptype'
import { formatTime } from '@/util/format'

const allData = ref<any[]>([])
const ipData = ref<any>(null)
const source = ref('ipapi.is')
const loading = ref(false)

const sourceOptions = [
    { label: 'ipapi.is',   value: 'ipapi.is'   },
    { label: 'ip-api.com', value: 'ip-api.com' },
    { label: 'ippure.com', value: 'ippure.com' },
]

function countryFlag(code?: string) {
    if (!code) return ''
    return code.toUpperCase().replace(/./g, c =>
        String.fromCodePoint(c.charCodeAt(0) + 127397)
    )
}

async function handleFetch() {
    loading.value = true
    try {
        const res = await fetchIP(source.value)
        ipData.value = res.data
        const idx = allData.value.findIndex((r: any) => r.Source === source.value && r.IPVersion === 'v4')
        if (idx >= 0) {
            allData.value[idx] = res.data
        } else {
            allData.value.push(res.data)
        }
    } finally {
        loading.value = false
    }
}

onMounted(async () => {
    const res = await getIP()
    if (res.data?.length) {
        allData.value = res.data
        ipData.value = res.data.find((r: any) => r.Source === source.value && r.IPVersion === 'v4') ?? null
    }
})

watch(source, (val) => {
    ipData.value = allData.value.find((r: any) => r.Source === val && r.IPVersion === 'v4') ?? null
})

watch(loading, (val) => {
    document.body.style.cursor = val ? 'wait' : ''
})
</script>

<style scoped>
.ip-detect-body {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 10px;
}

.ip-detect-item {
    display: flex;
    flex-direction: column;
    gap: 5px;
    min-width: 0;
}

.ip-detect-label {
    font-size: 12px;
    color: var(--color-text-dark);
}

.ip-detect-value {
    display: flex;
    flex-direction: row;
    align-items: center;
    gap: 5px;
    color: var(--color-text-light);
    font-weight: bold;
    font-size: 14px;
    min-width: 0;

    span {
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
    }

    .icon {
        font-size: 16px;
        flex-shrink: 0;
    }
}

@media (max-width: 1024px) {
    .ip-detect-body {
        grid-template-columns: repeat(2, 1fr);
    }
}

@media (max-width: 768px) {
    .ip-detect-body {
        grid-template-columns: repeat(1, 1fr);
    }
}
</style>