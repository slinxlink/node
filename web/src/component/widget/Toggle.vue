<template>
    <button
        type="button"
        role="switch"
        :aria-checked="value"
        :disabled="disabled"
        :class="{ on: value, small: size === 'small' }"
        @click="emit('update:value', !value)"
    >
        <span class="thumb"></span>
    </button>
</template>

<script setup lang="ts">
defineProps<{
    value: boolean
    disabled?: boolean
    size?: 'small'
}>()

const emit = defineEmits<{
    'update:value': [value: boolean]
}>()
</script>

<style scoped>
button {
    display: flex;
    align-items: center;
    flex-shrink: 0;
    width: 50px;
    height: 30px;
    padding: 5px;
    border: none;
    border-radius: 999px;
    outline: none;
    background-color: var(--background-200);
    cursor: pointer;
    transition: background-color 0.3s ease;

    &:hover:not(:disabled),
    &:focus-visible {
        background-color: var(--background-100);
    }

    &.on {
        background-color: var(--primary-300);

        &:hover:not(:disabled),
        &:focus-visible {
            background-color: var(--primary-500);
        }
    }

    &:disabled {
        opacity: 0.5;
        cursor: not-allowed;
    }

    .thumb {
        flex-shrink: 0;
        width: 20px;
        height: 20px;
        border-radius: 50%;
        background-color: var(--white);
        transition: transform 0.3s ease;
    }

    &.on .thumb {
        transform: translateX(20px);
    }

    &:active:not(:disabled) .thumb {
        transform: scale(0.9);
    }

    &.on:active:not(:disabled) .thumb {
        transform: translateX(20px) scale(0.9);
    }

    &.small {
        width: 36px;
        height: 22px;
        padding: 4px;

        .thumb {
            width: 14px;
            height: 14px;
        }

        &.on .thumb {
            transform: translateX(14px);
        }

        &.on:active:not(:disabled) .thumb {
            transform: translateX(14px) scale(0.9);
        }
    }
}
</style>
