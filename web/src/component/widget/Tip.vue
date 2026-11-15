<template>
    <span class="tip" @mouseenter="onHoverEnter" @mouseleave="onHoverLeave" @click="onClick">
        <slot>
            <Button variant="text" class="tip-btn" aria-label="帮助">
                <Icon name="help" />
            </Button>
        </slot>
        <Transition name="tip">
            <div v-if="visible && msg" class="anchor" :class="side">
                <div class="triangle" :class="side"></div>
                <div class="bubble" :class="align"><slot name="content"><BBCode :text="msg" /></slot></div>
            </div>
        </Transition>
    </span>
</template>

<script setup lang="ts">
import Button from '@/component/widget/Button.vue'
import Icon from '@/component/widget/Icon.vue'
import { BBCode } from '@/util/bbcode'

const props = withDefaults(
    defineProps<{
        msg: string // 支持 BBCode,换行用 [br];要排版得更复杂就用 content 插槽(msg 仍要给,为空不弹出)
        side?: 'top' | 'bottom'
        align?: 'left' | 'center' | 'right'
        lock?: boolean // 点击锁定，再点或点外面收起；默认问号按钮自动开启
    }>(),
    {
        side: 'top',
        align: 'center',
    },
)

const visible = ref(false)
const locked = ref(false)
let hideTimer: ReturnType<typeof setTimeout> | null = null

function onHoverEnter() {
    if (locked.value) return
    if (hideTimer) clearTimeout(hideTimer)
    visible.value = true
}

function onHoverLeave() {
    if (locked.value) return
    if (hideTimer) clearTimeout(hideTimer)
    hideTimer = setTimeout(() => {
        visible.value = false
    }, 50)
}

// 自定义内容要传 lock 才支持，免得影响里面的链接
const slots = useSlots()
const canLock = computed(() => props.lock || !slots.default)

function onClick() {
    if (canLock.value) toggle()
}

function toggle() {
    if (locked.value) {
        locked.value = false
        visible.value = false
    } else {
        locked.value = true
        visible.value = true
    }
}

function onClickOutside(e: MouseEvent) {
    if (!locked.value) return
    const el = e.target as HTMLElement
    if (!el.closest('.tip')) {
        locked.value = false
        visible.value = false
    }
}

onMounted(() => {
    document.addEventListener('click', onClickOutside)
})

onUnmounted(() => {
    document.removeEventListener('click', onClickOutside)
})
</script>

<style scoped>
.tip {
    position: relative;

    .anchor {
        position: absolute;
        left: 50%;
        width: 0;
        z-index: 89;
        pointer-events: none;

        &.top {
            bottom: calc(100% + 5px);

            .bubble {
                bottom: calc(100% + 5px);
            }
        }

        &.bottom {
            top: calc(100% + 5px);

            .bubble {
                top: calc(100% + 5px);
            }
        }

        .triangle {
            position: absolute;
            left: 0;
            transform: translateX(-50%);
            width: 0;
            height: 0;
            border-left: 5px solid transparent;
            border-right: 5px solid transparent;

            &.top {
                bottom: 0;
                border-top: 5px solid var(--background-100);
            }

            &.bottom {
                top: 0;
                border-bottom: 5px solid var(--background-100);
            }
        }

        .bubble {
            position: absolute;
            left: 0;
            padding: 6px 10px;
            border-radius: 5px;
            background-color: var(--background-100);
            box-shadow: var(--shadow);
            color: var(--text-100);
            font-size: 12px;
            text-align: left;
            white-space: pre;

            &.center {
                transform: translateX(-50%);
            }

            &.left {
                transform: translateX(-10px);
            }

            &.right {
                transform: translateX(calc(-100% + 10px));
            }
        }
    }

    .tip-enter-active,
    .tip-leave-active {
        transition: opacity 0.3s;
    }

    .tip-enter-from,
    .tip-leave-to {
        opacity: 0;
    }
}
</style>
