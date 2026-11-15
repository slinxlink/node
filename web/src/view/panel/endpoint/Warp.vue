<!-- view/panel/endpoint/Warp.vue -->
<template>
    <div class="warp">
        <!-- 还没有账户 -->
        <Button v-if="!warp" variant="primary" :disabled="loading" @click="register">
            <Icon name="add_link" />
            {{ loading ? '创建中...' : '创建 WARP 账户' }}
        </Button>

        <template v-else>
            <Field v-for="item in secrets" :key="item.label">
                <template #prefix>{{ item.label }}</template>
                <span class="mono">{{ item.value }}</span>
                <Copy :value="item.value" />
            </Field>

            <div class="buttons">
                <Button variant="border" :disabled="loading" @click="refresh">
                    <Icon name="refresh" />
                    刷新
                </Button>
                <Button variant="border" :disabled="loading" @click="register">
                    <Icon name="swap_horiz" />
                    更换 IP
                </Button>
                <Button variant="danger" :disabled="loading" @click="remove">
                    <Icon name="delete_forever" />
                    {{ confirming ? '再点一次确认删除' : '删除账户' }}
                </Button>
            </div>

            <Divider>设置</Divider>
            <Field>
                <template #prefix>自动更新 IP</template>
                <Input v-model:value="autoUpdate" type="number" :min="0" placeholder="0" />
            </Field>
            <div class="buttons">
                <Button variant="border" :disabled="loading" @click="update">修改自动更新天数</Button>
                <span class="hint">0 表示关闭，单位为天</span>
            </div>
            <Field>
                <template #prefix>WARP+ 许可证</template>
                <Input v-model:value="license" placeholder="XXXXXXXX-XXXXXXXX-XXXXXXXX" autocomplete="off" />
            </Field>
            <div class="buttons">
                <Button variant="border" :disabled="loading" @click="setLicense">修改许可证</Button>
            </div>

            <template v-if="data">
                <Divider>账户信息</Divider>
                <Field v-for="item in infos" :key="item.label">
                    <template #prefix>{{ item.label }}</template>
                    {{ item.value }}
                </Field>
                <div class="buttons">
                    <Button variant="primary" :disabled="loading" @click="apply">应用至端点</Button>
                    <Stamp :color="applied ? 'green' : 'red'">{{ applied ? '已应用' : '未应用' }}</Stamp>
                </div>
            </template>
        </template>
    </div>
</template>

<script setup lang="ts">
import Divider from '@/component/layout/Divider.vue'
import Field from '@/component/layout/Field.vue'
import Button from '@/component/widget/Button.vue'
import Copy from '@/component/widget/Copy.vue'
import Icon from '@/component/widget/Icon.vue'
import Input from '@/component/widget/Input.vue'
import Stamp from '@/component/widget/Stamp.vue'
import {
    getWarp,
    deleteWarp,
    registerWarp,
    refreshWarp,
    setWarpAutoUpdate,
    setWarpLicense,
} from '@/api/warp'
import type { Warp, WarpData } from '@/api/warp'
import { getEndpoints, createWarpEndpoint } from '@/api/endpoint'
import type { Endpoint } from '@/api/endpoint'
import { toast } from '@/composable/toast'
import { formatBytes } from '@/util/format'

// WARP 账户的管理：创建、刷新、更换 IP、设置，再把账户应用到 tag 为 warp 的端点
const emit = defineEmits<{
    applied: []
}>()

const warp = ref<Warp | null>(null)
const data = ref<WarpData | null>(null)
const endpoints = ref<Endpoint[]>([])
const loading = ref(false)
const license = ref('')
const autoUpdate = ref('0')
const confirming = ref(false)

let timer = 0

const secrets = computed(() => [
    { label: 'Access Token', value: warp.value?.AccessToken ?? '' },
    { label: '设备 ID', value: warp.value?.DeviceID ?? '' },
    { label: '许可证密钥', value: warp.value?.LicenseKey ?? '' },
    { label: '公钥', value: warp.value?.PublicKey ?? '' },
    { label: '私钥', value: warp.value?.PrivateKey ?? '' },
])

const infos = computed(() => {
    const d = data.value

    if (!d) return []

    return [
        { label: '设备名称', value: d.DeviceName },
        { label: '设备型号', value: d.DeviceModel },
        { label: '设备已启用', value: d.DeviceEnabled ? '是' : '否' },
        { label: '账户类型', value: d.AccountType },
        { label: '角色', value: d.Role },
        { label: 'WARP+ 数据', value: formatBytes(d.WarpPlusData) },
        { label: '配额', value: formatBytes(d.Quota) },
        { label: '使用', value: formatBytes(d.Usage) },
    ]
})

// 账户信息和 tag 为 warp 的端点逐项相同才算已应用
const applied = computed(() => {
    const endpoint = endpoints.value.find((item) => item.Tag === 'warp')
    const d = data.value

    if (!endpoint || !d) return false

    return (
        endpoint.PrivateKey === d.PrivateKey &&
        endpoint.PublicKey === d.PublicKey &&
        endpoint.PeerAddress === d.PeerAddress &&
        endpoint.PeerPort === d.PeerPort &&
        endpoint.PeerPublicKey === d.PeerPublicKey &&
        endpoint.Reserved === d.Reserved &&
        endpoint.Address === d.Address
    )
})

async function load() {
    warp.value = await getWarp()
    autoUpdate.value = String(warp.value?.AutoUpdate ?? 0)
    endpoints.value = await getEndpoints()
}

// 每个操作都是：置忙、调接口、失败弹提示、恢复
async function run(action: () => Promise<void>, fail: string) {
    loading.value = true

    try {
        await action()
    } catch (err) {
        toast.failure((err as { error?: string })?.error || fail)
    } finally {
        loading.value = false
    }
}

const register = () =>
    run(async () => {
        data.value = await registerWarp()
        await load()
    }, '创建失败')

const refresh = () =>
    run(async () => {
        data.value = await refreshWarp()
    }, '刷新失败')

// 弹窗同一时间只能开一个，确认不能再弹窗：第一次点变成确认，3 秒内再点才删
function remove() {
    if (!confirming.value) {
        confirming.value = true
        timer = window.setTimeout(() => (confirming.value = false), 3000)

        return
    }

    clearTimeout(timer)
    confirming.value = false

    run(async () => {
        await deleteWarp()
        warp.value = null
        data.value = null
        toast.success('删除成功')
    }, '删除失败')
}

const update = () =>
    run(async () => {
        warp.value = await setWarpAutoUpdate(Number(autoUpdate.value))
        toast.success('修改成功')
    }, '修改失败')

const setLicense = () =>
    run(async () => {
        data.value = await setWarpLicense(license.value)
        await load()
        toast.success('修改成功')
    }, '修改失败')

const apply = () =>
    run(async () => {
        if (!data.value) return

        await createWarpEndpoint(data.value)
        await load()
        emit('applied')
        toast.success('应用成功')
    }, '应用失败')

onMounted(load)
onUnmounted(() => clearTimeout(timer))
</script>

<style scoped>
.warp {
    display: flex;
    flex-direction: column;
    gap: 10px;
    width: 100%;
    text-align: left;

    :deep(.prefix) {
        flex-basis: 100px;
    }

    .mono {
        font-family: monospace;
        font-size: 12px;
        color: var(--text-200);
    }

    .hint {
        color: var(--text-300);
        font-size: 12px;
    }

    /* 按钮按内容宽度，从内容列(标题列 100px + 间距 10px)开始放 */
    .buttons {
        display: flex;
        align-items: center;
        gap: 10px;
        margin-left: 110px;

        button {
            width: auto;
            padding: 5px 10px;
        }
    }
}
</style>
