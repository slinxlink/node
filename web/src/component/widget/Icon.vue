<!-- component/widget/Icon.vue -->
<template>
    <i
        class="icon"
        :role="label ? 'img' : undefined"
        :aria-label="label"
        :aria-hidden="label ? undefined : true"
    >
        {{ name }}
    </i>
</template>

<script setup lang="ts">
defineProps<{
    name: string // 图标名，如 mail、appstore
    label?: string // 图标本身有意义时才传，读屏会念它
}>()
</script>

<style scoped>
/* 图标字体随项目打包，不依赖外部站点；没加载出来时退回缩小到 20% 的 Arial，避免显示成一大串图标名 */
@font-face {
    font-family: 'SLINX Icons';
    src: url('@/asset/font/icon.woff2') format('woff2');
    font-display: swap;
}

@font-face {
    font-family: 'SLINX Icons Fallback';
    src: local('Arial');
    size-adjust: 20%;
}

i {
    display: inline-block;
    white-space: nowrap;
    flex-shrink: 0; /* overflow 会让 flex 里的图标被压扁，必须不收缩 */
    width: 1em; /* 每个图标都是 1em 宽，字体加载前后宽度不变 */
    overflow: hidden;
    font-family: 'SLINX Icons', 'SLINX Icons Fallback';
    font-weight: normal;
    font-style: normal;
    font-size: 20px;
    line-height: 1;
    font-variation-settings: 'FILL' 1;
    font-feature-settings: 'liga';
    -webkit-font-smoothing: antialiased;
}
</style>
