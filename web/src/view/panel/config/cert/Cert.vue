<template>
    <div class="body">
        <!-- 证书 -->
        <button class="add-btn" @click="openCertDrawer('add')">
            <i class="icon">add</i>
            申请证书
        </button>
        <List :headers="['域名', '申请方式', '状态', '过期时间', '自动续签', '操作']">
            <tr v-for="cert in certs" :key="cert.ID">
                <td>{{ cert.Domain }}</td>
                <td>{{ formatMode(cert.Mode) }}</td>
                <td>
                    <div class="tags">
                        <span class="tag" :class="certStatus(cert).cls">{{ certStatus(cert).label }}</span>
                    </div>
                </td>
                <td class="muted">{{ formatTime(cert.ExpireAt) }}</td>
                <td><Toggle :model-value="cert.AutoRenew" :disabled="cert.Mode === 'manual'" @update:model-value="toggleRenew(cert)" /></td>
                <td class="actions">
                    <button class="icon-btn" @click="openCertDetail(cert)" title="详情"><i class="icon">info</i></button>
                    <button v-if="cert.Mode !== 'manual'" class="icon-btn" @click="openCertDrawer('reapply', cert)" title="重新申请"><i class="icon">refresh</i></button>
                    <button class="icon-btn" @click="openCertDrawer('edit', cert)" title="编辑"><i class="icon">edit</i></button>
                    <button class="icon-btn danger" @click="removeCert(cert.ID)" title="删除"><i class="icon">delete</i></button>
                </td>
            </tr>
            <tr v-if="certs.length === 0">
                <td colspan="6" class="empty">暂无证书</td>
            </tr>
        </List>

        <!-- ACME 账号 -->
        <button class="add-btn" @click="openAcmeDrawer">
            <i class="icon">add</i>
            添加 ACME 账号
        </button>
        <List class="acme-list" :headers="['邮箱', '服务商', '密钥', '操作']">
            <tr v-for="acme in acmes" :key="acme.ID">
                <td>{{ acme.Email }}</td>
                <td>{{ formatAcmeProvider(acme.Provider) }}</td>
                <td class="muted">{{ acme.PrivateKey ? '已配置' : '未注册' }}</td>
                <td class="actions">
                    <button class="icon-btn" @click="openAcmeDrawer(acme)" title="编辑"><i class="icon">edit</i></button>
                    <button class="icon-btn danger" @click="removeAcme(acme.ID)" title="删除"><i class="icon">delete</i></button>
                </td>
            </tr>
            <tr v-if="acmes.length === 0">
                <td colspan="4" class="empty">暂无 ACME 账号</td>
            </tr>
        </List>

        <!-- DNS 账号 -->
        <button class="add-btn" @click="openDnsDrawer">
            <i class="icon">add</i>
            添加 DNS 账号
        </button>
        <List class="dns-list" :headers="['备注', '服务商', '账号标识', '密钥', '操作']">
            <tr v-for="dns in dnsList" :key="dns.ID">
                <td>{{ dns.Name }}</td>
                <td>{{ formatDnsProvider(dns.Provider) }}</td>
                <td class="muted">{{ dns.Key }}</td>
                <td class="muted">{{ showSecret[dns.ID] ? dns.Secret : '••••••••' }}</td>
                <td class="actions">
                    <button class="icon-btn" @click="toggleSecret(dns.ID)">
                        <i class="icon">{{ showSecret[dns.ID] ? 'visibility_off' : 'visibility' }}</i>
                    </button>
                    <button class="icon-btn" @click="openDnsDrawer(dns)" title="编辑"><i class="icon">edit</i></button>
                    <button class="icon-btn danger" @click="removeDns(dns.ID)" title="删除"><i class="icon">delete</i></button>
                </td>
            </tr>
            <tr v-if="dnsList.length === 0">
                <td colspan="5" class="empty">暂无 DNS 账号</td>
            </tr>
        </List>
    </div>

    <Drawer
        v-model="showCertDrawer"
        :title="certFormMode === 'add' ? '申请证书' : certFormMode === 'edit' ? '编辑证书' : '重新申请'"
        :saveText="certFormMode === 'edit' ? '保存' : '申请'"
        :loading="certTask.loading"
        :done="certTask.done"
        @save="certFormMode === 'edit' ? updateCert() : addCert()"
    >
        <CertForm
            v-model="certForm"
            :acmes="acmes"
            :dns-list="dnsList"
            :form-mode="certFormMode"
        />
        <Task v-if="certTask.id" :task-id="certTask.id" @done="onCertTaskDone" @error="onCertTaskError" />
    </Drawer>
    <Drawer v-model="showCertDetailDrawer" title="证书详情" :footer="false">
        <CertDetail :cert="selectedCert" />
    </Drawer>
    <Drawer v-model="showAcmeDrawer" title="ACME 账号" :loading="acmeTask.loading" :done="acmeTask.done" @save="saveAcmeForm">
        <AcmeForm v-model="acmeForm" />
        <Task v-if="acmeTask.id" :task-id="acmeTask.id" @done="onAcmeTaskDone" @error="onAcmeTaskError" />
    </Drawer>
    <Drawer v-model="showDnsDrawer" title="DNS 账号" @save="saveDnsForm">
        <DnsForm v-model="dnsForm" />
    </Drawer>
</template>

<script setup lang="ts">
import List from '@/component/ui/List.vue'
import Toggle from '@/component/ui/Toggle.vue'
import Drawer from '@/component/Drawer.vue'
import CertDetail from '@/view/panel/config/cert//form/CertDetail.vue'
import CertForm from '@/view/panel/config/cert//form/Cert.vue'
import AcmeForm from '@/view/panel/config/cert/form/Acme.vue'
import DnsForm from '@/view/panel/config/cert/form/Dns.vue'
import Task from '@/component/widget/Task.vue'
import { formatTime } from '@/util/format'
import { getCert, saveCert, deleteCert, applyCert, getCertContent, getAcme, saveAcme, deleteAcme, getDns, saveDns, deleteDns } from '@/api/cert'

const modal = inject<any>('modal')

const certs = ref<any[]>([])
const acmes = ref<any[]>([])
const dnsList = ref<any[]>([])
const showSecret = ref<Record<number, boolean>>({})

const showAcmeDrawer = ref(false)
const showDnsDrawer = ref(false)
const acmeForm = ref<any>({ Provider: 'letsencrypt', Email: '', EabKid: '', EabHmac: '' })
const dnsForm = ref<any>({ Name: '', Provider: 'aliyun', Key: '', Secret: '' })

// ── 证书任务 ──────────────────────────────────────────────

const certTask = ref({
    loading: false,
    done: false,
    id: '',
})

function onCertTaskDone() {
    certTask.value.loading = false
    certTask.value.done = true
    loadCert()
}

function onCertTaskError() {
    certTask.value.loading = false
    certTask.value.done = false
    loadCert()
}

// ── ACME 任务 ─────────────────────────────────────────────

const acmeTask = ref({
    loading: false,
    done: false,
    id: '',
})

function onAcmeTaskDone() {
    acmeTask.value.loading = false
    acmeTask.value.done = true
    loadAcme()
}

function onAcmeTaskError() {
    acmeTask.value.loading = false
    acmeTask.value.done = false
    loadAcme()
}

// ── 加载 ──────────────────────────────────────────────────

onMounted(async () => {
    await Promise.all([loadCert(), loadAcme(), loadDns()])
})

async function loadCert() { certs.value = await getCert() }
async function loadAcme() { acmes.value = await getAcme() }
async function loadDns() { dnsList.value = await getDns() }

// ── 工具函数 ───────────────────────────────────────────────

function formatMode(mode: string) {
    const map: Record<string, string> = { dns: 'DNS验证', http: 'HTTP验证', manual: '手动上传' }
    return map[mode] ?? mode
}

function formatAcmeProvider(provider: string) {
    const map: Record<string, string> = { letsencrypt: "Let's Encrypt", zerossl: 'ZeroSSL' }
    return map[provider] ?? provider
}

function formatDnsProvider(provider: string) {
    const map: Record<string, string> = { aliyun: '阿里云', cloudflare: 'Cloudflare' }
    return map[provider] ?? provider
}

function certStatus(cert: any) {
    if (!cert.ExpireAt) return { label: '未知', cls: 'gray' }
    const days = (new Date(cert.ExpireAt).getTime() - Date.now()) / 86400000
    if (days < 0) return { label: '已过期', cls: 'red' }
    if (days < 10) return { label: '即将过期', cls: 'yellow' }
    return { label: '正常', cls: 'green' }
}

function toggleSecret(id: number) {
    showSecret.value[id] = !showSecret.value[id]
}

// ── 证书操作 ───────────────────────────────────────────────

const certFormMode = ref<'add' | 'edit' | 'reapply'>('add')
const showCertDrawer = ref(false)
const certForm = ref<any>({ Domain: '', Mode: 'dns', Acme: 0, Dns: 0, AutoRenew: true })

async function openCertDrawer(mode: 'add' | 'edit' | 'reapply', cert?: any) {
    certFormMode.value = mode
    certForm.value = cert ? { ...cert } : { Domain: '', Mode: 'dns', Acme: 0, Dns: 0, AutoRenew: true }
    certTask.value.loading = mode === 'reapply'
    certTask.value.done = false
    certTask.value.id = ''
    showCertDrawer.value = true
    if (cert?.ID && mode !== 'add') {
        const content = await getCertContent(cert.ID)
        certForm.value.CertContent = content.cert
        certForm.value.KeyContent = content.key
    }
    if (mode === 'reapply' && cert) {
        startApply(cert.ID)
    }
}

async function startApply(id: number) {
    certTask.value.id = ''
    await nextTick()
    try {
        const res = await applyCert(id)
        if (res.task_id) {
            certTask.value.loading = true
            certTask.value.done = false
            certTask.value.id = res.task_id
        }
    } catch (err: any) {
        modal.value?.show('error', err?.error)
        certTask.value.loading = false
    }
}

async function addCert() {
    try {
        let id = certForm.value.ID
        if (!id) {
            const saved = await saveCert(certForm.value)
            id = saved.ID
            certForm.value.ID = saved.ID
        }
        if (certForm.value.Mode === 'manual') {
            await loadCert()
            showCertDrawer.value = false
            modal.value?.show('success', '证书上传成功')
            return
        }
        await startApply(id)
    } catch (err: any) {
        modal.value?.show('error', err?.error)
    }
}

async function updateCert() {
    try {
        await saveCert(certForm.value)
        await loadCert()
        showCertDrawer.value = false
        modal.value?.show('success', '保存成功')
    } catch (err: any) {
        modal.value?.show('error', err?.error)
    }
}

async function toggleRenew(cert: any) {
    cert.AutoRenew = !cert.AutoRenew
    await saveCert(cert)
    await loadCert()
}

async function removeCert(id: number) {
    modal.value?.show('confirm', '确认删除该证书？', async () => {
        await deleteCert(id)
        await loadCert()
    })
}

const showCertDetailDrawer = ref(false)
const selectedCert = ref<any>(null)

function openCertDetail(cert: any) {
    selectedCert.value = cert
    showCertDetailDrawer.value = true
}

// ── ACME 操作 ──────────────────────────────────────────────

function openAcmeDrawer(acme?: any) {
    acmeForm.value = acme ? { ...acme } : { Provider: 'letsencrypt', Email: '', EabKid: '', EabHmac: '' }
    acmeTask.value.loading = false
    acmeTask.value.done = false
    acmeTask.value.id = ''
    showAcmeDrawer.value = true
}

async function saveAcmeForm() {
    try {
        const res = await saveAcme(acmeForm.value)
        if (res.task_id) {
            acmeTask.value.loading = true
            acmeTask.value.done = false
            acmeTask.value.id = res.task_id
        } else {
            await loadAcme()
            showAcmeDrawer.value = false
        }
    } catch (err: any) {
        modal.value?.show('error', err?.error)
    }
}

async function removeAcme(id: number) {
    modal.value?.show('confirm', '确认删除该 ACME 账号？', async () => {
        await deleteAcme(id)
        await loadAcme()
    })
}

// ── DNS 操作 ───────────────────────────────────────────────

function openDnsDrawer(dns?: any) {
    dnsForm.value = dns ? { ...dns } : { Provider: 'aliyun', Name: '', Key: '', Secret: '' }
    showDnsDrawer.value = true
}

async function saveDnsForm() {
    try {
        await saveDns(dnsForm.value)
        await loadDns()
        showDnsDrawer.value = false
        modal.value?.show('success', '保存成功')
    } catch (err: any) {
        modal.value?.show('error', err?.error)
    }
}

async function removeDns(id: number) {
    modal.value?.show('confirm', '确认删除该 DNS 账号？', async () => {
        await deleteDns(id)
        await loadDns()
    })
}
</script>

<style scoped>
.body {
    display: flex;
    flex-direction: column;
    gap: 20px;
}

.add-btn {
    margin-left: auto;
}

:deep(table) {
    td:nth-child(1), th:nth-child(1) {
        width: 200px;
        max-width: 200px;
        overflow: hidden;
        text-overflow: ellipsis;
    }
    td:nth-child(2), th:nth-child(2) {
        width: 110px;
        max-width: 110px;
        overflow: hidden;
        text-overflow: ellipsis;
    }
    td:nth-child(3), th:nth-child(3) {
        width: 200px;
        max-width: 250px;
        overflow: hidden;
        text-overflow: ellipsis;
    }
}

:deep(.acme-list table) {
    td:nth-last-child(2), th:nth-last-child(2) {
        width: auto;
        max-width: auto;
    }
}

:deep(.dns-list table) {
    td:nth-last-child(2), th:nth-last-child(2) {
        width: auto;
        max-width: auto;
    }
}
</style>