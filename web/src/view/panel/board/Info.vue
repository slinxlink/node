<!-- view/panel/board/Info.vue -->
<template>
    <div class="info">
        <div v-if="loading" class="empty">加载中...</div>
        <div v-else-if="users.length === 0" class="empty">无用户数据</div>
        <div v-else class="grid">
            <div v-for="item in users" :key="item.ID" class="card">
                <div class="top">
                    <span class="id">
                        用户 <strong>#{{ item.UserID }}</strong>
                    </span>
                    <Alive :ips="parseIPs(item.AliveIP)" />
                </div>
                <div class="key">
                    <div class="label">{{ protocol === 'hysteria' ? 'Password' : 'UUID' }}</div>
                    <div class="value">{{ protocol === 'hysteria' ? item.Passwd : item.UUID }}</div>
                </div>
                <div class="traffic">
                    <div class="item">
                        <Icon name="upload" />
                        <span class="label">上传</span>
                        <span class="value">{{ formatBytes(item.Upload) }}</span>
                    </div>
                    <div class="item">
                        <Icon name="download" />
                        <span class="label">下载</span>
                        <span class="value">{{ formatBytes(item.Download) }}</span>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import Icon from '@/component/widget/Icon.vue'
import Alive from '@/view/panel/board/Alive.vue'
import { getBoardUser } from '@/api/board'
import type { Board, BoardUser } from '@/api/board'
import { getInbounds } from '@/api/inbound'
import { formatBytes } from '@/util/format'

// 面板那边的用户：每人一张卡片，密钥、流量、在线 IP
const props = defineProps<{
    board: Board
}>()

const users = ref<BoardUser[]>([])
const protocol = ref('')
const loading = ref(true)

function parseIPs(text: string): string[] {
    try {
        const list = JSON.parse(text)

        return Array.isArray(list) ? list : []
    } catch {
        return []
    }
}

// hysteria 用密码认证，其余用 UUID
onMounted(async () => {
    const [list, inbounds] = await Promise.all([getBoardUser(props.board.ID), getInbounds()])

    users.value = list
    protocol.value = inbounds.find((item) => item.ID === props.board.Inbound)?.Protocol ?? ''
    loading.value = false
})
</script>

<style scoped>
.info {
    width: 100%;
    text-align: left;

    .grid {
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
        gap: 10px;
    }

    .card {
        display: flex;
        flex-direction: column;
        gap: 10px;
        padding: 14px 16px;
        border-radius: 12px;
        background-color: var(--background-200);

        .top {
            display: flex;
            align-items: center;
            justify-content: space-between;

            .id {
                color: var(--text-300);
                font-size: 12px;

                strong {
                    color: var(--text-100);
                    font-size: 14px;
                    font-weight: 500;
                }
            }
        }

        .key {
            padding: 8px 10px;
            border-radius: 8px;
            background-color: var(--background-300);

            .label {
                margin-bottom: 3px;
                color: var(--text-300);
                font-size: 12px;
            }

            .value {
                color: var(--text-100);
                font-family: monospace;
                font-size: 12px;
                line-height: 1.4;
                word-break: break-all;
            }
        }

        .traffic {
            display: flex;
            align-items: center;
            gap: 5px;

            .item {
                display: flex;
                flex: 1;
                align-items: center;
                gap: 5px;
                color: var(--text-300);
                font-size: 12px;

                .value {
                    color: var(--text-100);
                    font-weight: bold;
                }
            }
        }
    }

    .empty {
        padding: 40px 0;
        color: var(--text-300);
        font-size: 12px;
        text-align: center;
    }
}
</style>
