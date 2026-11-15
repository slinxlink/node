<!-- view/panel/board/Alive.vue -->
<template>
    <Status :on="ips.length > 0">
        <Tip v-if="ips.length > 0" :msg="ips.join('[br]')" side="bottom" lock>
            <Anchor variant="plain">在线 {{ ips.length }}</Anchor>
        </Tip>
        <span v-else class="off">离线</span>
    </Status>
</template>

<script setup lang="ts">
import Anchor from '@/component/widget/Anchor.vue'
import Status from '@/component/widget/Status.vue'
import Tip from '@/component/widget/Tip.vue'

// 面板用户的在线状态：圆点加文字；在线时只有「N 个 IP」这几个字能悬停，列出在线的 IP
defineProps<{
    ips: string[]
}>()
</script>

<style scoped>
.status {
    font-size: 12px;

    a.plain {
        font-size: 12px;
    }

    .off {
        color: var(--text-300);
    }
}
</style>
