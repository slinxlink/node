<!-- component/module/QRCode.vue -->
<template>
    <div class="qrcode" :style="{ width: `${SIZE}px` }">
        <div class="tools">
            <span class="name">{{ name }}</span>
            <Tip msg="下载图片">
                <Download :filename="`${name}.png`" :url="url" icon="archive" />
            </Tip>
            <Tip msg="复制链接">
                <Copy :value="value" />
            </Tip>
        </div>
        <img v-if="url" :src="url" :width="SIZE" :height="SIZE" alt="二维码" />
    </div>
</template>

<script setup lang="ts">
import Copy from '@/component/widget/Copy.vue'
import Download from '@/component/widget/Download.vue'
import Tip from '@/component/widget/Tip.vue'
import { buildQRCode } from '@/util/qrcode'

// 弹窗里的内容：上面一排是名称(靠左)和复制、下载按钮(靠右)，下面是二维码；value 是二维码里的内容，name 同时是下载的文件名
const props = defineProps<{
    name: string
    value: string
}>()

const SIZE = 220

const url = ref('')

onMounted(async () => {
    url.value = await buildQRCode({ value: props.value, size: SIZE })
})

onBeforeUnmount(() => {
    if (url.value) URL.revokeObjectURL(url.value)
})
</script>

<style scoped>
.qrcode {
    display: flex;
    flex-direction: column;
    gap: 10px;

    .tools {
        display: flex;
        align-items: center;
        gap: 10px;

        .name {
            flex: 1;
            min-width: 0;
            overflow: hidden;
            color: var(--text-100);
            text-align: left;
            text-overflow: ellipsis;
            white-space: nowrap;
        }
    }

    img {
        display: block;
        align-self: center;
        border-radius: 15px;
    }
}
</style>
