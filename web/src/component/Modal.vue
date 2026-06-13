<template>
    <dialog class="modal" :class="{ active: visible }" @click.self="mode !== 'warn' && close()">
        <div class="container">
            <button class="close-btn" :disabled="mode === 'warn'" @click="close">
                <i class="icon">close</i>
            </button>
            <div class="content">
                <i class="icon" :class="iconClass">{{ iconName }}</i>
                <p class="text">{{ message }}</p>
            </div>
            <div class="actions" :class="{ split: mode === 'confirm' }">
                <button v-if="mode === 'confirm'" class="cancel-btn" @click="close">取消</button>
                <button class="confirm-btn" :disabled="mode === 'warn'" @click="confirm">确认</button>
            </div>
        </div>
    </dialog>
</template>

<script setup lang="ts">
import { keyStack } from '@/composable/keyStack'

const visible = ref(false)
const mode = ref<'success' | 'error' | 'confirm' | 'warn'>('success')
const message = ref('')
const onConfirm = ref<(() => void) | null>(null)

const iconClass = computed(() => {
    if (mode.value === 'success') return 'c-yes'
    if (mode.value === 'error') return 'c-no'
    if (mode.value === 'warn') return 'c-reject'
    return 'c-reject'
})

const iconName = computed(() => {
    if (mode.value === 'success') return 'check_circle'
    if (mode.value === 'error') return 'cancel'
    if (mode.value === 'warn') return 'report'
    return 'report'
})

keyStack(() => visible.value, (e) => {
    if (mode.value === 'warn') return
    if (e.key === 'Escape') close()
    if (e.key === 'Enter') confirm()
})

function show(m: 'success' | 'error' | 'confirm' | 'warn', msg: string, callback?: () => void) {
    mode.value = m
    message.value = msg
    onConfirm.value = callback ?? null
    visible.value = true
}

function close() {
    visible.value = false
}

function confirm() {
    onConfirm.value?.()
    close()
}

defineExpose({ show })
</script>

<style scoped>
.modal {
    display: none;
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    margin: 0;
    padding: 0;
    z-index: 999;
    border: none;
    background-color: rgba(0, 0, 0, 0.5);
    backdrop-filter: blur(0);
    transition: backdrop-filter 0.3s ease;

    &.active {
        display: flex;
        align-items: center;
        justify-content: center;
        backdrop-filter: blur(5px);

        .container {
            opacity: 1;
        }
    }

    .container {
        position: relative;
        background: var(--color-bg-dark);
        border-radius: 10px;
        min-width: 250px;
        max-width: 80%;
        text-align: center;
        color: var(--color-text);
        font-size: var(--font-size);
        opacity: 0;
        transition: opacity 0.3s ease;

        .close-btn {
            position: absolute;
            top: 5px;
            right: 5px;
            border-radius: 5px;
            background-color: transparent;
            color: var(--color-text-dark);

            &:hover {
                background-color: var(--color-primary);
                color: var(--color-text-light);
            }
            &:active {
                background-color: var(--color-primary-dark);
            }
        }

        .content {
            padding: 20px;
            border-bottom: 1px solid var(--color-bg);

            .icon {
                font-size: 40px;
                margin-bottom: 10px;
            }

            .text {
                font-size: 14px;
            }
        }

        .actions {
            display: flex;

            .cancel-btn,
            .confirm-btn {
                flex: 1;
                display: flex;
                align-items: center;
                justify-content: center;
                padding: 10px;
                background-color: transparent;
                border: none;
                cursor: pointer;
                color: var(--color-text);

                &:hover {
                    background-color: var(--color-primary);
                    color: var(--color-text-light);
                }
                &:active {
                    background-color: var(--color-primary-dark);
                }
            }

            .cancel-btn {
                display: none;
                border-radius: 0 0 0 10px;
                border-right: 1px solid var(--color-bg);
            }

            .confirm-btn {
                border-radius: 0 0 10px 10px;
            }

            &.split .cancel-btn {
                display: flex;
            }

            &.split .confirm-btn {
                border-radius: 0 0 10px 0;
            }
        }
    }
}
</style>