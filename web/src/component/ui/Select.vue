<template>
    <div class="select" :class="{ disabled: props.disabled, small: props.small }" ref="wrapRef">
        <button ref="triggerRef" class="btn" @click="toggle" :disabled="props.disabled">
            <span>{{ options.find(o => o.value === modelValue)?.label ?? placeholder }}</span>
            <i class="icon">expand_more</i>
        </button>

        <Teleport to="body">
            <div v-if="open" ref="dropdownRef" class="dropdown" :class="{ small: props.small }" :style="dropdownStyle">
                <div class="scroll">
                    <div
                        v-for="option in options"
                        :key="option.value"
                        class="item"
                        :class="{ active: option.value === modelValue }"
                        @click="select(option.value)">
                        {{ option.label }}
                    </div>
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
const triggerRef = ref<HTMLElement>()
const dropdownRef = ref<HTMLElement>()
const dropdownStyle = ref<Record<string, string>>({})

const DROPDOWN_HEIGHT = 200
const GAP = 6

function toggle() {
    if (props.disabled) return
    open.value ? close() : show()
}

function show() {
    if (!triggerRef.value) return

    const rect = triggerRef.value.getBoundingClientRect()
    const below = window.innerHeight - rect.bottom - GAP
    const above = rect.top - GAP
    const flip = below < DROPDOWN_HEIGHT + GAP && above > below

    dropdownStyle.value = {
        left: `${rect.left}px`,
        width: props.small ? 'auto' : `${rect.width}px`,
        minWidth: `${rect.width}px`,
        ...(flip
            ? { bottom: `${window.innerHeight - rect.top + GAP}px` }
            : { top: `${rect.bottom + GAP}px` }),
    }

    open.value = true

    nextTick(() => {
        window.addEventListener('mousedown', onClickOutside, true)
        window.addEventListener('resize', close)
        window.addEventListener('scroll', onScroll, true)
    })
}

function close() {
    if (!open.value) return

    open.value = false
    window.removeEventListener('mousedown', onClickOutside, true)
    window.removeEventListener('resize', close)
    window.removeEventListener('scroll', onScroll, true)
}

function select(value: string | number) {
    emit('update:modelValue', value)
    close()
}

function onClickOutside(e: MouseEvent) {
    const target = e.target as Node

    if (triggerRef.value?.contains(target) || dropdownRef.value?.contains(target)) {
        return
    }

    close()
}

function onScroll(e: Event) {
    if (dropdownRef.value?.contains(e.target as Node)) return
    close()
}

keyStack(() => open.value, (e) => {
    if (e.key === 'Escape') close()
})

onBeforeUnmount(close)
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
    max-height: 200px;
    padding: 0;
    transition: max-height 0.25s ease;
    z-index: 9999;

    &.small {
        .item {
            padding: 8px 16px;
            font-size: var(--font-size-sm);
        }
    }

    .scroll {
        max-height: 200px;
        overflow-y: auto;
        display: flex;
        flex-direction: column;
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