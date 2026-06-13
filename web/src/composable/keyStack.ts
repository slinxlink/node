const stack: symbol[] = []
let handledEvent: KeyboardEvent | null = null

export function keyStack(isOpen: () => boolean, handler: (e: KeyboardEvent) => void) {
    const id = Symbol()
    let active = false

    const onKeydown = (e: KeyboardEvent) => {
        if (!active) return
        if (handledEvent === e) return
        if (stack[stack.length - 1] === id) {
            handledEvent = e
            handler(e)
        }
    }

    watch(isOpen, (val) => {
        if (val) {
            stack.push(id)
            active = false
            requestAnimationFrame(() => { active = true })
        } else {
            active = false
            const index = stack.indexOf(id)
            if (index !== -1) stack.splice(index, 1)
        }
    }, { immediate: true })

    onMounted(() => window.addEventListener('keyup', onKeydown))
    onUnmounted(() => {
        window.removeEventListener('keyup', onKeydown)
        const index = stack.indexOf(id)
        if (index !== -1) stack.splice(index, 1)
    })
}