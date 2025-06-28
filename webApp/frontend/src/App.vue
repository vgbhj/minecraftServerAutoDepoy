<script setup>
import Navbar from '@/components/Navbar.vue';
import { RouterView } from 'vue-router';
import { ref } from 'vue';
import AdminLogin from '@/components/AdminLogin.vue';

const hasPassword = ref(!!localStorage.getItem('adminPassword'));

function onLoginSuccess() {
  hasPassword.value = true;
}
</script>

<template>
  <Navbar v-if="hasPassword" />
  <RouterView v-if="hasPassword" />
  <AdminLogin v-else @success="onLoginSuccess" />
</template>