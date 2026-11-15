<template>
    <Box :title="core?.Name ?? ''">
        <template #action>
            <Status :on="running">{{ running ? '运行中' : '未启动' }}</Status>
        </template>
        <template #footer>
            <Button @click="running ? stop() : start()">
                <Icon name="power_settings_new" />
                {{ running ? '停止' : '启动' }}
            </Button>
            <Button @click="restart">
                <Icon name="refresh" />
                重启
            </Button>
            <Button disabled>
                <Icon name="commit" />
                {{ core?.Version }}
            </Button>
        </template>
    </Box>
</template>

<script setup lang="ts">
import Box from '@/component/layout/Box.vue'
import Button from '@/component/widget/Button.vue'
import Icon from '@/component/widget/Icon.vue'
import Status from '@/component/widget/Status.vue'
import { getCore, getCoreStatus, stopCore, startCore, restartCore } from '@/api/core'
import type { Core } from '@/api/core'
import { toast } from '@/composable/toast'

const core = ref<Core>()
const running = ref(false)

async function load() {
    const [c, s] = await Promise.all([getCore(), getCoreStatus()])
    core.value = c
    running.value = s.status === 'running'
}

// 操作失败时提示接口给的原因
async function run(action: () => Promise<unknown>) {
    try {
        await action()
        await load()
    } catch (err) {
        toast.failure((err as { error?: string })?.error || '操作失败')
    }
}

const start = () => run(startCore)
const stop = () => run(stopCore)
const restart = () => run(restartCore)

let timer: ReturnType<typeof setInterval>

onMounted(() => {
    load()
    timer = setInterval(load, 5000)
})

onUnmounted(() => clearInterval(timer))
</script>
