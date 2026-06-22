<template>
    <span class="download" @click="handle">
        <slot />
        <i class="icon">{{ props.icon ?? 'download' }}</i>
    </span>
</template>

<script setup lang="ts">
const props = defineProps<{
    filename: string
    icon?: string
    getCanvas?: () => HTMLCanvasElement | undefined
    content?: string
    contentType?: string
}>()

function handle() {
    if (props.getCanvas) {
        const canvas = props.getCanvas()
        if (!canvas) return
        const link = document.createElement('a')
        link.download = props.filename + '.png'
        link.href = canvas.toDataURL()
        link.click()
    } else if (props.content) {
        const blob = new Blob([props.content], { type: props.contentType ?? 'application/octet-stream' })
        const url = URL.createObjectURL(blob)
        const a = document.createElement('a')
        a.href = url
        a.download = props.filename
        a.click()
        URL.revokeObjectURL(url)
    }
}
</script>

<style scoped>
.download {
    display: inline-flex;
    align-items: center;
    gap: 5px;
    padding: 5px;
    border-radius: 5px;
    cursor: pointer;
    color: var(--color-text-light);

    .icon {
        font-size: 20px;
        color: var(--color-text-dark);
        transition: color 0.3s;
    }

    &:hover {
        background-color: var(--color-bg-light);

        .icon {
            color: var(--color-text-light);
        }
    }
}
</style>