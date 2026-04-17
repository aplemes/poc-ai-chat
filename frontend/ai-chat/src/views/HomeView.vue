<script setup lang="ts">
import RequestForm from '../components/RequestForm.vue'
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
  padding: 2rem 1rem;
  width: 100%;
}
</style>
