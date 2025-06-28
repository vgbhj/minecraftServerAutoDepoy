<script setup>
import { ref, onMounted, onBeforeUnmount } from 'vue'

const logs = ref('')
const error = ref('')
let ws

onMounted(() => {

// ws = new WebSocket("ws://localhost:8080/api/v1/console/stream");
  ws = new WebSocket(`ws://${window.location.host}/api/v1/console/stream`)
  ws.onmessage = (event) => {
    logs.value += event.data
  }
  ws.onerror = () => {
    error.value = 'Server is offline or console unavailable'
  }
  ws.onclose = () => {
    if (!logs.value) error.value = 'Server is offline or console unavailable'
  }
})
onBeforeUnmount(() => {
  if (ws) ws.close()
})
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
        style="height: calc(100vh - 7rem); width: calc(100vw - 14rem); overflow: hidden;"
      >
        <template v-if="!error">
          <div style="width:100%;height:100%;overflow:hidden;">
            <pre
              class="whitespace-pre-wrap break-words m-0 custom-scrollbar"
              style="overflow-y: scroll; max-width:100vw; max-height:100%;"
            >{{ logs.trim() ? logs : 'Waiting for console output...' }}</pre>
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
.custom-scrollbar {
  scrollbar-width: thin;
  scrollbar-color: transparent transparent;
}
.custom-scrollbar::-webkit-scrollbar {
  width: 8px;
  background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: transparent;
}
</style>