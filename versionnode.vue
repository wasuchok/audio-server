<template>
  <div class="container">
    <h1>🎧 ESP32 Audio Control</h1>

    <div class="btn-group">
      <button @click="sendCommand('play')" class="play-button">▶️ เล่นเสียง</button>
      <button @click="sendCommand('pause')" class="pause-button">⏸ หยุดชั่วคราว</button>
      <button @click="sendCommand('resume')" class="resume-button">⏯ เล่นต่อ</button>
      <button @click="sendCommand('stop')" class="stop-button">⏹ หยุดทั้งหมด</button>
      <button @click="start">🎤 เริ่มพูด</button>
      <button @click="stop">🚩 หยุด</button>
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

    <p v-if="connected">🟢 WebSocket Connected</p>
    <p v-if="connectedMic">🟢 WebSocket Mic Connected</p>
    <p v-else>🔴 Not Connected</p>
  </div>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue'

const connected = ref(false)
const connectedMic = ref(false)
const volume = ref(100) // 100%

let socket
let socketMic
let processor;
let audioContext;
let stream;

const start = async () => {
  sendCommand('pause')
  stream = await navigator.mediaDevices.getUserMedia({ audio: true });

  socketMic = new WebSocket('ws://localhost:7777');
  socketMic.binaryType = 'arraybuffer';

audioContext = new (window.AudioContext || window.webkitAudioContext)({
  latencyHint: 'interactive', // หรือ 'playback' / 'balanced' / number (seconds)
});
  const source = audioContext.createMediaStreamSource(stream);

  processor = audioContext.createScriptProcessor(512, 1, 1); 
  processor.onaudioprocess = (e) => {
    const input = e.inputBuffer.getChannelData(0);
    const pcm = float32ToInt16(input);
    if (socketMic.readyState === WebSocket.OPEN) {
      socketMic.send(pcm.buffer);
    }
  };

  source.connect(processor);
  processor.connect(audioContext.destination);
};

const stop = () => {
  if (processor) processor.disconnect();
  if (audioContext) audioContext.close();
  if (stream) stream.getTracks().forEach(t => t.stop());
  if (socketMic) socketMic.close();
};

function float32ToInt16(buffer) {
  const l = buffer.length;
  const result = new Int16Array(l);
  for (let i = 0; i < l; i++) {
    const s = Math.max(-1, Math.min(1, buffer[i]));
    result[i] = s < 0 ? s * 0x8000 : s * 0x7FFF;
  }
  return result;
}

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

// 📡 ส่งคำสั่ง mic-volume ทุกครั้งที่ volume เปลี่ยน
watch(volume, (newVal) => {
  sendCommand(`mic-volume-${newVal}`);
});
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

input[type="range"] {
  width: 200px;
  margin-top: 0.5rem;
}
</style>
