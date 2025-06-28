<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'

// Ограничения и типы для популярных server.properties (можно расширить)
const propertyMeta = {
  "allow-flight": { type: "boolean" },
  "allow-nether": { type: "boolean" },
  "broadcast-console-to-ops": { type: "boolean" },
  "broadcast-rcon-to-ops": { type: "boolean" },
  "difficulty": { type: "select", options: ["peaceful", "easy", "normal", "hard"] },
  "enable-command-block": { type: "boolean" },
  "enable-jmx-monitoring": { type: "boolean" },
  "enable-query": { type: "boolean" },
  "enable-rcon": { type: "boolean" },
  "enable-status": { type: "boolean" },
  "enforce-whitelist": { type: "boolean" },
  "entity-broadcast-range-percentage": { type: "number", min: 10, max: 1000 },
  "force-gamemode": { type: "boolean" },
  "function-permission-level": { type: "number", min: 1, max: 4 },
  "gamemode": { type: "select", options: ["survival", "creative", "adventure", "spectator"] },
  "generate-structures": { type: "boolean" },
  "generator-settings": { type: "text" },
  "hardcore": { type: "boolean" },
  "level-name": { type: "text", maxlength: 32 },
  "level-seed": { type: "text", maxlength: 32 },
  "level-type": { type: "select", options: ["default", "flat", "largebiomes", "amplified", "buffet", "customized"] },
  "max-chained-neighbor-updates": { type: "number", min: 0, max: 1000000 },
  "max-players": { type: "number", min: 1, max: 1000 },
  "max-tick-time": { type: "number", min: -1, max: 3600000 },
  "max-world-size": { type: "number", min: 1, max: 29999984 },
  "motd": { type: "text", maxlength: 59 },
  "network-compression-threshold": { type: "number", min: -1, max: 256 },
  "online-mode": { type: "boolean" },
  "op-permission-level": { type: "number", min: 1, max: 4 },
  "player-idle-timeout": { type: "number", min: 0, max: 32767 },
  "prevent-proxy-connections": { type: "boolean" },
  "pvp": { type: "boolean" },
  "query.port": { type: "number", min: 1, max: 65535 },
  "rate-limit": { type: "number", min: 0, max: 2147483647 },
  "rcon.password": { type: "text" },
  "rcon.port": { type: "number", min: 1, max: 65535 },
  "require-resource-pack": { type: "boolean" },
  "resource-pack": { type: "text" },
  "resource-pack-prompt": { type: "text" },
  "resource-pack-sha1": { type: "text", maxlength: 40 },
  "server-ip": { type: "text" },
  "server-port": { type: "number", min: 1, max: 65535 },
  "simulation-distance": { type: "number", min: 3, max: 32 },
  "spawn-animals": { type: "boolean" },
  "spawn-monsters": { type: "boolean" },
  "spawn-npcs": { type: "boolean" },
  "spawn-protection": { type: "number", min: 0, max: 16 },
  "sync-chunk-writes": { type: "boolean" },
  "text-filtering-config": { type: "text" },
  "use-native-transport": { type: "boolean" },
  "view-distance": { type: "number", min: 2, max: 32 },
  "white-list": { type: "boolean" },
}

const properties = ref({})
const loading = ref(true)
const error = ref('')
const success = ref('')

async function fetchProperties() {
  loading.value = true
  error.value = ''
  try {
    const res = await axios.get('/api/v1/server/properties')
    properties.value = res.data
  } catch (e) {
    error.value = 'Failed to load properties'
  }
  loading.value = false
}

async function saveProperties() {
  loading.value = true
  error.value = ''
  success.value = ''
  try {
    await axios.put('/api/v1/server/properties', properties.value)
    success.value = 'Properties updated successfully'
  } catch (e) {
    error.value = 'Failed to update properties'
  }
  loading.value = false
}

onMounted(fetchProperties)
</script>

<template>
  <div style="margin-left: 14rem;">
    <div class="min-h-screen bg-gray-900 flex items-center justify-center">
      <div class="w-full max-w-5xl bg-gray-800 rounded-xl shadow-lg p-8">
        <h2 class="text-3xl font-bold mb-8 text-white text-center tracking-wide">Server Properties</h2>
        <form @submit.prevent="saveProperties" v-if="!loading">
          <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            <div
              v-for="(value, key) in properties"
              :key="key"
              class="flex flex-col bg-gray-700 rounded-lg px-4 py-3 border border-gray-700 hover:border-green-500 transition"
            >
              <label :for="key" class="font-semibold text-gray-200 capitalize mb-2 break-all">{{ key }}</label>
              <!-- Typed rendering -->
              <template v-if="propertyMeta[key]?.type === 'number'">
                <input
                  :id="key"
                  type="number"
                  v-model="properties[key]"
                  :min="propertyMeta[key].min"
                  :max="propertyMeta[key].max"
                  class="px-3 py-2 rounded bg-gray-800 text-white border border-gray-600 focus:outline-none focus:ring-2 focus:ring-green-500 transition"
                />
              </template>
              <template v-else-if="propertyMeta[key]?.type === 'select'">
                <select
                  :id="key"
                  v-model="properties[key]"
                  class="px-3 py-2 rounded bg-gray-800 text-white border border-gray-600 focus:outline-none focus:ring-2 focus:ring-green-500 transition"
                >
                  <option v-for="opt in propertyMeta[key].options" :key="opt" :value="opt">{{ opt }}</option>
                </select>
              </template>
              <template v-else-if="propertyMeta[key]?.type === 'boolean'">
                <select
                  :id="key"
                  v-model="properties[key]"
                  class="px-3 py-2 rounded bg-gray-800 text-white border border-gray-600 focus:outline-none focus:ring-2 focus:ring-green-500 transition"
                >
                  <option value="true">true</option>
                  <option value="false">false</option>
                </select>
              </template>
              <template v-else>
                <input
                  :id="key"
                  type="text"
                  v-model="properties[key]"
                  :maxlength="propertyMeta[key]?.maxlength"
                  class="px-3 py-2 rounded bg-gray-800 text-white border border-gray-600 focus:outline-none focus:ring-2 focus:ring-green-500 transition"
                />
              </template>
            </div>
          </div>
          <div class="flex gap-4 mt-8 justify-end">
            <button
              type="submit"
              class="bg-green-600 hover:bg-green-700 text-white font-semibold px-8 py-2 rounded transition"
            >
              Save
            </button>
            <button
              type="button"
              class="bg-gray-600 hover:bg-gray-700 text-white font-semibold px-8 py-2 rounded transition"
              @click="fetchProperties"
            >
              Reset
            </button>
          </div>
          <div v-if="success" class="text-green-400 mt-4 text-center">{{ success }}</div>
          <div v-if="error" class="text-red-400 mt-4 text-center">{{ error }}</div>
        </form>
        <div v-else class="text-gray-300 text-center">Loading...</div>
      </div>
    </div>
  </div>
</template>