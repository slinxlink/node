type ToastType = 'success' | 'warn' | 'error' | 'loading'

export type ToastDetail = { level: 'warn' | 'error'; message: string }

const state = reactive({
    visible: false,
    type: 'success' as ToastType,
    message: '',
    detail: [] as ToastDetail[],
})

let timer: ReturnType<typeof setTimeout> | null = null

function show(type: ToastType, message: string, duration: number | null, detail: ToastDetail[] = []) {
    if (timer) clearTimeout(timer)

    state.type = type
    state.message = message
    state.detail = detail
    state.visible = true

    if (duration) {
        timer = setTimeout(() => {
            state.visible = false
        }, duration)
    }
}

// 带详情的要多留点时间读
function keep(detail?: ToastDetail[]) {
    if (detail?.length) return 5000
    return 3500
}

export { state as toastState }

export const toast = {
    success: (message: string) => show('success', message, 2000),
    warning: (message: string, detail?: ToastDetail[]) => show('warn', message, keep(detail), detail),
    failure: (message: string, detail?: ToastDetail[]) => show('error', message, keep(detail), detail),
    loading: (message: string) => show('loading', message, null),
    hide: () => {
        if (timer) clearTimeout(timer)
        state.visible = false
    },
}
