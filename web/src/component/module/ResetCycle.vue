<!-- component/module/ResetCycle.vue -->
<template>
    <Field>
        <template #prefix>
            <span>流量重置</span>
            <Tip msg="到点自动清零流量[br]重置时间为重置日 0:00" />
        </template>
        <Select v-model:value="cycle" :items="cycles" />
    </Field>
    <Field v-if="cycle === 'week'">
        <template #prefix>重置日</template>
        <Select v-model:value="day" :items="weekdays" />
    </Field>
    <Field v-if="cycle === 'month'">
        <template #prefix>
            <span>重置日</span>
            <Tip msg="1 - 31[br]当月没有这一天则在月底重置" />
        </template>
        <Input v-model:value="day" type="number" :min="1" :max="31" />
    </Field>
</template>

<script setup lang="ts">
import Field from '@/component/layout/Field.vue'
import Input from '@/component/widget/Input.vue'
import Select from '@/component/widget/Select.vue'
import Tip from '@/component/widget/Tip.vue'

const cycle = defineModel<string>('cycle', { required: true })
const day = defineModel<string>('day', { required: true })

const cycles = [
    { label: '从不', value: '' },
    { label: '每日', value: 'day' },
    { label: '每周', value: 'week' },
    { label: '每月', value: 'month' },
]

const weekdays = ['周一', '周二', '周三', '周四', '周五', '周六', '周日'].map((label, i) => ({
    label,
    value: String(i + 1),
}))

// 换周期时把重置日带到合法范围：周只有 1-7，其余不管默认 1
watch(cycle, (value) => {
    const number = Number(day.value)

    if (value === 'week' && (number < 1 || number > 7)) day.value = '1'
    if (value === 'month' && (number < 1 || number > 31)) day.value = '1'
})
</script>
