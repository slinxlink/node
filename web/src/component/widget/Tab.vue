<template>
    <div ref="root" class="tab" role="tablist">
        <Button
            v-for="item in items"
            :key="item.value"
            :ref="(el) => setButton(item.value, el)"
            variant="text"
            role="tab"
            :aria-selected="item.value === value"
            :selected="item.value === value"
            @click="emit('update:value', item.value)"
        >
            {{ item.label }}
        </Button>
        <span class="indicator" :class="{ ready }" :style="{ transform: `translateX(${x}px)` }" />
    </div>
</template>

<script setup lang="ts">
import type { ComponentPublicInstance } from 'vue'
import Button from '@/component/widget/Button.vue'

export interface TabItem {
    label: string
    value: string
}

const props = defineProps<{
    items: [TabItem, ...TabItem[]] // 至少一项
    value: string
}>()

const emit = defineEmits<{
    'update:value': [value: string]
}>()

const INDICATOR_WIDTH = 20

const root = ref<HTMLElement | null>(null)
const buttons = new Map<string, HTMLElement>()
const x = ref(0)
const ready = ref(false) // 第一次定位不要动画,不然会从左边滑过来

function setButton(value: string, el: Element | ComponentPublicInstance | null) {
    const element = el && '$el' in el ? el.$el : el

    if (element instanceof HTMLElement) buttons.set(value, element)
    else buttons.delete(value)
}

// 线移到选中按钮的中心;按钮在外层里的位置是相对外层内容算的,外层横向滚动时线跟着走
function place() {
    const button = buttons.get(props.value)

    if (button) {
        x.value = button.offsetLeft + button.offsetWidth / 2 - INDICATOR_WIDTH / 2
    }
}

let observer: ResizeObserver | null = null

onMounted(() => {
    place()
    requestAnimationFrame(() => (ready.value = true))

    // 字体加载完、窗口变化会让按钮宽度和位置变,重新定位
    observer = new ResizeObserver(place)
    observer.observe(root.value as HTMLElement)
    buttons.forEach((button) => observer?.observe(button))
})

onBeforeUnmount(() => observer?.disconnect())

watch(
    () => [props.value, props.items],
    () => nextTick(place),
)
</script>

<style scoped>
.tab {
    position: relative;
    display: flex;
    gap: 30px;
    min-width: 0;
    padding: 20px 0; /* 和 Section 默认标题一样上下各 20px,文字居中;留白放在外层,选中线画在里面不会被 overflow 裁掉 */
    overflow-x: auto;

    /* 按钮是 text 变体,这里只改选中和未选中的文字颜色,其余(悬停、按下、禁用)跟着 Button 走 */
    button.text {
        flex-shrink: 0;
        white-space: nowrap;
        line-height: 20px; /* 标签高度要正好 20px，加上外层上下各 20px 才是标题行的 60px；Button 默认行高是 1，不写这个文字会偏下 */

        &:not(.selected) {
            color: var(--text-300);

            &:hover,
            &:focus-visible {
                color: var(--primary-300);
            }
        }

        &.selected {
            background-color: transparent; /* Button 的选中态有底色,tab 只要文字变亮 */
            color: var(--text-100);
        }
    }

    /* 选中项下面一小段固定长度的粗线,离外层底边 2px */
    .indicator {
        position: absolute;
        left: 0;
        bottom: 2px;
        width: 20px;
        height: 4px;
        border-radius: 2px;
        background-color: var(--primary-300);

        &.ready {
            transition: transform 0.3s;
        }
    }
}
</style>
