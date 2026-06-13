<template>
    <div class="card md">
        <div class="header">
            <div class="title">{{ core.Name }}</div>
            <div class="status">
                <Status :status="status === 'running' ? 'yes' : 'no'" :animated="status === 'running'" />
                <span>{{ status === 'running' ? '运行中' : '未启动' }}</span>
            </div>
        </div>
        <div class="actions">
            <button class="action-btn" @click="status === 'running' ? handleStop() : handleStart()">
                <i class="icon">power_settings_new</i>
                {{ status === 'running' ? '停止' : '启动' }}
            </button>
            <button class="action-btn" @click="handleRestart">
                <i class="icon">refresh</i>
                重启
            </button>
            <button class="action-btn version">
                <i class="icon">commit</i>
                {{ core.Version }}
            </button>
        </div>
    </div>
</template>

<script setup lang="ts">
import Status from '@/component/ui/Status.vue'
import { getCore, getCoreStatus, stopCore, startCore, restartCore } from '@/api/core'

const core = ref<any>({})
const status = ref<string>('stopped')

async function fetch() {
    const [c, s] = await Promise.all([getCore(), getCoreStatus()])
    core.value = c
    status.value = s.status
}

async function handleStart() {
    await startCore()
    await fetch()
}

async function handleStop() {
    await stopCore()
    status.value = 'stopped'
}

async function handleRestart() {
    await restartCore()
    await fetch()
}

onMounted(() => {
    fetch()
    const timer = setInterval(fetch, 5000)
    onUnmounted(() => clearInterval(timer))
})
</script>

<style scoped>
.header {
    justify-content: space-between;

    .status {
        display: flex;
        align-items: center;
        gap: 8px;
        font-size: var(--font-size-sm);
        color: var(--color-text-dark);
    }
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
        gap: 6px;
        padding: 10px;
        background: transparent;
        border-right: 1px solid var(--color-bg);
        color: var(--color-text);
        font-size: var(--font-size-sm);

        &:first-child { border-radius: 0 0 0 20px; }
        &:nth-child(2) { border-radius: 0; }
        &:last-child { 
            border-right: none; 
            border-radius: 0 0 20px 0;
        }

        &:hover {
            background: var(--color-primary);
            color: var(--color-text-light);
        }
    }
}
</style>