<!-- components/ui/Tip.vue -->
<template>
    <div class="tip-wrap">
        <slot />
        <div class="tip-box">
            <slot name="content" />
        </div>
    </div>
</template>

<style scoped>
.tip-wrap {
    position: relative;
    display: inline-flex;
}
.tip-box {
    position: absolute;
    bottom: calc(100% + 8px);
    left: 50%;
    transform: translateX(-50%);
    background-color: var(--color-bg-light);
    border-radius: 8px;
    padding: 8px 12px;
    white-space: nowrap;
    font-size: 12px;
    color: var(--color-text-light);
    opacity: 0;
    visibility: hidden;
    transition: opacity 0.2s ease, visibility 0s 0.2s;
    z-index: 99;
}
.tip-wrap:hover .tip-box {
    opacity: 1;
    visibility: visible;
    transition: opacity 0.2s ease, visibility 0s;
}
</style>