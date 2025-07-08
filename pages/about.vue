<template>
  <div>
    <h2>🎧 Live Audio Stream</h2>
    <button @click="startAudio">Start Audio</button>
    <audio controls autoplay>
      <source src="http://localhost:4444/live" type="audio/mpeg" />
      Your browser does not support audio element.
    </audio>

  </div>
</template>

<script setup>
import { onMounted, onUnmounted } from 'vue'

let socket
let audioContext
let isStarted = false

const startAudio = () => {
  if (isStarted) return
  isStarted = true

  audioContext = new (window.AudioContext || window.webkitAudioContext)()
  socket = new WebSocket('ws://localhost:7777/ws/audio')
  socket.binaryType = 'arraybuffer'

  socket.onmessage = async (event) => {
    const arrayBuffer = event.data

    try {
      const audioBuffer = await audioContext.decodeAudioData(arrayBuffer.slice(0))
      const bufferSource = audioContext.createBufferSource()
      bufferSource.buffer = audioBuffer
      bufferSource.connect(audioContext.destination)
      bufferSource.start()
    } catch (err) {
      console.error('🎧 Audio decode error:', err)
    }
  }

  socket.onerror = (err) => {
    console.error('🔌 WebSocket error:', err)
  }

  socket.onclose = () => {
    console.log('🛑 WebSocket closed')
  }
}

onUnmounted(() => {
  if (socket) {
    socket.close()
  }
})
</script>
