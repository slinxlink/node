<!-- view/panel/user/Inbounds.vue -->
<template>
    <div class="inbounds">
        <Stamp v-for="item in items" :key="item.id" :color="item.color">
            <span class="line">
                {{ item.protocol }}:{{ item.port }}<template v-if="item.name"> - {{ item.name }}</template>
            </span>
        </Stamp>
        <span v-if="!items.length" class="none">-</span>
    </div>
</template>

<script setup lang="ts">
import Stamp from '@/component/widget/Stamp.vue'

// 用户关联的入站：一行一个，一个 Stamp 包住「协议:端口 - 名称」，颜色跟着协议；没有名称就只显示前面的
export interface Related {
    id: number
    name: string
    port: number
    protocol: string
    color: string
}

defineProps<{
    items: Related[]
}>()
</script>

<style scoped>
.inbounds {
    display: flex;
    flex-direction: column;
    gap: 5px;
    padding: 9px 0; /* 一行标签 22px，上下各 9px 正好 40px，单个入站时和别的单元格一样高 */

    /* 标签各占一行，宽度按内容；太长就在单元格宽度内截断(表格量列宽时 --cell-max 是 none，不截断) */
    min-width: 0;

    .stamp {
        align-self: flex-start;
        min-width: 0;
        width: max-content; /* 按内容本来的宽度，不然标签会先被压窄，表格量到的就不是真实宽度 */
        max-width: var(--cell-max, 100%);

        .line {
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
        }
    }

    .none {
        color: var(--text-200);
    }
}
</style>
