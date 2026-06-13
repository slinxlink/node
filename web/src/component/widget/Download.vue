<template>
    <span class="download" @click="handle">
        <slot />
        <i class="icon">image</i>
    </span>
</template>

<script setup lang="ts">
const props = defineProps<{
    filename: string
    getCanvas: () => HTMLCanvasElement | undefined
}>()

function handle() {
    const canvas = props.getCanvas()
    if (!canvas) return
    const link = document.createElement('a')
    link.download = props.filename + '.png'
    link.href = canvas.toDataURL()
    link.click()
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