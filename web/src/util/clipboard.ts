export async function copyText(text: string): Promise<boolean> {
    if (navigator.clipboard && window.isSecureContext) {
        try {
            await navigator.clipboard.writeText(text)
            return true
        } catch {
            // 继续走 fallback
        }
    }
    const ta = document.createElement('textarea')
    ta.value = text
    ta.setAttribute('readonly', '')
    ta.style.position = 'fixed'
    ta.style.top = '0'
    ta.style.left = '0'
    ta.style.opacity = '0'
    document.body.appendChild(ta)
    ta.focus()
    ta.setSelectionRange(0, text.length)
    let ok = false
    try {
        ok = document.execCommand('copy')
    } catch {
        ok = false
    }
    document.body.removeChild(ta)
    return ok
}
