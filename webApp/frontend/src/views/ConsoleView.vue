<script setup>
import { ref, onMounted, onBeforeUnmount, nextTick } from 'vue'
import axios from 'axios'

const logs = ref('')
const error = ref('')
let ws

const command = ref('')
const rconResponse = ref('')
const rconError = ref('')

const logsContainer = ref(null)

function scrollToBottom() {
  nextTick(() => {
    if (logsContainer.value) {
      logsContainer.value.scrollTop = logsContainer.value.scrollHeight
    }
  })
}

onMounted(() => {
  const adminPassword = localStorage.getItem('adminPassword')
  ws = new WebSocket(`ws://${window.location.host}/api/v1/console/stream?admin_password=${encodeURIComponent(adminPassword)}`)
  ws.onmessage = (event) => {
    logs.value += event.data
    scrollToBottom()
  }
  ws.onerror = () => {
    error.value = 'Server is offline or console unavailable'
  }
  ws.onclose = () => {
    if (!logs.value) error.value = 'Server is offline or console unavailable'
  }
  scrollToBottom()
})
onBeforeUnmount(() => {
  if (ws) ws.close()
})

async function sendRconCommand() {
  rconError.value = ''
  rconResponse.value = ''
  if (!command.value.trim()) {
    return
  }
  try {
    const res = await axios.post('/api/v1/console/rcon', { command: command.value })
    rconResponse.value = res.data.response
    command.value = ''
  } catch (e) {
    rconError.value = e.response?.data?.error || e.message
  }
}
</script>

<template>
  <div class="min-h-screen bg-gray-900 flex items-center justify-center">
    <div class="w-full bg-gray-800 rounded-none p-0" style="margin-left: 14rem;">
      <h2 class="text-2xl font-bold mb-4 text-white flex items-center gap-2 px-8 pt-8">
        <span class="pi pi-terminal text-green-400"></span>
        Server Console
      </h2>
      <div
        class="bg-black text-green-400 font-mono rounded-none px-8 pb-8 pt-4"
        style="height: calc(100vh - 13rem); width: calc(100vw - 14rem); overflow: hidden;"
      >
        <template v-if="!error">
          <div style="width:100%;height:100%;overflow:hidden;display:flex;flex-direction:column;height:100%;">
            <div
              ref="logsContainer"
              style="flex:1 1 0;overflow-y:auto;"
              class="console-scrollbar"
            >
              <pre
                class="whitespace-pre-wrap break-words m-0"
                style="overflow-y: visible; max-width:100vw; max-height:100%;"
              >{{ logs.trim() ? logs : 'Waiting for console output...' }}</pre>
              <div v-if="rconResponse" class="text-green-400 whitespace-pre-wrap">{{ rconResponse }}</div>
              <div v-if="rconError" class="text-red-400">{{ rconError }}</div>
            </div>
            <form
              @submit.prevent="sendRconCommand"
              class="flex items-center gap-2 mt-2"
              autocomplete="off"
              style="background: transparent; outline: none; box-shadow: none; border: none;"
            >
              <span class="text-green-400 font-mono select-none">$</span>
              <input
                v-model="command"
                type="text"
                placeholder="Enter RCON command"
                class="flex-1 px-3 py-2 rounded bg-black text-green-200 border-0 focus:outline-none font-mono"
                style="border-width:0;"
                autocomplete="off"
                @keyup.enter="sendRconCommand"
              />
            </form>
          </div>
        </template>
        <template v-else>
          <div class="text-red-400 text-center py-16">{{ error }}</div>
        </template>
      </div>
    </div>
  </div>
</template>

<style scoped>
.console-scrollbar {
  scrollbar-width: thin;
  scrollbar-color: transparent transparent;
}
.console-scrollbar::-webkit-scrollbar {
  width: 8px;
  background: transparent;
}
.console-scrollbar::-webkit-scrollbar-thumb {
  background: transparent;
}
</style>