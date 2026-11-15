<template>
    <Page class="core" title="核心配置">
        <Section :class="{ fill: tab !== 'settings' }">
            <template #action>
                <Tab v-model:value="tab" :items="tabs" />
            </template>

            <template v-if="tab === 'settings'">
                <div class="group">
                    <div class="gtitle">核心</div>
                    <div class="card">
                        <Field>
                            <template #prefix>
                                <hgroup>
                                    <strong>运行状态</strong>
                                    <span>{{ core?.Name || 'sing-box' }} {{ core?.Version || '-' }}</span>
                                </hgroup>
                            </template>
                            <Status :on="running">{{ running ? '运行中' : '未启动' }}</Status>
                            <Button class="action tight" variant="border" @click="running ? stop() : start()">
                                <Icon name="power_settings_new" />
                                {{ running ? '停止' : '启动' }}
                            </Button>
                            <Button class="action tight" variant="border" @click="restart">
                                <Icon name="refresh" />
                                重启
                            </Button>
                        </Field>
                    </div>
                </div>

                <Divider expand="0px" />

                <div class="group">
                    <div class="gtitle">
                        日志
                        <Stamp color="yellow">保存后需重启核心</Stamp>
                    </div>
                    <div class="card">
                        <Field>
                            <template #prefix>
                                <hgroup>
                                    <strong>启用</strong>
                                    <span>关闭后将无法记录用户在线 IP 等信息</span>
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
                            <Input v-model:value="form.LogPath" :prefix="`${dir}/`" placeholder="日志路径" autocomplete="off" />
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
                    <div class="gtitle">重置</div>
                    <div class="card">
                        <Field>
                            <template #prefix>
                                <hgroup>
                                    <strong>恢复默认</strong>
                                    <span>重置所有核心配置为默认值并重启核心</span>
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

            <JsonEditor v-else-if="tab === 'config'" :content="config" />
            <Log v-else :ws="coreLog" />
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
                        {{ needRestart ? '保存并重启核心' : '保存' }}
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
import JsonEditor from '@/component/module/JsonEditor.vue'
import Log from '@/component/module/Log.vue'
import Button from '@/component/widget/Button.vue'
import Icon from '@/component/widget/Icon.vue'
import Input from '@/component/widget/Input.vue'
import Select from '@/component/widget/Select.vue'
import Stamp from '@/component/widget/Stamp.vue'
import Status from '@/component/widget/Status.vue'
import Tab from '@/component/widget/Tab.vue'
import type { TabItem } from '@/component/widget/Tab.vue'
import Toggle from '@/component/widget/Toggle.vue'
import { getConfig } from '@/api/config'
import {
    coreLog,
    getCore,
    getCoreConfig,
    getCoreStatus,
    resetCore,
    restartCore,
    startCore,
    stopCore,
    updateCore,
} from '@/api/core'
import type { Core } from '@/api/core'
import { modal } from '@/composable/modal'
import { toast } from '@/composable/toast'

const tabs: [TabItem, ...TabItem[]] = [
    { label: '设置', value: 'settings' },
    { label: '高级配置', value: 'config' },
    { label: '日志', value: 'log' },
]

const levels = ['trace', 'debug', 'info', 'warn', 'error', 'fatal', 'panic'].map((value) => ({
    label: value.toUpperCase(),
    value,
}))

const tab = ref('settings')
const dir = ref('')
const config = ref('')
const running = ref(false)

const core = ref<Core>()
const form = ref({ LogEnable: true, LogPath: '', LogLevel: 'info' })

// 上次从服务器读到(或保存成功)时的样子；表单和它逐项比较，不一样才算改过：
// 改了又改回原样不算改，保存按钮也跟着变回不可点
const original = ref({ ...form.value })

type Key = keyof typeof form.value

// 分组的名字：保存条里列出的名字拼成「分组名 + 条目名」，如「日志等级」
const groups = {
    log: '日志',
}

// 每一项：属于哪个分组、条目名、保存后是否需要重启核心才生效
// 现在核心的设置保存时都会自动重启核心，所以都是 true；以后加了不用重启的项，写 false
const fields: [Key, keyof typeof groups, string, boolean][] = [
    ['LogEnable', 'log', '启用', true],
    ['LogPath', 'log', '路径', true],
    ['LogLevel', 'log', '等级', true],
]

// 和原样逐项比较，不一样才算改过；改了又改回原样不算改，保存条也跟着收起
const changes = computed(() =>
    fields
        .filter(([key]) => form.value[key] !== original.value[key])
        .map(([, group, item, restart]) => ({ label: `${groups[group]}${item}`, restart })),
)

const restartItems = computed(() => changes.value.filter((item) => item.restart))
const needRestart = computed(() => restartItems.value.length > 0)

const tip = computed(() =>
    needRestart.value
        ? `保存后核心将自动重启（${restartItems.value.map((item) => item.label).join('、')}）`
        : '保存后立即生效，不需要重启核心',
)

async function loadStatus() {
    running.value = (await getCoreStatus()).status === 'running'
}

async function load() {
    const [value, cfg] = await Promise.all([getCore(), getConfig()])

    dir.value = cfg.Dir
    core.value = value
    form.value = {
        LogEnable: value.LogEnable,
        LogPath: value.LogPath,
        LogLevel: value.LogLevel,
    }
    original.value = { ...form.value }
}

// 切到高级配置时才去取，每次切进来都取最新的
watch(tab, async (value) => {
    if (value === 'config') {
        config.value = await getCoreConfig()
    }
})

// 启动、停止、重启：失败提示接口给的原因，完成后刷新状态
async function run(action: () => Promise<unknown>, fail: string) {
    try {
        await action()
        await Promise.all([load(), loadStatus()])
    } catch (err) {
        toast.failure((err as { error?: string })?.error || fail)
    }
}

const start = () => run(startCore, '启动失败')
const stop = () => run(stopCore, '停止失败')
const restart = () => run(restartCore, '重启失败')

// 把表单写进服务器；成功返回 true
async function commit(): Promise<boolean> {
    if (!core.value) return false

    try {
        await updateCore({ ...core.value, ...form.value })
    } catch (err) {
        toast.failure((err as { error?: string })?.error || '保存失败')

        return false
    }

    toast.success('保存成功')
    original.value = { ...form.value }

    return true
}

// 操作条上的保存：改了需要重启核心的项，先确认(保存后核心自动重启，正在连着的客户端会短暂断开)；否则直接保存
async function submit() {
    if (!needRestart.value) {
        await commit()

        return
    }

    modal.show({
        title: '确认保存并重启核心？',
        detail: '保存后核心将自动重启',
        onConfirm: async () => ((await commit()) ? undefined : false),
    })
}

// 还原成原样
function undo() {
    form.value = { ...original.value }
}

function reset() {
    modal.show({
        title: '重置核心配置',
        detail: '确认将核心配置重置为默认值 ?',
        onConfirm: async () => {
            try {
                await resetCore()
            } catch (err) {
                toast.failure((err as { error?: string })?.error || '重置失败')

                return false
            }

            await Promise.all([load(), loadStatus()])
            toast.success('已重置')
        },
    })
}

let timer = 0

onMounted(() => {
    load()
    loadStatus()
    timer = window.setInterval(loadStatus, 5000)
})

onUnmounted(() => clearInterval(timer))
</script>

<style scoped>
.core {
    /* 左边标题(可带一行副标题)占一半，右边控件占一半 */
    :deep(.prefix) {
        flex: 0 0 calc(50% - 5px);
        color: var(--text-200); /* 标题；副标题和版本号用 text-300 */
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

    /* 保存条：和设置页一样，改了才出现，单独一张浮起来的卡片(照 mesh 的购买条)：
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

    /* 高级配置和日志是整页结构：容器高度由窗口决定，正好填满页头(80px)和页脚(80px)之间，不够 400px 时页面自己滚；里面的内容再撑满容器 */
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
