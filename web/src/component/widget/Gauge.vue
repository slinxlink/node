<!-- component/widget/Gauge.vue -->
<template>
    <div class="gauge">
        <div class="wrap">
            <svg viewBox="0 0 100 100">
                <circle
                    class="track"
                    cx="50"
                    cy="50"
                    r="40"
                    :stroke-dasharray="`${full} ${circ}`"
                    stroke-dashoffset="0"
                />
                <circle
                    class="fill"
                    cx="50"
                    cy="50"
                    r="40"
                    :stroke="color"
                    :stroke-dasharray="`${full} ${circ}`"
                    :stroke-dashoffset="offset"
                />
            </svg>
            <span class="pct">{{ pct.toFixed(1) }}%</span>
        </div>
        <strong>{{ name }}</strong>
        <div class="meta">
            <span>{{ detail }}</span>
            <Tip v-if="$slots.tip" :msg="name">
                <Icon name="memory" />
                <template #content>
                    <slot name="tip" />
                </template>
            </Tip>
        </div>
    </div>
</template>

<script setup lang="ts">
import Icon from '@/component/widget/Icon.vue'
import Tip from '@/component/widget/Tip.vue'

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

        svg {
            width: 100%;
            height: 100%;
            transform: rotate(-227.5deg);
        }

        circle {
            fill: none;
            stroke-width: 6;
            stroke-linecap: round;
        }

        .track {
            stroke: var(--background-100);
        }

        .fill {
            transition: stroke-dashoffset 0.6s ease;
        }

        .pct {
            position: absolute;
            top: 50%;
            left: 50%;
            transform: translate(-50%, -50%);
            color: var(--text-100);
            font-size: 15px;
            font-weight: bold;
        }
    }

    strong {
        color: var(--text-100);
        font-size: 13px;
    }

    .meta {
        display: flex;
        align-items: center;
        gap: 5px;
        color: var(--text-300);
        font-size: 12px;

        .icon {
            font-size: 14px;
            cursor: pointer;
        }
    }
}
</style>
