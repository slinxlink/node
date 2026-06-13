<template>
    <Nav />
    <div class="view">
        <router-view />
        <Footer />
    </div>
</template>

<script setup lang="ts">
    import Nav from './Nav.vue'
    import Footer from '../Footer.vue';
</script>