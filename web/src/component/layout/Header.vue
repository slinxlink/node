<!-- component/layout/Header.vue -->
<template>
    <header>
        <hgroup>
            <strong>{{ title }}</strong>
            <template v-if="detail">
                <Icon :name="icon" />
                <span>{{ detail }}</span>
            </template>
        </hgroup>
        <div v-if="motion.length" class="motion">
            <template v-for="(item, index) in motion" :key="index">
                <Tip v-if="item.tip" :msg="item.tip" side="bottom">
                    <Button @click="item.onClick">
                        <Icon :name="item.icon" />
                    </Button>
                </Tip>
                <Button v-else @click="item.onClick">
                    <Icon :name="item.icon" />
                </Button>
            </template>
        </div>
    </header>
</template>

<script setup lang="ts">
import Button from '@/component/widget/Button.vue'
import Icon from '@/component/widget/Icon.vue'
import Tip from '@/component/widget/Tip.vue'

export interface Motion {
    icon: string
    onClick: () => void
    tip?: string
}

withDefaults(
    defineProps<{
        title: string
        detail?: string
        icon?: string // 标题和副标题之间的图标，默认 double_arrow
        motion?: Motion[]
    }>(),
    {
        icon: 'double_arrow',
        motion: () => [],
    },
)
</script>

<style scoped>
header {
    display: flex;
    flex-shrink: 0; /* 页面内容超出可用高度时，标题栏不能被压扁 */
    align-items: center;
    width: 100%;
    height: 80px;
    max-width: 1200px;
    margin: 0 auto;
    padding: 20px;
    font-size: 24px;
    color: var(--text-100);

    hgroup {
        flex: 1;
        min-width: 0;
        display: flex;
        align-items: flex-end;
        gap: 5px;
        height: 1em;
        line-height: 1;

        strong {
            font-weight: normal;
        }

        i {
            color: var(--text-200);
            font-size: 14px;
            line-height: 1;
            font-variation-settings:
                'FILL' 1,
                'wght' 700,
                'opsz' 20;
        }

        span {
            color: var(--text-200);
            font-size: 14px;
        }
    }

    .motion {
        display: flex;
        flex-shrink: 0;
        align-items: center;
        gap: 10px;

        button {
            width: auto;
            padding: 10px;
            border-radius: 10px;
            box-shadow: var(--shadow);
        }
    }
}
</style>
