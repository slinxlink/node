// composable/modal.ts
import { shallowRef } from 'vue'
import type { Component, ComponentPublicInstance } from 'vue'
import { bind, unbind, type KeyBindings } from '@/util/keyboard'

type Instance = ComponentPublicInstance | null

// 传入的组件必须是带构造签名的 SFC，onConfirm 才能拿到它的实例类型
type ModalComponent = new (...args: never[]) => object

interface ModalOptions<T extends ModalComponent> {
    title: string
    detail?: string
    width?: string
    component?: T
    props?: Record<string, unknown>
    // 返回 false 或抛错时弹窗保持打开，可以是异步的
    onConfirm?: (instance: InstanceType<T>) => void | boolean | Promise<void | boolean>
    maskClosable?: boolean
    form?: boolean // 表单类弹窗：高度会随内容变化，贴着顶部显示，不居中
    loose?: boolean // 内容不裁剪：提示这类要伸出弹窗边缘的东西才不会被切掉；内容多了就不会滚动，只给内容固定的弹窗用
}

const visible = ref(false)
const title = ref('')
const detail = ref('')
const width = ref('')
const component = shallowRef<Component | null>(null)
const props = ref<Record<string, unknown>>({})
const actions = ref(true)
const maskClosable = ref(true)
const form = ref(false)
const loose = ref(false)
const componentRef = shallowRef<Instance>(null)
const onConfirmFn = shallowRef<((instance: Instance) => void | boolean | Promise<void | boolean>) | null>(
    null,
)
const busy = ref(false)

let token: KeyBindings | null = null

function show<T extends ModalComponent>(options: ModalOptions<T>) {
    title.value = options.title
    detail.value = options.detail ?? ''
    width.value = options.width ?? ''
    component.value = options.component ?? null
    props.value = options.props ?? {}
    actions.value = options.onConfirm != null
    maskClosable.value = options.maskClosable ?? true
    form.value = options.form ?? false
    loose.value = options.loose ?? false
    onConfirmFn.value = options.onConfirm
        ? (instance) => options.onConfirm?.(instance as InstanceType<T>)
        : null
    componentRef.value = null
    visible.value = true

    token = bind({
        Escape: close,
        Enter: confirm,
    })
}

function close() {
    if (busy.value) return
    visible.value = false
    component.value = null
    onConfirmFn.value = null
    componentRef.value = null
    unbind(token)
    token = null
}

async function confirm() {
    if (busy.value) return

    const cb = onConfirmFn.value

    if (cb) {
        busy.value = true

        try {
            if ((await cb(componentRef.value)) === false) return
        } finally {
            busy.value = false
        }
    }

    close()
}

export const modalState = {
    visible,
    title,
    detail,
    width,
    component,
    props,
    actions,
    maskClosable,
    form,
    loose,
    componentRef,
    busy,
}

export const modal = {
    show,
    close,
    confirm,
}
