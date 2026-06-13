<template>
    <div class="task" v-if="lines.length">
        <div class="output" ref="outputRef">
            <div class="line" v-for="(line, i) in lines" :key="i">
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
import { task } from '@/api/task'

const props = defineProps<{ taskId: string }>()
const emit = defineEmits<{ done: [], error: [] }>()

interface LogLine { time: string; level: string; message: string }

const lines = ref<LogLine[]>([])
const outputRef = ref<HTMLElement>()

watch(() => props.taskId, (id) => {
    if (!id) return
    lines.value = []
    const ws = task(id)
    ws.onmessage = (e) => {
        const line = JSON.parse(e.data) as LogLine
        if (line.level === 'DONE') {
            emit('done')
            ws.close()
            return
        }
        lines.value.push(line)
        if (line.level === 'ERROR') {
            emit('error')
            ws.close()
            return
        }
        nextTick(() => {
            outputRef.value?.scrollTo(0, outputRef.value.scrollHeight)
        })
    }
}, { immediate: true })
</script>

<style scoped>
.task {
    margin: 0 20px 20px;
    height: 500px;
    border-radius: 10px;
    overflow: hidden;
    background: var(--color-bg);

    .output {
        padding: 12px;
        max-height: 500px;
        overflow-y: auto;
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
    .task {
        .output {
            .line {
                flex-direction: column;
                gap: 0;
            }
        }
    }
}
</style>