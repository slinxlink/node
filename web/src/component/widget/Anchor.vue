<template>
    <component
        :is="type === 'route' ? RouterLink : 'a'"
        v-bind="disabled ? {} : type === 'route' ? { to: href } : { href }"
        :target="external ? '_blank' : target"
        :rel="external ? 'noopener noreferrer' : rel"
        :class="[variant, { disabled }]"
        :aria-disabled="disabled || undefined"
        :tabindex="disabled ? -1 : undefined"
        @click="onClick"
    >
        <slot />
    </component>
</template>

<script setup lang="ts">
import { RouterLink } from 'vue-router'

const props = defineProps<{
    href?: string // url 完整地址，route 站内路径
    type?: 'url' | 'route' // url（默认）、route
    target?: string
    rel?: string
    variant?: 'default' | 'primary' | 'text' | 'plain' | 'border'
    external?: boolean // 新标签页打开，自动带 rel
    disabled?: boolean // 禁用，不生成链接
}>()

const emit = defineEmits<{
    click: [e: MouseEvent]
}>()

function onClick(e: MouseEvent) {
    if (props.disabled) return e.preventDefault()
    emit('click', e)
}
</script>

<style scoped>
a {
    position: relative;
    display: flex;
    align-items: center;
    justify-content: center;
    width: 100%;
    padding: 10px 15px;
    border: none;
    border-radius: 999px;
    outline: none;
    background-color: var(--background-300);
    color: var(--text-100);
    font-size: 14px;
    text-align: center;
    cursor: pointer;
    text-decoration: none;
    -webkit-user-select: none; /* iOS Safari 需要前缀 */
    user-select: none;
    filter: blur(0); /* 临时处理 Chrome 背景色过渡闪回，待修复后删除 */
    transition:
        background-color 0.3s,
        color 0.3s,
        border-color 0.3s;

    &:hover,
    &:focus-visible {
        background-color: var(--background-100);
    }

    &:active {
        background-color: var(--primary-600);
    }

    &.primary {
        background-color: var(--primary-300);
        color: var(--white);

        &:hover,
        &:focus-visible {
            background-color: var(--primary-500);
        }

        &:active {
            background-color: var(--primary-600);
        }
    }

    &.plain {
        width: auto;
        padding: 0;
        background-color: transparent;
        color: var(--text-200);

        &:hover,
        &:focus-visible {
            background-color: transparent;
            color: var(--text-100);
        }

        &:active {
            background-color: transparent;
            color: var(--text-300);
        }
    }

    &.text {
        width: auto;
        padding: 0;
        background-color: transparent;
        color: var(--primary-100);

        &:hover,
        &:focus-visible {
            background-color: transparent;
            color: var(--primary-200);
        }

        &:active {
            background-color: transparent;
            color: var(--primary-400);
        }
    }

    &.border {
        background-color: var(--background-300);
        color: var(--text-200);
        border: 1px solid var(--background-100);

        &:hover,
        &:focus-visible {
            background-color: var(--background-300);
            border-color: var(--primary-200);
        }

        &:active {
            background-color: var(--primary-400);
            color: var(--white);
            border-color: var(--primary-400);
            --invert: 1;
        }
    }

    &.disabled {
        cursor: not-allowed;
        opacity: 0.5;
        filter: blur(0) grayscale(50%); /* blur(0) 同上 */
    }
}
</style>
