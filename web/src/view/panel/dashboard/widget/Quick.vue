<template>
    <div class="card md">
        <div class="header">
            <div class="title">一键生成</div>
        </div>
        <div class="desc">
            <p>自动生成 VLESS + Reality + XTLS 节点及订阅用户</p>
            <p>无需证书，开箱即用</p>
        </div>
        <div class="actions">
            <button class="action-btn" @click="handleQuick">
                立即生成
            </button>
        </div>
    </div>
</template>

<script setup lang="ts">
import { quick } from '@/api/quick'

const modal = inject<any>('modal')

async function handleQuick() {
    modal.value?.show('confirm', '将自动生成 VLESS + Reality + XTLS 入站及用户\n确认继续？', async () => {
        modal.value?.update('warn', '生成中...')
        try {
            await quick()
            modal.value?.update('success', '生成成功！')
        } catch (err: any) {
            modal.value?.update('error', err?.error ?? '生成失败')
        }
    })
}
</script>

<style scoped>
.desc {
    font-size: var(--font-size-sm);
    color: var(--color-text-dark);
    line-height: 1.5;
    flex: 1;
}

.actions {
    display: flex;
    margin-top: auto;
    margin-left: calc(-1 * var(--card-padding, 20px));
    margin-right: calc(-1 * var(--card-padding, 20px));
    margin-bottom: calc(-1 * var(--card-padding, 20px));
    border-top: 1px solid var(--color-bg);
    overflow: hidden;

    .action-btn {
        flex: 1;
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 10px;
        padding: 10px;
        background-color: transparent;
        color: var(--color-text);
        font-size: var(--font-size-sm);
        border-radius: 0 0 20px 20px;

        &:hover {
            background-color: var(--color-primary);
            color: var(--color-text-light);
        }
        &:active {
            background-color: var(--color-primary-dark);
        }
    }
}
</style>