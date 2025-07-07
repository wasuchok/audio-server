<template>
  <div class="container">
    <h1>🎧 ESP32 Audio Control</h1>

    <div class="btn-group">
      <button @click="sendCommand('play')" class="play-button">▶️ เล่นเสียง</button>
      <button @click="sendCommand('pause')" class="pause-button">⏸ หยุดชั่วคราว</button>
      <button @click="sendCommand('resume')" class="resume-button">⏯ เล่นต่อ</button>
      <button @click="sendCommand('stop')" class="stop-button">⏹ หยุดทั้งหมด</button>
    </div>

    <p v-if="connected">🟢 WebSocket Connected</p>
    <p v-else>🔴 Not Connected</p>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'

const connected = ref(false)
let socket

onMounted(() => {
  socket = new WebSocket('ws://localhost:8080')

  socket.addEventListener('open', () => {
    connected.value = true
    console.log('✅ WebSocket connected')
  })

  socket.addEventListener('close', () => {
    connected.value = false
    console.log('❌ WebSocket disconnected')
  })

  socket.addEventListener('error', (err) => {
    console.error('WebSocket error:', err)
  })
})

const sendCommand = (cmd) => {
  if (socket && socket.readyState === WebSocket.OPEN) {
    socket.send(cmd)
    console.log(`📤 Sent: ${cmd}`)
  } else {
    alert('WebSocket ยังไม่เชื่อมต่อ')
  }
}
</script>

<style scoped>
.container {
  padding: 2rem;
  text-align: center;
}

.btn-group {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 1rem;
  margin-bottom: 1rem;
}

button {
  font-size: 1rem;
  padding: 0.8rem 1.6rem;
  border: none;
  border-radius: 12px;
  cursor: pointer;
  color: white;
}

.play-button {
  background-color: #4ade80;
}
.play-button:hover {
  background-color: #22c55e;
}

.pause-button {
  background-color: #facc15;
  color: black;
}
.pause-button:hover {
  background-color: #eab308;
}

.resume-button {
  background-color: #60a5fa;
}
.resume-button:hover {
  background-color: #3b82f6;
}

.stop-button {
  background-color: #ef4444;
}
.stop-button:hover {
  background-color: #dc2626;
}
</style>
