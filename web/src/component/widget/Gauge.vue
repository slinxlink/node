<!-- components/widget/Gauge.vue -->
<template>
    <div class="gauge">
        <div class="wrap">
            <svg viewBox="0 0 100 100">
                <circle class="track" cx="50" cy="50" r="40"
                    :stroke-dasharray="`${full} ${circ}`"
                    stroke-dashoffset="0"/>
                <circle class="fill" cx="50" cy="50" r="40"
                    :stroke="color"
                    :stroke-dasharray="`${full} ${circ}`"
                    :stroke-dashoffset="offset"/>
            </svg>
            <div class="center">
                <span class="pct">{{ pct.toFixed(1) }}%</span>
            </div>
        </div>
        <span class="name">{{ name }}</span>
        <div class="meta">
            <span>{{ detail }}</span>
            <Tip v-if="$slots.tip">
                <i class="icon">memory</i>
                <template #content>
                    <slot name="tip" />
                </template>
            </Tip>
        </div>
    </div>
</template>

<script setup lang="ts">
    import Tip from '@/component/ui/Tip.vue'

    const props = defineProps<{
        name: string
        pct: number
        color: string
        detail?: string
    }>()

    const ARC = 275
    const R = 40
    const circ = 2 * Math.PI * R
    const full = circ * (ARC / 360)
    const offset = computed(() => full - (props.pct / 100) * full)
</script>

<style scoped>
.gauge {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 10px;

    .wrap {
        position: relative;
        width: 100px;
        height: 100px;
    }

    svg {
        width: 100px;
        height: 100px;
        transform: rotate(-227.5deg);
    }

    .track {
        fill: none;
        stroke: var(--color-bg-light);
        stroke-width: 6;
        stroke-linecap: round;
    }

    .fill {
        fill: none;
        stroke-width: 6;
        stroke-linecap: round;
        transition: stroke-dashoffset 0.6s ease;
    }

    .center {
        position: absolute;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
    }

    .pct {
        font-size: 15px;
        font-weight: bold;
        color: var(--color-text-light);
    }

    .name {
        font-size: 13px;
        font-weight: bold;
        color: var(--color-text-light);
    }

    .meta {
        display: flex;
        align-items: center;
        gap: 4px;
        font-size: 12px;
        color: var(--color-text-dark);

        .icon {
            font-size: 14px;
            color: var(--color-text-dark);
            cursor: pointer;
        }
    }
}
</style>