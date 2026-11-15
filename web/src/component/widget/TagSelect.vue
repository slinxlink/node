<!-- component/widget/TagSelect.vue -->
<template>
    <div class="tag-select" :class="{ disabled }">
        <Button ref="trigger" class="trigger" :class="{ trailing: $slots.trailing }" :disabled="disabled" @click="toggle">
            <span class="tags" :class="{ placeholder: !value.length }">
                <template v-if="value.length">
                    <span v-for="item in chosen" :key="String(item.value)" class="tag">{{ item.label }}</span>
                </template>
                <template v-else>{{ placeholder }}</template>
            </span>
            <Icon v-if="!$slots.trailing" name="unfold_more" />
        </Button>
        <!-- 右边的操作(比如刷新)放在里面，替掉展开图标；它是触发条的兄弟，点它不会展开列表 -->
        <span v-if="$slots.trailing" class="action">
            <slot name="trailing" />
        </span>

        <Teleport to="body">
            <Transition :css="false" @enter="onEnter" @leave="onLeave">
                <div v-if="open" ref="dropdown" class="dropdown" :class="{ flip }" :style="dropdownStyle">
                    <div class="scroll">
                        <Button
                            v-for="item in items"
                            :key="String(item.value)"
                            :class="{ on: has(item) }"
                            :disabled="item.disabled"
                            @click="toggleItem(item)"
                        >
                            <span class="label">{{ item.label }}</span>
                            <span class="check">
                                <Icon v-if="has(item)" name="check" />
                            </span>
                        </Button>
                    </div>
                </div>
            </Transition>
        </Teleport>
    </div>
</template>

<script setup lang="ts" generic="T">
import Button from '@/component/widget/Button.vue'
import Icon from '@/component/widget/Icon.vue'
import { bind, unbind, type KeyBindings } from '@/util/keyboard'

interface TagItem<T> {
    label: string
    value: T
    disabled?: boolean
}

// 紧凑的多选：平时像输入框，已选的项在里面排成小标签，一行放不下就自动换行；点开是勾选列表
const props = withDefaults(
    defineProps<{
        items: TagItem<T>[]
        value: T[]
        disabled?: boolean
        placeholder?: string
    }>(),
    { placeholder: '' },
)

const emit = defineEmits<{
    'update:value': [value: T[]]
}>()

const open = ref(false)
const flip = ref(false)
const trigger = ref<InstanceType<typeof Button> | null>(null)
const dropdown = ref<HTMLElement | null>(null)
const dropdownStyle = ref<Record<string, string>>({})

const DROPDOWN_HEIGHT = 300
const GAP = 5

let keyToken: KeyBindings | null = null

const has = (item: TagItem<T>) => props.value.includes(item.value)

// 标签按选择的先后排；选项里找不到的(比如旧数据)直接显示值本身
const chosen = computed(() =>
    props.value.map((v) => props.items.find((item) => item.value === v) ?? { label: String(v), value: v }),
)

function toggleItem(item: TagItem<T>) {
    emit(
        'update:value',
        has(item) ? props.value.filter((v) => v !== item.value) : [...props.value, item.value],
    )
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

    open.value = true

    nextTick(() => {
        keyToken = bind({
            Escape: () => close(),
            Enter: () => close(),
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

onBeforeUnmount(close)
</script>

<style scoped>
.tag-select {
    position: relative;
    width: 100%;

    .trigger {
        align-items: flex-start;
        justify-content: space-between;
        gap: 10px;
        height: auto;
        min-height: 34px;
        padding: 7px 7.5px 7px 7px;
        border-radius: 17px; /* 多行时仍是圆角矩形，不会变成大椭圆 */
        background-color: var(--background-200);

        &.trailing {
            padding-right: 35px;
        }

        &:hover,
        &:focus-visible {
            background-color: var(--background-100);
        }

        &:active {
            background-color: var(--background-100);
        }
    }

    /* 已选的标签排成一行，放不下就换行，输入框跟着变高 */
    .tags {
        display: flex;
        flex: 1;
        flex-wrap: wrap;
        align-items: center;
        gap: 5px;
        min-width: 0;
        min-height: 20px;
        line-height: 20px;
        text-align: left;

        &.placeholder {
            padding-left: 8px;
            color: var(--text-300);
        }

        .tag {
            flex-shrink: 0;
            padding: 4px 10px;
            border-radius: 999px;
            background-color: var(--background-100);
            color: var(--text-100);
            font-size: 12px;
            line-height: 1;
        }
    }

    .trigger:hover .tag,
    .trigger:focus-visible .tag {
        background-color: var(--background-200);
    }

    .icon {
        flex-shrink: 0;
        color: var(--text-300);
    }

    /* 固定在第一行的右边；多行时也不跟着居中 */
    .action {
        position: absolute;
        top: 7px;
        right: 10px; /* 和 Input、Textarea 的 trailing 同位置 */
        display: flex;
        align-items: center;
        justify-content: center;
        height: 20px;
        font-size: 20px;
    }

    &.disabled .tag {
        opacity: 0.5;
    }
}

.dropdown {
    z-index: 98; /* 和弹窗同层，下拉框挂在 body 最后，在弹窗里打开时盖在弹窗上面 */
    position: fixed;
    max-height: 300px;
    padding: 5px;
    background-color: var(--background-100);
    border-radius: 20px;
    box-shadow: var(--shadow);

    .scroll {
        display: flex;
        flex-direction: column;
        gap: 2px;
        max-height: 290px;
        overflow-y: auto;
        border-radius: 15px;

        button {
            justify-content: flex-start;
            gap: 5px;
            height: 34px;
            padding-right: 10px; /* 右边有勾，比默认的 15px 少一点 */
            background-color: var(--background-100);

            &:hover,
            &:focus-visible {
                background-color: var(--primary-500);
            }

            &:active {
                background-color: var(--primary-600);
            }

            &.on {
                background-color: var(--primary-400);
            }

            .check {
                display: flex;
                flex-shrink: 0;
                width: 20px;
            }

            .label {
                flex: 1;
                min-width: 0;
                overflow: hidden;
                text-align: left;
                text-overflow: ellipsis;
                white-space: nowrap;
            }
        }
    }
}
</style>
