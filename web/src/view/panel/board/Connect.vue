<!-- view/panel/board/Connect.vue -->
<template>
    <div class="form">
        <Field>
            <template #prefix>面板地址</template>
            <Input v-model:value="form.Host" placeholder="https://example.com" autocomplete="off" />
        </Field>
        <Field>
            <template #prefix>节点 ID</template>
            <Input v-model:value="form.NodeID" type="number" :min="0" />
        </Field>
        <Field>
            <template #prefix>通讯密钥</template>
            <Input v-model:value="form.Key" autocomplete="off" />
        </Field>
    </div>
</template>

<script setup lang="ts">
import Field from '@/component/layout/Field.vue'
import Input from '@/component/widget/Input.vue'
import { connectBoard } from '@/api/board'
import { toast } from '@/composable/toast'

const form = ref({
    Host: '',
    NodeID: '0',
    Key: '',
})

// 成功返回 true，失败返回 false（弹窗保持打开）；对接会连面板拉取信息并建好入站，检查都在后端
async function submit(): Promise<boolean> {
    try {
        await connectBoard({ ...form.value, NodeID: Number(form.value.NodeID) })
    } catch (err) {
        toast.failure((err as { error?: string })?.error || '对接失败')

        return false
    }

    toast.success('对接成功')

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
