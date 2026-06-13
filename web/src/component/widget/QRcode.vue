<template>
    <button class="btn" @click="show = true">
        <i class="icon">qr_code</i>
    </button>
    <Teleport to="body">
        <div v-if="show" class="mask" @click.self="show = false">
            <div class="box">
                <div class="header">
                    <span class="name">{{ name }}</span>
                    <Copy :value="value" />
                    <Download :filename="name" :getCanvas="() => canvas" />
                </div>
                <canvas ref="canvas" />
            </div>
        </div>
    </Teleport>
</template>

<script setup lang="ts">
import QRCodeLib from 'qrcode'
import Copy from '@/component/widget/Copy.vue'
import Download from '@/component/widget/Download.vue'
import { keyStack } from '@/composable/keyStack'

const props = defineProps<{
    name: string
    value: string
}>()

const show = ref(false)
const canvas = ref<HTMLCanvasElement>()

keyStack(() => show.value, (e) => {
    if (e.key === 'Escape') show.value = false
})

watch(show, (val) => {
    if (val) {
        nextTick(() => {
            QRCodeLib.toCanvas(canvas.value, props.value, {
                width: 300,
                margin: 2,
                color: { dark: '#000000', light: '#ffffff' }
            })
        })
    }
})
</script>

<style scoped>
.btn {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    padding: 5px;
    border-radius: 5px;
    background: none;
    color: var(--color-text-dark);
    transition: color 0.3s, background-color 0.3s;

    &:hover {
        background-color: var(--color-bg-light);
        color: var(--color-text-light);
    }

    .icon {
        font-size: 20px;
    }
}

.mask {
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.5);
    backdrop-filter: blur(5px);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 9999;
}

.box {
    background: var(--color-bg-dark);
    border-radius: 16px;
    padding: 20px;
    display: flex;
    flex-direction: column;
    gap: 12px;
}

.header {
    display: flex;
    align-items: center;
    gap: 5px;

    .name {
        flex: 1;
        font-size: var(--font-size);
        color: var(--color-text-light);
    }
}
</style>