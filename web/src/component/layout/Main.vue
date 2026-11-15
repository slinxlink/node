<template>
    <main>
        <slot />
    </main>
</template>

<style scoped>
main {
    display: flex;
    flex-direction: column;
    flex: 1 0 auto;
    width: 100%;
    min-width: 0;
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 20px;
    background-color: var(--background-200);
    z-index: 0;
    gap: 20px;
}
</style>
