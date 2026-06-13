<template>
    <div class="select" :class="{ disabled: props.disabled, small: props.small }" ref="wrapRef">
        <button class="btn" @click="toggle" :disabled="props.disabled">
            <span>{{ options.find(o => o.value === modelValue)?.label ?? placeholder }}</span>
            <i class="icon">expand_more</i>
        </button>
        <Teleport to="body">
            <div class="dropdown" :class="{ active: open, small: props.small }" :style="dropdownStyle">
                <div
                    v-for="option in options"
                    :key="option.value"
                    class="item"
                    :class="{ active: option.value === modelValue }"
                    @click="select(option.value)">
                    {{ option.label }}
                </div>
            </div>
        </Teleport>
    </div>
</template>

<script setup lang="ts">
import { keyStack } from '@/composable/keyStack'

const props = defineProps<{
    modelValue?: string | number
    options: Array<{ label: string, value: string | number }>
    placeholder?: string
    disabled?: boolean
    small?: boolean
}>()

const emit = defineEmits(['update:modelValue'])
const open = ref(false)
const wrapRef = ref<HTMLElement>()
const dropdownStyle = ref<any>({})

function updatePosition() {
    if (!wrapRef.value) return
    const rect = wrapRef.value.getBoundingClientRect()
    dropdownStyle.value = {
        top: rect.bottom + 6 + 'px',
        left: rect.left + 'px',
        width: props.small ? 'auto' : rect.width + 'px',
        minWidth: rect.width + 'px',
    }
}

function toggle() {
    if (!open.value) updatePosition()
    open.value = !open.value
}

function select(value: string | number) {
    emit('update:modelValue', value)
    open.value = false
}

function onClickOutside(e: MouseEvent) {
    if (wrapRef.value && !wrapRef.value.contains(e.target as Node)) {
        open.value = false
    }
}

keyStack(() => open.value, (e) => {
    if (e.key === 'Escape') open.value = false
})

onMounted(() => document.addEventListener('click', onClickOutside))
onUnmounted(() => document.removeEventListener('click', onClickOutside))
</script>

<style scoped>
.select {
    position: relative;
    width: 100%;

    &.small {
        display: inline-flex;
        width: auto;

        .btn {
            padding: 5px 10px 5px 15px;
            font-size: 12px;
            background-color: var(--color-bg-light);
            gap: 5px;

            &:hover {
                background-color: var(--color-bg);
            }

            .icon {
                font-size: 12px;
                color: var(--color-text-light);
            }
        }
    }

    .btn {
        display: flex;
        align-items: center;
        justify-content: space-between;
        width: 100%;
        padding: 8px 14px;
        border-radius: 999px;
        font-size: var(--font-size);
        background-color: var(--color-bg);
        color: var(--color-text-light);
        transition: background-color 0.3s;
        text-align: left;

        &:hover {
            background-color: var(--color-bg-light);
        }

        .icon {
            font-size: 18px;
            color: var(--color-text-dark);
            flex-shrink: 0;
        }
    }
}
</style>

<style>
.dropdown {
    position: fixed;
    background-color: var(--color-bg-light);
    border-radius: 10px;
    overflow: hidden;
    max-height: 0;
    transition: max-height 0.25s ease;
    z-index: 9999;



    &.small {
        .item {
            padding: 8px 16px;
            font-size: var(--font-size-sm);
        }
    }

    &.active {
        max-height: 200px;
        overflow-y: auto;
    }

    .item {
        padding: 10px 16px;
        font-size: var(--font-size);
        color: var(--color-text);
        cursor: pointer;
        white-space: nowrap;

        &:hover {
            background-color: var(--color-bg);
            color: var(--color-text-light);
        }

        &.active {
            color: var(--color-text-light);
        }
    }
}
</style>