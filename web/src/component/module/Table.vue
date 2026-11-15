<template>
    <div ref="tableEl" class="table">
        <div class="header">
            <Input v-model:value="keyword" placeholder="搜索" clearable>
                <template #leading>
                    <Icon name="search" />
                </template>
            </Input>
            <Select v-model:value="pageSize" :items="pageSizeOptions" />
        </div>

        <div v-if="rows.length === 0" class="empty">
            {{ placeholder }}
        </div>

        <template v-else>
            <div class="body">
                <div ref="scrollEl" class="scroll" :class="{ more }" @scroll="updateMore">
                    <div class="row head">
                        <div
                            v-for="(col, index) in columns"
                            :key="col.key"
                            class="cell"
                            :style="{ width: colWidth(index) + 'px' }"
                        >
                            <Button v-if="col.sort !== false" class="sort" variant="text" @click="toggleSort(col.key)">
                                <span class="label">{{ col.label }}</span>
                                <span class="sort-icons">
                                    <Icon name="arrow_drop_up" :class="{ active: sortKey === col.key && sortOrder === 'asc' }" />
                                    <Icon name="arrow_drop_down" :class="{ active: sortKey === col.key && sortOrder === 'desc' }" />
                                </span>
                            </Button>
                            <span v-else class="label">{{ col.label }}</span>
                        </div>
                        <div v-if="actions && inline" class="cell" :style="{ width: actionWidth + 'px' }">
                            <span class="label">操作</span>
                        </div>
                    </div>

                    <div
                        v-for="(row, rowIndex) in rows"
                        :key="String(row.id)"
                        class="row"
                        :class="{ hover: hovered === rowIndex }"
                        @mouseenter="hovered = rowIndex"
                        @mouseleave="hovered = -1"
                    >
                        <div
                            v-for="(col, index) in columns"
                            :key="col.key"
                            class="cell"
                            :class="{ multiline: col.multiline }"
                            :style="{ width: colWidth(index) + 'px' }"
                        >
                            <div class="inner">
                                <component
                                    v-if="col.render"
                                    :is="col.render(row).component"
                                    v-bind="col.render(row).props"
                                >{{ col.render(row).text }}</component>
                                <span v-else class="text">{{ row[col.key] }}</span>
                            </div>
                        </div>
                        <div v-if="actions && inline" class="cell action" :style="{ width: actionWidth + 'px' }">
                            <Tip v-for="(action, i) in visibleActions(row)" :key="i" :msg="action.tip ?? ''">
                                <Anchor v-if="action.href" type="route" :href="action.href(row)">
                                    <Icon :name="action.icon" />
                                </Anchor>
                                <Button v-else :class="{ alert: action.danger }" @click="action.onClick(row)">
                                    <Icon :name="action.icon" />
                                </Button>
                            </Tip>
                        </div>
                    </div>

                    <div
                        v-for="(col, index) in columns.slice(0, -1)"
                        :key="'resizer-' + col.key"
                        class="resizer"
                        :style="{ left: resizerLeft(index) + 'px' }"
                        @mousedown="startResize(index, $event)"
                    ></div>
                </div>

                <div v-if="actions && !inline" class="fixed">
                    <div class="cell action head">操作</div>
                    <div
                        v-for="(row, rowIndex) in rows"
                        :key="String(row.id)"
                        class="cell action"
                        :class="{ hover: hovered === rowIndex }"
                        :style="{ height: rowHeights[rowIndex] ? rowHeights[rowIndex] + 'px' : undefined }"
                        @mouseenter="hovered = rowIndex"
                        @mouseleave="hovered = -1"
                    >
                        <Tip v-for="(action, i) in visibleActions(row)" :key="i" :msg="action.tip ?? ''">
                            <Anchor v-if="action.href" type="route" :href="action.href(row)">
                                <Icon :name="action.icon" />
                            </Anchor>
                            <Button v-else :class="{ alert: action.danger }" @click="action.onClick(row)">
                                <Icon :name="action.icon" />
                            </Button>
                        </Tip>
                    </div>
                </div>
            </div>

            <div class="footer">
                <div class="pager">
                    <Button :disabled="page <= 1" @click="page--">
                        <Icon name="navigate_before" />
                    </Button>
                    <Button
                        v-for="p in pagerList"
                        :key="p"
                        :selected="p === page"
                        @click="page = p"
                    >
                        {{ p }}
                    </Button>
                    <Button :disabled="page >= totalPages" @click="page++">
                        <Icon name="navigate_next" />
                    </Button>
                </div>
                <div class="info">
                    第 {{ startIndex }} - {{ endIndex }} 项 · 共 {{ total }} 项
                </div>
            </div>
        </template>
    </div>
</template>

<script setup lang="ts">
import type { Component } from 'vue'
import Input from '@/component/widget/Input.vue'
import Select from '@/component/widget/Select.vue'
import Anchor from '@/component/widget/Anchor.vue'
import Button from '@/component/widget/Button.vue'
import Icon from '@/component/widget/Icon.vue'
import Tip from '@/component/widget/Tip.vue'

type Row = Record<string, unknown>

interface Column {
    key: string
    label: string
    sort?: boolean // 能否排序，默认 true
    width?: number // 宽屏下的默认列宽(px)，不写按内容自适应、最宽 220；手机上仍按自适应
    maxWidth?: number // 按内容自适应时的最宽(px)，默认 220，超过的内容截断
    multiline?: boolean // 单元格里的内容可以有多行：行跟着最高的单元格变高，列宽按最长的一行算，超过 220 照样截断(内容自己用 max-width: var(--cell-max, 100%) 截断，量宽度时这个变量是 none)
    collation?: 'number' | 'ip' | string[] // 排序时的比较方式：number 从显示内容里提取数字再比；ip 按地址的每一段数值比，IPv4 在 IPv6 前，解析不出来的排最后；数组按显示内容在数组里的先后排，不在数组里的排最后；不写则自动判断
    render?: (row: Row) => { component: Component; props?: Row; text?: string }
}

// href 跳转（渲染成链接，目前只支持站内路径），onClick 触发操作（渲染成按钮），二选一
type ActionItem = {
    icon: string
    tip?: string
    show?: (row: Row) => boolean
} & (
    | { href: (row: Row) => string; onClick?: never; danger?: never }
    | { onClick: (row: Row) => void; danger?: boolean; href?: never }
)

const props = withDefaults(defineProps<{
    columns: Column[]
    data: Row[]
    placeholder?: string
    actions?: ActionItem[]
}>(), {
    placeholder: '暂无数据',
})

function visibleActions(row: Row) {
    return (props.actions ?? []).filter((action) => !action.show || action.show(row))
}

const query = window.matchMedia('(max-width: 768px)')
const inline = ref(query.matches)
const onQuery = (e: MediaQueryListEvent) => {
    inline.value = e.matches
}
query.addEventListener('change', onQuery)
onBeforeUnmount(() => query.removeEventListener('change', onQuery))

const actionWidth = computed(() => {
    const most = Math.max(0, ...rows.value.map((row) => visibleActions(row).length))
    return Math.max(100, most * 36 + 20)
})

const keyword = ref('')
// 窄屏(手机)默认 10 条，桌面 20 条
const pageSize = ref(window.matchMedia('(max-width: 768px)').matches ? 10 : 20)
const pageSizeOptions = [
    { label: '5 条', value: 5 },
    { label: '10 条', value: 10 },
    { label: '20 条', value: 20 },
]
const page = ref(1)
const hovered = ref(-1) // 鼠标所在的行号;数据列和右边固定的操作列是两块,靠行号一起高亮
const sortKey = ref<string | null>(null)
const sortOrder = ref<'asc' | 'desc'>('asc')

const tableEl = ref<HTMLElement | null>(null)
const scrollEl = ref<HTMLElement | null>(null)
const more = ref(false)
const colWidths = ref<(number | null)[]>(props.columns.map(() => null))
const defaultColWidth = 150

// 右边还有内容没滚出来时为 true，用来显示右边缘的渐隐
function updateMore() {
    const el = scrollEl.value
    more.value = !!el && el.scrollLeft + el.clientWidth < el.scrollWidth - 1
}

function colWidth(index: number) {
    return colWidths.value[index] ?? defaultColWidth
}

function resizerLeft(index: number) {
    let left = 0
    for (let i = 0; i <= index; i++) {
        left += colWidth(i)
    }
    return left - 1
}

let resizingIndex = -1
let startX = 0
let startWidth = 0

// 单元格里的内容可以多行(行会变高)，右边固定的操作列和数据列是两块，要把每一行的高度量出来同步过去
const rowHeights = ref<number[]>([])

const multi = computed(() => props.columns.some((col) => col.multiline))

function syncRowHeights() {
    if (!multi.value) return

    nextTick(() => {
        const list = tableEl.value?.querySelectorAll('.scroll .row:not(.head)')
        rowHeights.value = Array.from(list ?? [], (el) => (el as HTMLElement).offsetHeight)
    })
}

function syncColWidths() {
    nextTick(() => {
        const dataRows = tableEl.value?.querySelectorAll('.row:not(.head)')
        const headRow = tableEl.value?.querySelector('.row.head')
        if (!dataRows || dataRows.length === 0 || !headRow) return

        const headLabels = headRow.querySelectorAll('.cell .label')
        const colCount = props.columns.length
        const maxWidths = Array.from({ length: colCount }, () => 0)

        // 量的是内容本来的宽度：先放开多行内容的截断，量完立刻恢复(同一帧内，不会闪)
        tableEl.value?.classList.add('measure')

        dataRows.forEach((row) => {
            const cells = row.querySelectorAll('.cell')
            cells.forEach((cell, i) => {
                if (i >= colCount) return
                const textEl = cell.querySelector('.text') as HTMLElement | null
                const contentWidth = textEl
                    ? textEl.scrollWidth + 25
                    : ((cell.querySelector('.inner')?.firstElementChild as HTMLElement)?.scrollWidth ?? 0) + 25
                if (contentWidth > (maxWidths[i] ?? 0)) maxWidths[i] = contentWidth
            })
        })

        tableEl.value?.classList.remove('measure')

        headLabels.forEach((headEl, i) => {
            if (i >= colCount) return
            const headWidth = (headEl as HTMLElement).scrollWidth + (props.columns[i]?.sort === false ? 20 : 35)
            if (headWidth > (maxWidths[i] ?? 0)) maxWidths[i] = headWidth
        })

        const wide = window.matchMedia('(min-width: 769px)').matches

        maxWidths.forEach((w, i) => {
            colWidths.value[i] = (wide ? props.columns[i]?.width : undefined) ?? Math.min(w, props.columns[i]?.maxWidth ?? 220)
        })
    })
}

function startResize(index: number, e: MouseEvent) {
    e.preventDefault() // 阻止按下后拖动触发文字框选，平时的框选不受影响
    resizingIndex = index
    startX = e.clientX
    startWidth = colWidth(index)
    document.body.style.cursor = 'col-resize'
    window.addEventListener('mousemove', onResize)
    window.addEventListener('mouseup', stopResize)
}

function onResize(e: MouseEvent) {
    if (resizingIndex === -1) return
    const delta = e.clientX - startX
    colWidths.value[resizingIndex] = Math.max(40, startWidth + delta)
}

function stopResize() {
    resizingIndex = -1
    document.body.style.cursor = ''
    window.removeEventListener('mousemove', onResize)
    window.removeEventListener('mouseup', stopResize)
}

function toggleSort(key: string) {
    if (sortKey.value === key) {
        sortOrder.value = sortOrder.value === 'asc' ? 'desc' : 'asc'
    } else {
        sortKey.value = key
        sortOrder.value = 'asc'
    }
}

function extractNumber(value: unknown): number {
    const match = String(value).replace(/,/g, '').match(/-?\d+(\.\d+)?/)

    return match ? Number(match[0]) : NaN
}

// IPv4 拆成 4 段，IPv6 展开成 8 段，都转成数字；解析不出来返回 null
function ipParts(value: unknown): { v6: boolean; parts: number[] } | null {
    const text = String(value).split('%')[0] ?? ''

    if (/^\d{1,3}(\.\d{1,3}){3}$/.test(text)) return { v6: false, parts: text.split('.').map(Number) }

    if (!/^[0-9a-f:]+$/i.test(text) || text.split('::').length > 2) return null

    const [head, tail] = text.split('::')
    const left = head ? head.split(':') : []
    const right = tail ? tail.split(':') : []
    const fill = text.includes('::') ? 8 - left.length - right.length : 0

    if (fill < 0 || (fill === 0 && !text.includes('::') && left.length !== 8)) return null

    return { v6: true, parts: [...left, ...Array(fill).fill('0'), ...right].map((part) => parseInt(part, 16)) }
}

const filtered = computed(() => {
    let list = props.data
    if (keyword.value) {
        list = list.filter((row) =>
            Object.values(row).some((v) => String(v).includes(keyword.value))
        )
    }
    if (sortKey.value) {
        const collation = props.columns.find((col) => col.key === sortKey.value)?.collation

        list = [...list].sort((a, b) => {
            const av = a[sortKey.value!]
            const bv = b[sortKey.value!]
            const sign = sortOrder.value === 'asc' ? 1 : -1

            if (Array.isArray(collation)) {
                const ai = collation.indexOf(String(av))
                const bi = collation.indexOf(String(bv))

                // 不在数组里的排最后，不受升降序影响
                if (ai < 0 || bi < 0) return Number(ai < 0) - Number(bi < 0)

                return (ai - bi) * sign
            }

            if (collation === 'ip') {
                const ap = ipParts(av)
                const bp = ipParts(bv)

                // 解析不出来的排最后，不受升降序影响
                if (!ap || !bp) return Number(!ap) - Number(!bp)
                if (ap.v6 !== bp.v6) return (Number(ap.v6) - Number(bp.v6)) * sign

                for (let i = 0; i < ap.parts.length; i++) {
                    if (ap.parts[i] !== bp.parts[i]) return (ap.parts[i]! - bp.parts[i]!) * sign
                }

                return 0
            }

            if (collation === 'number') {
                const an = extractNumber(av)
                const bn = extractNumber(bv)

                // 取不出数字的排最后，不受升降序影响
                if (isNaN(an) || isNaN(bn)) return Number(isNaN(an)) - Number(isNaN(bn))

                return (an - bn) * sign
            }

            const an = Number(av)
            const bn = Number(bv)
            const result = !isNaN(an) && !isNaN(bn) ? an - bn : String(av).localeCompare(String(bv))

            return result * sign
        })
    }
    return list
})

const total = computed(() => filtered.value.length)
const totalPages = computed(() => Math.max(1, Math.ceil(total.value / pageSize.value)))

const rows = computed(() => {
    const start = (page.value - 1) * pageSize.value
    return filtered.value.slice(start, start + pageSize.value)
})

const startIndex = computed(() => (total.value === 0 ? 0 : (page.value - 1) * pageSize.value + 1))
const endIndex = computed(() => Math.min(page.value * pageSize.value, total.value))

const pagerList = computed(() => {
    const p = page.value
    const tp = totalPages.value
    if (tp <= 5) {
        return Array.from({ length: tp }, (_, i) => i + 1)
    }
    let start = Math.max(1, p - 2)
    let end = start + 4
    if (end > tp) {
        end = tp
        start = end - 4
    }
    return Array.from({ length: 5 }, (_, i) => start + i)
})

watch(rows, () => {
    syncColWidths()
    syncRowHeights()
})
watch(colWidths, () => nextTick(updateMore), { deep: true })

// 搜索词或每页条数变了,结果的页数跟着变,回到第一页,不然停在后面的页上会是空的
watch([keyword, pageSize], () => (page.value = 1))

let observer: ResizeObserver | undefined

onMounted(() => {
    observer = new ResizeObserver(() => {
        updateMore()
        syncRowHeights()
    })
    if (scrollEl.value) observer.observe(scrollEl.value)
    // 创建时数据已经在了(比如切到另一个标签)，数据不会再变，列宽要在这里量一次，否则一直是默认宽度
    syncColWidths()
    syncRowHeights()
})

onBeforeUnmount(() => observer?.disconnect())
</script>

<style scoped>
.table {
    display: flex;
    flex-direction: column;
    width: 100%;
    min-width: 0;
    background-color: var(--background-300);
    border-radius: 20px;
    padding: 10px;

    /* 量列宽时临时放开多行内容的截断 */
    &.measure {
        --cell-max: none;
    }

    .header {
        display: flex;
        align-items: center;
        justify-content: space-between;
        gap: 10px;
        padding: 10px;

        .input {
            width: 300px;
        }

        .select {
            width: 100px;
        }
    }

    .empty {
        padding: 40px 0;
        font-size: 14px;
        color: var(--text-200);
        text-align: center;
    }

    .body {
        display: flex;
        flex-direction: row;
        align-items: flex-start;
        width: 100%;
        min-width: 0;
        padding: 0 10px;
    }

    .scroll {
        position: relative;
        display: flex;
        flex-direction: column;
        flex: 1;
        min-width: 0;
        overflow-x: auto;
        overflow-y: hidden; /* 只设 overflow-x 时，overflow-y 会被浏览器自动算成 auto，必须显式写 hidden */

        &.more {
            -webkit-mask-image: linear-gradient(to right, #000 calc(100% - 20px), transparent);
            mask-image: linear-gradient(to right, #000 calc(100% - 20px), transparent);
        }
    }

    .fixed {
        display: flex;
        flex-direction: column;
        flex-shrink: 0;
        border-left: 1px solid var(--background-200);
    }

    .row {
        display: flex;
        flex-direction: row;
        width: max-content;
        min-width: 100%;
        padding: 0;
        gap: 0;
        background-color: transparent;
        box-shadow: none;
        justify-content: flex-start;
        filter: blur(0); /* 避免 Chrome 的 hover 过渡结束闪回一帧 */
        transition: background-color 300ms ease;

        &.hover {
            background-color: var(--background-200);
        }

        &.head {
            font-size: 14px;
            position: relative;

            .cell {
                align-items: center;
                height: 40px;
            }

            &::after {
                content: '';
                position: absolute;
                left: 0;
                right: 0;
                bottom: 0;
                height: 1px;
                background-color: var(--background-200);
            }
        }
    }

    .cell {
        position: relative;
        display: flex;
        align-items: flex-start;
        padding: 0 10px;
        flex-shrink: 0;
        height: 40px;
        white-space: nowrap;

        &.multiline {
            height: auto;
            min-height: 40px;
        }

        /* 内容至少占 40px 并在里面居中；内容多行时变高，同一行的其他单元格贴着顶部，也就是都显示在第一行 */
        .inner {
            display: flex;
            flex: 1;
            align-items: center;
            min-width: 0;
            min-height: 40px;
        }

        .text {
            min-width: 0;
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
            font-size: 14px;
        }

        /* 只管表头的排序按钮，单元格里放的其他按钮、开关不受影响 */
        .sort {
            justify-content: flex-start;
            gap: 0;
            width: 100%;
            border-radius: 0; /* 按钮自带的大圆角会把表头文字切掉一块 */

            .sort-icons {
                display: flex;
                flex-direction: column;
                align-items: center;

                .icon {
                    font-size: 14px;
                    line-height: 0.5;
                    color: var(--text-300);

                    &.active {
                        color: var(--text-100);
                    }
                }
            }
        }

        &.action {
            min-width: 100px;
            width: auto;
            justify-content: flex-start;
            padding: 5px 10px;
            gap: 2px;
            filter: blur(0);
            transition: background-color 300ms ease;

            &.hover {
                background-color: var(--background-200);
            }

            &.head {
                align-items: center;
                height: 40px;
                padding: 0 10px;
                font-size: 14px;
                position: relative;

                &::after {
                    content: '';
                    position: absolute;
                    left: 0;
                    right: 0;
                    bottom: 0;
                    height: 1px;
                    background-color: var(--background-200);
                }
            }

            button,
            a {
                width: auto;
                padding: 5px;
                border-radius: 5px;
                background-color: transparent;
                color: var(--text-200);

                &:hover {
                    background-color: var(--background-100);
                    color: var(--text-100);
                }

                &:active {
                    background-color: var(--background-200);
                }

                &.alert {
                    color: var(--red);
                }
            }
        }
    }

    .resizer {
        position: absolute;
        top: 0;
        bottom: 0;
        width: 6px;
        margin-left: -3px;
        cursor: col-resize;
        z-index: 1;

        &::after {
            content: '';
            position: absolute;
            top: 0;
            bottom: 0;
            left: 50%;
            margin-left: -1px;
            width: 0;
            border-left: 2px dashed var(--text-300);
            opacity: 0;
            transition: opacity 0.3s ease;
        }

        &:hover::after {
            opacity: 1;
        }
    }

    .footer {
        display: flex;
        align-items: center;
        justify-content: space-between;
        padding: 10px;

        .pager {
            display: flex;
            flex-direction: row;
            gap: 2px;

            button {
                width: 25px;
                height: 25px;
                color: var(--text-200);
                border-radius: 5px;
                background-color: transparent;
                padding: 10px;
                font-size: 12px;

                &:hover {
                    background-color: var(--background-100);
                    color: var(--text-100);

                }
                
                &:active {
                    background-color: var(--primary-600);

                }

                &.selected {
                    color: var(--text-100);
                }
            }
        }

        .info {
            font-size: 12px;
            color: var(--text-200);
        }
    }
}
</style>
