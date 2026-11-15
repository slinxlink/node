<template>
    <Page class="inbound" title="入站" :motion="motion">
        <Table :columns="columns" :data="rows" placeholder="暂无入站协议" :actions="actions" />
    </Page>
</template>

<script setup lang="ts">
import type { Component } from 'vue'
import type { Motion } from '@/component/layout/Header.vue'
import Page from '@/component/layout/Page.vue'
import Table from '@/component/module/Table.vue'
import Traffic from '@/component/module/Traffic.vue'
import Toggle from '@/component/widget/Toggle.vue'
import Stamp from '@/component/widget/Stamp.vue'
import Form from '@/view/panel/inbound/Form.vue'
import { display } from '@/view/panel/inbound/display'
import { getInbounds, deleteInbound, resetInboundTraffic, toggleInbound } from '@/api/inbound'
import type { Inbound } from '@/api/inbound'
import { modal } from '@/composable/modal'
import { bb } from '@/util/bbcode'
import { formatTime } from '@/util/format'
import { toast } from '@/composable/toast'

type Row = Record<string, unknown>

const inbounds = ref<Inbound[]>([])

// 表格按 id 认行，把接口的字段整理成小写的一行
const rows = computed<Row[]>(() =>
    inbounds.value.map((item) => ({
        id: item.ID,
        name: item.Name || '-',
        port: item.Port,
        protocol: item.Protocol,
        transport: display.transport(item.Protocol, item.Transport),
        tls: display.tls(item.TLSType) || '-',
        used: item.Upload + item.Download,
        limit: item.TrafficLimit,
        enable: item.Enable,
        created: formatTime(item.CreatedAt),
    })),
)

async function load() {
    inbounds.value = await getInbounds()
}

async function toggle(id: number) {
    try {
        await toggleInbound(id)
        await load()
    } catch (err) {
        toast.failure((err as { error?: string })?.error || '操作失败')
    }
}

function remove(row: Row) {
    const item = inbounds.value.find((i) => i.ID === row.id)
    const label = item ? `${item.Protocol}:${item.Port}${item.Name ? ` - ${item.Name}` : ''}` : ''

    modal.show({
        title: '删除入站',
        detail: bb`确认删除入站 ${label} ?`,
        onConfirm: async () => {
            try {
                await deleteInbound(row.id as number)
                await load()
            } catch (err) {
                toast.failure((err as { error?: string })?.error || '删除失败')
                return false
            }
        },
    })
}

function reset(row: Row) {
    const item = inbounds.value.find((i) => i.ID === row.id)
    const label = item ? `${item.Protocol}:${item.Port}${item.Name ? ` - ${item.Name}` : ''}` : ''

    modal.show({
        title: '重置流量',
        detail: bb`确认重置入站 ${label} 的流量 ?`,
        onConfirm: async () => {
            try {
                await resetInboundTraffic(row.id as number)
                await load()
            } catch (err) {
                toast.failure((err as { error?: string })?.error || '重置失败')
                return false
            }
        },
    })
}

const columns = [
    {
        key: 'enable',
        label: '启用',
        collation: ['true', 'false'],
        render: (row: Row) => ({
            component: Toggle as Component,
            props: { value: row.enable, size: 'small', 'onUpdate:value': () => toggle(row.id as number) },
        }),
    },
    { key: 'name', label: '备注' },
    { key: 'port', label: '端口' },
    {
        key: 'protocol',
        label: '协议',
        collation: ['vless', 'vmess', 'trojan', 'hysteria', 'tuic', 'anytls'],
        render: (row: Row) => ({
            component: Stamp as Component,
            props: { color: display.protocolColor(row.protocol as string) },
            text: row.protocol as string,
        }),
    },
    { key: 'transport', label: '传输', collation: ['TCP', 'WS', 'UDP'] },
    // 没加密的入站显示短横线
    { key: 'tls', label: '安全', collation: ['TLS', 'Reality'] },
    {
        key: 'used',
        label: '流量',
        width: 260,
        collation: 'number' as const,
        render: (row: Row) => ({
            component: Traffic as Component,
            props: { used: row.used, limit: row.limit },
        }),
    },
    { key: 'created', label: '创建时间' },
]

function add() {
    modal.show({
        title: '添加入站',
        width: '800px',
        form: true,
        component: Form,
        onConfirm: async (form) => {
            if (!(await form.submit())) return false

            await load()
        },
    })
}

function edit(row: Row) {
    const target = inbounds.value.find((item) => item.ID === row.id)
    if (!target) return

    modal.show({
        title: '编辑入站',
        width: '800px',
        form: true,
        component: Form,
        props: { inbound: target },
        onConfirm: async (form) => {
            if (!(await form.submit())) return false

            await load()
        },
    })
}

const motion: Motion[] = [{ icon: 'add', tip: '添加入站', onClick: add }]

const actions = [
    { icon: 'restart_alt', tip: '重置流量', onClick: reset },
    { icon: 'edit', tip: '编辑', onClick: edit },
    { icon: 'delete', tip: '删除', danger: true, onClick: remove },
]

onMounted(load)
</script>
