<template>
    <main>
        <div class="page">
            <Tab v-model="activeTab" :tabs="tabs">
                <div v-if="activeTab === 'settings'" class="body">
                    <RestartBar :dirty="dirty" :on-save="save" :on-restart="restart" @update:dirty="dirty = $event" />
                    <Section title="常规">
                        <Row title="当前面板地址" subtitle="备注">
                            <label>{{ panelUrl }}</label>
                        </Row>
                        <Row title="域名" subtitle="选择证书域名，设置后自动启用HTTPS">
                            <Select v-model="Config.Domain" :options="domainOptions" placeholder="选择域名" />
                        </Row>
                        <Row title="端口" subtitle="重启面板生效「范围1024 ~ 65535」">
                            <Input v-model="Config.Port" type="number" :min="1" :max="65535" placeholder="端口" />
                        </Row>
                        <Row title="路径" subtitle="必须以 '/' 开头">
                            <Input v-model="Config.Path" placeholder="路径" />
                        </Row>
                    </Section>
                    <Section title="日志">
                        <Row title="启用" subtitle="关闭后面板将不记录任何日志">
                            <Toggle v-model="Config.LogEnable" />
                        </Row>
                        <Row title="路径" subtitle="日志文件仅支持保存在面板根目录下">
                            <Input v-model="Config.LogPath" prefix="/etc/slinx/" placeholder="日志路径" />
                        </Row>
                        <Row title="等级" subtitle="等级越高，记录的日志越少">
                            <Select v-model="Config.LogLevel" :options="logLevelOptions" placeholder="日志等级" />
                        </Row>
                    </Section>
                    <Section title="管理员凭据">
                        <form @submit.prevent>
                            <Row title="用户名" subtitle="&nbsp;">
                                <Input :modelValue="Config.Username" placeholder="用户名" autocomplete="username" disabled />
                            </Row>
                            <Row title="原密码" subtitle="&nbsp;">
                                <Input v-model="credentials.old_password" autocomplete="current-password" placeholder="原密码" type="password" />
                            </Row>
                        </form>
                        <form @submit.prevent>
                            <Row title="新用户名" subtitle="留空则不修改">
                                <Input v-model="credentials.new_username" placeholder="新用户名" autocomplete="username" />
                            </Row>
                            <Row title="新密码" subtitle="留空则不修改">
                                <Input v-model="credentials.new_password" autocomplete="new-password" placeholder="新密码" type="password" />
                            </Row>
                            <Row title="">
                                <button @click="submitCredentials">确定</button>
                            </Row>
                        </form>
                    </Section>
                    <Section title="面板对接">
                        <Row title="启用" subtitle="关闭后停止所有面板对接与同步">
                            <Toggle v-model="Config.BoardEnable" />
                        </Row>
                    </Section>
                </div>
                <div v-if="activeTab === 'certs'" class="body">
                    <Cert />
                </div>
                <div v-if="activeTab === 'log'" class="log-wrap">
                    <Log :ws="log" />
                </div>
            </Tab>
        </div>
    </main>
</template>

<script setup lang="ts">
// Vue 组件
import Tab from '@/component/ui/Tab.vue'
import Section from '@/component/ui/Section.vue'
import Row from '@/component/ui/Row.vue'
import Input from '@/component/ui/Input.vue'
import Select from '@/component/ui/Select.vue'
import Toggle from '@/component/ui/Toggle.vue'
import RestartBar from '@/component/widget/RestartBar.vue'
import Log from '@/component/widget/Log.vue'
import Cert from '@/view/panel/config/cert/Cert.vue'

// API
import { getConfig, updateConfig, Log as log } from '@/api/config'
import { getCert } from '@/api/cert'
import { changeCredentials } from '@/api/auth'
import { restartPanel } from '@/api/bootstrap'

// Store
import { loadConfig } from '@/store/config'

const modal = inject<any>('modal')

// ── 状态 ──────────────────────────────────────────────────

const activeTab = ref('settings')
const tabs = [
    { key: 'settings', label: '设置' },
    { key: 'certs',    label: '证书' },
    { key: 'log',      label: '日志' },
]

const Config = ref<any>({ LogEnable: false, BoardEnable: false })
const certs = ref<any[]>([])
const credentials = ref({ old_password: '', new_username: '', new_password: '' })
const dirty = ref(false)
const loaded = ref(false)

// ── 计算属性 ───────────────────────────────────────────────

const domainOptions = computed(() => [
    { label: '无', value: '' },
    ...certs.value.map(c => ({ label: c.Domain, value: c.Domain }))
])

const panelUrl = computed(() => {
    if (!Config.value.Port) return ''
    const host = Config.value.Domain
        ? `https://${Config.value.Domain}`
        : `http://${Config.value.IPv4}`
    return `${host}:${Config.value.Port}${Config.value.Path}`
})

const logLevelOptions = [
    { label: 'INFO', value: 'info' },
    { label: 'WARN', value: 'warn' },
    { label: 'ERROR', value: 'error' },
]

// ── 生命周期 ───────────────────────────────────────────────

onMounted(async () => {
    Config.value = await getConfig()
    certs.value = await getCert()
    loaded.value = true
})

watch(Config, () => {
    if (loaded.value) dirty.value = true
}, { deep: true })

// ── 操作 ──────────────────────────────────────────────────

async function save() {
    try {
        await updateConfig(Config.value)
        await loadConfig()
        modal.value?.show('success', '保存成功')
    } catch (err: any) {
        modal.value?.show('error', err?.error)
        throw err
    }
}

async function restart() {
    try {
        modal.value?.show('warn', '准备重启...')
        const config = await getConfig()
        const host = config.Domain ? `https://${config.Domain}` : `http://${config.IPv4}`
        const url = `${host}:${config.Port}${config.Path}`

        let count = 3
        const jump = () => {
            localStorage.removeItem('token')
            window.location.href = url
        }
        modal.value?.show('success', `面板重启中，${count}秒后自动跳转...`)
        const timer = setInterval(() => {
            count--
            if (count <= 0) {
                clearInterval(timer)
                jump()
            } else {
                modal.value?.show('success', `面板重启中，${count}秒后自动跳转...`)
            }
        }, 1000)

        await restartPanel()
    } catch (err: any) {
        modal.value?.show('error', err?.error)
    }
}

async function submitCredentials() {
    try {
        await changeCredentials(credentials.value)
        await loadConfig()
        credentials.value = { old_password: '', new_username: '', new_password: '' }
        modal.value?.show('success', '修改成功')
    } catch (err: any) {
        modal.value?.show('error', err?.error)
    }
}
</script>

<style scoped>
form {
    gap: 10px;
}
button {
    margin-left: auto;
    font-size: var(--font-size-sm);
    padding: 10px 20px;
    line-height: 1;
}
.log-wrap {
    height: calc(100vh - 180px);
    display: flex;
    flex-direction: column;
}
</style>