<!-- component/layout/Field.vue -->
<!-- 只管三栏布局：prefix 左、main 中（默认插槽）、suffix 右；里面放什么、样式怎么调，都由调用方决定 -->
<template>
    <div class="field">
        <div v-if="$slots.prefix" class="prefix">
            <slot name="prefix" />
        </div>
        <div class="main">
            <slot />
        </div>
        <div v-if="$slots.suffix" class="suffix">
            <slot name="suffix" />
        </div>
    </div>
</template>

<style scoped>
.field {
    display: flex;
    align-items: center;
    gap: 10px;
    min-height: 40px;
    font-size: 14px;

    .prefix {
        display: flex;
        flex-shrink: 0;
        align-items: center;
        gap: 10px;
        min-width: 60px;
        color: var(--text-300);
    }

    .main {
        display: flex;
        flex: 1;
        align-items: center;
        gap: 10px;
        min-width: 0;
        overflow-wrap: anywhere;
        color: var(--text-100);
    }

    .suffix {
        display: flex;
        flex-shrink: 0;
        align-items: center;
        gap: 10px;
        margin-left: auto;
        color: var(--text-300);
    }
}
</style>
