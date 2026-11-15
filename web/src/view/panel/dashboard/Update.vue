<!-- view/panel/dashboard/Update.vue -->
<template>
    <div class="update">
        <div class="version">
            <div class="row">
                <span>当前版本</span>
                <strong>{{ currentVersion }}</strong>
            </div>
            <div class="row">
                <span>最新版本</span>
                <strong class="new">{{ latestVersion }}</strong>
            </div>
        </div>

        <div class="changelog">
            <Markdown>
                <div v-html="rendered" />
            </Markdown>
        </div>

        <p class="notice">
            <Icon name="info" />
            更新完成后面板将自动重启
        </p>
    </div>
</template>

<script setup lang="ts">
import { Marked } from 'marked'
import Markdown from '@/component/layout/Markdown.vue'
import Icon from '@/component/widget/Icon.vue'

const props = defineProps<{
    currentVersion: string
    latestVersion: string
    changelog: string
}>()

const escape = (text: string) =>
    text.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;')

// 更新说明来自 GitHub Release，里面的原始 HTML 一律当文字显示，不能执行
const marked = new Marked({
    renderer: {
        html: ({ text }) => escape(text),
    },
})

const rendered = computed(() => marked.parse(props.changelog, { async: false }))
</script>

<style scoped>
.update {
    display: flex;
    flex-direction: column;
    gap: 10px;
    width: 100%;
    font-size: 14px;
    text-align: left;

    /* 版本对比：两行，左边标题右边版本号，和设置页的行一样；新版本用绿色 */
    .version {
        padding: 4px 15px;
        border-radius: 15px;
        background-color: var(--background-200);

        .row {
            display: flex;
            align-items: center;
            justify-content: space-between;
            min-height: 48px;
            border-bottom: 1px solid var(--background-300);

            &:last-child {
                border-bottom: none;
            }

            span {
                color: var(--text-300);
            }

            strong {
                color: var(--text-100);
                font-weight: 500;

                &.new {
                    color: var(--green);
                }
            }
        }
    }

    /* 更新内容：放在和上面一样的容器里，太长就在里面滚；里面 markdown 的版式由 Markdown 组件管 */
    .changelog {
        max-height: 300px;
        padding: 15px;
        overflow-y: auto;
        border-radius: 15px;
        background-color: var(--background-200);
    }

    /* 提示条：样式照设置页的黄色提示，34px 高，内容垂直居中 */
    .notice {
        display: flex;
        align-items: center;
        gap: 6px;
        height: 34px;
        padding: 0 9px;
        border-radius: 12px;
        background-color: color-mix(in srgb, var(--yellow) 12%, var(--background-300));
        color: var(--yellow);
        font-size: 12px;
        line-height: 1;
        white-space: nowrap;

        .icon {
            font-size: 16px;
        }
    }
}
</style>
