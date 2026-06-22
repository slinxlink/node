<template>
    <main>
        <div class="page">
            <div class="body">
                <button class="add-btn" @click="openCreate">
                    <i class="icon">add</i>
                    添加对接
                </button>

                <List :headers="['ID', '备注', '关联入站', '启用', '操作']">
                    <tr v-for="b in boards" :key="b.ID">
                        <td class="muted">{{ b.ID }}</td>
                        <td>{{ b.Name || '-' }}</td>
                        <td>
                            <div class="tags">
                                <template v-for="ib in getInboundTags(b.Inbound)" :key="ib.id">
                                    <span class="tag" :class="ib.color">{{ ib.protocol }}</span>
                                    <span class="tag" :class="ib.color">{{ ib.port }}</span>
                                </template>
                            </div>
                        </td>
                        <td>
                            <Toggle :model-value="b.Enable" @update:model-value="toggle(b)" />
                        </td>
                        <td class="actions">
                            <button class="icon-btn" @click="openInfo(b)" title="信息">
                                <i class="icon">info</i>
                            </button>
                            <button class="icon-btn" @click="openEdit(b)" title="编辑">
                                <i class="icon">edit</i>
                            </button>
                            <button class="icon-btn danger" @click="remove(b.ID)" title="删除">
                                <i class="icon">delete</i>
                            </button>
                        </td>
                    </tr>
                    <tr v-if="boards.length === 0">
                        <td colspan="5" class="empty">暂无对接</td>
                    </tr>
                </List>
            </div>
        </div>
    </main>
    <Drawer v-model="showDrawer" :title="drawerTitle" @save="handleSave">
        <Form v-model="defaultBoard" :inbounds="inbounds" />
    </Drawer>
    <User v-model="showUserDrawer" :board="activeBoard" />
</template>

<script setup lang="ts">
import Toggle from '@/component/ui/Toggle.vue'
import Drawer from '@/component/Drawer.vue'
import List from '@/component/ui/List.vue'
import Form from '@/view/panel/board/form/Form.vue'
import User from '@/view/panel/board/widget/User.vue'
import { getBoards, saveBoard, deleteBoard, toggleBoard } from '@/api/board'
import { getInbounds } from '@/api/inbound'

const modal = inject<any>('modal')

// ── 状态 ──────────────────────────────────────────────────

const boards = ref<any[]>([])
const inbounds = ref<any[]>([])
const showDrawer = ref(false)

const baseBoard = () => ({
    Enable: true,
    Name: '',
    Host: '',
    NodeID: 0,
    Key: '',
    Inbound: 0,
    Type: 'SLINX',
    SyncInterval: 60,
})

const defaultBoard = ref<any>(baseBoard())

// ── 生命周期 ───────────────────────────────────────────────

onMounted(() => load())

// ── 工具函数 ───────────────────────────────────────────────

function getInboundTags(inboundID: number) {
    const colorMap: Record<string, string> = {
        vless: 'primary',
        vmess: 'green',
        hysteria: 'blue',
        trojan: 'purple',
    }
    const ib = inbounds.value.find(i => i.ID === inboundID)
    return ib ? [{ id: ib.ID, port: ib.Port, protocol: ib.Protocol, color: colorMap[ib.Protocol] ?? 'gray' }] : []
}

// ── 列表操作 ───────────────────────────────────────────────

async function load() {
    const [b, ib] = await Promise.all([getBoards(), getInbounds()])
    boards.value = b
    inbounds.value = ib
}

async function toggle(b: any) {
    await toggleBoard(b.ID)
    await load()
}

async function remove(id: number) {
    modal.value?.show('confirm', '确认删除该对接？', async () => {
        await deleteBoard(id)
        await load()
    })
}

// ── Drawer 操作 ────────────────────────────────────────────

const drawerTitle = ref('添加对接')

async function openCreate() {
    defaultBoard.value = baseBoard()
    drawerTitle.value = '添加对接'
    showDrawer.value = true
}

function openEdit(b: any) {
    defaultBoard.value = { ...b }
    drawerTitle.value = '编辑对接'
    showDrawer.value = true
}

async function handleSave() {
    const data = { ...defaultBoard.value }
    try {
        await saveBoard(data)
        await load()
        showDrawer.value = false
    } catch (err: any) {
        const msg = err?.error
        modal.value?.show('error', msg)
    }
}

const showUserDrawer = ref(false)
const activeBoard = ref<any>(null)

function openInfo(b: any) {
    activeBoard.value = b
    showUserDrawer.value = true
}
</script>