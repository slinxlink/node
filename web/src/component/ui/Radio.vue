<template>
    <div class="radio" :class="{ disabled }" ref="radioRef">
        <div class="slider" :style="sliderStyle" />
        <button
            v-for="(option, index) in options"
            :key="option.value"
            :ref="el => btnRefs[index] = el as HTMLElement"
            :class="{ active: modelValue === option.value }"
            @click="!disabled && $emit('update:modelValue', option.value)"
        >
            {{ option.label }}
        </button>
    </div>
</template>

<script setup lang="ts">
const props = defineProps<{
    modelValue?: string
    options: Array<{ label: string, value: string }>
    disabled?: boolean
}>()
defineEmits(['update:modelValue'])

const btnRefs = ref<HTMLElement[]>([])
const btnWidths = ref<number[]>([])

onMounted(() => {
    const observer = new ResizeObserver(() => {
        btnWidths.value = btnRefs.value.map(el => el?.offsetWidth ?? 0)
    })
    watch(btnRefs, (els) => {
        els.forEach(el => el && observer.observe(el))
    }, { immediate: true })
    onUnmounted(() => observer.disconnect())
})

const sliderStyle = computed(() => {
    const index = props.options.findIndex(o => o.value === props.modelValue)
    if (index === -1 || !btnRefs.value[index] || !btnWidths.value[index]) return { opacity: 0 }
    const btn = btnRefs.value[index]
    const padding = 5
    return {
        opacity: 1,
        width: btn.offsetWidth + 'px',
        transform: `translateX(${btn.offsetLeft - padding}px)`,
    }
})
</script>

<style scoped>
.radio {
    display: flex;
    gap: 1px;
    background-color: var(--color-bg);
    border-radius: 999px;
    padding: 5px;
    position: relative;

    &.disabled {
        opacity: 0.5;
        cursor: not-allowed;

        button {
            cursor: not-allowed;
        }
    }

    .slider {
        position: absolute;
        top: 5px;
        left: 5px;
        height: calc(100% - 10px);
        background-color: var(--color-primary-light);
        border-radius: 999px;
        transition: transform 0.3s ease, width 0.3s ease;
        pointer-events: none;
    }

    button {
        position: relative;
        z-index: 1;
        padding: 5px 15px;
        border: none;
        border-radius: 999px;
        font-size: var(--font-size-sm);
        color: var(--color-text-dark);
        background: none;
        cursor: pointer;
        transition: color 0.3s;

        &:hover {
            color: var(--color-text-light);
        }

        &.active {
            color: var(--color-text-light);
        }
    }
}
</style>