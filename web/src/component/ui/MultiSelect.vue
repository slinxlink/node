<template>
    <div class="multiselect-wrap" ref="wrapRef" :class="{ disabled }">
        <div class="multiselect-input" @click="!disabled && toggle()">
            <div class="multiselect-tags">
                <span class="multiselect-tag" v-for="val in modelValue" :key="val">
                    {{ options.find(o => o.value === val)?.label }}
                    <i class="icon" v-if="!disabled" @click.stop="remove(val)">close</i>
                </span>
                <span class="multiselect-placeholder" v-if="!modelValue?.length">{{ placeholder }}</span>
            </div>
        </div>
        <Teleport to="body">
            <div class="multiselect-dropdown" :class="{ active: open }" :style="dropdownStyle">
                <div
                    v-for="option in options"
                    :key="option.value"
                    class="multiselect-item"
                    :class="{ active: modelValue?.includes(option.value) }"
                    @click="select(option.value)">
                    {{ option.label }}
                    <i class="icon" v-if="modelValue?.includes(option.value)">check</i>
                </div>
            </div>
        </Teleport>
    </div>
</template>

<script setup lang="ts">
const props = defineProps<{
    modelValue?: string[]
    options: Array<{ label: string, value: string }>
    placeholder?: string
    disabled?: boolean
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
        width: rect.width + 'px',
    }
}

function toggle() {
    if (!open.value) updatePosition()
    open.value = !open.value
}

function select(value: string) {
    const current = props.modelValue ?? []
    if (current.includes(value)) {
        emit('update:modelValue', current.filter(v => v !== value))
    } else {
        emit('update:modelValue', [...current, value])
    }
}

function remove(value: string) {
    emit('update:modelValue', (props.modelValue ?? []).filter(v => v !== value))
}

function onClickOutside(e: MouseEvent) {
    if (wrapRef.value && !wrapRef.value.contains(e.target as Node)) {
        open.value = false
    }
}

onMounted(() => document.addEventListener('click', onClickOutside))
onUnmounted(() => document.removeEventListener('click', onClickOutside))
</script>

<style scoped>
.multiselect-wrap {
    position: relative;
    width: 100%;

    &.disabled .multiselect-input {
        opacity: 0.5;
        cursor: not-allowed;
    }
}

.multiselect-input {
    min-height: 36px;
    padding: 4px 40px 4px 14px;
    border-radius: 18px;
    background-color: var(--color-bg);
    cursor: pointer;
    transition: background-color 0.3s;
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    align-content: center;

    &:hover {
        background-color: var(--color-bg-light);
    }
}

.multiselect-tags {
    display: flex;
    flex-wrap: wrap;
    gap: 4px;
    align-items: center;
}

.multiselect-tag {
    display: flex;
    align-items: center;
    gap: 2px;
    padding: 2px 8px;
    border-radius: 999px;
    background-color: var(--color-bg-light);
    color: var(--color-text-light);
    font-size: var(--font-size-sm);

    .icon {
        font-size: 14px;
        color: var(--color-text-dark);
        cursor: pointer;

        &:hover {
            color: var(--color-text-light);
        }
    }
}

.multiselect-placeholder {
    color: var(--color-text-dark);
    font-size: var(--font-size);
}
</style>

<style>
.multiselect-dropdown {
    position: fixed;
    background-color: var(--color-bg-light);
    border-radius: 10px;
    overflow: hidden;
    max-height: 0;
    transition: max-height 0.25s ease;
    z-index: 9999;

    &.active {
        max-height: 200px;
        overflow-y: auto;
    }
}

.multiselect-item {
    padding: 10px 16px;
    font-size: var(--font-size);
    color: var(--color-text);
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: space-between;

    &:hover {
        background-color: var(--color-bg);
        color: var(--color-text-light);
    }

    &.active {
        background-color: var(--color-primary-light);
        color: var(--color-text-light);

        .icon {
            color: var(--color-primary-dark);
        }
    }
}
</style>