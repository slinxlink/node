<template>
    <Box title="一键生成">
        <div class="desc">
            <p>自动生成 VLESS + Reality + XTLS 节点及订阅用户</p>
            <p>无需证书，开箱即用</p>
        </div>
        <template #footer>
            <Button @click="confirm">立即生成</Button>
        </template>
    </Box>
</template>

<script setup lang="ts">
import Box from '@/component/layout/Box.vue'
import Button from '@/component/widget/Button.vue'
import { quick } from '@/api/quick'
import { modal } from '@/composable/modal'
import { toast } from '@/composable/toast'

function confirm() {
    modal.show({
        title: '一键生成',
        detail: '将自动生成 VLESS + Reality + XTLS 入站及用户',
        onConfirm: async () => {
            toast.loading('生成中...')
            try {
                await quick()
            } catch (err) {
                toast.failure((err as { error?: string })?.error || '生成失败')
                return false
            }
            toast.success('生成成功')
        },
    })
}
</script>

<style scoped>
.desc {
    color: var(--text-300);
    font-size: 12px;
    line-height: 1.5;
}
</style>
