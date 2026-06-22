<template>
    <div class="link">
        <div class="left">
            <span class="tag" :class="color">{{ label }}</span>
            <span class="name">{{ name }}</span>
        </div>
        <div class="right">
            <Download v-if="download" :filename="name" :content="download" content-type="application/json" icon="data_object" />
            <Copy :value="value" />
            <QRcode :name="name" :value="value" />
        </div>
    </div>
</template>

<script setup lang="ts">
import Copy from '@/component/widget/Copy.vue'
import QRcode from '@/component/widget/QRcode.vue'
import Download from '@/component/widget/Download.vue'

defineProps<{
    label: string
    color?: string
    name: string
    value: string
    download?: string
}>()
</script>

<style scoped>
.link {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 10px;
    border-radius: 10px;
    background-color: var(--color-bg);

    .left {
        display: flex;
        align-items: center;
        gap: 10px;
    }

    .right {
        display: flex;
        align-items: center;
        gap: 5px;
    }

    .name {
        font-size: var(--font-size-sm);
        color: var(--color-text-light);
    }
}
</style>