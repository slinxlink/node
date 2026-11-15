<template>
    <button
        :type="type"
        :disabled="disabled"
        :tabindex="tabindex"
        :class="[variant, { selected }]"
        @click="onClick"
    >
        <slot />
    </button>
</template>

<script setup lang="ts">
withDefaults(
    defineProps<{
        type?: 'button' | 'submit' | 'reset'
        variant?: 'default' | 'primary' | 'text' | 'border' | 'danger'
        selected?: boolean
        disabled?: boolean
        tabindex?: number
    }>(),
    {
        type: 'button',
        variant: 'default',
    },
)

const emit = defineEmits<{
    click: [e: MouseEvent]
}>()

function onClick(e: MouseEvent) {
    emit('click', e)
}
</script>

<style scoped>
button {
    position: relative;
    display: flex;
    align-items: center;
    justify-content: center;
    width: 100%;
    padding: 10px 15px;
    gap: 5px;
    text-align: center;
    color: var(--text-100);
    font-size: 14px;
    line-height: 1;
    background-color: var(--background-300);
    border: none;
    border-radius: 999px;
    outline: none;
    -webkit-user-select: none; /* iOS Safari 需要前缀 */
    user-select: none;
    cursor: pointer;
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

    &.selected {
        background-color: var(--primary-600);
        cursor: default;
        pointer-events: none;
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

    &.text {
        width: auto;
        padding: 0;
        background-color: transparent;
        color: var(--text-200);

        &:hover,
        &:focus-visible {
            background-color: transparent;
            color: var(--primary-500);
        }

        &:active {
            background-color: transparent;
            color: var(--primary-600);
        }
    }

    &.border {
        background-color: var(--background-300);
        color: var(--text-200);
        border: 1px solid;
        border-color: var(--background-100);

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
        &.selected {
            background-color: var(--background-300);
            border-color: var(--primary-100);
        }
    }

    &.danger {
        background-color: var(--background-300);
        color: var(--red);
        border: 1px solid;
        border-color: var(--red);

        &:hover,
        &:focus-visible {
            background-color: color-mix(in srgb, var(--red) 25%, var(--background-300));
        }

        &:active {
            background-color: var(--red);
            color: var(--white);
        }
    }

    &:disabled {
        cursor: not-allowed;
        opacity: 0.5;
        filter: blur(0) grayscale(50%); /* blur(0) 同上 */
    }
}
</style>
