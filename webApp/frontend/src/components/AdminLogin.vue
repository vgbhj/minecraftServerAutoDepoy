<template>
    <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-gray-900 via-gray-800 to-blue-900">
      <div class="bg-gray-800 rounded-lg shadow-lg p-8 w-full max-w-sm">
        <h2 class="text-2xl font-bold text-center text-gray-100 mb-6">Admin Login</h2>
        <input
          v-model="adminPassword"
          type="password"
          placeholder="Admin password"
          class="w-full border border-gray-700 bg-gray-900 text-gray-100 px-3 py-2 rounded mb-4 focus:outline-none focus:ring-2 focus:ring-blue-600"
        />
        <button
          @click="login"
          class="w-full px-4 py-2 rounded font-bold transition bg-blue-600 text-white hover:bg-blue-700"
        >
          Login
        </button>
        <div v-if="error" class="mt-4 text-center text-red-400">{{ error }}</div>
      </div>
    </div>
  </template>
  
  <script setup>
  import { ref } from 'vue';
  
  const adminPassword = ref('');
  const error = ref('');
  
  function login() {
    fetch('/api/v1/server/status', {
      headers: { 'X-Admin-Password': adminPassword.value }
    }).then(res => {
      if (res.status === 200) {
        localStorage.setItem('adminPassword', adminPassword.value);
        location.reload();
      } else {
        error.value = 'Неверный пароль';
      }
    });
  }
  </script>