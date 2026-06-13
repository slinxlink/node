<template>
    <div class="update">
        <div class="version">
            <div class="item">
                <span class="label">当前版本</span>
                <span class="value">{{ currentVersion }}</span>
            </div>
            <i class="icon arrow">arrow_forward</i>
            <div class="item">
                <span class="label">最新版本</span>
                <span class="value new">{{ latestVersion }}</span>
            </div>
        </div>
        <div class="changelog" v-html="renderedChangelog" />
    </div>
</template>

<script setup lang="ts">
import { marked } from 'marked'

const props = defineProps<{
    currentVersion: string
    latestVersion: string
    changelog: string
}>()

const renderedChangelog = computed(() => marked(props.changelog))
</script>

<style scoped>
.update {
    display: flex;
    flex-direction: column;
    gap: 20px;
    padding: 20px;
}

.version {
    display: flex;
    flex-direction: row;
    align-items: center;
    gap: 10px;
    padding: 10px;
    background: var(--color-bg);
    border-radius: 10px;

    .item {
        display: flex;
        flex-direction: column;
        gap: 5px;
        font-size: var(--font-size-sm);
        text-align: center;

        .label {
            color: var(--color-text-dark);
        }

        .value {
            font-weight: bold;
            color: var(--color-text-light);

            &.new {
                color: var(--color-green);
            }
        }
    }
}

.arrow {
    color: var(--color-text-dark);
    font-size: var(--font-size-sm);
}

.changelog {
    font-size: var(--font-size-sm);
    color: var(--color-text);
    line-height: 1.5;

    &:deep(h2) {
        font-size: var(--font-size);
        color: var(--color-text-light);
        margin-bottom: 10px;
    }

    &:deep(ul) {
        padding-left: 20px;
        display: flex;
        flex-direction: column;
        gap: 5px;
        font-size: var(--font-size-sm);
    }

    &:deep(li) {
        color: var(--color-text);
    }
}
</style>