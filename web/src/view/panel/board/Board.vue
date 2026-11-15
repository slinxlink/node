<template>
    <Page class="board" title="面板对接" :motion="motion">
        <Table :columns="columns" :data="rows" placeholder="暂无对接" :actions="actions" />
    </Page>
</template>

<script setup lang="ts">
import type { Component } from 'vue'
import type { Motion } from '@/component/layout/Header.vue'
import Page from '@/component/layout/Page.vue'
import Table from '@/component/module/Table.vue'
import Toggle from '@/component/widget/Toggle.vue'
import Stamp from '@/component/widget/Stamp.vue'
import Connect from '@/view/panel/board/Connect.vue'
import Form from '@/view/panel/board/Form.vue'
import Info from '@/view/panel/board/Info.vue'
import { display } from '@/view/panel/inbound/display'
import { getBoards, deleteBoard, toggleBoard } from '@/api/board'
import type { Board } from '@/api/board'
import { getInbounds } from '@/api/inbound'
import type { Inbound } from '@/api/inbound'
import { modal } from '@/composable/modal'
import { bb } from '@/util/bbcode'
import { toast } from '@/composable/toast'
import { formatTime } from '@/util/format'

type Row = Record<string, unknown>

const boards = ref<Board[]>([])
const inbounds = ref<Inbound[]>([])

// 关联的入站写成「协议:端口」；入站被删了就显示短横线
function inboundText(id: number) {
    const item = inbounds.value.find((i) => i.ID === id)

    return item ? `${item.Protocol}:${item.Port}` : '-'
}

// 表格按 id 认行，把接口的字段整理成小写的一行
const rows = computed<Row[]>(() =>
    boards.value.map((item) => ({
        id: item.ID,
        name: item.Name || '-',
        host: item.Host,
        inbound: inboundText(item.Inbound),
        protocol: inbounds.value.find((i) => i.ID === item.Inbound)?.Protocol ?? '',
        enable: item.Enable,
        created: formatTime(item.CreatedAt),
    })),
)

async function load() {
    const [list, inboundList] = await Promise.all([getBoards(), getInbounds()])

    boards.value = list
    inbounds.value = inboundList
}

async function toggle(id: number) {
    try {
        await toggleBoard(id)
        await load()
    } catch (err) {
        toast.failure((err as { error?: string })?.error || '操作失败')
    }
}

function remove(row: Row) {
    modal.show({
        title: '删除对接',
        detail: bb`确认删除对接 ${row.name as string} ?`,
        onConfirm: async () => {
            try {
                await deleteBoard(row.id as number)
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
    { key: 'name', label: '备注' },
    { key: 'host', label: '面板地址' },
    {
        key: 'inbound',
        label: '关联入站',
        render: (row: Row) => ({
            component: Stamp as Component,
            props: { color: display.protocolColor(row.protocol as string) },
            text: row.inbound as string,
        }),
    },
    { key: 'created', label: '创建时间' },
]

function add() {
    modal.show({
        title: '添加对接',
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

function connect() {
    modal.show({
        title: '一键对接',
        width: '600px',
        form: true,
        component: Connect,
        onConfirm: async (form) => {
            if (!(await form.submit())) return false

            await load()
        },
    })
}

function edit(row: Row) {
    const target = boards.value.find((item) => item.ID === row.id)
    if (!target) return

    modal.show({
        title: '编辑对接',
        width: '600px',
        form: true,
        component: Form,
        props: { inbounds: inbounds.value, board: target },
        onConfirm: async (form) => {
            if (!(await form.submit())) return false

            await load()
        },
    })
}

// 面板那边的用户，只看，没有确认按钮
function info(row: Row) {
    const target = boards.value.find((item) => item.ID === row.id)
    if (!target) return

    modal.show({
        title: '面板用户',
        width: '800px',
        form: true,
        component: Info,
        props: { board: target },
    })
}

const motion: Motion[] = [
    { icon: 'bolt', tip: '一键对接', onClick: connect },
    { icon: 'add', tip: '添加对接', onClick: add },
]

const actions = [
    { icon: 'info', tip: '信息', onClick: info },
    { icon: 'edit', tip: '编辑', onClick: edit },
    { icon: 'delete', tip: '删除', danger: true, onClick: remove },
]

onMounted(load)
</script>
