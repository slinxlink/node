<template>
  <router-view />
  <Modal ref="modal" />
</template>

<script setup lang="ts">
import Modal from '@/component/Modal.vue'

const modal = ref()
provide('modal', modal)
</script>