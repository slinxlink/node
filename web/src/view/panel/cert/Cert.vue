<template>
    <Page class="cert" title="证书" :motion="motion">
        <Section>
            <template #action>
                <Tab v-model:value="tab" :items="tabs" />
            </template>

            <Table
                v-if="tab === 'cert'"
                key="cert"
                :columns="columns"
                :data="rows"
                placeholder="暂无证书"
                :actions="actions"
            />
            <Table
                v-else-if="tab === 'acme'"
                key="acme"
                :columns="acmeColumns"
                :data="acmeRows"
                placeholder="暂无 ACME 账号"
                :actions="acmeActions"
            />
            <Table
                v-else
                key="dns"
                :columns="dnsColumns"
                :data="dnsRows"
                placeholder="暂无 DNS 账号"
                :actions="dnsActions"
            />
        </Section>
    </Page>
</template>

<script setup lang="ts">
import type { Component } from 'vue'
import type { Motion } from '@/component/layout/Header.vue'
import Page from '@/component/layout/Page.vue'
import Section from '@/component/layout/Section.vue'
import Table from '@/component/module/Table.vue'
import Stamp from '@/component/widget/Stamp.vue'
import Tab from '@/component/widget/Tab.vue'
import type { TabItem } from '@/component/widget/Tab.vue'
import Toggle from '@/component/widget/Toggle.vue'
import Detail from '@/view/panel/cert/Detail.vue'
import AcmeForm from '@/view/panel/cert/AcmeForm.vue'
import DnsForm from '@/view/panel/cert/DnsForm.vue'
import Form from '@/view/panel/cert/Form.vue'
import { display } from '@/view/panel/cert/display'
import { deleteAcme, deleteCert, deleteDns, getAcme, getCert, getDns, toggleCertRenew } from '@/api/cert'
import type { Acme, Cert, DnsAccount } from '@/api/cert'
import { modal } from '@/composable/modal'
import { bb } from '@/util/bbcode'
import { toast } from '@/composable/toast'
import { formatTime } from '@/util/format'

type Row = Record<string, unknown>

const tabs: [TabItem, ...TabItem[]] = [
    { label: '证书', value: 'cert' },
    { label: 'ACME', value: 'acme' },
    { label: 'DNS', value: 'dns' },
]

const tab = ref('cert')
const certs = ref<Cert[]>([])
const acmes = ref<Acme[]>([])
const dnsList = ref<DnsAccount[]>([])

// 表格按 id 认行，把接口的字段整理成小写的一行
const rows = computed<Row[]>(() =>
    certs.value.map((item) => {
        const status = display.status(item.ExpireAt)

        return {
            id: item.ID,
            domain: item.Domain,
            mode: display.mode(item.Mode),
            status: status.text,
            statusColor: status.color,
            expire: formatTime(item.ExpireAt),
            renew: item.AutoRenew,
            manual: item.Mode === 'manual',
        }
    }),
)

async function loadCert() {
    certs.value = await getCert()
}

async function load() {
    const [certList, acmeList, dnsAccounts] = await Promise.all([getCert(), getAcme(), getDns()])

    certs.value = certList
    acmes.value = acmeList
    dnsList.value = dnsAccounts
}

function find(row: Row) {
    return certs.value.find((item) => item.ID === row.id)
}

// 手动上传的证书没有自动续签，开关是灰的
async function toggleRenew(row: Row) {
    try {
        await toggleCertRenew(row.id as number)
        await loadCert()
    } catch (err) {
        toast.failure((err as { error?: string })?.error || '操作失败')
    }
}

function remove(row: Row) {
    modal.show({
        title: '删除证书',
        detail: bb`确认删除证书 ${row.domain as string} ?`,
        onConfirm: async () => {
            try {
                await deleteCert(row.id as number)
                await loadCert()
            } catch (err) {
                toast.failure((err as { error?: string })?.error || '删除失败')
                return false
            }
        },
    })
}

const columns = [
    {
        key: 'renew',
        label: '自动续签',
        collation: ['true', 'false'],
        render: (row: Row) => ({
            component: Toggle as Component,
            props: {
                value: row.renew,
                size: 'small',
                disabled: row.manual,
                'onUpdate:value': () => toggleRenew(row),
            },
        }),
    },
    { key: 'domain', label: '域名' },
    { key: 'mode', label: '申请方式' },
    {
        key: 'status',
        label: '状态',
        collation: ['正常', '即将过期', '已过期', '未知'],
        render: (row: Row) => ({
            component: Stamp as Component,
            props: { color: row.statusColor },
            text: row.status as string,
        }),
    },
    { key: 'expire', label: '过期时间' },
]

// 添加、重新申请、编辑共用一个表单，只是模式不同；申请的过程输出会显示在弹窗里
function open(mode: 'add' | 'edit' | 'reapply', row?: Row) {
    const target = row ? find(row) : undefined
    if (row && !target) return

    const title = { add: '申请证书', edit: '编辑证书', reapply: '重新申请' }[mode]

    modal.show({
        title,
        width: '800px',
        form: true,
        component: Form,
        props: { mode, cert: target, acmes: acmes.value, dnsList: dnsList.value, onChanged: loadCert },
        onConfirm: async (form) => {
            if (!(await form.submit())) return false
        },
    })
}

function info(row: Row) {
    const target = find(row)
    if (!target) return

    modal.show({
        title: '证书详情',
        width: '800px',
        form: true,
        component: Detail,
        props: { cert: target, acmes: acmes.value, dnsList: dnsList.value },
    })
}

// ── ACME 账号 ──────────────────────────────────────────────

const acmeRows = computed<Row[]>(() =>
    acmes.value.map((item) => ({
        id: item.ID,
        email: item.Email,
        provider: display.acmeProvider(item.Provider),
        key: item.PrivateKey ? '已配置' : '未注册',
    })),
)

async function loadAcme() {
    acmes.value = await getAcme()
}

// 新账号要向服务商注册，注册的过程输出会显示在弹窗里
function openAcme(row?: Row) {
    const target = row ? acmes.value.find((item) => item.ID === row.id) : undefined
    if (row && !target) return

    modal.show({
        title: target ? '编辑 ACME 账号' : '添加 ACME 账号',
        width: '600px',
        form: true,
        component: AcmeForm,
        props: { acme: target, onChanged: loadAcme },
        onConfirm: async (form) => {
            if (!(await form.submit())) return false
        },
    })
}

function removeAcmeAccount(row: Row) {
    modal.show({
        title: '删除 ACME 账号',
        detail: bb`确认删除 ACME 账号 ${row.email as string} ?`,
        onConfirm: async () => {
            try {
                await deleteAcme(row.id as number)
                await loadAcme()
            } catch (err) {
                toast.failure((err as { error?: string })?.error || '删除失败')
                return false
            }
        },
    })
}

const acmeColumns = [
    { key: 'email', label: '邮箱' },
    { key: 'provider', label: '服务商' },
    { key: 'key', label: '密钥' },
]

const acmeActions = [
    { icon: 'edit', tip: '编辑', onClick: (row: Row) => openAcme(row) },
    { icon: 'delete', tip: '删除', danger: true, onClick: removeAcmeAccount },
]

// ── DNS 账号 ───────────────────────────────────────────────

// 哪些账号的密钥正在显示，默认都用点遮住
const revealed = ref<Record<number, boolean>>({})

const dnsRows = computed<Row[]>(() =>
    dnsList.value.map((item) => ({
        id: item.ID,
        name: item.Name || '-',
        provider: display.dnsProvider(item.Provider),
        key: item.Key,
        secret: revealed.value[item.ID] ? item.Secret : '••••••••',
        revealed: !!revealed.value[item.ID],
    })),
)

async function loadDns() {
    dnsList.value = await getDns()
}

function openDns(row?: Row) {
    const target = row ? dnsList.value.find((item) => item.ID === row.id) : undefined
    if (row && !target) return

    modal.show({
        title: target ? '编辑 DNS 账号' : '添加 DNS 账号',
        width: '600px',
        form: true,
        component: DnsForm,
        props: { dns: target, onChanged: loadDns },
        onConfirm: async (form) => {
            if (!(await form.submit())) return false
        },
    })
}

function removeDnsAccount(row: Row) {
    modal.show({
        title: '删除 DNS 账号',
        detail: bb`确认删除 DNS 账号 ${row.name as string} ?`,
        onConfirm: async () => {
            try {
                await deleteDns(row.id as number)
                await loadDns()
            } catch (err) {
                toast.failure((err as { error?: string })?.error || '删除失败')
                return false
            }
        },
    })
}

function toggleSecret(row: Row) {
    revealed.value = { ...revealed.value, [row.id as number]: !row.revealed }
}

const dnsColumns = [
    { key: 'name', label: '备注' },
    { key: 'provider', label: '服务商' },
    { key: 'key', label: '账号标识' },
    { key: 'secret', label: '密钥', maxWidth: 520 },
]

const dnsActions = [
    { icon: 'visibility', tip: '显示密钥', show: (row: Row) => !row.revealed, onClick: toggleSecret },
    { icon: 'visibility_off', tip: '隐藏密钥', show: (row: Row) => !!row.revealed, onClick: toggleSecret },
    { icon: 'edit', tip: '编辑', onClick: (row: Row) => openDns(row) },
    { icon: 'delete', tip: '删除', danger: true, onClick: removeDnsAccount },
]

// 每个标签有自己的操作按钮
const motion = computed<Motion[]>(() => {
    if (tab.value === 'acme') return [{ icon: 'add', tip: '添加 ACME 账号', onClick: () => openAcme() }]
    if (tab.value === 'dns') return [{ icon: 'add', tip: '添加 DNS 账号', onClick: () => openDns() }]

    return [{ icon: 'add', tip: '申请证书', onClick: () => open('add') }]
})

const actions = [
    { icon: 'info', tip: '详情', onClick: info },
    {
        icon: 'refresh',
        tip: '重新申请',
        show: (row: Row) => !row.manual,
        onClick: (row: Row) => open('reapply', row),
    },
    { icon: 'edit', tip: '编辑', onClick: (row: Row) => open('edit', row) },
    { icon: 'delete', tip: '删除', danger: true, onClick: remove },
]

onMounted(load)
</script>

<style scoped>
.cert {
    /* 表格已经自带容器和内边距，放进标签容器里把它们去掉，只留容器的 */
    :deep(.content) {
        padding: 10px;
    }

    :deep(.table) {
        padding: 0;
        border-radius: 0;
        background-color: transparent;
    }
}
</style>
