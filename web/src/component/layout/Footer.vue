<!-- component/layout/Footer.vue -->
<template>
    <footer class="footer">
        <div class="copyright">
            <span>&copy; 2026 - {{ year }} Powered by</span>
            <Anchor variant="text" href="https://github.com/slinxlink" external>SLINX</Anchor>
        </div>
        <ThemeButton v-if="theme" />
    </footer>
</template>

<script setup lang="ts">
import Anchor from '@/component/widget/Anchor.vue'
import ThemeButton from '@/component/module/ThemeButton.vue'

// 主题按钮：没有页头的页面(登录页)没地方放，才放在这里
defineProps<{
    theme?: boolean
}>()

const year = new Date().getFullYear()
</script>

<style scoped>
.footer {
    display: flex;
    flex-shrink: 0;
    align-items: center;
    gap: 10px;
    width: 100%;
    max-width: 1200px;
    min-height: 80px;
    margin: 0 auto;
    padding: 0 20px;
    color: var(--text-300);
    font-size: 12px;

    .copyright {
        display: flex;
        align-items: center;
        gap: 5px;
        height: 40px;

        a {
            width: auto;
            font-size: 12px;
        }
    }

    .tip {
        margin-left: auto;
    }
}
</style>
