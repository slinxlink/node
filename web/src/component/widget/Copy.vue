<template>
    <span class="copy" @click="handle">
        <slot />
        <i class="icon" :class="{ sm: props.size === 'sm' }">{{ copied ? 'check' : 'content_copy' }}</i>
    </span>
</template>

<script setup lang="ts">
const props = defineProps<{
    value: string
    size?: 'sm'
}>()

const copied = ref(false)

async function handle() {
    await navigator.clipboard.writeText(props.value)
    copied.value = true
    setTimeout(() => copied.value = false, 2000)
}
</script>

<style scoped>
.copy {
    display: inline-flex;
    align-items: center;
    gap: 5px;
    padding: 5px;
    border-radius: 5px;
    cursor: pointer;
    color: var(--color-text-light);

    .icon {
        font-size: 20px;
        color: var(--color-text-dark);
        transition: color 0.3s;

        &.sm {
            font-size: 14px;
        }
    }

    &:hover {
        background-color: var(--color-bg-light);

        .icon {
            color: var(--color-text-light);
        }
    }
}
</style>