<!-- components/UnlockDetect.vue -->
<template>
<div class="detect">
        <div class="header">
            <div class="title-wrap">
                <div class="title">解锁检测</div>
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

        <div v-for="group in groups" :key="group.label" class="unlock-group">
            <div class="unlock-group-title">{{ group.label }}</div>
            <div class="unlock-grid">
                <div v-for="platform in group.platforms" :key="platform" class="unlock-card">
                    <div class="unlock-card-header">
                        <img :src="getPlatformImage(platform)" :alt="platform" class="unlock-logo" />
                        <span class="unlock-name">{{ nameMap[platform] }}</span>
                        <Tip v-if="noteMap[platform]">
                            <i class="icon note-icon">help</i>
                            <template #content>{{ noteMap[platform] }}</template>
                        </Tip>
                    </div>
                    <div class="unlock-status">
                        <i v-if="getStatus(platform)" class="icon" :class="statusClass(getStatus(platform))">
                            {{ statusIcon(getStatus(platform)) }}
                        </i>
                        <span :class="statusClass(getStatus(platform))">
                            {{ statusText(getStatus(platform)) }}
                        </span>
                        <span v-if="getRegion(platform)" class="unlock-region">
                            {{ countryFlag(getRegion(platform)) }} {{ getRegion(platform) }}
                        </span>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import Tip from '@/component/ui/Tip.vue'
import { getUnlock, fetchUnlock } from '@/api/detect'
import { formatTime } from '@/util/format'

const loading = ref(false)
const data = ref<any[]>([])
const platformImages = import.meta.glob('@/asset/image/platform/*.webp', { eager: true, import: 'default' })

function getPlatformImage(platform: string): string {
    return (platformImages[`/src/asset/image/platform/${platform}.webp`] as string) ?? ''
}

const groups = [
    {
        label: '流媒体',
        platforms: ['netflix', 'disney', 'prime_video', 'youtube_premium', 'dazn', 'tvbanywhere', 'iqiyi']
    },
    {
        label: 'AI 工具',
        platforms: ['claude', 'chatgpt', 'gemini']
    },
    {
        label: '其他',
        platforms: ['apple', 'bing', 'google_play', 'steam', 'reddit', 'wikipedia']
    }
]

const nameMap: Record<string, string> = {
    'apple':           'Apple',
    'bing':            'Bing',
    'google_play':     'Google Play Store',
    'steam':           'Steam',
    'reddit':          'Reddit',
    'wikipedia':       'Wikipedia Editability',
    'netflix':         'Netflix',
    'disney':          'Disney+',
    'prime_video':     'Prime Video',
    'youtube_premium': 'YouTube Premium',
    'dazn':            'DAZN',
    'tvbanywhere':     'TVBAnywhere+',
    'iqiyi':           'iQiyi Oversea',
    'claude':          'Claude',
    'chatgpt':         'ChatGPT',
    'gemini':          'Gemini',
}

const noteMap: Record<string, string> = {
    'bing':    '受限表示该 IP 被 Bing 标记为风险 IP',
    'netflix': '受限表示仅可观看 Netflix 自制内容',
    'iqiyi':   '受限表示仅海外通用线路，无本地化内容',
    'chatgpt': '受限表示仅 Web 或仅 Mobile App 可用',
}

function getStatus(platform: string): string {
    const record = data.value.find((r: any) => r.Platform === platform)
    return record?.Status ?? ''
}

function getRegion(platform: string): string {
    const record = data.value.find((r: any) => r.Platform === platform)
    return record?.Region ?? ''
}

function statusIcon(status: string): string {
    if (status === 'true') return 'check_circle'
    if (status === 'false') return 'cancel'
    if (status === 'reject') return 'do_not_disturb_on'
    return ''
}

function statusText(status: string): string {
    if (status === 'true') return '可用'
    if (status === 'false') return '不可用'
    if (status === 'reject') return '受限'
    return '-'
}

function statusClass(status: string): string {
    if (status === 'true') return 'c-yes'
    if (status === 'false') return 'c-no'
    if (status === 'reject') return 'c-reject'
    return 'c-empty'
}

function countryFlag(code?: string) {
    if (!code) return ''
    return code.toUpperCase().replace(/./g, c =>
        String.fromCodePoint(c.charCodeAt(0) + 127397)
    )
}

async function handleFetch() {
    loading.value = true
    try {
        const res = await fetchUnlock()
        data.value = res.data ?? []
    } finally {
        loading.value = false
    }
}

onMounted(async () => {
    const res = await getUnlock()
    data.value = res.data ?? []
})


watch(loading, (val) => {
    document.body.style.cursor = val ? 'wait' : ''
})
</script>

<style scoped>
.unlock-group {
    display: flex;
    flex-direction: column;
    gap: 10px;
}

.unlock-group-title {
    font-size: 12px;
    color: var(--color-text-dark);
}

.unlock-grid {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 10px;
}

.unlock-card {
    display: flex;
    flex-direction: column;
    gap: 8px;
    padding: 12px;
    background-color: var(--color-bg);
    border-radius: 12px;
    min-width: 0;
}

.unlock-card-header {
    display: flex;
    align-items: center;
    gap: 8px;
    min-width: 0;
}

.unlock-logo {
    width: 20px;
    height: 20px;
    border-radius: 4px;
    flex-shrink: 0;
}

.unlock-name {
    font-size: 13px;
    font-weight: bold;
    color: var(--color-text-light);
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
}

.note-icon {
    font-size: 14px;
    color: var(--color-text-dark);
    flex-shrink: 0;
    cursor: help;
}

.unlock-status {
    display: flex;
    align-items: center;
    gap: 6px;
    font-size: 12px;
}

.unlock-region {
    color: var(--color-text);
    font-size: 12px;
}

@media (max-width: 1024px) {
    .unlock-grid {
        grid-template-columns: repeat(3, 1fr);
    }
}

@media (max-width: 768px) {
    .unlock-grid {
        grid-template-columns: repeat(2, 1fr);
    }
}
</style>