<!-- component/layout/Markdown.vue，Markdown 内容的版式(更新说明用):12px 小字号、列表紧凑;和 mesh 的 markdown/Wiki.vue 是同一份样式,保持一致;类名叫 markdown。里面放一个 div 包着渲染出来的 HTML(v-html),第一块、最后一块的上下边距靠 div > :first-child 去掉 -->
<template>
    <article class="markdown">
        <slot />
    </article>
</template>

<style scoped>
.markdown {
    font-size: 12px;
    line-height: 1.6;
    color: var(--text-200);

    /* 标题当小节标签用:新增、优化、修复这类 */
    & :deep(h1),
    & :deep(h2),
    & :deep(h3) {
        margin: 14px 0 4px;
        color: var(--text-100);
        font-size: 14px;
        font-weight: 700;

        &:first-child {
            margin-top: 0;
        }
    }

    & :deep(a[href]) {
        color: var(--primary-100);
        text-decoration: none;
        transition: color 0.3s;

        &:hover,
        &:focus-visible {
            color: var(--primary-200);
        }

        &:active {
            color: var(--primary-400);
        }
    }

    & :deep(p) {
        margin: 0.5em 0;
    }

    /* 第一块内容顶边对齐、最后一块贴着容器底边:段间距只管块与块之间,不管最上面和最下面 */
    & :deep(div > :first-child) {
        margin-top: 0;
    }

    & :deep(div > :last-child) {
        margin-bottom: 0;
    }

    & :deep(ul) {
        margin: 0.3em 0;
        padding-left: 20px;
    }

    /* 第一层列表项用标题色，更深的层级用正文色 */
    & :deep(li) {
        margin: 2px 0;
        color: var(--text-100);
    }

    & :deep(li li) {
        color: var(--text-200);
    }

    /* 层级只靠缩进区分 */
    & :deep(li > ul) {
        margin: 1px 0;
        padding-left: 18px;
    }
}
</style>
