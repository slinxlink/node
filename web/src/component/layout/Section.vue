<!-- component/layout/Section.vue -->
<template>
    <section>
        <template v-if="title || $slots.action || motion?.length">
            <div class="head">
                <strong v-if="title">{{ title }}</strong>
                <div v-if="$slots.action" class="action">
                    <slot name="action" />
                </div>
                <div v-if="motion?.length" class="motion">
                    <template v-for="(item, index) in motion" :key="index">
                        <Tip v-if="item.tip" :msg="item.tip">
                            <Button variant="text" @click="item.onClick">
                                <Icon :name="item.icon" />
                            </Button>
                        </Tip>
                        <Button v-else variant="text" @click="item.onClick">
                            <Icon :name="item.icon" />
                        </Button>
                    </template>
                </div>
            </div>
            <Divider />
        </template>
        <div class="content" :class="{ double: columns === 2 }">
            <slot />
        </div>
    </section>
</template>

<script setup lang="ts">
import Divider from '@/component/layout/Divider.vue'
import type { Motion } from '@/component/layout/Header.vue'
import Button from '@/component/widget/Button.vue'
import Icon from '@/component/widget/Icon.vue'
import Tip from '@/component/widget/Tip.vue'

// columns 为 2 时内容排成两列，窗口宽度 768 以下回到单列
withDefaults(
    defineProps<{
        title?: string
        columns?: 1 | 2
        motion?: Motion[] // 标题行右边的操作按钮
    }>(),
    { columns: 1 },
)
</script>

<style scoped>
section {
    display: flex;
    flex-direction: column;
    min-width: 0;
    border-radius: 20px;
    background-color: var(--background-300);
    box-shadow: var(--shadow);
    font-size: 14px;

    /* 容器被固定高度时，标题行和分隔线不能被压扁(分隔线只有 1px，会缩到 0 看不见) */
    .head,
    hr {
        flex-shrink: 0;
    }

    .head {
        display: flex;
        align-items: center;
        justify-content: space-between;
        gap: 10px;

        strong {
            padding: 20px;
            color: var(--text-300);
        }

        /* 标题行里的操作区(如切换 tab),占满剩下的宽度 */
        .action {
            display: flex;
            flex: 1;
            align-items: flex-end;
            min-width: 0;
            height: 60px; /* 固定高度,里面的内容晚到(如 tab 等数据)时标题行不跳 */
            padding: 0 20px;
        }

        /* 标题行最右边的操作按钮(和页面标题栏的 motion 一样：图标、提示、点击函数)，透明背景 */
        .motion {
            display: flex;
            flex-shrink: 0;
            align-items: center;
            gap: 10px;
            padding: 0 12px;

            button {
                width: auto;
                padding: 10px;
                border-radius: 10px;
            }
        }
    }

    /* 内容是空的(比如列表还没有条目)，内容块和标题下面的分隔线都不显示，只剩标题行 */
    &:has(> .content:empty) > hr {
        display: none;
    }

    .content:empty {
        display: none;
    }

    .content {
        display: flex;
        flex-direction: column;
        gap: 10px;
        padding: 20px;
        color: var(--text-100);

        &.double {
            display: grid;
            grid-template-columns: repeat(2, minmax(0, 1fr));
            align-items: start;

            @media (max-width: 768px) {
                grid-template-columns: minmax(0, 1fr);
            }
        }
    }
}
</style>
