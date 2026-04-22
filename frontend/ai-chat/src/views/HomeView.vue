<script setup lang="ts">
import RequestForm from '@/components/RequestForm.vue'
</script>

<template>
  <main class="home">
    <RequestForm />
  </main>
</template>

<style scoped>
.home {
  display: flex;
  justify-content: center;
  align-items: flex-start;
  padding: 2.5rem 1.5rem 5rem;
  width: 100%;
  min-height: 100vh;
  background: #f1f5f9;
}
</style>
