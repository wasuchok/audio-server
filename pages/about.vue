<template>
  <div>
    <p>🔊 Streaming audio...</p>
    <audio ref="player" controls autoplay muted></audio>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'

onMounted(() => {
  let audioCtx = new (window.AudioContext || window.webkitAudioContext)()
  const sampleRate = 44100
  const bufferQueue = []

  let isPlaying = false

  const socket = new WebSocket('ws://103.52.109.49:7777/ws/stream-client')
 
  socket.binaryType = 'arraybuffer'

  socket.onmessage = (event) => {
    const int16Array = new Int16Array(event.data)
    const float32Array = new Float32Array(int16Array.length)

    for (let i = 0; i < int16Array.length; i++) {
      float32Array[i] = int16Array[i] / 32768
    }

    const audioBuffer = audioCtx.createBuffer(1, float32Array.length, sampleRate)
    audioBuffer.copyToChannel(float32Array, 0)

    bufferQueue.push(audioBuffer)
    if (!isPlaying) playFromQueue()
  }

  function playFromQueue() {
    if (bufferQueue.length === 0) {
      isPlaying = false
      return
    }

    isPlaying = true
    const buffer = bufferQueue.shift()
    const source = audioCtx.createBufferSource()
    source.buffer = buffer
    source.connect(audioCtx.destination)
    source.onended = playFromQueue
    source.start()
  }
})
</script>
