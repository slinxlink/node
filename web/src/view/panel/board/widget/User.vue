<template>
    <Drawer v-model="show" title="机场用户" :footer="false">
        <div class="user-list">
            <div v-if="loading" class="empty">加载中...</div>
            <div v-else-if="users.length === 0" class="empty">无用户数据</div>
            <div v-else class="grid">
                <div class="card" v-for="u in users" :key="u.ID">
                    <div class="top">
                        <span class="id">用户 <strong>#{{ u.UserID }}</strong></span>
                        <div class="status">
                            <Status :status="u.AliveIP > 0 ? 'yes' : 'no'" :animated="u.AliveIP > 0" />
                            <span v-if="u.AliveIP > 0" class="ip">{{ u.AliveIP }} IP</span>
                        </div>
                    </div>
                    <div class="key">
                        <div class="label">{{ protocol === 'hysteria' ? 'Password' : 'UUID' }}</div>
                        <div class="value">{{ protocol === 'hysteria' ? u.Passwd : u.UUID }}</div>
                    </div>
                    <div class="traffic">
                        <div class="item">
                            <i class="icon">upload</i>
                            <span class="label">上传</span>
                            <span class="value">{{ formatBytes(u.Upload) }}</span>
                        </div>
                        <div class="item">
                            <i class="icon">download</i>
                            <span class="label">下载</span>
                            <span class="value">{{ formatBytes(u.Download) }}</span>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </Drawer>
</template>

<script setup lang="ts">
import Drawer from '@/component/Drawer.vue'
import Status from '@/component/ui/Status.vue'
import { getBoardUser } from '@/api/board'
import { getInbounds } from '@/api/inbound'
import { formatBytes } from '@/util/format'

const props = defineProps<{
    modelValue: boolean
    board: any
}>()

const emit = defineEmits(['update:modelValue'])

const show = computed({
    get: () => props.modelValue,
    set: (v) => emit('update:modelValue', v),
})

const users = ref<any[]>([])
const loading = ref(false)
const protocol = ref('')

watch(() => props.modelValue, async (val) => {
    if (!val || !props.board) return
    loading.value = true
    const [us, inbounds] = await Promise.all([
        getBoardUser(props.board.ID),
        getInbounds(),
    ])
    users.value = us
    protocol.value = inbounds.find((i: any) => i.ID === props.board.Inbound)?.Protocol ?? ''
    loading.value = false
})
</script>

<style scoped>
.user-list {
    padding: 20px;

    .grid {
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
        gap: 10px;

        @media (max-width: 768px) {
            grid-template-columns: 1fr;
        }

        .card {
            background: var(--color-bg);
            border-radius: 12px;
            padding: 14px 16px;
            display: flex;
            flex-direction: column;
            gap: 10px;

            .top {
                display: flex;
                align-items: center;
                justify-content: space-between;

                .id {
                    font-size: var(--font-size-sm);
                    color: var(--color-text-dark);

                    strong {
                        font-size: var(--font-size);
                        font-weight: 500;
                        color: var(--color-text);
                    }
                }

                .status {
                    display: flex;
                    align-items: center;
                    gap: 6px;

                    .ip {
                        font-size: var(--font-size-sm);
                        color: var(--color-green);
                    }
                }
            }

            .key {
                background: var(--color-bg-dark);
                border-radius: 8px;
                padding: 8px 10px;

                .label {
                    font-size: var(--font-size-sm);
                    color: var(--color-text-dark);
                    margin-bottom: 3px;
                }

                .value {
                    font-size: var(--font-size-sm);
                    font-family: monospace;
                    color: var(--color-text);
                    word-break: break-all;
                    line-height: 1.4;
                }
            }

            .traffic {
                display: flex;
                gap: 5px;
                align-items: center;

                .item {
                    flex: 1;
                    display: flex;
                    align-items: center;
                    gap: 5px;

                    .icon {
                        font-size: var(--font-size-sm);
                        color: var(--color-text-dark);
                    }

                    .label {
                        font-size: var(--font-size-sm);
                        color: var(--color-text-dark);
                    }

                    .value {
                        font-size: var(--font-size-sm);
                        font-weight: bold;
                        color: var(--color-text);
                    }
                }
            }
        }
    }

    .empty {
        text-align: center;
        color: var(--color-text-dark);
        font-size: var(--font-size-sm);
        padding: 40px 0;
    }
}


</style>