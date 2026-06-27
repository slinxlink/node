<template>
    <main>
        <div class="page">
            <div class="body">
                <button class="add-btn" @click="openCreate">
                    <i class="icon">add</i>
                    添加入站
                </button>

                <List :headers="['ID', '备注', '端口', '协议', '启用', '操作']">
                    <tr v-for="ib in inbounds" :key="ib.ID">
                        <td class="muted">{{ ib.ID }}</td>
                        <td>{{ ib.Name || '-' }}</td>
                        <td>{{ ib.Port }}</td>
                        <td>
                            <div class="tags">
                                <span class="tag primary">{{ ib.Protocol }}</span>
                                <span class="tag green">
                                    {{ ib.Protocol === 'hysteria' || ib.Protocol === 'tuic' ? 'UDP' : (ib.Transport === 'websocket' ? 'WS' : 'TCP') }}
                                </span>
                                <span v-if="ib.TLSType !== 'none'" class="tag blue">{{ formatTLS(ib.TLSType) }}</span>
                            </div>
                        </td>
                        <td>
                            <Toggle :model-value="ib.Enable" @update:model-value="toggle(ib)" />
                        </td>
                        <td class="actions">
                            <button class="icon-btn" @click="openEdit(ib)" title="编辑">
                                <i class="icon">edit</i>
                            </button>
                            <button class="icon-btn danger" @click="remove(ib.ID)" title="删除">
                                <i class="icon">delete</i>
                            </button>
                        </td>
                    </tr>
                    <tr v-if="inbounds.length === 0">
                        <td colspan="6" class="empty">暂无入站协议</td>
                    </tr>
                </List>
            </div>
        </div>
    </main>
    <Drawer v-model="showDrawer" :title="drawerTitle" @save="handleSave">
        <Form v-model="defaultInbound" />
    </Drawer>
</template>

<script setup lang="ts">
import Toggle from '@/component/ui/Toggle.vue'
import Drawer from '@/component/Drawer.vue'
import List from '@/component/ui/List.vue'
import Form from '@/view/panel/inbound/form/Form.vue'
import { getInbounds, saveInbound, deleteInbound, toggleInbound } from '@/api/inbound'
import { generatePort, generateRealityTarget, generateRealityKeyPair, generateShortIDs } from '@/api/generate'
import { getCert } from '@/api/cert'

const modal = inject<any>('modal')

// ── 状态 ──────────────────────────────────────────────────

const inbounds = ref<any[]>([])
const showDrawer = ref(false)

const baseInbound = () => ({
    Enable: true,
    Name: '',
    Protocol: 'vless',
    Port: 0,

    Transport: 'raw',
    WsPath: '/',
    WsHost: '',
    WsPingInterval: 0,

    UDPTimeout: 60,
    MasqueradeEnabled: false,
    MasqueradeType: 'default',
    MasqueradeURL: '',
    RewriteHost: false,
    IgnoreTLSVerify: false,
    MasqueradePath: '',
    MasqueradeCode: 200,
    MasqueradeBody: '',
    ObfsType: '',
    ObfsPassword: '',
    ObfsMinPacketSize: 512,
    ObfsMaxPacketSize: 1200,
    
    TuicCongestionControl: 'bbr',
    TuicAuthTimeout: 3,
    TuicZeroRTT: false,
    TuicHeartbeat: 10,
    TuicUDPRelayMode: 'native',

    HopEnabled: false,
    HopPort: '10000-11000',
    HopInterval: '5-10',

    AnyTLSIdleSessionCheckInterval: 30,
    AnyTLSIdleSessionTimeout: 30,
    AnyTLSMinIdleSession: 1,

    TLSType: 'none',
    ServerName: '',
    CipherSuites: '',
    ALPN: ['h2', 'http/1.1'],
    TLSMinVersion: '1.2',
    TLSMaxVersion: '1.3',
    UTLS: 'chrome',
    Insecure: false,
    Certs: [1],
    ECHEnabled: false,
    ECHKey: '',
    ECHConfig: '',

    RealityServerName: '',
    RealityServer: '',
    RealityServerPort: 443,
    RealityMaxTimeDiff: 0,
    RealityShortIDs: [],
    RealityPublicKey: '',
    RealityPrivateKey: '',
    Flow: 'xtls-rprx-vision',
})

const generatedDefaults = ref<any>({})
const defaultInbound = ref<any>(baseInbound())

// ── 生命周期 ───────────────────────────────────────────────

onMounted(() => {
    load()
})

// ── 工具函数 ───────────────────────────────────────────────

function formatTLS(val: string) {
    const map: Record<string, string> = {
        'tls': 'TLS',
        'reality': 'Reality',
    }
    return map[val] ?? val
}

// ── 列表操作 ───────────────────────────────────────────────

async function load() {
    inbounds.value = await getInbounds()
}

async function toggle(ib: any) {
    await toggleInbound(ib.ID)
    await load()
}

async function remove(id: number) {
    modal.value?.show('confirm', '确认删除该入站？', async () => {
        await deleteInbound(id)
        await load()
    })
}
// ── Drawer 操作 ────────────────────────────────────────────
const drawerTitle = ref('添加入站')

async function openCreate() {
    drawerTitle.value = '添加入站'
    defaultInbound.value = baseInbound()
    showDrawer.value = true

    const [portRes, targetRes, keyPairRes, shortIDsRes, certsRes] = await Promise.all([
        generatePort(),
        generateRealityTarget(),
        generateRealityKeyPair(),
        generateShortIDs(),
        getCert(),
    ])
    const firstCertID = certsRes.length > 0 ? certsRes[0].ID : 0
    defaultInbound.value = {
        ...defaultInbound.value,
        Port: portRes.port,
        RealityServerName: targetRes.domain,
        RealityServer: targetRes.domain,
        RealityServerPort: 443,
        RealityPrivateKey: keyPairRes.private_key,
        RealityPublicKey: keyPairRes.public_key,
        RealityShortIDs: shortIDsRes.short_ids,
        Certs: [firstCertID],
    }
}

function openEdit(ib: any) {
    defaultInbound.value = {
        ...defaultInbound.value,
        ...ib,
        ALPN: ib.ALPN ? ib.ALPN.split(',') : [],
        RealityShortIDs: ib.RealityShortIDs ? JSON.parse(ib.RealityShortIDs) : [],
        Certs: ib.Certs ? JSON.parse(ib.Certs) : [0],
    }
    drawerTitle.value = '编辑入站'
    showDrawer.value = true
}

async function handleSave() {
    const data = { ...defaultInbound.value }

    if (Array.isArray(data.ALPN)) {
        data.ALPN = data.ALPN.join(',')
    }

    if (Array.isArray(data.RealityShortIDs)) {
        data.RealityShortIDs = JSON.stringify(data.RealityShortIDs)
    }

    if (Array.isArray(data.Certs)) {
        data.Certs = JSON.stringify(data.Certs)
    }

    try {
        await saveInbound(data)
        await load()
        showDrawer.value = false
        modal.value?.show('success', '保存成功')
    } catch (err: any) {
        const msg = err?.error
        modal.value?.show('error', msg)
    }
}
</script>