<template>
    <RouterView />
    <Toast />
    <Modal />
</template>

<script setup lang="ts">
import Toast from '@/component/module/Toast.vue'
import Modal from '@/component/module/Modal.vue'
import { initTheme } from '@/composable/theme'

onMounted(() => {
    initTheme()
})
</script>
