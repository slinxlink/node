<!-- widget/Warp.vue -->
<template>
    <div class="warp">
        <!-- 状态一：空状态 -->
        <div v-if="!warp" class="empty-state">
            <button :disabled="loading" @click="register">
                <i class="icon">add_link</i>
                {{ loading ? '创建中...' : '创建 WARP 账户' }}
            </button>
        </div>

        <!-- 状态二/三：账户信息已存在 -->
        <template v-else>
            <FormRow title="Access Token">
                <span>{{ warp.AccessToken }}</span>
            </FormRow>
            <FormRow title="设备 ID">
                <span>{{ warp.DeviceID }}</span>
            </FormRow>
            <FormRow title="许可证密钥">
                <span>{{ warp.LicenseKey }}</span>
            </FormRow>
            <FormRow title="公钥">
                <span>{{ warp.PublicKey }}</span>
            </FormRow>
            <FormRow title="私钥">
                <span>{{ warp.PrivateKey }}</span>
            </FormRow>
            <FormRow>
                <button :disabled="loading" @click="refresh">
                    <i class="icon">refresh</i>
                    {{ loading ? '处理中...' : '刷新' }}
                </button>
                <button :disabled="loading" @click="register">
                    <i class="icon">refresh</i>
                    {{ loading ? '处理中...' : '更换 IP' }}
                </button>
                <button class="danger" :disabled="loading" @click="remove">
                    <i class="icon">delete_forever</i>
                    删除账户
                </button>
            </FormRow>

            <div class="divider">设置</div>
            <FormRow class="quarter" title="自动更新 IP">
                <Input v-model.number="warp.AutoUpdate" type="number" placeholder="0" />
                <span class="dark">0 表示关闭，单位为天</span>
            </FormRow>
            <FormRow>
                <button :disabled="loading" @click="update">修改自动更新天数</button>
            </FormRow>
            <FormRow title="WARP+许可证">
                <Input v-model="license" placeholder="XXXXXXXX-XXXXXXXX-XXXXXXXX" />
            </FormRow>
            <FormRow>
                <button :disabled="loading" @click="set">修改许可证</button>
            </FormRow>

            <template v-if="data">
                <div class="divider">账户信息</div>
                <FormRow title="设备名称">
                    <span>{{ data.DeviceName }}</span>
                </FormRow>
                <FormRow title="设备型号">
                    <span>{{ data.DeviceModel }}</span>
                </FormRow>
                <FormRow title="设备已启用">
                    <span>{{ data.DeviceEnabled }}</span>
                </FormRow>
                <FormRow title="账户类型">
                    <span>{{ data.AccountType }}</span>
                </FormRow>
                <FormRow title="角色">
                    <span>{{ data.Role }}</span>
                </FormRow>
                <FormRow title="WARP+ 数据">
                    <span>{{ data.WarpPlusData }}</span>
                </FormRow>
                <FormRow title="配额">
                    <span>{{ data.Quota }}</span>
                </FormRow>
                <FormRow title="使用">
                    <span>{{ data.Usage }}</span>
                </FormRow>
                <FormRow>
                    <button class="btn" :disabled="loading" @click="create">应用至端点</button>
                    <span class="tag" :class="check.cls">{{ check.label }}</span>
                </FormRow>
            </template>
        </template>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted, inject } from 'vue'
import FormRow from '@/component/ui/FormRow.vue'
import Input from '@/component/ui/Input.vue'
import { getWarp, deleteWarp, registerWarp, refreshWarp, setWarpAutoUpdate, setWarpLicense } from '@/api/warp'
import { getEndpoints, createWarpEndpoint } from '@/api/endpoint'

const modal = inject<any>('modal')

const emit = defineEmits(['applied'])

const warp = ref<any>(null)
const data = ref<any>(null)
const endpoints = ref<any[]>([])
const loading = ref(false)
const license = ref('')

const load = async () => {
    const res = await getWarp()
    warp.value = res.data

    const endpointRes = await getEndpoints()
    endpoints.value = endpointRes.data
}

const register = async () => {
    loading.value = true
    try {
        const res = await registerWarp()
        data.value = res.data
        await load()
    } catch (err: any) {
        modal.value?.show('error', err?.error)
    } finally {
        loading.value = false
    }
}

const refresh = async () => {
    loading.value = true
    try {
        const res = await refreshWarp()
        data.value = res.data
    } catch (err: any) {
        modal.value?.show('error', err?.error)
    } finally {
        loading.value = false
    }
}

const remove = () => {
    modal.value?.show('confirm', '确定要删除当前 WARP 账户吗？', async () => {
        loading.value = true
        try {
            await deleteWarp()
            warp.value = null
            data.value = null
            modal.value?.show('success', '删除成功')
        } catch (err: any) {
            modal.value?.show('error', err?.error)
        } finally {
            loading.value = false
        }
    })
}

const update = async () => {
    loading.value = true
    try {
        const res = await setWarpAutoUpdate(warp.value.AutoUpdate)
        warp.value = res.data
        modal.value?.show('success', '修改成功')
    } catch (err: any) {
        modal.value?.show('error', err?.error)
    } finally {
        loading.value = false
    }
}

const set = async () => {
    loading.value = true
    try {
        const res = await setWarpLicense(license.value)
        data.value = res.data
        await load()
        modal.value?.show('success', '修改成功')
    } catch (err: any) {
        modal.value?.show('error', err?.error)
    } finally {
        loading.value = false
    }
}

const create = async () => {
    loading.value = true
    try {
        await createWarpEndpoint(data.value)
        await load()
        emit('applied')
        modal.value?.show('success', '应用成功')
    } catch (err: any) {
        modal.value?.show('error', err?.error)
    } finally {
        loading.value = false
    }
}

const check = computed(() => {
    const endpoint = endpoints.value.find((e: any) => e.Tag === 'warp')
    if (!endpoint || !data.value) return { label: '未应用', cls: 'red' }

    const matched =
        endpoint.PrivateKey === data.value.PrivateKey &&
        endpoint.PublicKey === data.value.PublicKey &&
        endpoint.PeerAddress === data.value.PeerAddress &&
        endpoint.PeerPort === data.value.PeerPort &&
        endpoint.PeerPublicKey === data.value.PeerPublicKey &&
        endpoint.Reserved === data.value.Reserved &&
        endpoint.Address === data.value.Address

    return matched ? { label: '已应用', cls: 'green' } : { label: '未应用', cls: 'red' }
})

onMounted(() => {
    load()
})

watch(loading, (val) => {
    document.body.style.cursor = val ? 'wait' : ''
})
</script>

<style scoped>
.warp {
    display: flex;
    flex-direction: column;
    gap: 10px;
    padding: 20px;

    .action-row {
        display: flex;
        flex-direction: row;
        gap: 10px;
        margin-left: 110px;
    }

    .dark {
        color: var(--color-text-dark);
    }

    button {
        display: flex;
        flex-direction: row;
        gap: 5px;
        padding: 9px 20px;
        font-size: var(--font-size-sm);
        align-items: center;

        .icon {
            font-size: var(--font-size-sm);
        }

        &.danger {
            border: 1px solid var(--color-red);
            background-color: var(--color-bg-dark);
            color: var(--color-red);

            &:hover {
                border-color: var(--color-primary-light);
                color: var(--color-primary-light);
            }
        }
    }
}

.empty-state {
    display: flex;
    padding: 20px;
}

.full-width {
    width: 100%;
    margin-top: 12px;
}

.quarter :deep(.input-wrap),
.quarter :deep(.select) {
    max-width: 25%;
}
</style>