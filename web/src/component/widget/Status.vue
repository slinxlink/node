<!-- component/widget/Status.vue -->
<template>
    <span class="status" :class="{ on }">
        <slot />
    </span>
</template>

<script setup lang="ts">
defineProps<{
    on?: boolean // 在线、运行中：绿点加呼吸光晕；否则是灰点
}>()
</script>

<style scoped>
.status {
    display: flex;
    align-items: center;
    gap: 6px;

    &::before {
        content: '';
        flex-shrink: 0;
        width: 8px;
        height: 8px;
        border-radius: 50%;
        background-color: var(--text-300);
    }

    &.on::before {
        background-color: var(--green);
        animation: breathe 2.5s ease-in-out infinite alternate;
    }
}

@keyframes breathe {
    from {
        box-shadow: 0 0 0 3px color-mix(in srgb, var(--green) 30%, transparent);
    }

    to {
        box-shadow: 0 0 0 5px color-mix(in srgb, var(--green) 10%, transparent);
    }
}
</style>
