<!-- view/panel/user/Expire.vue -->
<template>
    <span :class="{ red: level === 'expired', orange: level === 'soon' }">{{ text }}</span>
</template>

<script setup lang="ts">
defineProps<{
    text: string
    level: 'expired' | 'soon' | ''
}>()
</script>

<style scoped>
span {
    min-width: 0;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    font-size: 14px;
}
</style>
