<!-- view/panel/cert/DnsForm.vue -->
<template>
    <div class="form">
        <Field>
            <template #prefix>服务商</template>
            <Select v-model:value="form.Provider" :items="providers" />
        </Field>
        <Field>
            <template #prefix>备注</template>
            <Input v-model:value="form.Name" autocomplete="off" />
        </Field>
        <Field>
            <template #prefix>{{ keyLabel }}</template>
            <Input v-model:value="form.Key" autocomplete="off" />
        </Field>
        <Field>
            <template #prefix>{{ secretLabel }}</template>
            <Input v-model:value="form.Secret" autocomplete="off" />
        </Field>
    </div>
</template>

<script setup lang="ts">
import Field from '@/component/layout/Field.vue'
import Input from '@/component/widget/Input.vue'
import Select from '@/component/widget/Select.vue'
import { saveDns } from '@/api/cert'
import type { DnsAccount } from '@/api/cert'
import { toast } from '@/composable/toast'

const props = defineProps<{
    dns?: DnsAccount // 不传是添加，传了是编辑
}>()

const emit = defineEmits<{
    changed: []
}>()

const providers = [
    { label: '阿里云', value: 'aliyun' },
    { label: 'Cloudflare', value: 'cloudflare' },
]

const form = ref({
    ID: undefined as number | undefined,
    Provider: 'aliyun',
    Name: '',
    Key: '',
    Secret: '',
})

// 两家服务商的两项凭据叫法不一样
const keyLabel = computed(() => (form.value.Provider === 'cloudflare' ? 'Email' : 'Access key'))
const secretLabel = computed(() => (form.value.Provider === 'cloudflare' ? 'API Token' : 'Secret key'))

onMounted(() => {
    const dns = props.dns

    if (!dns) return

    form.value = {
        ID: dns.ID,
        Provider: dns.Provider,
        Name: dns.Name,
        Key: dns.Key,
        Secret: dns.Secret,
    }
})

// 成功返回 true，失败返回 false（弹窗保持打开）
async function submit(): Promise<boolean> {
    try {
        await saveDns({ ...form.value })
    } catch (err) {
        toast.failure((err as { error?: string })?.error || '保存失败')

        return false
    }

    toast.success('保存成功')
    emit('changed')

    return true
}

defineExpose({ submit })
</script>

<style scoped>
.form {
    display: flex;
    flex-direction: column;
    gap: 10px;
    width: 100%;
    text-align: left;

    :deep(.prefix) {
        flex-basis: 100px;
    }
}
</style>
