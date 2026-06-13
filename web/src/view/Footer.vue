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

<style>
footer {
    width: 100%;
    margin-top: auto;
    padding: 20px;
    color: var(--color-text-dark);
    font-size: var(--font-size-sm);
    display: flex;
    flex-direction: row;
    gap: var(--gap-size-half);

    .switch {
        padding: 10px;
        border-radius: 50%;
        background-color: var(--color-bg-dark);
        box-shadow: var(--box-shadow);
        margin-left: auto;

        &:hover {
            background-color: var(--color-primary);
        }
        &:active {
            background-color: var(--color-primary-dark);
        }
    }
}
</style>