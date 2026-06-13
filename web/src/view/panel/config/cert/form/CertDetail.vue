<template>
    <div class="form" v-if="cert">
        <div class="form-row">
            <span class="form-label">域名</span>
            <span>{{ cert.Domain }}</span>
        </div>
        <div class="form-row">
            <span class="form-label">申请方式</span>
            <span>{{ formatMode(cert.Mode) }}</span>
        </div>
        <div class="form-row">
            <span class="form-label">状态</span>
            <span class="tag" :class="certStatus(cert).cls">{{ certStatus(cert).label }}</span>
        </div>
        <div class="form-row">
            <span class="form-label">过期时间</span>
            <span>{{ formatTime(cert.ExpireAt) }}</span>
        </div>
        <div class="form-row" v-if="cert.Mode !== 'manual'">
            <span class="form-label">自动续签</span>
            <Toggle :model-value="cert.AutoRenew" disabled />
        </div>
        <div class="form-row" v-if="cert.Mode !== 'manual' && acme">
            <span class="form-label">ACME</span>
            <span>{{ acme.Email }} [{{ formatAcmeProvider(acme.Provider) }}]</span>
        </div>
        <div class="form-row" v-if="cert.Mode === 'dns' && dns">
            <span class="form-label">DNS</span>
            <span>{{ dns.Name }} [{{ formatDnsProvider(dns.Provider) }}]</span>
        </div>
        <div class="form-row">
            <span class="form-label">证书内容</span>
            <Text :model-value="content.cert" :rows="6" readonly />
        </div>
        <div class="form-row">
            <span class="form-label">私钥内容</span>
            <Text :model-value="content.key" :rows="6" readonly />
        </div>
    </div>
</template>

<script setup lang="ts">
import Toggle from '@/component/ui/Toggle.vue'
import Text from '@/component/ui/Text.vue'
import { getCertContent, getAcme, getDns } from '@/api/cert'
import { formatTime } from '@/util/format'

const props = defineProps<{ cert: any }>()

const content = ref({ cert: '', key: '' })
const acmes = ref<any[]>([])
const dnsList = ref<any[]>([])

const acme = computed(() => acmes.value.find(a => a.ID === props.cert?.Acme))
const dns = computed(() => dnsList.value.find(d => d.ID === props.cert?.Dns))

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

onMounted(async () => {
    const [a, d] = await Promise.all([getAcme(), getDns()])
    acmes.value = a
    dnsList.value = d
})

watch(() => props.cert, async (val) => {
    if (val?.ID) {
        content.value = await getCertContent(val.ID)
    }
}, { immediate: true })
</script>

<style scoped>
.form {
    display: flex;
    flex-direction: column;
    gap: 14px;
    padding: 16px 20px;
}

.form-row {
    display: flex;
    align-items: center;
    gap: 16px;
}

.form-label {
    width: 100px;
    flex-shrink: 0;
    font-size: var(--font-size-sm);
    color: var(--color-text-dark);
    text-align: right;
}
</style>