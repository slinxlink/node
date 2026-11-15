import type { Motion } from '@/component/layout/Header.vue'

type Theme = 'light' | 'dark' | 'system'

const KEY = 'theme'

function readSaved(): Theme | null {
    try {
        const value = localStorage.getItem(KEY)
        return value === 'light' || value === 'dark' || value === 'system' ? value : null
    } catch {
        return null
    }
}

function saveTheme(value: Theme) {
    try {
        localStorage.setItem(KEY, value)
    } catch {
        // 存储被禁用时忽略
    }
}

function getSystemTheme(): 'light' | 'dark' {
    try {
        return window.matchMedia('(prefers-color-scheme: light)').matches ? 'light' : 'dark'
    } catch {
        return 'dark'
    }
}

export const theme = ref<Theme>(readSaved() ?? 'system')

// 实际生效的主题，system 已解析成 light 或 dark
export const resolved = ref<'light' | 'dark'>('dark')

function applyTheme() {
    resolved.value = theme.value === 'system' ? getSystemTheme() : theme.value
    document.body.classList.remove('light', 'dark')
    document.body.classList.add(resolved.value)
}

export function toggleTheme() {
    const order: Theme[] = ['system', 'dark', 'light']
    const next = order[(order.indexOf(theme.value) + 1) % order.length] ?? 'system'
    saveTheme(next)

    const change = () => {
        theme.value = next
        applyTheme()
    }

    // 支持就整页淡入淡出，等 Vue 更新完再截切换后的画面
    if (
        !document.startViewTransition ||
        window.matchMedia('(prefers-reduced-motion: reduce)').matches
    )
        return change()
    document.startViewTransition(async () => {
        change()
        await nextTick()
    })
}

export function initTheme() {
    applyTheme()

    window.matchMedia('(prefers-color-scheme: light)').addEventListener('change', () => {
        if (theme.value === 'system') applyTheme()
    })
}

// 当前主题对应的图标和提示
export const themeState = computed(() => {
    switch (theme.value) {
        case 'light':
            return { icon: 'light_mode', tip: '浅色模式' }
        case 'dark':
            return { icon: 'dark_mode', tip: '深色模式' }
        default:
            return { icon: 'contrast', tip: '跟随系统' }
    }
})

// 放进页头右上角的主题切换按钮；面板各页的 Page 传这个
export const themeMotion = computed<Motion[]>(() => [
    { ...themeState.value, onClick: toggleTheme },
])
