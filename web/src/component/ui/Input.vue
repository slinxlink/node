<template>
    <div class="input-wrap">
        <label v-if="label" class="input-label">{{ label }}</label>
        <div class="input-inner">
            <span v-if="prefix" ref="prefixRef" class="prefix">{{ prefix }}</span>
            <input
                class="input"
                :type="isNumber ? 'text' : type"
                :placeholder="placeholder"
                :value="modelValue"
                :readonly="readonly"
                :disabled="disabled"
                :autocomplete="autocomplete"
                :style="prefix ? { paddingLeft: prefixWidth } : {}"
                @input="onInput"
                @keydown="isNumber ? onKeydown($event) : null"
            />
            <div v-if="isNumber" class="spin-btns">
                <button class="spin-btn" @click="step(1)"><i class="icon">expand_less</i></button>
                <button class="spin-btn" @click="step(-1)"><i class="icon">expand_more</i></button>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
const props = defineProps<{
    modelValue?: string | number
    label?: string
    placeholder?: string
    type?: string
    min?: number
    max?: number
    readonly?: boolean
    disabled?: boolean
    autocomplete?: string
    prefix?: string
}>()
const emit = defineEmits(['update:modelValue'])

const isNumber = computed(() => props.type === 'number')
const prefixRef = ref<HTMLElement>()
const prefixWidth = ref('14px')

onMounted(async () => {
    await nextTick()
    if (prefixRef.value) {
        prefixWidth.value = `${prefixRef.value.offsetWidth + 14}px`
    }
})

watch(() => props.prefix, async () => {
    await nextTick()
    if (prefixRef.value) {
        prefixWidth.value = `${prefixRef.value.offsetWidth + 14}px`
    }
})
function onKeydown(e: KeyboardEvent) {
    const allowed = ['Backspace', 'Delete', 'ArrowLeft', 'ArrowRight', 'Tab']
    if (allowed.includes(e.key)) return
    if (!/^[0-9]$/.test(e.key)) e.preventDefault()
}

function onInput(e: Event) {
    const input = e.target as HTMLInputElement
    const raw = input.value
    if (isNumber.value) {
        if (raw === '') {
            emit('update:modelValue', 0)
            return
        }
        const num = Number(raw)
        if (isNaN(num)) {
            input.value = String(props.modelValue ?? 0)
            return
        }
        if (props.min !== undefined && num < props.min) {
            input.value = String(props.modelValue ?? props.min)
            return
        }
        if (props.max !== undefined && num > props.max) {
            input.value = String(props.modelValue ?? props.max)
            return
        }
        emit('update:modelValue', num)
    } else {
        emit('update:modelValue', raw)
    }
}

function step(dir: number) {
    let val = Number(props.modelValue ?? 0) + dir
    if (props.min !== undefined) val = Math.max(props.min, val)
    if (props.max !== undefined) val = Math.min(props.max, val)
    emit('update:modelValue', val)
}
</script>

<style scoped>
.input-wrap {
    display: flex;
    flex-direction: column;
    gap: 6px;
    width: 100%;
}

.input-label {
    font-size: var(--font-size-sm);
    color: var(--color-text-dark);
}

.input-inner {
    position: relative;
    display: flex;
    align-items: center;
}

.prefix {
    position: absolute;
    left: 14px;
    font-size: var(--font-size);
    color: var(--color-text-dark);
    white-space: nowrap;
    pointer-events: none;
    user-select: none;
}

.input {
    border: none;
    outline: none;
    text-align: left;
    font-size: var(--font-size);
    line-height: var(--font-size);
    color: var(--color-text-light);
    box-shadow: inset var(--box-shadow);
    transition: background-color 0.3s;
    background-color: var(--color-bg);
    border-radius: 999px;
    padding: 8px 14px;
    width: 100%;
    box-sizing: border-box;

    &:disabled {
        opacity: 0.5;
        cursor: not-allowed;
    }

    &:focus {
        background-color: var(--color-bg-light);
    }

    &:-webkit-autofill,
    &:-webkit-autofill:hover,
    &:-webkit-autofill:focus,
    &:-webkit-autofill:active {
        -webkit-box-shadow: 0 0 0 1000px var(--color-bg) inset;
        -webkit-text-fill-color: var(--color-text-light);
        caret-color: var(--color-text-light);
        transition: background-color 5000s ease-in-out 0s;
    }
}

.spin-btns {
    position: absolute;
    right: 0px;
    display: flex;
    flex-direction: column;
    gap: 0;
}

.spin-btn {
    background: none;
    border: none;
    cursor: pointer;
    padding: 0px 5px;
    color: var(--color-text-dark);
    display: flex;
    align-items: center;
    line-height: 1;

    &:hover {
        color: var(--color-text-light);
    }

    .icon {
        font-size: 16px;
        margin: -2px 0;
    }
}
</style>