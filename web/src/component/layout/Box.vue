<!-- component/layout/Box.vue -->
<template>
    <section class="box">
        <div v-if="title || $slots.action" class="head">
            <strong v-if="title">{{ title }}</strong>
            <div v-if="$slots.action" class="action">
                <slot name="action" />
            </div>
        </div>
        <slot />
        <template v-if="$slots.footer">
            <span class="fill"></span>
            <div class="footer">
                <slot name="footer" />
            </div>
        </template>
    </section>
</template>

<script setup lang="ts">
defineProps<{
    title?: string
}>()
</script>

<style scoped>
.box {
    display: flex;
    flex-direction: column;
    gap: 10px;
    min-width: 0;
    padding: 20px;
    border-radius: 20px;
    background-color: var(--background-300);
    box-shadow: var(--shadow);
    color: var(--text-200);
    font-size: 14px;

    .head {
        display: flex;
        align-items: center;
        gap: 5px;
        height: 20px;

        strong {
            color: var(--text-300);
        }

        /* 标题行右边放状态、下拉之类的 */
        .action {
            display: flex;
            align-items: center;
            gap: 5px;
            margin-left: auto;
            color: var(--text-300);
            font-size: 12px;
        }
    }

    /* 撑开内容和底部按钮之间的空白：卡片被拉高时按钮沉到底；它两边各有一个间距，所以按钮离上面至少 20px */
    .fill {
        flex: 1;
    }

    /* 贴着底边的一排按钮：撑满卡片宽度，按钮之间用细线隔开，悬停整格变主题色 */
    .footer {
        display: flex;
        margin: 0 -20px -20px;
        border-top: 1px solid var(--background-200);
        border-radius: 0 0 20px 20px;
        overflow: hidden;

        &:deep(button),
        &:deep(a) {
            flex: 1;
            gap: 5px;
            min-height: 40px;
            padding: 10px;
            border: none;
            border-right: 1px solid var(--background-200);
            border-radius: 0;
            background-color: transparent;
            font-size: 12px;

            &:last-child {
                border-right: none;
            }

            &:hover,
            &:focus-visible {
                background-color: var(--primary-300);
                color: var(--white);
            }

            &:active {
                background-color: var(--primary-600);
                color: var(--white);
            }

            &:disabled {
                opacity: 1;
                pointer-events: none;
            }
        }
    }
}
</style>
