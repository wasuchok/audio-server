<template>
  <audio ref="audio" controls autoplay />
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount } from 'vue'
let ws
const audio = ref(null)
let mediaSource, sourceBuffer

onMounted(() => {
  mediaSource = new MediaSource()
  audio.value.src = URL.createObjectURL(mediaSource)
  ws = new WebSocket('ws://localhost:7777/ws/stream-client')
  ws.binaryType = 'arraybuffer'
  mediaSource.addEventListener('sourceopen', () => {
    sourceBuffer = mediaSource.addSourceBuffer('audio/mpeg')
    ws.onmessage = (event) => {
      if (sourceBuffer && !sourceBuffer.updating) {
        sourceBuffer.appendBuffer(new Uint8Array(event.data))
      }
    }
  })
})

onBeforeUnmount(() => {
  ws && ws.close()
})
</script>