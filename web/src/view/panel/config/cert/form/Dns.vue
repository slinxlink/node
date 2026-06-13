<template>
    <div class="form">
        <div class="form-row">
            <span class="form-label">服务商</span>
            <Select v-model="form.Provider" :options="providerOptions" />
        </div>
        <div class="form-row">
            <span class="form-label">备注</span>
            <Input v-model="form.Name" />
        </div>
        <div class="form-row">
            <span class="form-label">{{ keyLabel }}</span>
            <Input v-model="form.Key" />
        </div>
        <div class="form-row">
            <span class="form-label">{{ secretLabel }}</span>
            <Input v-model="form.Secret" />
        </div>
    </div>
</template>

<script setup lang="ts">
import Input from '@/component/ui/Input.vue'
import Select from '@/component/ui/Select.vue'

const form = defineModel<any>({ default: () => ({
    Provider: 'aliyun',
    Name: '',
    Key: '',
    Secret: '',
}) })

const providerOptions = [
    { label: '阿里云', value: 'aliyun' },
    { label: 'Cloudflare', value: 'cloudflare' },
]

const keyLabel = computed(() => {
    const map: Record<string, string> = {
        aliyun: 'Access key',
        cloudflare: 'Email',
    }
    return map[form.value.Provider] ?? 'Key'
})

const secretLabel = computed(() => {
    const map: Record<string, string> = {
        aliyun: 'Secret key',
        cloudflare: 'API Token',
    }
    return map[form.value.Provider] ?? 'Secret'
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