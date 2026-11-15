<template>
    <Page class="detect" title="线路检测">
        <div class="grid">
            <IP />
            <Unlock />
            <BackRoute />
        </div>
    </Page>
</template>

<script setup lang="ts">
import Page from '@/component/layout/Page.vue'
import IP from '@/view/panel/detect/IP.vue'
import Unlock from '@/view/panel/detect/Unlock.vue'
import BackRoute from '@/view/panel/detect/BackRoute.vue'
</script>

<style scoped>
.grid {
    display: flex;
    flex-direction: column;
    gap: 20px;
    padding-bottom: 20px;

    /* 标题行做成 60px 高(和设置、核心页容器的标题行一样)：顶到卡片上沿，标题和右上角的控件在里面垂直居中 */
    &:deep(.head) {
        height: 60px;
        margin: -20px -20px -10px; /* 下面的 -10px 抵掉卡片内部的 10px 间距，标题行下面不留空 */
        padding: 0 13px 0 20px; /* 左边是标题，留 20px；右边是控件，离右边和离上边一样是 13px */
    }

    /* 右上角的下拉和检测按钮：都是标准尺寸，34px 高 */
    &:deep(.action .select) {
        width: 130px;
    }

    &:deep(.action .run) {
        width: auto;
        height: 34px;
        padding: 0 10px 0 5px;
        background-color: var(--background-200);

        &:hover:not(:disabled),
        &:focus-visible {
            background-color: var(--background-100);
        }

        &:active:not(:disabled) {
            background-color: var(--primary-600);
        }

        @media (max-width: 768px) {
            padding: 0 6px;

            span {
                display: none;
            }
        }
    }
}
</style>
