<template>
    <span class="stamp" :style="{ '--stamp': `var(--${color})` }">
        <slot />
    </span>
</template>

<script setup lang="ts">
withDefaults(
    defineProps<{
        color?: string // color.css 里的令牌名，如 green、primary-500、text-200
    }>(),
    {
        color: 'text-200',
    },
)
</script>

<style scoped>
.stamp {
    display: inline-flex;
    align-items: center;
    flex-shrink: 0;
    padding: 5px 10px;
    border-radius: 999px;
    background-color: color-mix(in srgb, var(--stamp) 30%, transparent);
    color: var(--stamp);
    font-size: 12px;
    line-height: 1;
    white-space: nowrap;
}
</style>
