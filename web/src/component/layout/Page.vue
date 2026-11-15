<!-- component/layout/Page.vue -->
<template>
    <div class="page">
        <Header v-if="title" :title="title" :detail="detail" :icon="icon" :motion="motion" />
        <Main>
            <div v-if="split" class="split">
                <div class="side"><slot name="left" /></div>
                <div class="side"><slot name="right" /></div>
            </div>
            <slot v-else />
        </Main>
        <Footer v-if="footer" :theme="!title" />
    </div>
</template>

<script setup lang="ts">
import Header from '@/component/layout/Header.vue'
import type { Motion } from '@/component/layout/Header.vue'
import Main from '@/component/layout/Main.vue'
import Footer from '@/component/layout/Footer.vue'

// 没有页头的页面(登录页)主题按钮放在底部栏里
withDefaults(
    defineProps<{
        title?: string
        detail?: string
        icon?: string
        motion?: Motion[]
        footer?: boolean
        split?: boolean
    }>(),
    { footer: true },
)
</script>

<style scoped>
.page {
    display: flex;
    flex-direction: column;
    min-width: 0;
    min-height: 0;
    overflow-x: hidden;
    overflow-y: auto;

    /* 左右分栏：两栏各自往下排，窄屏叠成一栏 */
    .split {
        display: grid;
        grid-template-columns: repeat(2, minmax(0, 1fr));
        align-items: start;
        gap: 20px;

        @media (max-width: 768px) {
            grid-template-columns: minmax(0, 1fr);
        }
    }

    .side {
        display: flex;
        flex-direction: column;
        gap: 20px;
        min-width: 0;
    }
}
</style>
