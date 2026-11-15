<template>
    <Page class="user" title="用户" :motion="motion">
        <Table :columns="columns" :data="rows" placeholder="暂无用户" :actions="actions" />
    </Page>
</template>

<script setup lang="ts">
import type { Component } from 'vue'
import type { Motion } from '@/component/layout/Header.vue'
import Page from '@/component/layout/Page.vue'
import Table from '@/component/module/Table.vue'
import Traffic from '@/component/module/Traffic.vue'
import Toggle from '@/component/widget/Toggle.vue'
import QRCode from '@/component/module/QRCode.vue'
import Form from '@/view/panel/user/Form.vue'
import Info from '@/view/panel/user/Info.vue'
import Alive from '@/view/panel/user/Alive.vue'
import Expire from '@/view/panel/user/Expire.vue'
import Inbounds from '@/view/panel/user/Inbounds.vue'
import type { Related } from '@/view/panel/user/Inbounds.vue'
import { getConfig } from '@/api/config'
import type { Config } from '@/api/config'
import { getInbounds } from '@/api/inbound'
import type { Inbound } from '@/api/inbound'
import { getUsers, deleteUser, resetUserTraffic, toggleUser } from '@/api/user'
import type { User } from '@/api/user'
import { modal } from '@/composable/modal'
import { bb } from '@/util/bbcode'
import { expireLevel, formatExpire } from '@/util/format'
import { toast } from '@/composable/toast'
import { display } from '@/view/panel/inbound/display'

type Row = Record<string, unknown>

const users = ref<User[]>([])
const inbounds = ref<Inbound[]>([])
const config = ref<Config>()

// 用户关联的入站 ID 换成协议、端口、名称；找不到的(入站已被删)略过
function related(json: string): Related[] {
    let ids: number[] = []

    try {
        ids = JSON.parse(json || '[]')
    } catch {
        return []
    }

    return ids.flatMap((id) => {
        const item = inbounds.value.find((i) => i.ID === id)

        if (!item) return []

        return [
            {
                id,
                name: item.Name,
                port: item.Port,
                protocol: item.Protocol,
                color: display.protocolColor(item.Protocol),
            },
        ]
    })
}

function aliveIPs(json: string): string[] {
    try {
        return JSON.parse(json || '[]')
    } catch {
        return []
    }
}

// 表格按 id 认行，把接口的字段整理成小写的一行
const rows = computed<Row[]>(() =>
    users.value.map((item) => ({
        id: item.ID,
        name: item.Name || '-',
        token: item.Token,
        inbounds: related(item.Inbounds),
        alive: aliveIPs(item.AliveIP),
        used: item.Upload + item.Download,
        limit: item.TrafficLimit,
        expire: formatExpire(item.ExpireAt),
        expireLevel: expireLevel(item.ExpireAt),
        enable: item.Enable,
    })),
)

async function load() {
    const [u, i, c] = await Promise.all([getUsers(), getInbounds(), getConfig()])
    users.value = u
    inbounds.value = i
    config.value = c
}

// 订阅地址：有域名用 https，没有就用 IP
function subUrl(token: string) {
    const c = config.value
    if (!c) return ''

    const host = c.Domain ? `https://${c.Domain}` : `http://${c.IPv4}`

    return `${host}:${c.SubPort}${c.SubPath}/${token}`
}

function qrcode(row: Row) {
    // 没起名字的用户用 Token 当标题和文件名
    const name = row.name === '-' ? (row.token as string) : (row.name as string)

    modal.show({
        title: '',
        loose: true,
        component: QRCode,
        props: { name, value: subUrl(row.token as string) },
    })
}

async function toggle(id: number) {
    try {
        await toggleUser(id)
        await load()
    } catch (err) {
        toast.failure((err as { error?: string })?.error || '操作失败')
    }
}

function remove(row: Row) {
    // 没起名字的用户用 Token 代替
    const name = row.name === '-' ? (row.token as string) : (row.name as string)

    modal.show({
        title: '删除用户',
        detail: bb`确认删除用户 ${name} ?`,
        onConfirm: async () => {
            try {
                await deleteUser(row.id as number)
                await load()
            } catch (err) {
                toast.failure((err as { error?: string })?.error || '删除失败')
                return false
            }
        },
    })
}

function reset(row: Row) {
    const name = row.name === '-' ? (row.token as string) : (row.name as string)

    modal.show({
        title: '重置流量',
        detail: bb`确认重置用户 ${name} 的流量 ?`,
        onConfirm: async () => {
            try {
                await resetUserTraffic(row.id as number)
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
    { key: 'name', label: '名称' },
    {
        // 一行一个入站，行会跟着变高；列宽按最长的一行算，超过上限就截断
        key: 'inbounds',
        label: '关联入站',
        sort: false,
        multiline: true,
        render: (row: Row) => ({
            component: Inbounds as Component,
            props: { items: row.inbounds },
        }),
    },
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
    {
        key: 'alive',
        label: '在线IP',
        sort: false,
        multiline: true,
        render: (row: Row) => ({
            component: Alive as Component,
            props: { items: row.alive },
        }),
    },
    {
        key: 'expire',
        label: '有效期',
        render: (row: Row) => ({
            component: Expire as Component,
            props: { text: row.expire, level: row.expireLevel },
        }),
    },
]

function add() {
    modal.show({
        title: '添加用户',
        width: '600px',
        form: true,
        component: Form,
        props: { inbounds: inbounds.value },
        onConfirm: async (form) => {
            if (!(await form.submit())) return false

            await load()
        },
    })
}

function edit(row: Row) {
    const target = users.value.find((item) => item.ID === row.id)
    if (!target) return

    modal.show({
        title: '编辑用户',
        width: '600px',
        form: true,
        component: Form,
        props: { inbounds: inbounds.value, user: target },
        onConfirm: async (form) => {
            if (!(await form.submit())) return false

            await load()
        },
    })
}

function info(row: Row) {
    const target = users.value.find((item) => item.ID === row.id)
    if (!target) return

    modal.show({
        title: '',
        width: '600px',
        form: true,
        component: Info,
        props: { user: target, inbounds: inbounds.value },
    })
}

const motion: Motion[] = [{ icon: 'add', tip: '添加用户', onClick: add }]

const actions = [
    { icon: 'qr_code', tip: '二维码', onClick: qrcode },
    { icon: 'info', tip: '信息', onClick: info },
    { icon: 'restart_alt', tip: '重置流量', onClick: reset },
    { icon: 'edit', tip: '编辑', onClick: edit },
    { icon: 'delete', tip: '删除', danger: true, onClick: remove },
]

onMounted(load)
</script>
