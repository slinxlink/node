<template>
    <div class="form">
        <div class="form-row">
            <span class="form-label">申请方式</span>
            <Select v-model="form.Mode" :options="modeOptions" :disabled="props.formMode === 'reapply'" />
        </div>
        <div class="form-row" v-if="form.Mode !== 'manual'">
            <span class="form-label">域名</span>
            <Input v-model="form.Domain" placeholder="example.com" :disabled="props.formMode !== 'add'" />
        </div>
        <div class="form-row" v-if="form.Mode !== 'manual'">
            <span class="form-label">ACME 账号</span>
            <Select v-model="form.Acme" :options="acmeOptions" placeholder="选择 ACME 账号" :disabled="props.formMode === 'reapply'" />
        </div>
        <div class="form-row" v-if="form.Mode === 'dns'">
            <span class="form-label">DNS 账号</span>
            <Select v-model="form.Dns" :options="dnsOptions" placeholder="选择 DNS 账号" :disabled="props.formMode === 'reapply'" />
        </div>
        <div class="form-row" v-if="form.Mode !== 'manual'">
            <span class="form-label">自动续签</span>
            <Toggle v-model="form.AutoRenew" :disabled="props.formMode === 'reapply'" />
        </div>
        <template v-if="form.Mode === 'manual'">
            <div class="form-row">
                <span class="form-label">公钥</span>
                <Text v-model="form.CertContent" placeholder="-----BEGIN CERTIFICATE-----" :rows="6" :disabled="props.formMode === 'reapply'" />
            </div>
            <div class="form-row">
                <span class="form-label">私钥</span>
                <Text v-model="form.KeyContent" placeholder="-----BEGIN PRIVATE KEY-----" :rows="6" :disabled="props.formMode === 'reapply'" />
            </div>
        </template>
    </div>
</template>

<script setup lang="ts">
import Input from '@/component/ui/Input.vue'
import Select from '@/component/ui/Select.vue'
import Text from '@/component/ui/Text.vue'
import Toggle from '@/component/ui/Toggle.vue'

const form = defineModel<any>()

const props = defineProps<{
    acmes: any[]
    dnsList: any[]
    formMode: 'add' | 'edit' | 'reapply'
}>()

const modeOptions = [
    { label: 'DNS 验证', value: 'dns' },
    { label: 'HTTP 验证', value: 'http' },
    { label: '手动上传', value: 'manual' },
]

const acmeOptions = computed(() =>
    props.acmes.map(a => ({
        label: `${a.Email} [${formatAcmeProvider(a.Provider)}]`,
        value: a.ID
    }))
)

const dnsOptions = computed(() =>
    props.dnsList.map(d => ({
        label: `${d.Name} [${formatDnsProvider(d.Provider)}]`,
        value: d.ID
    }))
)

function formatAcmeProvider(provider: string) {
    const map: Record<string, string> = { letsencrypt: "Let's Encrypt", zerossl: 'ZeroSSL' }
    return map[provider] ?? provider
}

function formatDnsProvider(provider: string) {
    const map: Record<string, string> = { aliyun: '阿里云', cloudflare: 'Cloudflare' }
    return map[provider] ?? provider
}

watch(() => form.value.Mode, (val) => {
    if (val === 'manual') {
        form.value.AutoRenew = false
    }
})
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