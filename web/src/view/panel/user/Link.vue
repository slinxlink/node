<!-- view/panel/user/Link.vue -->
<template>
    <div class="link">
        <div class="row">
            <div class="left">
                <Stamp :color="color">{{ label }}</Stamp>
                <span class="name">{{ name }}</span>
            </div>
            <div class="right">
                <Tip v-if="download" msg="下载 json 配置">
                    <Download :filename="`${name}.json`" :content="download" content-type="application/json" icon="file_save" />
                </Tip>
                <Tip msg="复制链接">
                    <Copy :value="value" />
                </Tip>
                <Tip msg="二维码">
                    <Button variant="text" aria-label="二维码" @click="show">
                        <Icon name="qr_code" />
                    </Button>
                </Tip>
            </div>
        </div>
        <!-- 二维码单独弹出，盖在信息弹窗上面；弹窗同一时间只能开一个，所以这里自己带一层遮罩 -->
        <Teleport to="body">
            <Transition name="mask">
                <div v-if="open" class="mask" @click.self="close">
                    <div class="box">
                        <QRCode :name="name" :value="value" />
                    </div>
                </div>
            </Transition>
        </Teleport>
    </div>
</template>

<script setup lang="ts">
import QRCode from '@/component/module/QRCode.vue'
import Button from '@/component/widget/Button.vue'
import Copy from '@/component/widget/Copy.vue'
import Download from '@/component/widget/Download.vue'
import Icon from '@/component/widget/Icon.vue'
import Stamp from '@/component/widget/Stamp.vue'
import Tip from '@/component/widget/Tip.vue'
import { bind, unbind, type KeyBindings } from '@/util/keyboard'

// 一条链接：协议或类型的标签、名称，右边是下载配置(可选)、复制、二维码；点二维码弹出大图
defineProps<{
    label: string
    color?: string
    name: string
    value: string
    download?: string // 有值才显示下载按钮，内容是要存成文件的 JSON
}>()

const open = ref(false)

let token: KeyBindings | null = null

function show() {
    open.value = true
    token = bind({ Escape: close })
}

function close() {
    open.value = false
    unbind(token)
    token = null
}

onBeforeUnmount(close)
</script>

<style scoped>
.link {
    display: flex;
    flex-direction: column;
    gap: 10px;
    padding: 5px 10px 5px 5px;
    border-radius: 999px;
    background-color: var(--background-200);

    .row {
        display: flex;
        align-items: center;
        justify-content: space-between;
        gap: 10px;

        .left {
            display: flex;
            align-items: center;
            gap: 10px;
            min-width: 0;

            .name {
                overflow: hidden;
                color: var(--text-100);
                font-size: 12px;
                text-overflow: ellipsis;
                white-space: nowrap;
            }
        }

        .right {
            display: flex;
            align-items: center;
            flex-shrink: 0;
            gap: 5px;
        }
    }
}

.mask {
    position: fixed;
    inset: 0;
    z-index: 98;
    display: flex;
    align-items: center;
    justify-content: center;
    background-color: rgba(0, 0, 0, 0.5);
    -webkit-backdrop-filter: blur(5px);
    backdrop-filter: blur(5px);

    &.mask-enter-active,
    &.mask-leave-active {
        transition: opacity 0.3s;
    }

    &.mask-enter-from,
    &.mask-leave-to {
        opacity: 0;
    }

    .box {
        padding: 20px;
        border-radius: 20px;
        background-color: var(--background-300);
        box-shadow: var(--shadow);
    }
}
</style>
