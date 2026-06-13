<!-- Row.vue -->
<template>
    <div class="row">
        <div class="left">
            <div class="title">{{ title }}</div>
            <div class="subtitle" v-if="subtitle">{{ subtitle }}</div>
        </div>
        <div class="right">
            <slot />
        </div>
    </div>
</template>

<script setup lang="ts">
defineProps<{ title: string; subtitle?: string }>()
</script>

<style scoped>
.row {
    display: flex;
    align-items: center;
    padding: 5px 0;

    &:last-child { border-bottom: none; }

    .left {
        flex: 1;
        display: flex;
        flex-direction: column;

        .title { color: var(--color-text-light); }
        .subtitle { font-size: var(--font-size-sm); color: var(--color-text-dark); }
    }

    .right {
        flex: 1;
        display: flex;
    }
}
</style>