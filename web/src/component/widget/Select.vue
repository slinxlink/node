<template>
    <div class="select" :class="{ disabled, small, custom: $slots.item }">
        <Button ref="trigger" class="trigger" :disabled="disabled" @click="toggle">
            <span class="display" :class="{ placeholder: !selected }">
                <slot v-if="selected && $slots.item" name="item" :item="selected" />
                <template v-else>{{ selectedLabel }}</template>
            </span>
            <Icon name="unfold_more" />
        </Button>

        <Teleport to="body">
            <Transition :css="false" @enter="onEnter" @leave="onLeave">
                <div v-if="open" ref="dropdown" class="dropdown" :class="{ flip, small, custom: $slots.item }" :style="dropdownStyle">
                    <div ref="scroll" class="scroll">
                        <Button v-for="(item, index) in items" :key="index" :ref="(el) => setItemRef(el, index)" :selected="resolveValue(item) === value" :disabled="item.disabled" @click="select(item)">
                            <span class="display">
                                <slot name="item" :item="item">{{ item.label }}</slot>
                            </span>
                        </Button>
                    </div>
                </div>
            </Transition>
        </Teleport>
    </div>
</template>

<script setup lang="ts" generic="T">
import type { ComponentPublicInstance } from 'vue'
import Button from '@/component/widget/Button.vue'
import Icon from '@/component/widget/Icon.vue'
import { initial } from '@/util/initial'
import { bind, unbind, type KeyBindings } from '@/util/keyboard'

interface SelectItem<T> {
    label: string
    value?: T
    disabled?: boolean
}

const props = withDefaults(defineProps<{
    items: SelectItem<T>[]
    value: T | null
    placeholder?: string
    disabled?: boolean
    small?: boolean // 小号：触发按钮和展开的选项都是 26px 高、12px 字，放在标题行之类窄的地方
}>(), {
    placeholder: '',
})

const emit = defineEmits<{
    'update:value': [value: T]
}>()

function resolveValue(item: SelectItem<T>): T {
    return item.value ?? (item.label as unknown as T)
}

const selected = computed(() => props.items.find((i) => resolveValue(i) === props.value))

const selectedLabel = computed(() => selected.value?.label ?? props.placeholder)

const open = ref(false)
const flip = ref(false)
const trigger = ref<InstanceType<typeof Button> | null>(null)
const dropdown = ref<HTMLElement | null>(null)
const scroll = ref<HTMLElement | null>(null)
const dropdownStyle = ref<Record<string, string>>({})

const DROPDOWN_HEIGHT = 300
const GAP = 5

let itemEls: (HTMLElement | null)[] = []
let keyToken: KeyBindings | null = null

function setItemRef(el: Element | ComponentPublicInstance | null, index: number) {
    itemEls[index] = (el && '$el' in el ? el.$el : el) as HTMLElement | null
}

// 展开和收起只做高度变化；上下内边距一起动，不然收到最小也会剩 10px
function grow(box: HTMLElement, from: number, to: number, done: () => void) {
    const height = box.offsetHeight
    const heights = [`${height * from}px`, `${height * to}px`]
    const paddings = [`${5 * from}px`, `${5 * to}px`]

    box.style.overflow = 'hidden'

    const animation = box.animate(
        [
            { height: heights[0], paddingTop: paddings[0], paddingBottom: paddings[0] },
            { height: heights[1], paddingTop: paddings[1], paddingBottom: paddings[1] },
        ],
        { duration: 300, easing: 'ease' },
    )

    animation.onfinish = () => {
        box.style.overflow = ''
        done()
    }
}

function onEnter(el: Element, done: () => void) {
    grow(el as HTMLElement, 0, 1, done)
}

function onLeave(el: Element, done: () => void) {
    grow(el as HTMLElement, 1, 0, done)
}

function toggle() {
    if (props.disabled) return
    if (open.value) close()
    else show()
}

function show() {
    const el = trigger.value?.$el as HTMLElement | undefined
    if (!el) return

    const rect = el.getBoundingClientRect()
    const below = window.innerHeight - rect.bottom - GAP
    const above = rect.top - GAP
    flip.value = below < DROPDOWN_HEIGHT + GAP && above > below

    dropdownStyle.value = {
        left: `${rect.left}px`,
        width: `${rect.width}px`,
        ...(flip.value
            ? { bottom: `${window.innerHeight - rect.top + GAP}px` }
            : { top: `${rect.bottom + GAP}px` }),
    }

    itemEls = []
    open.value = true

    nextTick(() => {
        reveal()
        keyToken = bind({
            Escape: () => close(),
            ArrowUp: () => move(-1),
            ArrowDown: () => move(1),
            Enter: () => close(),
            alnum: onChar,
        })
        window.addEventListener('mousedown', onClickOutside, true)
        window.addEventListener('resize', close)
        window.addEventListener('scroll', onScroll, true)
    })
}

function close() {
    if (!open.value) return
    open.value = false
    unbind(keyToken)
    keyToken = null
    window.removeEventListener('mousedown', onClickOutside, true)
    window.removeEventListener('resize', close)
    window.removeEventListener('scroll', onScroll, true)
}

function onClickOutside(e: MouseEvent) {
    const target = e.target as Node
    const triggerEl = trigger.value?.$el as HTMLElement | undefined
    if (triggerEl?.contains(target) || dropdown.value?.contains(target)) return
    close()
}

function onScroll(e: Event) {
    if (dropdown.value?.contains(e.target as Node)) return
    close()
}

function move(step: number) {
    const currentIndex = props.items.findIndex((i) => resolveValue(i) === props.value)
    const next = (currentIndex === -1 ? 0 : currentIndex) + step
    if (next < 0 || next >= props.items.length) return
    emit('update:value', resolveValue(props.items[next]!))
    nextTick(reveal)
}

// 按字母或数字键：从当前项往下找下一个首字相同的，找到底回到第一个
function onChar(e: KeyboardEvent) {
    const letter = e.key.toUpperCase()
    const total = props.items.length
    const current = props.items.findIndex((i) => resolveValue(i) === props.value)

    for (let step = 1; step <= total; step++) {
        const index = (current + step + total) % total
        const item = props.items[index]!

        if (item.disabled || initial(item.label) !== letter) continue

        if (index !== current) {
            emit('update:value', resolveValue(item))
            nextTick(reveal)
        }

        return
    }
}

function reveal() {
    const container = scroll.value
    const index = props.items.findIndex((i) => resolveValue(i) === props.value)
    const target = itemEls[index]
    if (!container || !target) return
    const containerRect = container.getBoundingClientRect()
    const targetRect = target.getBoundingClientRect()
    const delta = (targetRect.top - containerRect.top) - (containerRect.height / 2 - targetRect.height / 2)
    const from = container.scrollTop
    const to = Math.max(0, Math.min(from + delta, container.scrollHeight - container.clientHeight))
    animateScroll(container, from, to, 300)
}

function animateScroll(el: HTMLElement, from: number, to: number, duration: number) {
    if (from === to) return
    const start = performance.now()

    function step(now: number) {
        const elapsed = now - start
        const progress = Math.min(elapsed / duration, 1)
        const eased = 1 - Math.pow(1 - progress, 3)
        el.scrollTop = from + (to - from) * eased
        if (progress < 1) {
            requestAnimationFrame(step)
        }
    }

    requestAnimationFrame(step)
}

function select(item: SelectItem<T>) {
    emit('update:value', resolveValue(item))
    close()
}

onBeforeUnmount(close)
</script>

<style scoped>
.select {
    position: relative;
    width: 100%;

    .trigger {
        justify-content: space-between;
        background-color: var(--background-200);
        padding: 10px 7.5px 10px 15px;
        height: 34px;

        &:hover,
        &:focus-visible {
            background-color: var(--background-100);
        }
        &:active {
            background-color: var(--primary-600);
        }
    }

    &.small .trigger {
        height: 26px;
        padding: 0 5px 0 12px;
        font-size: 12px;
    }

    .icon {
        color: var(--text-300);
    }

    .display {
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;

        &.placeholder {
            color: var(--text-300);
        }
    }

    /* 有自定义内容(如旗帜图)时，图在左、文字在右；图占 24px，左边留 5px */
    &.custom {
        .trigger {
            padding-left: 5px;
        }

        .display {
            display: flex;
            align-items: center;
            gap: 10px;

            &.placeholder {
                padding-left: 10px;
            }

            &:deep(img) {
                flex-shrink: 0;
                width: 24px;
                height: 24px;
            }
        }
    }
}

.dropdown {
    z-index: 98; /* 和弹窗同层,下拉框挂在 body 最后,在弹窗里打开时盖在弹窗上面 */
    position: fixed;
    max-height: 300px;
    padding: 5px;
    background-color: var(--background-100);
    border-radius: 20px;
    box-shadow: var(--shadow);

    .scroll {
        max-height: 290px;
        overflow-y: auto;
        display: flex;
        flex-direction: column;
        gap: 2px;
        border-radius: 15px;

        button {
            height: 34px;
            justify-content: start;
            background-color: var(--background-100);

            &:hover,
            &:focus-visible {
                background-color: var(--primary-500);
            }
            &:active {
                background-color: var(--primary-600);
            }
            &.selected {
                background-color: var(--primary-400);
            }
        }

        .display {
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
        }
    }

    &.small .scroll button {
        height: 26px;
        padding: 0 12px;
        font-size: 12px;
    }

    &.custom .scroll button {
        padding: 5px 10px 5px 5px;

        .display {
            display: flex;
            align-items: center;
            gap: 10px;

            &:deep(img) {
                flex-shrink: 0;
                width: 24px;
                height: 24px;
            }
        }
    }
}
</style>