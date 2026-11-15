<!-- component/module/JsonEditor.vue -->
<template>
    <div class="json">
        <div class="toolbar">
            <span class="hint">当前版本只支持查看</span>
            <Choice v-model:value="section" :items="sections" />
        </div>
        <div ref="editor" class="editor" />
    </div>
</template>

<script setup lang="ts">
import { EditorView, basicSetup } from 'codemirror'
import { json } from '@codemirror/lang-json'
import { oneDark } from '@codemirror/theme-one-dark'
import { Compartment } from '@codemirror/state'
import Choice from '@/component/widget/Choice.vue'
import { resolved } from '@/composable/theme'

// 只读的 JSON 查看器：可以只看配置里的某一块(入站、出站、端点、路由)，跟着主题切换亮暗
const props = defineProps<{
    content: string
}>()

const editor = ref<HTMLElement>()
const section = ref('all')

const sections = [
    { label: '全部', value: 'all' },
    { label: '入站', value: 'inbounds' },
    { label: '出站', value: 'outbounds' },
    { label: '端点', value: 'endpoints' },
    { label: '路由', value: 'route' },
]

let view: EditorView | null = null
const themeSlot = new Compartment()

// 当前要显示的文字：选了某一块就取出来重新排版，取不到或解析不了就是空对象
const text = computed(() => {
    if (!props.content) return '{}'
    if (section.value === 'all') return props.content

    try {
        const part = JSON.parse(props.content)[section.value]

        return part ? JSON.stringify(part, null, 2) : '{}'
    } catch {
        return '{}'
    }
})

watch(text, (value) => {
    view?.dispatch({
        changes: { from: 0, to: view.state.doc.length, insert: value },
    })
})

watch(resolved, (value) => {
    view?.dispatch({
        effects: themeSlot.reconfigure(value === 'dark' ? oneDark : []),
    })
})

onMounted(() => {
    if (!editor.value) return

    view = new EditorView({
        doc: text.value,
        extensions: [
            basicSetup,
            json(),
            themeSlot.of(resolved.value === 'dark' ? oneDark : []),
            EditorView.editable.of(false),
            EditorView.theme({
                '&': { flex: '1', minHeight: '0' },
                '.cm-scroller': { overflow: 'auto' },
                '.cm-content': { fontSize: '12px', fontFamily: 'monospace' },
            }),
        ],
        parent: editor.value,
    })
})

onUnmounted(() => view?.destroy())
</script>

<style scoped>
.json {
    display: flex;
    flex: 1; /* 放在固定高度的容器里：撑满剩下的高度 */
    flex-direction: column;
    gap: 20px;
    min-height: 0;

    .toolbar {
        display: flex;
        align-items: center;
        justify-content: space-between;
        gap: 10px;

        .hint {
            color: var(--text-300);
            font-size: 12px;
        }
    }

    .editor {
        display: flex;
        flex: 1;
        flex-direction: column;
        min-height: 0;
        overflow: hidden;
        border-radius: 15px;
    }
}
</style>
