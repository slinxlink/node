<template>
    <div class="form-row">
        <span class="label">{{ label }}</span>
        <div class="content">
            <slot />
        </div>
    </div>
</template>

<script setup lang="ts">
defineProps<{
    label?: string
}>()
</script>

<style scoped>
.form-row {
    display: flex;
    align-items: center;
    gap: 10px;
    min-height: 30px;

    .label {
        width: 100px;
        flex-shrink: 0;
        font-size: var(--font-size-sm);
        color: var(--color-text-dark);
        text-align: right;
    }

    .content {
        display: flex;
        flex-direction: row;
        align-items: center;
        gap: 5px;
        flex: 1;
        min-width: 0;
        font-size: var(--font-size-sm);
        color: var(--color-text-light);

        &:deep(> *:first-child) {
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
            min-width: 0;
        }
    }
}
</style>