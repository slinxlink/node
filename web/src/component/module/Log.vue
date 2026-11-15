<!-- component/module/Log.vue -->
<template>
    <div class="log">
        <div class="toolbar">
            <label class="auto">
                <Toggle v-model:value="autoScroll" />
                自动滚动
            </label>
            <div class="right">
                <Select v-model:value="level" :items="levels" />
                <Button class="clear" @click="clear">清空</Button>
            </div>
        </div>
        <div ref="output" class="output">
            <div v-for="(line, i) in shown" :key="i" class="line">
                <div class="meta">
                    <span class="time">{{ line.time }}</span>
                    <span class="level" :class="line.level.toLowerCase()">{{ line.level }}</span>
                </div>
                <span class="message">{{ line.message }}</span>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import Button from '@/component/widget/Button.vue'
import Select from '@/component/widget/Select.vue'
import Toggle from '@/component/widget/Toggle.vue'

// 实时日志：连上 WebSocket 一行行收，最多留 500 行，断开 3 秒后自动重连；可按级别过滤、清空、自动滚到底
const props = defineProps<{
    ws: () => WebSocket
}>()

interface Line {
    time: string
    level: string
    message: string
}

const MAX = 500

const lines = ref<Line[]>([])
const autoScroll = ref(true)
const level = ref('ALL')
const output = ref<HTMLElement>()

const levels = [
    { label: '全部', value: 'ALL' },
    { label: 'INFO', value: 'INFO' },
    { label: 'WARN', value: 'WARN' },
    { label: 'ERROR', value: 'ERROR' },
]

const shown = computed(() =>
    level.value === 'ALL' ? lines.value : lines.value.filter((item) => item.level === level.value),
)

let socket: WebSocket | null = null
let timer = 0
let closed = false

function connect() {
    socket = props.ws()

    socket.onmessage = (e) => {
        lines.value.push(JSON.parse(e.data) as Line)

        if (lines.value.length > MAX) {
            lines.value.shift()
        }

        if (autoScroll.value) {
            nextTick(() => output.value?.scrollTo(0, output.value.scrollHeight))
        }
    }

    socket.onclose = () => {
        if (!closed) timer = window.setTimeout(connect, 3000)
    }
}

function clear() {
    lines.value = []
}

onMounted(connect)

onUnmounted(() => {
    closed = true
    clearTimeout(timer)
    socket?.close()
})
</script>

<style scoped>
.log {
    display: flex;
    flex: 1; /* 放在固定高度的容器里：撑满剩下的高度，里面的日志自己滚 */
    flex-direction: column;
    gap: 20px;
    min-height: 0;

    .toolbar {
        display: flex;
        align-items: center;
        justify-content: space-between;
        gap: 10px;
        color: var(--text-300);
        font-size: 12px;

        .auto {
            display: flex;
            align-items: center;
            gap: 8px;
            cursor: pointer;
        }

        .right {
            display: flex;
            align-items: center;
            gap: 10px;
        }

        .right .select {
            width: 100px;
        }

        .clear {
            width: auto;
            background-color: var(--background-200);

            &:hover:not(:disabled),
            &:focus-visible {
                background-color: var(--background-100);
            }

            &:active:not(:disabled) {
                background-color: var(--primary-600);
            }
        }
    }

    .output {
        flex: 1;
        min-height: 0;
        padding: 15px;
        overflow-y: auto;
        border-radius: 15px;
        background-color: var(--background-200);
        font-family: monospace;
        font-size: 12px;
        line-height: 1.8;

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
}
</style>
