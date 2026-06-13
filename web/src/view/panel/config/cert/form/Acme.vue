<template>
    <div class="form">
        <div class="form-row">
            <span class="form-label">服务商</span>
            <Select v-model="form.Provider" :options="providerOptions" />
        </div>
        <div class="form-row">
            <span class="form-label">邮箱</span>
            <Input v-model="form.Email" placeholder="example@email.com" />
        </div>
        <template v-if="form.Provider === 'zerossl'">
            <div class="form-row">
                <span class="form-label">EAB Kid</span>
                <Input v-model="form.EabKid" placeholder="EAB Kid" />
            </div>
            <div class="form-row">
                <span class="form-label">EAB Hmac</span>
                <Input v-model="form.EabHmac" placeholder="EAB Hmac Key" />
            </div>
        </template>
    </div>
</template>

<script setup lang="ts">
import Input from '@/component/ui/Input.vue'
import Select from '@/component/ui/Select.vue'

const form = defineModel<any>({ default: () => ({
    Provider: 'letsencrypt',
    Email: '',
    EabKid: '',
    EabHmac: '',
}) })

const providerOptions = [
    { label: "Let's Encrypt", value: 'letsencrypt' },
    { label: 'ZeroSSL',       value: 'zerossl' },
]
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