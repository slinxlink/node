<template>
    <hr v-if="!$slots.default" :style="expand ? { '--expand': expand } : undefined" />
    <div v-else :style="{ '--expand': expand ?? '20px' }">
        <slot />
    </div>
</template>

<script setup lang="ts">
defineProps<{
    expand?: string // 线向两边多伸出的长度，抵掉容器的 padding；负数向内缩
}>()
</script>

<style scoped>
/* 没有文字时只画一条线 */
hr {
    width: calc(100% + 2 * var(--expand, 0px));
    margin-inline: calc(-1 * var(--expand, 0px));
    height: 1px;
    border: none;
    background-color: var(--background-200);
}

div {
    display: flex;
    align-items: center;
    width: 100%;
    color: var(--text-300);
    font-size: 14px;
    text-align: center;

    &::before,
    &::after {
        content: '';
        flex: 1;
        border-bottom: 1px solid var(--background-200);
        margin: 0 calc(-1 * var(--expand));
    }

    &::before {
        margin-right: 10px;
    }

    &::after {
        margin-left: 10px;
    }
}
</style>
