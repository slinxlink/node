<template>
    <div class="json">
        <div class="header">
            <div class="title-wrap">
                <span class="title">高级配置</span>
                <span class="subtitle">当前版本只支持查看功能</span>
            </div>
            <RadioGroup v-model="activeSection" :options="sections" />
        </div>
        <div class="editor" ref="editorRef" />
    </div>
</template>

<script setup lang="ts">
import RadioGroup from '@/component/ui/Radio.vue'
import { EditorView, basicSetup } from 'codemirror'
import { json } from '@codemirror/lang-json'
import { oneDark } from '@codemirror/theme-one-dark'
import { Compartment } from '@codemirror/state'
import { isDark } from '@/composable/theme'

const props = defineProps<{
    content: string
}>()

const editorRef = ref<HTMLElement>()
const activeSection = ref('all')
let view: EditorView | null = null
const themeCompartment = new Compartment()

const sections = [
    { label: '全部', value: 'all' },
    { label: '入站', value: 'inbounds' },
    { label: '出站', value: 'outbounds' },
    { label: '端点', value: 'endpoints' },
    { label: '路由', value: 'route' },
]

const displayContent = computed(() => {
    if (!props.content) return '{}'
    if (activeSection.value === 'all') return props.content
    try {
        const parsed = JSON.parse(props.content)
        const section = parsed[activeSection.value]
        return section ? JSON.stringify(section, null, 2) : '{}'
    } catch {
        return '{}'
    }
})

watch(displayContent, (val) => {
    if (view) {
        view.dispatch({
            changes: { from: 0, to: view.state.doc.length, insert: val }
        })
    }
})

watch(isDark, (dark) => {
    view?.dispatch({
        effects: themeCompartment.reconfigure(dark ? oneDark : [])
    })
})

onMounted(() => {
    if (!editorRef.value) return
    view = new EditorView({
        doc: displayContent.value,
        extensions: [
            basicSetup,
            json(),
            themeCompartment.of(isDark.value ? oneDark : []),
            EditorView.editable.of(false),
            EditorView.theme({
                '&': { height: '100%' },
                '.cm-scroller': { overflow: 'auto' },
                '.cm-content': { fontSize: '12px', fontFamily: 'monospace' },
            }),
        ],
        parent: editorRef.value,
    })
})

onUnmounted(() => {
    view?.destroy()
})
</script>

<style scoped>
.json {
    display: flex;
    flex-direction: column;
    height: calc(100vh - 180px);
    padding: 20px;
    background-color: var(--color-bg-dark);
    border-radius: 20px;
    gap: 20px;

    .header {
        display: flex;
        flex-direction: row;
        align-items: center;
        justify-content: space-between;
        gap: 10px;

        @media (max-width: 768px) {
            flex-direction: column;
            align-items: flex-start;
        }

        .title-wrap {
            display: flex;
            flex-direction: column;
            gap: 5px;

            .title {
                font-size: var(--font-size);
                font-weight: bold;
                color: var(--color-text-light);
            }
            .subtitle {
                font-size: var(--font-size-sm);
                color: var(--color-text-dark);
            }
        }
    }

    .editor {
        flex: 1;
        overflow: hidden;
        border-radius: 10px;
    }
}
</style>