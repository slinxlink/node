<template>
    <div class="textarea-wrap">
        <textarea
            v-bind="$attrs"
            class="textarea"
            :class="{ 'has-trailing': $slots.trailing }"
            :placeholder="placeholder"
            :value="value"
            :maxlength="maxlength"
            :readonly="readonly"
            :disabled="disabled"
            @input="onInput"
        />
        <!-- 右上角的操作(比如刷新)，位置和 Input 的 trailing 一样 -->
        <span v-if="$slots.trailing" class="trailing">
            <slot name="trailing" />
        </span>
    </div>
</template>

<script setup lang="ts">
defineOptions({ inheritAttrs: false })

defineProps<{
    value: string
    placeholder?: string
    maxlength?: number
    disabled?: boolean
    readonly?: boolean
}>()

const emit = defineEmits<{
    'update:value': [value: string]
}>()

function onInput(e: Event) {
    emit('update:value', (e.target as HTMLTextAreaElement).value)
}
</script>

<style scoped>
.textarea-wrap {
    position: relative;
    width: 100%;

    .trailing {
        position: absolute;
        top: 7px;
        right: 10px;
        display: flex;
        align-items: center;
        justify-content: center;
        height: 20px;
        font-size: 20px;
    }
}

.textarea {
    display: block;
    width: 100%;
    height: 100px;
    min-height: 60px;
    max-height: 400px;
    /* 上面的 10px 用和背景同色的边框来留：滚动条画在边框里面，所以顶端避开圆角；下面和左右还是内边距，右下角的拖动把手不受影响 */
    padding: 0 15px 10px;
    border: solid var(--background-200);
    border-width: 10px 0 0;
    border-radius: 15px;
    outline: none;
    resize: vertical; /* 只能上下拖,浏览器自带右下角的把手 */
    background-color: var(--background-200);
    color: var(--text-100);
    caret-color: var(--primary-300);
    font-family: inherit;
    font-size: 14px;
    line-height: 1.5;
    transition:
        background-color 0.3s,
        border-color 0.3s;

    &.has-trailing {
        padding-right: 40px;
    }

    &:hover,
    &:focus {
        border-color: var(--background-100);
        background-color: var(--background-100);
    }

    &:disabled {
        opacity: 0.5;
        cursor: not-allowed;
        filter: grayscale(50%);
    }
}
</style>
