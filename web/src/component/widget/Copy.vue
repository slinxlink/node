<template>
    <Button variant="text" aria-label="复制" @click="copy">
        <span class="icons">
            <Icon name="content_copy" :class="{ hidden: copied }" />
            <Icon name="check" :class="{ hidden: !copied }" />
        </span>
    </Button>
</template>

<script setup lang="ts">
import { toast } from '@/composable/toast'
import Button from '@/component/widget/Button.vue'
import Icon from '@/component/widget/Icon.vue'
import { copyText } from '@/util/clipboard'

// 点一下复制 value,成功后图标变成对勾,2 秒后恢复;位置和大小由使用的地方定
const props = defineProps<{
    value: string
}>()

const copied = ref(false)

async function copy() {
    if (!(await copyText(props.value))) {
        toast.failure('复制失败')

        return
    }

    copied.value = true

    setTimeout(() => (copied.value = false), 2000)
}
</script>

<style scoped>
/* 两个图标叠在同一格里,切换时 300ms 交叉淡入淡出 */
.icons {
    display: grid;

    i {
        grid-area: 1 / 1;
        transition: opacity 0.3s;

        &.hidden {
            opacity: 0;
        }
    }
}
</style>
