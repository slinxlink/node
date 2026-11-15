<template>
    <Page class="endpoint" title="端点" :motion="motion">
        <Table :columns="columns" :data="rows" placeholder="暂无端点" :actions="actions" />
    </Page>
</template>

<script setup lang="ts">
import type { Component } from 'vue'
import type { Motion } from '@/component/layout/Header.vue'
import Page from '@/component/layout/Page.vue'
import Table from '@/component/module/Table.vue'
import Toggle from '@/component/widget/Toggle.vue'
import Stamp from '@/component/widget/Stamp.vue'
import Form from '@/view/panel/endpoint/Form.vue'
import Route from '@/view/panel/endpoint/Route.vue'
import Warp from '@/view/panel/endpoint/Warp.vue'
import { getEndpoints, deleteEndpoint, toggleEndpoint } from '@/api/endpoint'
import type { Endpoint } from '@/api/endpoint'
import { modal } from '@/composable/modal'
import { bb } from '@/util/bbcode'
import { toast } from '@/composable/toast'
import { formatTime } from '@/util/format'

type Row = Record<string, unknown>

const endpoints = ref<Endpoint[]>([])

// 表格按 id 认行，把接口的字段整理成小写的一行
const rows = computed<Row[]>(() =>
    endpoints.value.map((item) => ({
        id: item.ID,
        tag: item.Tag,
        type: item.Type,
        address: `${item.PeerAddress}:${item.PeerPort}`,
        enable: item.Enable,
        created: formatTime(item.CreatedAt),
    })),
)

async function load() {
    endpoints.value = await getEndpoints()
}

async function toggle(id: number) {
    try {
        await toggleEndpoint(id)
        await load()
    } catch (err) {
        toast.failure((err as { error?: string })?.error || '操作失败')
    }
}

function remove(row: Row) {
    modal.show({
        title: '删除端点',
        detail: bb`确认删除端点 ${row.tag as string} ?`,
        onConfirm: async () => {
            try {
                await deleteEndpoint(row.id as number)
                await load()
            } catch (err) {
                toast.failure((err as { error?: string })?.error || '删除失败')
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
    { key: 'tag', label: '标签' },
    {
        key: 'type',
        label: '协议',
        render: (row: Row) => ({
            component: Stamp as Component,
            props: { color: row.type === 'wireguard' ? 'primary-300' : 'green' },
            text: row.type as string,
        }),
    },
    { key: 'address', label: '地址' },
    { key: 'created', label: '创建时间' },
]

function add() {
    modal.show({
        title: '添加端点',
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
    const target = endpoints.value.find((item) => item.ID === row.id)
    if (!target) return

    modal.show({
        title: '编辑端点',
        width: '800px',
        form: true,
        component: Form,
        props: { endpoint: target },
        onConfirm: async (form) => {
            if (!(await form.submit())) return false

            await load()
        },
    })
}

// 没有确认按钮，账户的操作都在里面；应用到端点后列表要刷新
function warp() {
    modal.show({
        title: 'Cloudflare WARP',
        width: '600px',
        form: true,
        component: Warp,
        props: { onApplied: load },
    })
}

const motion: Motion[] = [
    { icon: 'cloud', tip: 'WARP', onClick: warp },
    { icon: 'add', tip: '添加端点', onClick: add },
]

function route(row: Row) {
    const target = endpoints.value.find((item) => item.ID === row.id)
    if (!target) return

    modal.show({
        title: '路由设置',
        width: '600px',
        form: true,
        component: Route,
        props: { endpoint: target },
        onConfirm: async (form) => {
            if (!(await form.submit())) return false
        },
    })
}

const actions = [
    { icon: 'alt_route', tip: '路由', onClick: route },
    { icon: 'edit', tip: '编辑', onClick: edit },
    { icon: 'delete', tip: '删除', danger: true, onClick: remove },
]

onMounted(load)
</script>
