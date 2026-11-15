<template>
    <img :src="src" :alt="alt" :class="{ border }" draggable="false" />
</template>

<script setup lang="ts">
withDefaults(
    defineProps<{
        src: string
        alt?: string
        border?: boolean // 加一圈 1px 的描边,给本身带白边的图用(比如日本国旗,浅色模式下会和背景融在一起)
    }>(),
    {
        alt: '',
    },
)
</script>

<style scoped>
.border {
    box-sizing: border-box;
    border: 1px solid var(--background-200);
    border-radius: 50%; /* 带描边的都是圆形国旗,描边跟着圆走,不然表格里没设圆角的会画出一个方框 */
}

/* 单色黑图跟随主题反色：--invert 为 1 变白，0 保持原色；按钮按下等场景覆盖 --invert 即可 */
.invert {
    filter: invert(var(--invert));
    transition: filter 0.3s;
}
</style>
