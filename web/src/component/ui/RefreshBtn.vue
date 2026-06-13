<template>
    <div class="input-refresh-wrap">
        <slot />
        <button class="refresh-btn" @click.stop="handleClick">
            <i class="icon">refresh</i>
        </button>
    </div>
</template>

<script setup lang="ts">
const emit = defineEmits(['click'])

function handleClick() {
    emit('click')
}
</script>

<style scoped>
.input-refresh-wrap {
    position: relative;
    width: 100%;
    display: flex;
    align-items: center;

    :deep(.input),
    :deep(.select .btn),
    :deep(.multiselect-input) {
        padding-right: 36px;
    }
}

.refresh-btn {
    position: absolute;
    right: 4px;
    top: 50%;
    transform: translateY(-50%);
    z-index: 10;
    background: none;
    border: none;
    cursor: pointer;
    padding: 4px;
    border-radius: 999px;
    color: var(--color-text-dark);
    display: flex;
    align-items: center;
    transition: color 0.2s, background 0.2s;

    &:hover {
        background: var(--color-bg-light);
        color: var(--color-text-light);
    }

    .icon {
        font-size: 16px;
    }
}
</style>