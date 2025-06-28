<script setup>
import { RouterLink } from 'vue-router';
import banner from '@/assets/img/banner.png';
</script>

<template>
  <section class="flex flex-col justify-center items-center min-h-screen bg-gray-900">
    <!-- <img :src="banner" alt="Banner" class="mb-8 w-full max-w-2xl rounded-lg shadow-lg" /> -->
    <i class="pi pi-exclamation-triangle text-yellow-500 text-7xl mb-5"></i>
    <h1 class="text-5xl font-bold mb-4 text-white">404 Not Found</h1>
    <p class="text-xl mb-5 text-gray-300">This page does not exist</p>
    <RouterLink
      to="/"
      class="text-white bg-green-700 hover:bg-green-900 rounded-md px-5 py-2 mt-4 font-semibold transition"
    >
      Go Back
    </RouterLink>
  </section>
</template>