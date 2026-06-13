<template>
    <div class="sub">
        <Info :data="data" />
    </div>
    <Footer />
</template>

<script setup lang="ts">
import Info from '@/view/sub/widget/Info.vue'
import Footer from '@/view/Footer.vue';
const data = (window as any).__SUB_DATA__ ?? null
</script>


<style scoped>
.sub {
    display: flex;
    flex-direction: column;
    margin: 20px auto 0;
    min-height: calc(100vh - 100px);
    width: 100%;
    align-items: center;
    
}
</style>

<style>
html,
body,
#sub {
    flex-direction: column;
}
</style>