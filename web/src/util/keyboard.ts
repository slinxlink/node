// util/keyboard.ts

export interface KeyBindings {
    [key: string]: (e: KeyboardEvent) => void
}

const stack: KeyBindings[] = []
let listening = false

// 键组：绑定里写组名，一项就覆盖一批按键；先精确匹配 e.key，没命中才查键组；组里的键不响应带 Ctrl、Cmd、Alt 的组合
const groups: Record<string, (e: KeyboardEvent) => boolean> = {
    alnum: (e) => /^[a-z0-9]$/i.test(e.key) && !e.ctrlKey && !e.metaKey && !e.altKey,
}

function find(top: KeyBindings | undefined, e: KeyboardEvent) {
    if (!top) return undefined

    const exact = top[e.key]

    if (exact) return exact

    for (const [name, match] of Object.entries(groups)) {
        if (top[name] && match(e)) return top[name]
    }

    return undefined
}

function onKeydown(e: KeyboardEvent) {
    const handler = find(stack[stack.length - 1], e)
    if (!handler) return
    e.preventDefault()
    handler(e)
}

function ensureListening() {
    if (listening) return
    listening = true
    window.addEventListener('keydown', onKeydown)
}

function teardownIfIdle() {
    if (stack.length > 0) return
    listening = false
    window.removeEventListener('keydown', onKeydown)
}

export function bind(bindings: KeyBindings): KeyBindings {
    stack.push(bindings)
    ensureListening()
    return bindings
}

export function unbind(token: KeyBindings | null | undefined) {
    if (!token) return
    const index = stack.lastIndexOf(token)
    if (index !== -1) stack.splice(index, 1)
    teardownIfIdle()
}
