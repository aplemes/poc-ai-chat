<script setup lang="ts">
import { RouterView } from 'vue-router'
import FloatingChat from './components/FloatingChat.vue'
</script>

<template>
  <RouterView />
  <FloatingChat />
</template>

<style scoped>
.logo {
  display: block;
  margin: 0 auto 2rem;
}
</style>
