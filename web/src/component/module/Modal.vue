<!-- component/module/Modal.vue -->
<template>
    <Transition name="modal">
        <div v-if="modalState.visible.value" class="mask" :class="{ form: modalState.form.value }" @click.self="onMaskClick">
            <div class="modal" :class="{ loose: modalState.loose.value }" :style="{ width: modalState.width.value || undefined }">
                <div class="content">
                    <p v-if="modalState.title.value" class="title">{{ modalState.title.value }}</p>
                    <p v-if="modalState.detail.value" class="detail"><BBCode :text="modalState.detail.value" /></p>
                    <component
                        :is="modalState.component.value"
                        v-if="modalState.component.value"
                        v-bind="modalState.props.value"
                        :ref="setComponentRef"
                    />
                </div>
                <div class="actions" v-if="modalState.actions.value">
                    <Button class="cancel-btn" @click="modal.close">取消</Button>
                    <Button class="confirm-btn" :disabled="modalState.busy.value" @click="modal.confirm">确认</Button>
                </div>
            </div>
        </div>
    </Transition>
</template>

<script setup lang="ts">
import type { ComponentPublicInstance } from 'vue'
import { modalState, modal } from '@/composable/modal'
import { BBCode } from '@/util/bbcode'
import Button from '@/component/widget/Button.vue'

function onMaskClick() {
    if (modalState.maskClosable.value) modal.close()
}

function setComponentRef(el: Element | ComponentPublicInstance | null) {
    modalState.componentRef.value = el as ComponentPublicInstance | null
}
</script>

<style scoped>
.mask {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    background-color: rgba(0, 0, 0, 0.5);
    -webkit-backdrop-filter: blur(5px);
    backdrop-filter: blur(5px);
    z-index: 98;
}

/* 表单类弹窗：贴着顶部留 20px，高度随内容变化时顶边不动，最高撑到离底部 20px */
.mask.form {
    align-items: flex-start;
    padding: 20px;

    .modal {
        max-width: 100%;
        max-height: 100%;
    }
}

.modal {
    position: relative;
    max-width: 80vw;
    max-height: 80vh;
    min-width: 200px;
    background: var(--background-300);
    border-radius: 20px;
    box-shadow: var(--shadow);
    color: var(--text-100);
    display: flex;
    flex-direction: column;

    .content {
        /* 圆角自己设，不靠外层去裁：有底部按钮时只有上面两个角，没有按钮时四个角都是 */
        border-radius: 20px 20px 0 0;

        &:last-child {
            border-radius: 20px;
        }

        display: flex;
        flex-direction: column;
        gap: 15px;
        align-items: center;
        padding: 20px;
        overflow: auto;

        .title {
            font-size: 14px;
            text-align: center;
        }

        .detail {
            font-size: 14px;
            color: var(--text-300);
            text-align: center;
        }
    }

    /* 内容不裁剪：提示等可以伸出弹窗边缘 */
    &.loose .content {
        overflow: visible;
    }

    .actions {
        display: flex;
        border-top: 1px solid var(--background-200);

        button {
            border-radius: 0;
            width: 100%;
            height: 40px;

            &:first-child {
                border-radius: 0 0 0 20px;
            }

            &:last-child {
                border-radius: 0 0 20px 0;
            }
        }

        .cancel-btn {
            border-right: 1px solid var(--background-200);
        }
    }
}

.modal-enter-active,
.modal-leave-active {
    transition: opacity 0.3s;
}

.modal-enter-from,
.modal-leave-to {
    opacity: 0;
}
</style>
