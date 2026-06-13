<template>
    <div class="tab">
        <div class="header">
            <button
                v-for="tab in tabs"
                :key="tab.key"
                :class="{ active: modelValue === tab.key }"
                @click="$emit('update:modelValue', tab.key)">
                {{ tab.label }}
            </button>
        </div>
        <div class="content">
            <slot />
        </div>
    </div>
</template>

<script setup lang="ts">
defineProps<{
    tabs: { key: string; label: string }[]
    modelValue: string
}>()
defineEmits(['update:modelValue'])
</script>

<style scoped>
.tab {
    display: flex;
    flex-direction: column;
    gap: 20px;

    .header {
        display: flex;
        gap: 5px;

        button {
            padding: 10px 20px;
            border-radius: 999px;
            font-size: var(--font-size-sm);
            color: var(--color-text-dark);
            background-color: transparent;

            &:hover {
                background-color: var(--color-primary);
                color: var(--color-text-light);
            }
            &:active {
                background-color: var(--color-primary);
            }
            &.active {
                background-color: var(--color-bg-dark);
                color: var(--color-text-light);
            }
        }
    }
}
</style>