<!-- component/module/Task.vue -->
<template>
    <div v-if="lines.length" ref="output" class="task">
        <div v-for="(line, i) in lines" :key="i" class="line">
            <div class="meta">
                <span class="time">{{ line.time }}</span>
                <span class="level" :class="line.level.toLowerCase()">{{ line.level }}</span>
            </div>
            <span class="message">{{ line.message }}</span>
        </div>
    </div>
</template>

<script setup lang="ts">
import { task } from '@/api/task'

// 后台任务的实时输出(申请证书、注册 ACME 账号)：给任务号就连上去，一行行显示；
// 收到 DONE 就是完成，收到 ERROR 就是失败，两种情况都会断开
const props = defineProps<{
    taskId: string
}>()

const emit = defineEmits<{
    done: []
    error: []
}>()

interface Line {
    time: string
    level: string
    message: string
}

const lines = ref<Line[]>([])
const output = ref<HTMLElement>()

let socket: WebSocket | null = null

function connect(id: string) {
    socket?.close()
    lines.value = []

    if (!id) return

    socket = task(id)

    socket.onmessage = (e) => {
        const line = JSON.parse(e.data) as Line

        if (line.level === 'DONE') {
            emit('done')
            socket?.close()

            return
        }

        lines.value.push(line)
        nextTick(() => output.value?.scrollTo(0, output.value.scrollHeight))

        if (line.level === 'ERROR') {
            emit('error')
            socket?.close()
        }
    }
}

watch(() => props.taskId, connect, { immediate: true })

onUnmounted(() => socket?.close())
</script>

<style scoped>
.task {
    max-height: 300px;
    padding: 15px;
    overflow-y: auto;
    border-radius: 15px;
    background-color: var(--background-200);
    font-family: monospace;
    font-size: 12px;
    line-height: 1.8;
    text-align: left;

    .line {
        display: flex;
        gap: 10px;

        @media (max-width: 768px) {
            flex-direction: column;
            gap: 0;
        }
    }

    .meta {
        display: flex;
        flex-shrink: 0;
        gap: 10px;
    }

    .time {
        color: var(--text-300);
        white-space: nowrap;
    }

    .level {
        font-weight: bold;
        white-space: nowrap;

        &.info {
            color: var(--green);
        }

        &.warn {
            color: var(--yellow);
        }

        &.error {
            color: var(--red);
        }
    }

    .message {
        color: var(--text-100);
        word-break: break-all;
    }
}
</style>
