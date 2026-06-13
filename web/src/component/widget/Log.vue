<template>
    <div class="log">
        <div class="toolbar">
            <label class="checkbox">
                <input type="checkbox" v-model="autoScroll" />
                自动滚动
            </label>
            <div class="toolbar-right">
                <Select v-model="filterLevel" :options="levelOptions" small/>
                <button @click="clear">清空</button>
            </div>
        </div>
        <div class="output" ref="outputRef">
            <template v-for="(line, i) in filteredLines" :key="i">
                <div class="line">
                    <div class="meta">
                        <span class="time">{{ line.time }}</span>
                        <span class="level" :class="line.level.toLowerCase()">{{ line.level }}</span>
                    </div>
                    <span class="message">{{ line.message }}</span>
                </div>
            </template>
        </div>
    </div>
</template>

<script setup lang="ts">
import Select from '@/component/ui/Select.vue'

const props = defineProps<{
    ws: () => WebSocket
}>()

interface LogLine {
    time: string
    level: string
    message: string
}

const lines = ref<LogLine[]>([])
const autoScroll = ref(true)
const filterLevel = ref('ALL')
const outputRef = ref<HTMLElement>()
let socket: WebSocket

const levelOptions = [
    { label: '全部', value: 'ALL' },
    { label: 'INFO', value: 'INFO' },
    { label: 'WARN', value: 'WARN' },
    { label: 'ERROR', value: 'ERROR' },
]

const filteredLines = computed(() => {
    if (filterLevel.value === 'ALL') return lines.value
    return lines.value.filter(l => l.level === filterLevel.value)
})

function connect() {
    socket = props.ws()
    socket.onmessage = (e) => {
        const line = JSON.parse(e.data) as LogLine
        lines.value.push(line)
        if (lines.value.length > 500) {
            lines.value.shift()
        }
        if (autoScroll.value) {
            nextTick(() => {
                outputRef.value?.scrollTo(0, outputRef.value.scrollHeight)
            })
        }
    }
    socket.onclose = () => {
        setTimeout(connect, 3000)
    }
}

function clear() {
    lines.value = []
}

onMounted(() => connect())
onUnmounted(() => socket?.close())
</script>

<style scoped>
.log {
    display: flex;
    flex-direction: column;
    gap: 20px;
    height: 100%;
    background-color: var(--color-bg-dark);
    padding: 20px;
    border-radius: 20px;

    .toolbar {
        display: flex;
        align-items: center;
        justify-content: space-between;
        font-size: var(--font-size-sm);
        color: var(--color-text-light);

        label {
            background-color: var(--color-bg-light);
            border-radius: 999px;
            padding: 5px 10px;
        }

        .toolbar-right {
            display: flex;
            align-items: center;
            gap: 5px;
        }

        button {
            padding: 5px 15px;
            border-radius: 999px;
            font-size: var(--font-size-sm);
            background-color: var(--color-bg-light);

            &:hover {
                background-color: var(--color-primary);
            }

            &:active {
                background-color: var(--color-primary-dark);
            }
        }
    }

    .output {
        flex: 1;
        overflow-y: auto;
        background: var(--color-bg);
        border-radius: 12px;
        padding: 12px;
        font-family: monospace;
        font-size: 12px;
        line-height: 1.8;

        .line {
            display: flex;
            gap: 10px;

            .meta {
                display: flex;
                flex-direction: row;
                flex-wrap: nowrap;
                gap: 10px;

                .time {
                    color: var(--color-text-dark);
                    white-space: nowrap;
                    flex-shrink: 0;
                }

                .level {
                    white-space: nowrap;
                    flex-shrink: 0;
                    font-weight: bold;

                    &.info  { color: var(--color-green); }
                    &.warn  { color: var(--color-yellow); }
                    &.error { color: var(--color-red); }
                }

            }

            .message {
                color: var(--color-text-light);
                word-break: break-all;
            }
        }
    }
}

@media (max-width: 768px) {
    .log {
        .output {
            .line {
                flex-direction: column;
                gap: 0;
            }
        }
    }
}
</style>