<template>
    <div class="drawer" :class="{ active: modelValue }">
        <div class="container">
            <div class="header">
                <span class="title">{{ title }}</span>
                <button class="close-btn" :disabled="loading" @click="$emit('update:modelValue', false)">
                    <i class="icon">close</i>
                </button>
            </div>
            <div class="drawer-body">
                <slot />
            </div>
            <div class="footer" v-if="footer !== false">
                <button class="save-btn" :disabled="loading" @click="done ? $emit('update:modelValue', false) : $emit('save')">
                    {{ loading ? '处理中...' : done ? '完成' : (saveText || '保存') }}
                </button>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { keyStack } from '@/composable/keyStack'

const props = withDefaults(defineProps<{
    modelValue: boolean
    title: string
    saveText?: string
    loading?: boolean
    done?: boolean
    footer?: boolean
}>(), {
    footer: true
})

const emit = defineEmits(['update:modelValue', 'save'])

keyStack(() => props.modelValue, (e) => {
    if (e.key === 'Escape' && !props.loading) emit('update:modelValue', false)
    if (e.key === 'Enter' && !props.loading) props.done ? emit('update:modelValue', false) : emit('save')
})

watch(() => props.loading, (val) => {
    document.body.style.cursor = val ? 'wait' : ''
})
</script>

<style scoped>
.drawer {
    display: none;
    backdrop-filter: blur(0);
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    z-index: 999;
    overflow-y: auto;
    background-color: rgba(0, 0, 0, 0.5);
    transition: backdrop-filter 0.3s ease;

    &.active {
        display: block;
        backdrop-filter: blur(5px);

        .container {
            opacity: 1;
        }
    }

    .container {
        position: relative;
        top: auto;
        left: auto;
        transform: none;
        margin: 20px auto 40px;
        width: min(800px, 90vw);
        background-color: var(--color-bg-dark);
        border-radius: 20px;
        box-shadow: 0 0.25rem 0.5rem rgba(0, 0, 0, 0.3);
        color: var(--color-text);
        font-size: var(--font-size);
        opacity: 0;
        transition: opacity 0.3s ease;

        .header {
            display: flex;
            align-items: center;
            justify-content: space-between;
            padding: 16px 20px;
            border-bottom: 1px solid var(--color-bg);

            .title {
                font-size: var(--font-size-md);
                font-weight: 600;
                color: var(--color-text-light);
            }

            .close-btn {
                background: none;
                border: none;
                cursor: pointer;
                padding: 4px;
                border-radius: 4px;
                color: var(--color-text-dark);
                display: flex;
                align-items: center;
                transition: color 0.3s, background 0.3s;

                &:hover {
                    background: var(--color-bg-light);
                    color: var(--color-text-light);
                }

                &:disabled {
                    opacity: 0.5;
                    cursor: not-allowed;
                    pointer-events: none;
                }
            }
        }

        .drawer-body {
            padding: 0;
            overflow: visible;
        }

        .footer {
            display: flex;
            justify-content: flex-end;
            padding: 14px 20px;
            border-top: 1px solid var(--color-bg);

            .save-btn {
                padding: 9px 20px;
                border: none;
                border-radius: 999px;
                color: var(--color-text-light);

                &:disabled {
                    background-color: var(--color-bg-dark);
                    border: 1px solid var(--color-bg-light);
                    cursor: not-allowed;
                    pointer-events: none;
                }
            }
        }
    }
}
</style>