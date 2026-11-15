<template>
    <div class="input">
        <div class="wrap">
            <span v-if="hasLeading" class="leading">
                <slot name="leading" />
            </span>
            <span v-if="prefix" ref="prefixEl" class="prefix">{{ prefix }}</span>
            <input
                v-bind="$attrs"
                :type="inputType"
                :placeholder="placeholder"
                :value="value"
                :autocomplete="autocomplete"
                :maxlength="maxlength"
                :readonly="readonly"
                :disabled="disabled"
                :inputmode="isNumber ? 'decimal' : undefined"
                :class="{ 'has-trailing': hasTrailing, 'has-leading': hasLeading, 'has-prefix': !!prefix, 'has-spin': isNumber, center }"
                :style="prefix ? { paddingLeft: `${15 + prefixWidth}px` } : undefined"
                @input="onInput"
                @keydown="onKeydown"
                @blur="onBlur"
            />
            <span v-if="hasTrailing" class="trailing">
                <slot name="trailing">
                    <Button
                        v-if="type === 'password'"
                        variant="text"
                        class="eyes"
                        :tabindex="-1"
                        @click="showPassword = !showPassword"
                    >
                        <Icon name="visibility" class="eye" :class="{ show: !showPassword }" />
                        <Icon name="visibility_off" class="eye" :class="{ show: showPassword }" />
                    </Button>
                    <div v-else-if="isNumber" class="spin">
                        <Button variant="text" :tabindex="-1" :disabled="disabled || readonly" @click="spin(-1)">
                            <Icon name="remove" />
                        </Button>
                        <Button variant="text" :tabindex="-1" :disabled="disabled || readonly" @click="spin(1)">
                            <Icon name="add" />
                        </Button>
                    </div>
                    <Button
                        v-else-if="clearable && value && !disabled && !readonly"
                        variant="text"
                        :tabindex="-1"
                        @click="clear"
                    >
                        <Icon name="cancel" />
                    </Button>
                </slot>
            </span>
        </div>
        <Transition name="error">
            <span v-if="error" class="error" :class="{ center }">{{ error }}</span>
        </Transition>
    </div>
</template>

<script setup lang="ts">
import Button from '@/component/widget/Button.vue'
import Icon from '@/component/widget/Icon.vue'

defineOptions({ inheritAttrs: false })

const props = defineProps<{
    value: string
    type?: string
    placeholder?: string
    autocomplete?: string
    maxlength?: number
    disabled?: boolean
    readonly?: boolean
    validate?: (value: string) => string
    center?: boolean
    clearable?: boolean
    min?: number // type 为 number 时的下限
    max?: number // type 为 number 时的上限
    step?: number // type 为 number 时每次增减的量，默认 1
    prefix?: string // 输入框里的文字前缀，字号和输入的内容一样，宽度按文字，输入的内容紧接在它后面，如路径前面的目录
}>()

const emit = defineEmits<{
    'update:value': [value: string]
}>()

const slots = useSlots()

const error = ref('')
const showPassword = ref(false)

const prefixEl = ref<HTMLElement>()
const prefixWidth = ref(0)

// 前缀的宽度决定输入内容从哪里开始：挂载后、前缀变了都重新量
async function measure() {
    await nextTick()
    prefixWidth.value = prefixEl.value?.offsetWidth ?? 0
}

onMounted(measure)
watch(() => props.prefix, measure)

const isNumber = computed(() => props.type === 'number')

// number 用文本框自己限制输入，浏览器原生的数字框会带难看的上下箭头
const inputType = computed(() => {
    if (isNumber.value) return 'text'
    if (props.type !== 'password') return props.type
    return showPassword.value ? 'text' : 'password'
})

const hasTrailing = computed(() => {
    return !!slots.trailing || props.type === 'password' || isNumber.value || (props.clearable && !!props.value && !props.disabled && !props.readonly)
})

const hasLeading = computed(() => {
    return !!slots.leading
})

function decimals(value: number | string) {
    return (String(value).split('.')[1] ?? '').length
}

function onKeydown(e: KeyboardEvent) {
    if (!isNumber.value || e.metaKey || e.ctrlKey) return

    if (e.key === 'ArrowUp') return spin(1, e)
    if (e.key === 'ArrowDown') return spin(-1, e)

    // Backspace、方向键、Tab、Enter 这类功能键都放行
    if (e.key.length > 1) return

    const minus = e.key === '-' && props.min !== undefined && props.min < 0
    const dot = e.key === '.' && !props.value.includes('.')

    if (!/^[0-9]$/.test(e.key) && !minus && !dot) e.preventDefault()
}

function onInput(e: Event) {
    const input = e.target as HTMLInputElement
    const val = input.value.trim()

    if (isNumber.value && val !== '') {
        const num = Number(val)
        const outside =
            (props.min !== undefined && num < props.min) ||
            (props.max !== undefined && num > props.max)

        // 不是数字或超出范围，退回上一个值
        if (!/^-?\d*\.?\d*$/.test(val) || outside) {
            input.value = props.value
            return
        }
    }

    emit('update:value', val)
    if (error.value) check()
}

function spin(direction: number, e?: Event) {
    e?.preventDefault()

    const size = props.step ?? 1
    let num = Number(props.value) + direction * size

    if (Number.isNaN(num)) num = direction * size
    if (props.min !== undefined) num = Math.max(props.min, num)
    if (props.max !== undefined) num = Math.min(props.max, num)

    emit('update:value', num.toFixed(Math.max(decimals(size), decimals(props.value))))
    if (error.value) check()
}

function onBlur() {
    check()
}

function clear() {
    emit('update:value', '')
    if (error.value) check()
}

function check(): boolean {
    if (!props.validate) {
        error.value = ''
        return true
    }

    error.value = props.validate(props.value.trim())
    return !error.value
}

defineExpose({ check })

const formRegister = inject<((input: { check: () => boolean }) => void) | null>('form-register', null)

onMounted(() => {
    formRegister?.({ check })
})
</script>

<style scoped>
.input {
    display: flex;
    flex-direction: column;
    gap: 5px;
    width: 100%;

    .wrap {
        position: relative;

        .leading {
            position: absolute;
            top: 50%;
            left: 10px;
            transform: translateY(-50%);
            display: flex;
            align-items: center;
            justify-content: center;
            font-size: 20px;
        }

        /* 文字前缀：和输入的内容同样大小，颜色淡一档，叠在输入框左边，不挡点击 */
        .prefix {
            position: absolute;
            top: 50%;
            left: 15px;
            transform: translateY(-50%);
            color: var(--text-300);
            font-size: 14px;
            line-height: 1;
            white-space: nowrap;
            pointer-events: none;
            user-select: none;
        }

        input {
            display: flex;
            width: 100%;
            height: 34px;
            padding: 10px 15px;
            border: none;
            border-radius: 999px;
            outline: none;
            background-color: var(--background-200);
            color: var(--text-100);
            caret-color: var(--primary-300);
            font-size: 14px;
            line-height: 1;
            text-align: left;
            transition:
                background-color 0.3s,
                padding 0.3s;

            &.center {
                text-align: center;
            }

            &.has-trailing {
                padding-right: 35px;
            }

            &.has-leading {
                padding-left: 35px;
            }

            /* 左内边距是量出来的，不要过渡，否则挂载时会从 15px 滑到位 */
            &.has-prefix {
                transition: background-color 0.3s;
            }

            &.has-spin {
                padding-right: 70px;
            }

            &.center.has-trailing {
                padding-left: 35px;
                padding-right: 35px;
            }

            &:hover {
                background-color: var(--background-100);
            }

            &:focus {
                background-color: var(--background-100);
            }

            &:disabled {
                opacity: 0.5;
                cursor: not-allowed;
                filter: grayscale(50%);
            }
        }

        .trailing {
            position: absolute;
            top: 50%;
            right: 10px;
            transform: translateY(-50%);
            display: flex;
            align-items: center;
            justify-content: center;
            font-size: 20px;

            .spin {
                display: flex;
                gap: 5px;
            }

            .eyes {
                position: relative;
                width: 20px;
                height: 20px;

                .eye {
                    position: absolute;
                    top: 0;
                    left: 0;
                    opacity: 0;
                    transform: scale(0.8);
                    transition:
                        opacity 0.3s,
                        transform 0.3s;

                    &.show {
                        opacity: 1;
                        transform: scale(1);
                    }
                }
            }
        }
    }

    .error {
        padding: 0 15px;
        overflow: hidden;
        color: var(--red);
        font-size: 12px;
        text-overflow: ellipsis;
        white-space: nowrap;

        &.center {
            text-align: center;
        }
    }
}

.error-enter-active,
.error-leave-active {
    transition: all 0.3s;
}

.error-enter-from,
.error-leave-to {
    opacity: 0;
    height: 0;
}

.error-enter-to,
.error-leave-from {
    opacity: 1;
    height: 1.2em;
}
</style>
