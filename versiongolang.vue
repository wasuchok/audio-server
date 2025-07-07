<template>
  <div class="container">
    <h1>🎧 ESP32 Audio Control</h1>

    <div class="btn-group">
      <button @click="sendCommand('play')" class="play-button">▶️ เล่นเสียง</button>
      <button @click="sendCommand('pause')" class="pause-button">⏸ หยุดชั่วคราว</button>
      <button @click="sendCommand('resume')" class="resume-button">⏯ เล่นต่อ</button>
      <button @click="sendCommand('stop')" class="stop-button">⏹ หยุดทั้งหมด</button>
      <button @click="startMic" class="mic-button">🎤 เริ่มพูด</button>
      <button @click="stopMic" class="stop-mic-button">🚩 หยุด</button>
    </div>

    <div class="volume-control">
      <label for="volume">🔉 Mic Volume: {{ volume }}%</label>
      <input
        id="volume"
        type="range"
        min="0"
        max="100"
        step="5"
        v-model="volume"
      />
    </div>

    <p v-if="connectedControl">🟢 Control WebSocket Connected</p>
    <p v-else>🔴 Control WebSocket Disconnected</p>

    <p v-if="connectedMic">🟢 Mic WebSocket Connected</p>
    <p v-else>🔴 Mic WebSocket Disconnected</p>
  </div>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue'

const volume = ref(100)

const connectedControl = ref(false)
const connectedMic = ref(false)

let socketControl
let socketMic
let processor
let audioContext
let stream

function float32ToInt16(buffer, factor) {
  const l = buffer.length
  const result = new Int16Array(l)
  for (let i = 0; i < l; i++) {
    const s = Math.max(-1, Math.min(1, buffer[i]))
    result[i] = s < 0 ? s * 0x8000 * factor : s * 0x7FFF * factor
  }
  return result
}

const startMic = async () => {
  try {
    stream = await navigator.mediaDevices.getUserMedia({ audio: true })

    socketMic = new WebSocket('ws://103.52.109.49:7777/ws/mic')
    socketMic.binaryType = 'arraybuffer'

    socketMic.onopen = () => {
      connectedMic.value = true
      console.log('🎙️ Mic WebSocket connected')
    }

    socketMic.onclose = () => {
      connectedMic.value = false
      console.log('❌ Mic WebSocket disconnected')
    }

    audioContext = new (window.AudioContext || window.webkitAudioContext)()
    const source = audioContext.createMediaStreamSource(stream)

    processor = audioContext.createScriptProcessor(512, 1, 1)
    processor.onaudioprocess = (e) => {
      const input = e.inputBuffer.getChannelData(0)
      const pcm = float32ToInt16(input, volume.value / 100)
      if (socketMic && socketMic.readyState === WebSocket.OPEN) {
        socketMic.send(pcm.buffer)
      }
    }

    source.connect(processor)
    processor.connect(audioContext.destination)

    console.log('🎤 Microphone started')
  } catch (err) {
    console.error('❌ Failed to access microphone:', err)
    alert('ไม่สามารถเข้าถึงไมโครโฟนได้')
  }
}

const stopMic = () => {
  if (processor) processor.disconnect()
  if (audioContext) audioContext.close()
  if (stream) stream.getTracks().forEach(t => t.stop())
  if (socketMic) socketMic.close()

  console.log('🎤 Microphone stopped')
}

onMounted(() => {
  socketControl = new WebSocket('ws://103.52.109.49:7777/ws/control') // ✅ สำหรับคำสั่ง play/pause/resume/stop

  socketControl.onopen = () => {
    connectedControl.value = true
    console.log('✅ Control WebSocket connected')
  }

  socketControl.onclose = () => {
    connectedControl.value = false
    console.log('❌ Control WebSocket disconnected')
  }

  socketControl.onerror = (err) => {
    console.error('WebSocket error:', err)
  }
})

const sendCommand = (cmd) => {
  if (socketControl && socketControl.readyState === WebSocket.OPEN) {
    socketControl.send(cmd)
    console.log(`📤 Sent: ${cmd}`)
  } else {
    alert('❌ Control WebSocket ยังไม่เชื่อมต่อ')
  }
}

watch(volume, (newVal) => {
  sendCommand(`mic-volume-${newVal}`)
})
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
  margin-bottom: 1.5rem;
}

.volume-control {
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

.mic-button {
  background-color: #8b5cf6;
}
.mic-button:hover {
  background-color: #7c3aed;
}

.stop-mic-button {
  background-color: #9ca3af;
}
.stop-mic-button:hover {
  background-color: #6b7280;
}

input[type="range"] {
  width: 200px;
  margin-top: 0.5rem;
}
</style>
