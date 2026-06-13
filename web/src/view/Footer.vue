<template>
  <footer>
      <button class="switch" @click="toggleTheme">
        <span class="icon">{{ isDark ? 'light_mode' : 'dark_mode' }}</span>
      </button>
  </footer>
</template>

<script setup lang="ts">
import { isDark } from '@/composable/theme'

onMounted(() => {
    if (!isDark.value) document.body.classList.add('light-mode')
})

const toggleTheme = () => {
    isDark.value = !isDark.value
    document.body.classList.toggle('light-mode')
    localStorage.setItem('theme', isDark.value ? 'dark' : 'light')
}
</script>