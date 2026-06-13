<template>
    <div class="text-wrap">
        <label v-if="label" class="text-label">{{ label }}</label>
        <textarea
            class="textarea"
            :placeholder="placeholder"
            :value="modelValue"
            :rows="rows || 4"
            :readonly="readonly"
            @input="$emit('update:modelValue', ($event.target as HTMLTextAreaElement).value)"
        />
    </div>
</template>

<script setup lang="ts">
defineProps<{
    modelValue?: string
    label?: string
    placeholder?: string
    rows?: number
    readonly?: boolean
}>()
defineEmits(['update:modelValue'])
</script>

<style scoped>
.text-wrap {
    display: flex;
    flex-direction: column;
    gap: 6px;
    width: 100%;
}

.text-label {
    font-size: var(--font-size-sm);
    color: var(--color-text-dark);
}

.textarea {
    border: none;
    outline: none;
    resize: vertical;
    font-size: var(--font-size);
    color: var(--color-text-light);
    background-color: var(--color-bg);
    border-radius: 10px;
    padding: 10px 14px;
    width: 100%;
    box-sizing: border-box;
    transition: background-color 0.3s;
    font-family: inherit;
    line-height: 1.5;

    &:focus {
        background-color: var(--color-bg-light);
    }
}
</style>