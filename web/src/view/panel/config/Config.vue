<template>
    <Page class="config" title="面板设置">
        <Section :class="{ fill: tab === 'log' }">
            <template #action>
                <Tab v-model:value="tab" :items="tabs" />
            </template>

            <template v-if="tab === 'settings'">
                <div class="group">
                    <div class="gtitle">
                        {{ groups.panel }}
                        <Stamp color="yellow">保存后需重启面板</Stamp>
                    </div>
                    <div class="card">
                    <Field>
                        <template #prefix>当前面板地址</template>
                        <span class="mono">{{ panelUrl || '-' }}</span>
                        <Copy v-if="panelUrl" :value="panelUrl" />
                    </Field>
                    <Field>
                        <template #prefix>
                            <hgroup>
                                <strong>域名</strong>
                                <span>选择证书域名，设置后自动启用 HTTPS</span>
                            </hgroup>
                        </template>
                        <Select v-model:value="form.Domain" :items="domains" placeholder="选择域名" />
                    </Field>
                    <Field>
                        <template #prefix>
                            <hgroup>
                                <strong>端口</strong>
                                <span>范围 1024 ~ 65535</span>
                            </hgroup>
                        </template>
                        <Input v-model:value="form.Port" type="number" :min="1" :max="65535" placeholder="端口" />
                    </Field>
                    <Field>
                        <template #prefix>
                            <hgroup>
                                <strong>路径</strong>
                                <span>必须以 / 开头</span>
                            </hgroup>
                        </template>
                        <Input v-model:value="form.Path" placeholder="路径" autocomplete="off" />
                    </Field>
                    </div>
                </div>

                <Divider expand="0px" />

                <div class="group">
                    <div class="gtitle">
                        {{ groups.sub }}
                        <Stamp color="yellow">保存后需重启面板</Stamp>
                    </div>
                    <div class="card">
                    <Field>
                        <template #prefix>
                            <hgroup>
                                <strong>启用</strong>
                                <span>关闭后订阅服务停止</span>
                            </hgroup>
                        </template>
                        <Toggle v-model:value="form.SubEnable" />
                    </Field>
                    <Field>
                        <template #prefix>
                            <hgroup>
                                <strong>端口</strong>
                                <span>范围 1024 ~ 65535</span>
                            </hgroup>
                        </template>
                        <Input v-model:value="form.SubPort" type="number" :min="1" :max="65535" placeholder="2096" />
                    </Field>
                    <Field>
                        <template #prefix>
                            <hgroup>
                                <strong>路径</strong>
                                <span>必须以 / 开头</span>
                            </hgroup>
                        </template>
                        <Input v-model:value="form.SubPath" placeholder="/link" autocomplete="off" />
                    </Field>
                    </div>
                </div>

                <Divider expand="0px" />

                <div class="group">
                    <div class="gtitle">
                        {{ groups.ruleset }}
                    </div>
                    <div class="card">
                    <Field>
                        <template #prefix>
                            <hgroup>
                                <strong>自动更新</strong>
                                <span>每周一自动更新规则集</span>
                            </hgroup>
                        </template>
                        <Toggle v-model:value="form.RulesetAutoUpdate" />
                    </Field>
                    <Field>
                        <template #prefix>
                            <hgroup>
                                <strong>手动更新</strong>
                                <span>立即下载并更新规则集</span>
                            </hgroup>
                        </template>
                        <Button class="action tight" variant="border" @click="refreshRuleset">
                            <Icon name="refresh" />
                            更新
                        </Button>
                    </Field>
                    </div>
                </div>

                <Divider expand="0px" />

                <div class="group">
                    <div class="gtitle">
                        {{ groups.log }}
                    </div>
                    <div class="card">
                    <Field>
                        <template #prefix>
                            <hgroup>
                                <strong>启用</strong>
                                <span>关闭后面板将不记录任何日志</span>
                            </hgroup>
                        </template>
                        <Toggle v-model:value="form.LogEnable" />
                    </Field>
                    <Field>
                        <template #prefix>
                            <hgroup>
                                <strong>路径</strong>
                                <span>日志文件仅支持保存在面板根目录下</span>
                            </hgroup>
                        </template>
                        <Input v-model:value="form.LogPath" :prefix="`${config?.Dir ?? ''}/`" placeholder="日志路径" autocomplete="off" />
                    </Field>
                    <Field>
                        <template #prefix>
                            <hgroup>
                                <strong>等级</strong>
                                <span>等级越高，记录的日志越少</span>
                            </hgroup>
                        </template>
                        <Select v-model:value="form.LogLevel" :items="levels" placeholder="日志等级" />
                    </Field>
                    </div>
                </div>

                <Divider expand="0px" />

                <div class="group">
                    <div class="gtitle">
                        {{ groups.system }}
                    </div>
                    <div class="card">
                    <Field>
                        <template #prefix>
                            <hgroup>
                                <strong>BBR 加速</strong>
                                <span>TCP 拥塞控制优化</span>
                            </hgroup>
                        </template>
                        <Toggle v-model:value="form.BBR" />
                    </Field>
                    </div>
                </div>

                <Divider expand="0px" />

                <div class="group">
                    <div class="gtitle">
                        {{ groups.board }}
                    </div>
                    <div class="card">
                    <Field>
                        <template #prefix>
                            <hgroup>
                                <strong>启用</strong>
                                <span>关闭后停止所有面板对接与同步</span>
                            </hgroup>
                        </template>
                        <Toggle v-model:value="form.BoardEnable" />
                    </Field>
                    </div>
                </div>

                <Divider expand="0px" />

                <div class="group">
                    <div class="gtitle">
                        {{ groups.reset }}
                    </div>
                    <div class="card">
                    <Field>
                        <template #prefix>
                            <hgroup>
                                <strong>恢复默认</strong>
                                <span>重置所有面板配置为默认值并重启核心</span>
                            </hgroup>
                        </template>
                        <Button class="action tight" variant="danger" @click="reset">
                            <Icon name="restart_alt" />
                            恢复默认
                        </Button>
                    </Field>
                    </div>
                </div>
            </template>

            <form v-else-if="tab === 'admin'" class="credentials" @submit.prevent>
                <Field>
                    <template #prefix>用户名</template>
                    <Input :value="config?.Username ?? ''" autocomplete="username" disabled />
                </Field>
                <Field>
                    <template #prefix>原密码</template>
                    <Input v-model:value="credentials.old_password" type="password" placeholder="原密码" autocomplete="current-password" />
                </Field>
                <Field>
                    <template #prefix>
                        <hgroup>
                            <strong>新用户名</strong>
                            <span>留空则不修改</span>
                        </hgroup>
                    </template>
                    <Input v-model:value="credentials.new_username" placeholder="新用户名" autocomplete="username" />
                </Field>
                <Field>
                    <template #prefix>
                        <hgroup>
                            <strong>新密码</strong>
                            <span>留空则不修改</span>
                        </hgroup>
                    </template>
                    <Input v-model:value="credentials.new_password" type="password" placeholder="新密码" autocomplete="new-password" />
                </Field>
                <Field>
                    <template #prefix></template>
                    <Button class="action" variant="border" @click="submitCredentials">提交修改</Button>
                </Field>
            </form>

            <Log v-else-if="tab === 'log'" :ws="panelLog" />
        </Section>

        <Transition name="bar">
            <div v-if="tab === 'settings' && changes.length" class="bar">
                <div class="in">
                    <span class="dot"></span>
                    <div class="msg">
                        <span>已修改 {{ changes.length }} 项</span>
                        <small>{{ tip }}</small>
                    </div>
                    <Button class="action" @click="undo">还原</Button>
                    <Button class="action tight" variant="primary" @click="submit">
                        <Icon name="save" />
                        {{ needRestart ? '保存并重启面板' : '保存' }}
                    </Button>
                </div>
            </div>
        </Transition>
    </Page>
</template>

<script setup lang="ts">
import Divider from '@/component/layout/Divider.vue'
import Field from '@/component/layout/Field.vue'
import Page from '@/component/layout/Page.vue'
import Section from '@/component/layout/Section.vue'
import Log from '@/component/module/Log.vue'
import Button from '@/component/widget/Button.vue'
import Copy from '@/component/widget/Copy.vue'
import Icon from '@/component/widget/Icon.vue'
import Input from '@/component/widget/Input.vue'
import Select from '@/component/widget/Select.vue'
import Stamp from '@/component/widget/Stamp.vue'
import Tab from '@/component/widget/Tab.vue'
import type { TabItem } from '@/component/widget/Tab.vue'
import Toggle from '@/component/widget/Toggle.vue'
import { changeCredentials } from '@/api/auth'
import { restartPanel } from '@/api/bootstrap'
import { getCert } from '@/api/cert'
import type { Cert } from '@/api/cert'
import { Log as panelLog, getConfig, resetConfig, updateConfig } from '@/api/config'
import type { Config } from '@/api/config'
import { refreshRuleset as fetchRuleset } from '@/api/ruleset'
import { modal } from '@/composable/modal'
import { toast } from '@/composable/toast'
import { loadConfig } from '@/store/config'

const tabs: [TabItem, ...TabItem[]] = [
    { label: '设置', value: 'settings' },
    { label: '管理员', value: 'admin' },
    { label: '日志', value: 'log' },
]

const levels = [
    { label: 'INFO', value: 'info' },
    { label: 'WARN', value: 'warn' },
    { label: 'ERROR', value: 'error' },
]

const tab = ref('settings')
const config = ref<Config>()
const certs = ref<Cert[]>([])

// 端口用字符串存着，保存时再转数字
const form = ref({
    Domain: '',
    Port: '',
    Path: '',
    SubEnable: false,
    SubPort: '',
    SubPath: '',
    RulesetAutoUpdate: false,
    LogEnable: false,
    LogPath: '',
    LogLevel: 'info',
    BBR: false,
    BoardEnable: false,
})

const credentials = ref({ old_password: '', new_username: '', new_password: '' })

// 上次从服务器读到(或保存成功)时的样子；表单和它逐项比较，不一样才算改过：
// 改了又改回原样不算改，底部的操作条也跟着收起
const original = ref({ ...form.value })

type Key = keyof typeof form.value

// 分组的名字：页面上的分组标题和保存条里列出的名字都用它
const groups = {
    panel: '面板',
    sub: '订阅',
    ruleset: '规则集',
    log: '日志',
    system: '系统',
    board: '面板对接',
    reset: '重置',
}

// 每一项：属于哪个分组、条目名、保存后是否需要重启面板才生效；
// 保存条里列出的名字自动拼成「分组名 + 条目名」，如「面板端口」「订阅路径」
// 域名、端口、路径、订阅这几项是启动时才读的；日志、BBR、规则集、面板对接保存时就已经生效
const fields: [Key, keyof typeof groups, string, boolean][] = [
    ['Domain', 'panel', '域名', true],
    ['Port', 'panel', '端口', true],
    ['Path', 'panel', '路径', true],
    ['SubEnable', 'sub', '启用', true],
    ['SubPort', 'sub', '端口', true],
    ['SubPath', 'sub', '路径', true],
    ['RulesetAutoUpdate', 'ruleset', '自动更新', false],
    ['LogEnable', 'log', '启用', false],
    ['LogPath', 'log', '路径', false],
    ['LogLevel', 'log', '等级', false],
    ['BBR', 'system', 'BBR 加速', false],
    ['BoardEnable', 'board', '启用', false],
]

const changes = computed(() =>
    fields
        .filter(([key]) => form.value[key] !== original.value[key])
        .map(([, group, item, restart]) => ({ label: `${groups[group]}${item}`, restart })),
)

const restartItems = computed(() => changes.value.filter((item) => item.restart))
const needRestart = computed(() => restartItems.value.length > 0)

const tip = computed(() =>
    needRestart.value
        ? `其中 ${restartItems.value.length} 项需要重启面板才能生效（${restartItems.value.map((item) => item.label).join('、')}）`
        : '保存后立即生效，不需要重启',
)

const domains = computed(() => [
    { label: '无', value: '' },
    ...certs.value.map((item) => ({ label: item.Domain, value: item.Domain })),
])

const panelUrl = computed(() => {
    if (!form.value.Port || !config.value) return ''

    const host = form.value.Domain ? `https://${form.value.Domain}` : `http://${config.value.IPv4}`

    return `${host}:${form.value.Port}${form.value.Path}`
})

function fill(value: Config) {
    config.value = value
    form.value = {
        Domain: value.Domain,
        Port: String(value.Port),
        Path: value.Path,
        SubEnable: value.SubEnable,
        SubPort: String(value.SubPort),
        SubPath: value.SubPath,
        RulesetAutoUpdate: value.RulesetAutoUpdate,
        LogEnable: value.LogEnable,
        LogPath: value.LogPath,
        LogLevel: value.LogLevel,
        BBR: value.BBR,
        BoardEnable: value.BoardEnable,
    }
    original.value = { ...form.value }
}

async function load() {
    const [value, list] = await Promise.all([getConfig(), getCert()])

    certs.value = list
    fill(value)
}

// 把表单写进服务器；成功返回 true
async function commit(): Promise<boolean> {
    if (!config.value) return false

    try {
        await updateConfig({
            ...config.value,
            ...form.value,
            Port: Number(form.value.Port),
            SubPort: Number(form.value.SubPort),
        })
        await loadConfig()
    } catch (err) {
        toast.failure((err as { error?: string })?.error || '保存失败')

        return false
    }

    toast.success('保存成功')
    original.value = { ...form.value }

    return true
}

// 操作条上的保存：只改了立即生效的项就直接保存；改了需要重启的项，确认之后保存并重启面板
async function submit() {
    if (!needRestart.value) {
        await commit()

        return
    }

    modal.show({
        title: '确认保存并重启面板？',
        detail: '保存后面板会重启，页面将自动跳转到新地址',
        onConfirm: async () => {
            if (!(await commit())) return false

            await restart()
        },
    })
}

// 还原成原样
function undo() {
    form.value = { ...original.value }
}

// 重启面板：端口、路径可能变了，3 秒后清掉登录状态，跳到新地址
async function restart() {
    try {
        const value = await getConfig()
        const host = value.Domain ? `https://${value.Domain}` : `http://${value.IPv4}`
        const url = `${host}:${value.Port}${value.Path}`

        toast.loading('准备重启...')
        await restartPanel()

        let count = 3
        toast.loading(`面板重启中，${count} 秒后自动跳转...`)

        const timer = setInterval(() => {
            count--

            if (count <= 0) {
                clearInterval(timer)
                localStorage.removeItem('token')
                window.location.href = url

                return
            }

            toast.loading(`面板重启中，${count} 秒后自动跳转...`)
        }, 1000)
    } catch (err) {
        toast.failure((err as { error?: string })?.error || '重启失败')
    }
}

async function refreshRuleset() {
    toast.loading('规则集更新中，请稍候...')

    try {
        await fetchRuleset()
    } catch (err) {
        toast.failure((err as { error?: string })?.error || '更新失败')

        return
    }

    toast.success('规则集更新完成')
}

async function submitCredentials() {
    try {
        await changeCredentials(credentials.value)
        await loadConfig()
    } catch (err) {
        toast.failure((err as { error?: string })?.error || '修改失败')

        return
    }

    credentials.value = { old_password: '', new_username: '', new_password: '' }
    fill(await getConfig())
    toast.success('修改成功')
}

function reset() {
    modal.show({
        title: '确认重置所有面板配置为默认值？',
        detail: '此操作将重置面板端口、路径、用户名及密码为随机值[br]重置后将无法通过当前地址访问面板[br]请前往服务器终端命令获取新的登录信息',
        onConfirm: async () => {
            try {
                await resetConfig()
            } catch (err) {
                toast.failure((err as { error?: string })?.error || '重置失败')

                return false
            }

            await load()
            toast.success('已重置')
        },
    })
}

onMounted(load)
</script>

<style scoped>
.config {
    /* 左边标题(可带一行副标题)占一半，右边控件占一半 */
    :deep(.prefix) {
        flex: 0 0 calc(50% - 5px);
        color: var(--text-200); /* 标题；副标题用 text-300 */
    }

    hgroup {
        display: flex;
        flex-direction: column;
        gap: 2px;

        strong {
            font-weight: normal;
        }

        span {
            color: var(--text-300);
            font-size: 12px;
        }
    }

    .mono {
        font-family: monospace;
        font-size: 12px;
        color: var(--text-200);
    }

    /* 分组：分割线(不带字)隔开，标题写在线下面，下面是这一组的行 */
    .group {
        display: flex;
        flex-direction: column;
        gap: 10px;
    }

    .gtitle {
        display: flex;
        align-items: center;
        gap: 10px;
        color: var(--text-100);
    }

    .card {
        display: flex;
        flex-direction: column;
        gap: 10px;
    }

    hr {
        margin-block: 10px;
    }

    /* 右半边里的按钮：靠左，按内容宽度，34px 高(图标 20px 比字高，光靠行高压不下来，直接定高) */
    .action {
        width: auto;
        height: 34px;
        padding: 0 15px;

        /* 带图标的小按钮，左 5px 右 10px */
        &.tight {
            padding: 0 10px 0 5px;
        }
    }

    .credentials {
        display: flex;
        flex-direction: column;
        gap: 10px;
    }

    /* 保存条：照 mesh 的购买条做，改了才出现，单独一张浮起来的卡片(不在容器里面)：
       外层只负责粘在页面底部，离窗口底边 40px；这 40px 是粘住的位置，不是占位的内边距，滚到最底后条自然落在页脚上面 */
    .bar {
        position: sticky;
        bottom: 40px;
        z-index: 5;
        flex-shrink: 0;

        .in {
            display: flex;
            align-items: center;
            gap: 10px;
            padding: 12px;
            border: 1px solid var(--background-100);
            border-radius: 20px;
            background-color: var(--background-300);
            box-shadow: 0 6px 30px rgba(0, 0, 0, 0.3);
        }

        .dot {
            flex-shrink: 0;
            width: 8px;
            height: 8px;
            margin-left: 8px;
            border-radius: 50%;
            background-color: var(--yellow);
        }

        .msg {
            display: flex;
            flex: 1;
            flex-direction: column;
            min-width: 0;
            color: var(--text-100);

            small {
                color: var(--text-200);
                font-size: 12px;
            }
        }
    }

    .bar-enter-active,
    .bar-leave-active {
        transition:
            opacity 0.25s,
            transform 0.25s;
    }

    .bar-enter-from,
    .bar-leave-to {
        opacity: 0;
        transform: translateY(100%);
    }

    /* 日志是整页结构：容器高度由窗口决定，正好填满页头(80px)和页脚(80px)之间，不够 400px 时页面自己滚；日志再撑满容器 */
    .fill {
        height: calc(100vh - 160px);
        min-height: 400px;

        &:deep(.content) {
            flex: 1;
            min-height: 0;
        }
    }
}
</style>
