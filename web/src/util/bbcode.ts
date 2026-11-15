import { defineComponent, h, type VNode } from 'vue'

// 弹窗说明文字用的 BBCode：[b] [i] [u] [s] [color=令牌名] [size=10 到 24] [br]
// 不拼 HTML，解析成节点再渲染成元素，所以没有 XSS；不认识或者不合规的标签当普通文字显示
// 用户填的内容要放进去时用 bb 标签模板：bb`确认删除 ${name}`，插进去的变量永远是纯文本，不会被当成标签

const STYLES = ['b', 'i', 'u', 's']
const COLORS = ['red', 'orange', 'yellow', 'green', 'cyan', 'blue', 'purple', 'rainbow']

type Node =
    | { type: 'text'; text: string }
    | { type: 'br' }
    | { type: 'style'; tag: string; children: Node[] }
    | { type: 'color'; name: string; children: Node[] }
    | { type: 'size'; value: number; children: Node[] }

interface Frame {
    name: string
    make: (children: Node[]) => Node
    parent: Node[]
    children: Node[]
}

type Token =
    | { kind: 'br' }
    | { kind: 'close'; name: string }
    | { kind: 'open'; name: string; make: (children: Node[]) => Node }

function read(token: string): Token | null {
    if (token === 'br') return { kind: 'br' }

    const close = token.match(/^\/(b|i|u|s|color|size)$/)
    if (close) return { kind: 'close', name: close[1]! }

    if (STYLES.includes(token)) {
        return { kind: 'open', name: token, make: (children) => ({ type: 'style', tag: token, children }) }
    }

    const color = token.match(/^color=([a-z]+)$/)
    if (color && COLORS.includes(color[1]!)) {
        const name = color[1]!
        return { kind: 'open', name: 'color', make: (children) => ({ type: 'color', name, children }) }
    }

    const size = token.match(/^size=(\d{2})$/)
    if (size && Number(size[1]) >= 10 && Number(size[1]) <= 24) {
        const value = Number(size[1])
        return { kind: 'open', name: 'size', make: (children) => ({ type: 'size', value, children }) }
    }

    return null
}

function parse(source: string): Node[] {
    const root: Node[] = []
    const stack: Frame[] = []
    let current = root
    let text = ''

    const flush = () => {
        if (text) current.push({ type: 'text', text })
        text = ''
    }

    // 把栈顶的标签收起来，放进它的上一层
    const pop = () => {
        const frame = stack.pop()!
        frame.parent.push(frame.make(frame.children))
        current = stack.length ? stack[stack.length - 1]!.children : root
    }

    for (let i = 0; i < source.length; i++) {
        const char = source[i]!

        if (char === '\\' && (source[i + 1] === '[' || source[i + 1] === '\\')) {
            text += source[i + 1]
            i++
            continue
        }

        const end = char === '[' ? source.indexOf(']', i) : -1
        const token = end > i ? read(source.slice(i + 1, end)) : null

        if (!token) {
            text += char
            continue
        }

        if (token.kind === 'br') {
            flush()
            current.push({ type: 'br' })
        } else if (token.kind === 'open') {
            flush()
            const frame: Frame = { name: token.name, make: token.make, parent: current, children: [] }
            stack.push(frame)
            current = frame.children
        } else {
            const at = stack.map((frame) => frame.name).lastIndexOf(token.name)

            if (at < 0) {
                text += char
                continue
            }

            flush()

            while (stack.length > at) pop()
        }

        i = end
    }

    flush()

    while (stack.length) pop()

    return root
}

function render(nodes: Node[]): (VNode | string)[] {
    return nodes.map((node) => {
        if (node.type === 'text') return node.text
        if (node.type === 'br') return h('br')

        const children = render(node.children)

        if (node.type === 'color') {
            // rainbow 不是颜色令牌,是 color.css 里的彩虹渐变文字类
            return node.name === 'rainbow' ? h('span', { class: 'rainbow' }, children) : h('span', { style: { color: `var(--${node.name})` } }, children)
        }
        if (node.type === 'size') return h('span', { style: { fontSize: `${node.value}px` } }, children)

        const tag = { b: 'strong', i: 'em', u: 'u', s: 's' }[node.tag]!

        return h(tag, children)
    })
}

export const BBCode = defineComponent({
    props: {
        text: { type: String, required: true },
    },
    setup: (props) => () => render(parse(props.text)),
})

// 标签模板：写死的部分按 BBCode 解析，插进去的变量转成纯文本
export const bb = (strings: TemplateStringsArray, ...values: unknown[]): string =>
    strings.reduce((result, part, i) => {
        const value = i < values.length ? String(values[i]).replace(/[\\[]/g, '\\$&') : ''

        return result + part + value
    }, '')
