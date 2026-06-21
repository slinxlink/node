<template>
    <div class="form-row">
        <div class="title-warp">
            <span class="title">{{ title }}</span>
            <span class="subtitle">{{ subtitle }}</span>
        </div>
        <div class="content">
            <slot />
        </div>
    </div>
</template>

<script setup lang="ts">
defineProps<{
    title?: string
    subtitle?: string
}>()
</script>

<style scoped>
.form-row {
    display: flex;
    align-items: center;
    gap: 10px;
    min-height: 30px;

    .title-warp {
        width: 100px;
        flex-shrink: 0;
        font-size: var(--font-size-sm);
        color: var(--color-text-dark);
        text-align: right;
        display: flex;
        flex-direction: column;
        gap: 0px;
            color: var(--color-text-dark);

        .subtitle {
            font-size: 8px;
        }
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