<script setup>
import { ref } from 'vue'

const serverIP = ref('')
const username = ref('')
const password = ref('')
const loading = ref(false)
const message = ref('')
const error = ref('')
const output = ref('')

async function deployServer() {
  loading.value = true
  message.value = ''
  error.value = ''
  output.value = ''
  try {
    const response = await fetch('/api/v1/server/deploy', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        server_ip: serverIP.value,
        username: username.value,
        password: password.value,
      }),
    })
    const data = await response.json()
    if (!response.ok) {
      error.value = data.error || 'Deployment error'
      output.value = data.details || ''
    } else {
      message.value = data.message
      output.value = data.output
    }
  } catch (err) {
    error.value = err.message
  } finally {
    loading.value = false
  }
}
</script>

<template>
    <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-gray-900 via-gray-800 to-blue-900">
        <form
        @submit.prevent="deployServer"
        class="max-w-md w-full p-4 rounded text-gray-100"
      >
        <div class="flex justify-center mb-8">
          <img src="@/assets/img/banner.png" alt="Banner" class="h-25" />
        </div>
        <div class="text-center mb-8 text-xl font-semibold">
          Deploy your Minecraft server via SSH
        </div>
        <div class="mb-4">
          <input
            v-model="serverIP"
            type="text"
            class="w-full border border-gray-700 bg-gray-800 text-gray-100 px-2 py-1 rounded"
            required
            placeholder="Server IP"
          />
        </div>
        <div class="mb-4">
          <input
            v-model="username"
            type="text"
            class="w-full border border-gray-700 bg-gray-800 text-gray-100 px-2 py-1 rounded"
            required
            placeholder="Username"
          />
        </div>
        <div class="mb-4">
          <input
            v-model="password"
            type="password"
            class="w-full border border-gray-700 bg-gray-800 text-gray-100 px-2 py-1 rounded"
            required
            placeholder="Password"
          />
        </div>
        <button
          type="submit"
          :disabled="loading"
          class="w-full px-4 py-2 rounded font-bold transition bg-blue-600 text-white hover:bg-blue-700"
        >
          {{ loading ? 'Deploying...' : 'Deploy server' }}
        </button>
        <div v-if="message" class="mt-4 text-center text-green-400">{{ message }}</div>
        <div v-if="error" class="mt-4 text-center text-red-400">{{ error }}</div>
        <pre
          v-if="output"
          class="mt-4 p-2 rounded text-xs overflow-x-auto bg-gray-800 text-gray-300"
        >{{ output }}</pre>
      </form>
    </div>
</template>