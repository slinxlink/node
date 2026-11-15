<!-- component/module/ThemeButton.vue -->
<template>
    <Tip :msg="themeState.tip">
        <Button :aria-label="themeState.tip" @click="toggleTheme">
            <Icon :name="themeState.icon" />
        </Button>
    </Tip>
</template>

<script setup lang="ts">
import Button from '@/component/widget/Button.vue'
import Icon from '@/component/widget/Icon.vue'
import Tip from '@/component/widget/Tip.vue'
import { themeState, toggleTheme } from '@/composable/theme'
</script>

<style scoped>
button {
    width: auto;
    padding: 10px;
    box-shadow: var(--shadow);
}
</style>
