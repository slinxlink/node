<template>
    <div class="toggle" :class="{ active: modelValue, disabled: disabled }" @click="$emit('update:modelValue', !modelValue)">
        <div class="thumb" />
    </div>
</template>

<script setup lang="ts">
defineProps<{ modelValue: boolean, disabled?: boolean }>()
defineEmits(['update:modelValue'])
</script>

<style scoped>
.toggle {
    width: 42px;
    height: 24px;
    border-radius: 12px;
    background: var(--color-bg-light);
    position: relative;
    cursor: pointer;
    transition: background 0.25s ease;
    flex-shrink: 0;

    &.active {
        background: var(--color-primary-light);
    }

    &.disabled {
        opacity: 0.5;
        cursor: not-allowed;
        pointer-events: none;
    }

    .thumb {
        position: absolute;
        top: 3px;
        left: 3px;
        width: 18px;
        height: 18px;
        border-radius: 50%;
        background: #fff;
        transition: left 0.25s ease;
        box-shadow: 0 1px 3px rgba(0,0,0,0.3);
    }

    &.active .thumb {
        left: 21px;
    }
}
</style>
