<template>
    <main>
        <div class="page">
            <Tab v-model="activeTab" :tabs="tabs">
                <div v-if="activeTab === 'settings'" class="body">
                    <RestartBar 
                        :dirty="dirty" 
                        :on-save="save" 
                        :on-restart="restart" 
                        tip="此处的所有更改保存后将自动重启核心生效"
                        restart-text="重启核心"
                        @update:dirty="dirty = $event" 
                    />
                    <CoreCard />
                    <Section title="日志">
                        <Row title="启用" subtitle="关闭后面板将不记录任何日志">
                            <Toggle v-model="Core.LogEnable" />
                        </Row>
                        <Row title="路径" subtitle="日志文件仅支持保存在面板根目录下">
                            <Input v-model="Core.LogPath" :prefix="dir + '/'" placeholder="日志路径" />
                        </Row>
                        <Row title="等级" subtitle="等级越高，记录的日志越少">
                            <Select v-model="Core.LogLevel" :options="logLevelOptions" placeholder="日志等级" />
                        </Row>
                    </Section>
                    <Section title="重置" :default-open="false">
                        <Row title="恢复默认" subtitle="重置所有核心配置为默认值并重启核心">
                            <button class="reset" @click="reset">重置</button>
                        </Row>
                    </Section>
                </div>
                <div v-if="activeTab === 'config'" class="config-wrap">
                    <JsonEditor :content="coreConfig" />
                </div>
                <div v-if="activeTab === 'log'" class="log-wrap">
                    <Log :ws="log" />
                </div>
            </Tab>
        </div>
    </main>
</template>

<script setup lang="ts">
import Tab from '@/component/ui/Tab.vue'
import Section from '@/component/ui/Section.vue'
import Row from '@/component/ui/Row.vue'
import Input from '@/component/ui/Input.vue'
import Select from '@/component/ui/Select.vue'
import Toggle from '@/component/ui/Toggle.vue'
import Log from '@/component/widget/Log.vue'
import JsonEditor from '@/component/widget/JsonEditor.vue'
import RestartBar from '@/component/widget/RestartBar.vue'
import CoreCard from '@/component/widget/CoreCard.vue'
import { getConfig } from '@/api/config'
import { Log as log, getCoreConfig, getCore, updateCore, resetCore, restartCore } from '@/api/core'

const modal = inject<any>('modal')

const Core = ref<any>({ LogEnable: true, LogLevel: 'info' })
const dirty = ref(false)
const loaded = ref(false)

const activeTab = ref('settings')
const tabs = [
    { key: 'settings', label: '设置' },
    { key: 'config',   label: '高级配置' },
    { key: 'log',      label: '日志' },
]

const coreConfig = ref<string>('')

const logLevelOptions = [
    { label: 'TRACE', value: 'trace' },
    { label: 'DEBUG', value: 'debug' },
    { label: 'INFO',  value: 'info' },
    { label: 'WARN',  value: 'warn' },
    { label: 'ERROR', value: 'error' },
    { label: 'FATAL', value: 'fatal' },
    { label: 'PANIC', value: 'panic' },
]

const dir = ref('')

watch(Core, () => {
    if (loaded.value) dirty.value = true
}, { deep: true })

watch(activeTab, async (val) => {
    if (val === 'config') {
        coreConfig.value = await getCoreConfig()
    }
})

onMounted(() => {
    getConfig().then(cfg => dir.value = cfg.Dir)
    getCore().then(core => {
        Core.value = core
        nextTick(() => loaded.value = true)
    })
})

async function save() {
    try {
        modal.value?.show('warn', '保存中...')
        await updateCore(Core.value)
        modal.value?.show('success', '保存成功')
    } catch (err: any) {
        modal.value?.show('error', err?.error)
        throw err
    }
}

async function restart() {
    try {
        modal.value?.show('warn', '重启中...')
        await restartCore()
        modal.value?.show('success', '核心重启成功')
    } catch (err: any) {
        modal.value?.show('error', err?.error)
    }
}

async function reset() {
    modal.value?.show('confirm', '确认重置核心配置为默认值？', async () => {
        await resetCore()
        Core.value = await getCore()
        modal.value?.show('success', '已重置')
    })
}
</script>

<style scoped>
.log-wrap {
    height: calc(100vh - 180px);
    display: flex;
    flex-direction: column;
}

:deep(.card.md) {
    width: 100%;
}

.reset {
    padding: 10px 40px;
    font-size: var(--font-size-sm);
    border: 1px solid var(--color-red);
    background-color: var(--color-bg-dark);
    color: var(--color-red);

    &:hover {
        border-color: var(--color-primary-light);
        color: var(--color-primary-light);
    }
}
</style>