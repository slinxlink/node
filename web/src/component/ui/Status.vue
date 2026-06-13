<template>
    <span class="status" :class="[colorClass, { animated: animated }]">
        <span class="circle"></span>
        <span class="circle"></span>
        <span class="circle"></span>
    </span>
</template>

<script setup lang="ts">
const props = defineProps<{
    status: 'yes' | 'no' | 'reject'
    animated?: boolean
}>()

const colorClass = computed(() => ({
    'yes':    'status-green',
    'no':     'status-red',
    'reject': 'status-orange',
}[props.status]))
</script>

<style scoped>
.status-green  { --status-color: #2FB344; }
.status-red    { --status-color: #D63939; }
.status-orange { --status-color: #E0A030; }

.status {
    --size: 20px;
    --circle-size: .75rem;

    display: block;
    position: relative;
    width: var(--size);
    height: var(--size);

    .circle {
        position: absolute;
        left: 50%;
        top: 50%;
        translate: -50% -50%;
        width: var(--circle-size);
        height: var(--circle-size);
        border-radius: 100rem;
        background: var(--status-color, #667382);

        &:nth-child(1) { z-index: 3; }
        &:nth-child(2) { z-index: 2; opacity: .1; }
        &:nth-child(3) { z-index: 1; opacity: .3; }
    }

    &.animated {
        .circle {
            &:nth-child(1) { animation: 2s linear 1s infinite backwards status-pulsate-main; }
            &:nth-child(2) { animation: 2s linear 1s infinite backwards status-pulsate-secondary; }
            &:nth-child(3) { animation: 2s linear 1s infinite backwards status-pulsate-tertiary; }
        }
    }
}

@keyframes status-pulsate-main {
    40% { transform: scale(1.25, 1.25); }
    60% { transform: scale(1.25, 1.25); }
}

@keyframes status-pulsate-secondary {
    10%  { transform: scale(1, 1); }
    30%  { transform: scale(3, 3); }
    80%  { transform: scale(3, 3); }
    100% { transform: scale(1, 1); }
}

@keyframes status-pulsate-tertiary {
    25%  { transform: scale(1, 1); }
    80%  { transform: scale(3, 3); opacity: 0; }
    100% { transform: scale(3, 3); opacity: 0; }
}
</style>