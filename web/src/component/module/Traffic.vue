<!-- component/module/Traffic.vue -->
<template>
    <div class="traffic">
        <span class="used">{{ formatBytes(used) }}</span>
        <div class="bar" :class="{ unlimited, over }">
            <div :style="{ width: `${shown}%` }"></div>
        </div>
        <span class="limit">
            <Icon v-if="unlimited" name="all_inclusive" />
            <template v-else>{{ formatBytes(limit) }}</template>
        </span>
    </div>
</template>

<script setup lang="ts">
import Icon from '@/component/widget/Icon.vue'
import { formatBytes } from '@/util/format'

const props = defineProps<{
    used: number
    limit: number
}>()

const unlimited = computed(() => !props.limit)
const over = computed(() => !unlimited.value && props.used >= props.limit)

const percent = computed(() => {
    if (unlimited.value) return 100
    return Math.min(((props.used || 0) / props.limit) * 100, 100)
})

const shown = ref(0)

onMounted(() => requestAnimationFrame(() => (shown.value = percent.value)))

watch(percent, (value) => (shown.value = value))
</script>

<style scoped>
.traffic {
    display: flex;
    align-items: center;
    gap: 10px;
    width: 100%;
    font-size: 12px;
    color: var(--text-200);

    .used,
    .limit {
        flex-shrink: 0;
        width: 62px;
        white-space: nowrap;
    }

    .used {
        text-align: right;
    }

    .limit {
        display: flex;
        align-items: center;

        .icon {
            font-size: 16px;
        }
    }

    .bar {
        flex: 1;
        min-width: 50px;
        height: 10px;
        border-radius: 5px;
        background-color: var(--background-200);
        overflow: hidden;

        div {
            position: relative;
            height: 100%;
            border-radius: 5px;
            background-color: var(--blue);
            overflow: hidden;
            transition: width 0.8s ease-out;

            &::after {
                content: '';
                position: absolute;
                inset: 0;
                background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.35), transparent);
                transform: translateX(-100%);
                animation: sweep 2.4s ease-in-out infinite;
            }
        }

        &.unlimited div {
            background-color: var(--green);
        }

        &.over div {
            background-color: var(--red);
        }
    }
}

@keyframes sweep {
    from {
        transform: translateX(-100%);
    }

    60%,
    to {
        transform: translateX(100%);
    }
}

@media (prefers-reduced-motion: reduce) {
    .bar div::after {
        animation: none;
    }
}
</style>
