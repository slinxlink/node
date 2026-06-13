<template>
    <main>
        <div class="page">
            <div class="body">
                <button class="add-btn" @click="openCreate">
                    <i class="icon">add</i>
                    添加用户
                </button>

                <List :headers="['用户', '关联入站', '启用', '操作']">
                    <tr v-for="u in users" :key="u.ID">
                        <td>
                            <div class="user-info">
                                <span>{{ u.Name || '-' }}</span>
                                <span class="token">{{ u.Token }}</span>
                            </div>
                        </td>
                        <td>
                            <div class="tags">
                                <span
                                    v-for="ib in getInboundTags(u.Inbounds)"
                                    :key="ib.id"
                                    class="tag"
                                    :class="ib.color"
                                >
                                    {{ ib.port }}
                                </span>
                            </div>
                        </td>
                        <td>
                            <Toggle :model-value="u.Enable" @update:model-value="toggle(u)" />
                        </td>
                        <td class="actions">
                            <button class="icon-btn" title="二维码">
                                <i class="icon">qr_code</i>
                            </button>
                            <button class="icon-btn" title="信息">
                                <i class="icon">info</i>
                            </button>
                            <button class="icon-btn" @click="openEdit(u)" title="编辑">
                                <i class="icon">edit</i>
                            </button>
                            <button class="icon-btn danger" @click="remove(u.ID)" title="删除">
                                <i class="icon">delete</i>
                            </button>
                        </td>
                    </tr>
                    <tr v-if="users.length === 0">
                        <td colspan="4" class="empty">暂无用户</td>
                    </tr>
                </List>
            </div>
        </div>
    </main>
    <Drawer v-model="showDrawer" :title="drawerTitle" @save="handleSave">
        <Form v-model="defaultUser" :inbounds="inbounds" />
    </Drawer>
</template>

<script setup lang="ts">
import Toggle from '@/component/ui/Toggle.vue'
import Drawer from '@/component/Drawer.vue'
import List from '@/component/ui/List.vue'
import Form from '@/view/panel/user/form/Form.vue'
import { getUsers, saveUser, deleteUser, toggleUser } from '@/api/user'
import { getInbounds } from '@/api/inbound'
import { generateToken, generateUUID, generatePassword } from '@/api/generate'

const modal = inject<any>('modal')

// ── 状态 ──────────────────────────────────────────────────

const users = ref<any[]>([])
const inbounds = ref<any[]>([])
const showDrawer = ref(false)

const baseUser = () => ({
    Enable: true,
    Name: '',
    Token: '',
    UUID: '',
    Password: '',
    Inbounds: [],
})

const generatedDefaults = ref<any>({})
const defaultUser = ref<any>(baseUser())

// ── 生命周期 ───────────────────────────────────────────────

onMounted(() => {
    load()
    
    Promise.all([
        generateToken(),
        generateUUID(),
        generatePassword(),
    ]).then(([tokenRes, uuidRes, passwordRes]) => {
        generatedDefaults.value = {
            Token: tokenRes.token,
            UUID: uuidRes.uuid,
            Password: passwordRes.password,
        }
        defaultUser.value = { ...baseUser(), ...generatedDefaults.value }
    })
})

// ── 工具函数 ───────────────────────────────────────────────

function getInboundTags(inboundsJSON: string) {
    const colorMap: Record<string, string> = {
        vless: 'primary',
        vmess: 'green',
        hysteria: 'blue',
    }
    try {
        const ids: number[] = JSON.parse(inboundsJSON || '[]')
        return ids.map(id => {
            const ib = inbounds.value.find(i => i.ID === id)
            return ib ? { id, port: ib.Port, color: colorMap[ib.Protocol] ?? 'gray' } : null
        }).filter((ib): ib is { id: number, port: number, color: string } => ib !== null)
    } catch {
        return []
    }
}

// ── 列表操作 ───────────────────────────────────────────────


async function load() {
    const [u, ib] = await Promise.all([getUsers(), getInbounds()])
    users.value = u
    inbounds.value = ib
}

async function toggle(u: any) {
    await toggleUser(u.ID)
    await load()
}

async function remove(id: number) {
    modal.value?.show('confirm', '确认删除该用户？', async () => {
        await deleteUser(id)
        await load()
    })
}

// ── Drawer 操作 ────────────────────────────────────────────

const drawerTitle = ref('添加用户')

async function openCreate() {
    const [tokenRes, uuidRes, passwordRes] = await Promise.all([
        generateToken(),
        generateUUID(),
        generatePassword(),
    ])
    generatedDefaults.value = {
        Token: tokenRes.token,
        UUID: uuidRes.uuid,
        Password: passwordRes.password,
    }
    defaultUser.value = { ...baseUser(), ...generatedDefaults.value }
    drawerTitle.value = '添加用户'
    showDrawer.value = true
}
function openEdit(u: any) {
    defaultUser.value = {
        ...u,
        Inbounds: u.Inbounds ? JSON.parse(u.Inbounds).map(String) : [],
    }
    drawerTitle.value = '编辑用户'
    showDrawer.value = true
}

async function handleSave() {
    const data = { ...defaultUser.value }

    if (Array.isArray(data.Inbounds)) {
        data.Inbounds = JSON.stringify(data.Inbounds.map(Number))
    }

    try {
        await saveUser(data)
        await load()
        showDrawer.value = false
    } catch (err: any) {
        const msg = err?.error
        modal.value?.show('error', msg)
    }
}
</script>

<style scoped>
.user-info {
    display: flex;
    flex-direction: column;
    gap: 2px;

    .token {
        font-size: var(--font-size-sm);
        color: var(--color-text-dark);
    }
}

:deep(th),
:deep(td) {
    &:first-child {
        padding-right: 25px !important;
    }
}
</style>