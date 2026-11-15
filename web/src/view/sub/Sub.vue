<template>
    <Page class="sub">
        <div class="box">
            <Info v-if="data" :user="data.user" :inbounds="data.inbounds" :data="data" />
            <p v-else class="empty">订阅不存在或已停用</p>
        </div>
    </Page>

    <Toast />
</template>

<script setup lang="ts">
import Page from '@/component/layout/Page.vue'
import Toast from '@/component/module/Toast.vue'
import { initTheme } from '@/composable/theme'
import Info from '@/view/panel/user/Info.vue'

// 订阅页：用户打开订阅链接看到的页面，信息是后端直接写进页面的，不用登录、不请求接口(链接里的数据都带着了)
const data = window.__SUB_DATA__ ?? null

onMounted(initTheme)
</script>

<style scoped>
.sub {
    .box {
        width: 100%;
        max-width: 800px;
        min-width: 0;
        margin: 20px auto 0; /* 贴着顶部留 20px，不居中；底部栏仍在最下面 */
        padding: 20px;
        border-radius: 20px;
        background-color: var(--background-300);
        box-shadow: var(--shadow);
        color: var(--text-200);
        font-size: 14px;
    }

    .empty {
        padding: 40px 0;
        color: var(--text-300);
        text-align: center;
    }
}
</style>
