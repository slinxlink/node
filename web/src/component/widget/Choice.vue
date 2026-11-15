<template>
    <div ref="root" class="choice" :class="{ disabled }">
        <span :class="{ ready }" :style="slider" />
        <button
            v-for="item in items"
            :key="String(item.value)"
            type="button"
            :class="{ active: item.value === value }"
            :disabled="disabled"
            @click="emit('update:value', item.value)"
        >
            {{ item.label }}
        </button>
    </div>
</template>

<script setup lang="ts" generic="T">
interface ChoiceItem<T> {
    label: string
    value: T
}

const props = defineProps<{
    items: ChoiceItem<T>[]
    value: T | null
    disabled?: boolean
}>()

const emit = defineEmits<{
    'update:value': [value: T]
}>()

const root = ref<HTMLElement>()
const slider = ref({ width: '0px', transform: 'translateX(0)', opacity: 0 })

// 第一次定位不要动画，免得滑块从左边滑过来
const ready = ref(false)

// 滑块盖在选中的按钮上；按钮的位置和宽度随字体、内容变化，所以每次都重新量
function measure() {
    const button = root.value?.querySelector<HTMLElement>('button.active')

    slider.value = button
        ? { width: `${button.offsetWidth}px`, transform: `translateX(${button.offsetLeft}px)`, opacity: 1 }
        : { ...slider.value, opacity: 0 }
}

// 根元素宽度由按钮撑开，按钮变宽变窄它就跟着变
const observer = new ResizeObserver(measure)

onMounted(() => {
    measure()
    observer.observe(root.value!)
    requestAnimationFrame(() => (ready.value = true))
})

onBeforeUnmount(() => observer.disconnect())

watch(() => props.value, measure, { flush: 'post' })
</script>

<style scoped>
.choice {
    position: relative;
    display: flex;
    width: fit-content;
    height: 34px;
    padding: 5px;
    border-radius: 999px;
    background-color: var(--background-200);

    &.disabled {
        opacity: 0.5;
    }

    > span {
        position: absolute;
        top: 5px;
        left: 0;
        height: calc(100% - 10px);
        border-radius: 999px;
        background-color: var(--primary-300);
        pointer-events: none;

        &.ready {
            transition:
                transform 0.3s ease-in-out,
                width 0.3s ease-in-out;
        }
    }

    button {
        position: relative;
        z-index: 1;
        padding: 0 15px;
        border: none;
        border-radius: 999px;
        outline: none;
        background: none;
        color: var(--text-300);
        font-size: 14px;
        line-height: 1;
        white-space: nowrap;
        cursor: pointer;
        transition: color 0.3s ease-in-out;

        &:hover,
        &:focus-visible,
        &.active {
            color: var(--text-100);
        }

        &.active {
            cursor: default;
        }

        &:disabled {
            cursor: not-allowed;
        }
    }
}
</style>
