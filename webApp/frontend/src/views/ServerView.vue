<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'

const status = ref('Loading...')
const version = ref('...')
const serverIp = ref('...')
const loading = ref(false)
const error = ref('')

async function fetchStatus() {
  try {
    const res = await axios.get('/api/v1/server/status')
    status.value = res.data.status === 'UP' ? 'Online' : 'Offline'
  } catch (e) {
    status.value = 'Error'
  }
}

async function fetchVersion() {
  try {
    const res = await axios.get('/api/v1/minecraft/current')
    version.value = res.data.version
  } catch (e) {
    version.value = 'Error'
  }
}

async function fetchIp() {
  try {
    const res = await axios.get('/api/v1/server/ip')
    serverIp.value = res.data.ip
  } catch (e) {
    serverIp.value = 'Error'
  }
}

async function serverAction(action) {
  loading.value = true
  error.value = ''
  try {
    await axios.post(`/api/v1/server/${action}`)
    await fetchStatus()
  } catch (e) {
    error.value = 'Error: ' + (e.response?.data?.error || e.message)
  }
  loading.value = false
}

onMounted(() => {
  fetchStatus()
  fetchVersion()
  fetchIp()
})
</script>

<template>
  <div style="margin-left: 14rem;">
    <div class="min-h-screen flex items-center justify-center bg-gray-900">
      <div class="max-w-xl mx-auto mt-10 p-8 bg-gray-800 rounded-lg">
        <h2 class="text-2xl font-bold mb-6 text-white flex items-center gap-2">
          Server status:
          <span :class="status==='Online' ? 'text-green-400' : status==='Offline' ? 'text-red-400' : 'text-yellow-400'">
            {{ status }}
          </span>
        </h2>
        <div class="flex flex-wrap gap-4 mb-6">
          <button
            class="bg-green-500 hover:bg-green-600 text-white font-semibold py-2 px-6 rounded transition disabled:opacity-50 flex items-center gap-2"
            :disabled="loading || status==='Online'"
            @click="serverAction('start')"
          >
            <span class="pi pi-caret-right"></span>
            Start
          </button>
          <button
            class="bg-red-500 hover:bg-red-600 text-white font-semibold py-2 px-6 rounded transition disabled:opacity-50 flex items-center gap-2"
            :disabled="loading || status==='Offline'"
            @click="serverAction('stop')"
          >
            <span class="pi pi-power-off"></span>
            Stop
          </button>
          <button
            class="bg-blue-500 hover:bg-blue-600 text-white font-semibold py-2 px-6 rounded transition disabled:opacity-50 flex items-center gap-2"
            :disabled="loading || status==='Offline'"
            @click="serverAction('restart')"
          >
            <span class="pi pi-refresh"></span>
            Restart
          </button>
        </div>
        <div v-if="error" class="text-red-400 mb-4">{{ error }}</div>
        <div class="mt-6 text-white space-y-2">
          <div>
            <span class="font-semibold">Server IP: </span>
            <span class="text-lg">{{ serverIp }}</span>
          </div>
          <div>
            <span class="font-semibold">Server version: </span>
            <span class="text-lg">{{ version }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>