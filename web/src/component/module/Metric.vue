<!-- component/module/Metric.vue -->
<template>
    <Box :title="title">
        <div class="items">
            <div v-for="item in items" :key="item.label" class="item">
                <span class="label">{{ item.label }}</span>
                <span class="value">
                    <Icon :name="item.icon" />
                    {{ item.value }}
                </span>
            </div>
        </div>
    </Box>
</template>

<script setup lang="ts">
import Box from '@/component/layout/Box.vue'
import Icon from '@/component/widget/Icon.vue'

// 一张卡片两个数：标题、每项的名字、图标和值
export interface MetricItem {
    label: string
    icon: string
    value: string | number
}

defineProps<{
    title: string
    items: [MetricItem, MetricItem]
}>()
</script>

<style scoped>
.items {
    display: grid;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    gap: 10px;
}

.item {
    display: flex;
    flex-direction: column;
    gap: 5px;

    .label {
        color: var(--text-300);
        font-size: 12px;
    }

    .value {
        display: flex;
        align-items: center;
        gap: 5px;
        color: var(--text-100);
        font-weight: bold;
    }
}
</style>
