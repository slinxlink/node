<!-- Endpoint.vue -->
<template>
    <main>
        <div class="page">
            <div class="body">
                <button class="add-btn" @click="create">
                    <i class="icon">add</i>
                    添加端点
                </button>
                <button class="add-btn" @click="showWarp = true">
                    WARP
                </button>

                <List :headers="['标签', '协议', '地址', '启用', '操作']">
                    <tr v-for="ep in endpoints" :key="ep.ID">
                        <td>{{ ep.Tag }}</td>
                        <td>
                            <span class="tag" :class="ep.Type === 'wireguard' ? 'primary' : 'green'">
                                {{ ep.Type }}
                            </span>
                        </td>
                        <td>{{ ep.PeerAddress }}:{{ ep.PeerPort }}</td>
                        <td>
                            <Toggle :model-value="ep.Enable" @update:model-value="toggle(ep)" />
                        </td>
                        <td class="actions">
                            <button class="icon-btn" @click="openRoute(ep)" title="路由">
                                <i class="icon">alt_route</i>
                            </button>
                            <button class="icon-btn" @click="edit(ep)" title="编辑">
                                <i class="icon">edit</i>
                            </button>
                            <button class="icon-btn danger" @click="remove(ep.ID)" title="删除">
                                <i class="icon">delete</i>
                            </button>
                        </td>
                    </tr>
                    <tr v-if="endpoints.length === 0">
                        <td colspan="5" class="empty">暂无端点</td>
                    </tr>
                </List>
            </div>
        </div>
    </main>
    <Drawer v-model="showWarp" title="Cloudflare WARP" :footer="false">
        <Warp v-if="showWarp" @applied="load" />
    </Drawer>
    <Drawer v-model="showDrawer" :title="drawerTitle" @save="save">
        <Form v-model="defaultEndpoint" />
    </Drawer>



    <Drawer v-model="showRoute" title="路由设置" @save="routeRef?.save()">
        <Route v-if="showRoute" ref="routeRef" :endpoint="routeEndpoint" />
    </Drawer>
</template>

<script setup lang="ts">
import Drawer from '@/component/Drawer.vue'
import List from '@/component/ui/List.vue'
import Toggle from '@/component/ui/Toggle.vue'
import Warp from '@/view/panel/endpoint/widget/Warp.vue'
import Form from '@/view/panel/endpoint/form/Form.vue'
import { getEndpoints, saveEndpoint, deleteEndpoint, toggleEndpoint } from '@/api/endpoint'
import { generateWireguardKeyPair } from '@/api/generate'

const modal = inject<any>('modal')

// ── 状态 ──────────────────────────────────────────────────

const endpoints = ref<any[]>([])
const showWarp = ref(false)
const showDrawer = ref(false)
const generatedDefaults = ref<any>({})
const defaultEndpoint = ref<any>(baseEndpoint())
const drawerTitle = ref('添加端点')

// ── 工具函数 ───────────────────────────────────────────────

function baseEndpoint() {
    return {
        Enable: true,
        Tag: '',
        Type: 'wireguard',
        MTU: 1408,
        Address: ['172.16.0.2/32', 'fd00::2/128'],
        PrivateKey: '',
        PublicKey: '',
        PeerAddress: '',
        PeerPort: 2408,
        PeerPublicKey: '',
        AllowedIPs: '0.0.0.0/0,::/0',
        Reserved: '[0,0,0]',
        UDPTimeout: 0,
        Workers: 0,
    }
}

// ── 生命周期 ───────────────────────────────────────────────

onMounted(() => {
    load()

    generateWireguardKeyPair().then((res) => {
        generatedDefaults.value = { PrivateKey: res.private_key, PublicKey: res.public_key }
        defaultEndpoint.value = { ...baseEndpoint(), ...generatedDefaults.value }
    })
})

// ── 列表操作 ───────────────────────────────────────────────

async function load() {
    const res = await getEndpoints()
    endpoints.value = res.data
}

async function toggle(ep: any) {
    await toggleEndpoint(ep.ID)
    await load()
}

function remove(id: number) {
    modal.value?.show('confirm', '确认删除该端点？', async () => {
        await deleteEndpoint(id)
        await load()
    })
}

// ── Drawer 操作 ────────────────────────────────────────────

async function create() {
    const res = await generateWireguardKeyPair()
    generatedDefaults.value = { PrivateKey: res.private_key, PublicKey: res.public_key }
    defaultEndpoint.value = { ...baseEndpoint(), ...generatedDefaults.value }
    drawerTitle.value = '添加端点'
    showDrawer.value = true
}

function edit(ep: any) {
    defaultEndpoint.value = { ...defaultEndpoint.value, ...ep, Address: JSON.parse(ep.Address || '[]') }
    drawerTitle.value = '编辑端点'
    showDrawer.value = true
}

async function save() {
    try {
        const payload = { ...defaultEndpoint.value, Address: JSON.stringify(defaultEndpoint.value.Address) }
        await saveEndpoint(payload)
        await load()
        showDrawer.value = false
    } catch (err: any) {
        modal.value?.show('error', err?.error)
    }
}






import Route from '@/view/panel/endpoint/widget/Route.vue'

const showRoute = ref(false)
const routeEndpoint = ref<any>(null)
const routeRef = ref<any>(null)

function openRoute(ep: any) {
    routeEndpoint.value = ep
    showRoute.value = true
}
</script>