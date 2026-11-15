<!-- view/panel/detect/Unlock.vue -->
<template>
    <Box title="解锁检测">
        <template #action>
            <span v-if="data.length">{{ formatTime(data[0]?.UpdatedAt ?? '') }}</span>
            <Button class="run" :disabled="loading" @click="detect">
                <Icon :name="loading ? 'hourglass_empty' : 'travel_explore'" />
                <span>{{ loading ? '检测中' : '开始检测' }}</span>
            </Button>
        </template>

        <div v-for="group in groups" :key="group.label" class="group">
            <div class="title">{{ group.label }}</div>
            <div class="grid">
                <div v-for="platform in group.platforms" :key="platform" class="card">
                    <div class="head">
                        <img :src="image(platform)" :alt="platform" class="logo" />
                        <span class="name">{{ names[platform] }}</span>
                        <Tip v-if="notes[platform]" :msg="notes[platform]">
                            <Icon name="help" class="help" />
                        </Tip>
                    </div>
                    <div class="status" :style="{ color: result(platform).color }">
                        <Icon v-if="result(platform).icon" :name="result(platform).icon" />
                        <span>{{ result(platform).text }}</span>
                        <span v-if="region(platform)" class="region">
                            {{ countryFlag(region(platform)) }} {{ region(platform) }}
                        </span>
                    </div>
                </div>
            </div>
        </div>
    </Box>
</template>

<script setup lang="ts">
import Box from '@/component/layout/Box.vue'
import Button from '@/component/widget/Button.vue'
import Icon from '@/component/widget/Icon.vue'
import Tip from '@/component/widget/Tip.vue'
import { getUnlock, fetchUnlock } from '@/api/detect'
import type { Unlock } from '@/api/detect'
import { toast } from '@/composable/toast'
import { countryFlag } from '@/util/country'
import { formatTime } from '@/util/format'

const data = ref<Unlock[]>([])
const loading = ref(false)

const images = import.meta.glob('@/asset/image/platform/*.webp', {
    eager: true,
    import: 'default',
})

function image(platform: string): string {
    return (images[`/src/asset/image/platform/${platform}.webp`] as string) ?? ''
}

const groups = [
    {
        label: '流媒体',
        platforms: ['netflix', 'disney', 'prime_video', 'youtube_premium', 'dazn', 'tvbanywhere', 'iqiyi'],
    },
    {
        label: 'AI 工具',
        platforms: ['claude', 'chatgpt', 'gemini'],
    },
    {
        label: '其他',
        platforms: ['apple', 'bing', 'google_play', 'steam', 'reddit', 'wikipedia'],
    },
]

const names: Record<string, string> = {
    apple: 'Apple',
    bing: 'Bing',
    google_play: 'Google Play Store',
    steam: 'Steam',
    reddit: 'Reddit',
    wikipedia: 'Wikipedia Editability',
    netflix: 'Netflix',
    disney: 'Disney+',
    prime_video: 'Prime Video',
    youtube_premium: 'YouTube Premium',
    dazn: 'DAZN',
    tvbanywhere: 'TVBAnywhere+',
    iqiyi: 'iQiyi Oversea',
    claude: 'Claude',
    chatgpt: 'ChatGPT',
    gemini: 'Gemini',
}

const notes: Record<string, string> = {
    bing: '受限表示该 IP 被 Bing 标记为风险 IP',
    netflix: '受限表示仅可观看 Netflix 自制内容',
    iqiyi: '受限表示仅海外通用线路，无本地化内容',
    chatgpt: '受限表示仅 Web 或仅 Mobile App 可用',
}

function record(platform: string) {
    return data.value.find((item) => item.Platform === platform)
}

function region(platform: string): string {
    return record(platform)?.Region ?? ''
}

// 状态：true 可用、false 不可用、reject 受限，没有记录就是短横线
function result(platform: string) {
    const status = record(platform)?.Status

    if (status === 'true') return { icon: 'check_circle', text: '可用', color: 'var(--green)' }
    if (status === 'false') return { icon: 'cancel', text: '不可用', color: 'var(--red)' }
    if (status === 'reject') return { icon: 'do_not_disturb_on', text: '受限', color: 'var(--yellow)' }

    return { icon: '', text: '-', color: 'var(--text-300)' }
}

async function detect() {
    loading.value = true

    try {
        data.value = await fetchUnlock()
    } catch (err) {
        toast.failure((err as { error?: string })?.error || '检测失败')
    } finally {
        loading.value = false
    }
}

onMounted(async () => {
    data.value = await getUnlock()
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

.grid {
    display: grid;
    grid-template-columns: repeat(4, minmax(0, 1fr));
    gap: 10px;

    @media (max-width: 1024px) {
        grid-template-columns: repeat(3, minmax(0, 1fr));
    }

    @media (max-width: 768px) {
        grid-template-columns: repeat(2, minmax(0, 1fr));
    }
}

.card {
    display: flex;
    flex-direction: column;
    gap: 8px;
    min-width: 0;
    padding: 12px;
    border-radius: 12px;
    background-color: var(--background-200);

    .head {
        display: flex;
        align-items: center;
        gap: 8px;
        min-width: 0;
    }

    .logo {
        flex-shrink: 0;
        width: 20px;
        height: 20px;
        border-radius: 4px;
    }

    .name {
        overflow: hidden;
        color: var(--text-100);
        font-weight: bold;
        text-overflow: ellipsis;
        white-space: nowrap;
    }

    .help {
        flex-shrink: 0;
        color: var(--text-300);
        font-size: 14px;
        cursor: help;
    }

    .status {
        display: flex;
        align-items: center;
        gap: 6px;
        font-size: 12px;
    }

    .region {
        color: var(--text-200);
    }
}
</style>
