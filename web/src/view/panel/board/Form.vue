<!-- view/panel/board/Form.vue -->
<template>
    <div class="form">
        <Field>
            <template #prefix>启用</template>
            <Toggle v-model:value="form.Enable" />
        </Field>
        <Field>
            <template #prefix>备注</template>
            <Input v-model:value="form.Name" autocomplete="off" />
        </Field>
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
        <Field>
            <template #prefix>面板类型</template>
            <Select v-model:value="form.Type" :items="types" />
        </Field>
        <Field>
            <template #prefix>同步间隔 (s)</template>
            <Input v-model:value="form.SyncInterval" type="number" :min="0" />
        </Field>
        <Field>
            <template #prefix>关联入站</template>
            <Select v-model:value="form.Inbound" :items="items" />
        </Field>
    </div>
</template>

<script setup lang="ts">
import Field from '@/component/layout/Field.vue'
import Input from '@/component/widget/Input.vue'
import Select from '@/component/widget/Select.vue'
import Toggle from '@/component/widget/Toggle.vue'
import type { Board } from '@/api/board'
import { saveBoard } from '@/api/board'
import type { Inbound } from '@/api/inbound'
import { toast } from '@/composable/toast'

const props = defineProps<{
    inbounds: Inbound[]
    board?: Board // 不传是添加，传了是编辑
}>()

const types = [{ label: 'SLINX', value: 'SLINX' }]

// 数字都用字符串存着，提交时再转；Inbound 是入站 ID
const form = ref({
    Enable: true,
    Name: '',
    Host: '',
    NodeID: '0',
    Key: '',
    Inbound: 0,
    Type: 'SLINX',
    SyncInterval: '60',
})

const items = computed(() =>
    props.inbounds.map((item) => ({
        label: `${item.Protocol}:${item.Port}${item.Name ? ` - ${item.Name}` : ''}`,
        value: item.ID,
    })),
)

onMounted(() => {
    const board = props.board

    if (!board) return

    form.value = {
        Enable: board.Enable,
        Name: board.Name,
        Host: board.Host,
        NodeID: String(board.NodeID),
        Key: board.Key,
        Inbound: board.Inbound,
        Type: board.Type,
        SyncInterval: String(board.SyncInterval),
    }
})

// 成功返回 true，失败返回 false（弹窗保持打开）
async function submit(): Promise<boolean> {
    // 编辑时把原来的整条带上：后端整条覆盖
    const data = {
        ...props.board,
        ...form.value,
        NodeID: Number(form.value.NodeID),
        SyncInterval: Number(form.value.SyncInterval),
    }

    try {
        await saveBoard(data)
    } catch (err) {
        toast.failure((err as { error?: string })?.error || '保存失败')

        return false
    }

    toast.success('保存成功')

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
