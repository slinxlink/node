<template>
    <Button variant="text" aria-label="下载" @click="download">
        <Icon :name="icon ?? 'download'" />
    </Button>
</template>

<script setup lang="ts">
import Button from '@/component/widget/Button.vue'
import Icon from '@/component/widget/Icon.vue'

// 点一下下载：url 是已经生成好的文件地址(比如二维码图片)，content 是要存成文件的文本，二选一；位置和大小由使用的地方定
const props = defineProps<{
    filename: string
    icon?: string // 图标名，默认 download
    url?: string
    content?: string
    contentType?: string
}>()

function download() {
    const link = document.createElement('a')
    link.download = props.filename

    if (props.url) {
        link.href = props.url
        link.click()

        return
    }

    if (props.content === undefined) return

    const blob = new Blob([props.content], { type: props.contentType ?? 'application/octet-stream' })
    link.href = URL.createObjectURL(blob)
    link.click()
    URL.revokeObjectURL(link.href)
}
</script>
