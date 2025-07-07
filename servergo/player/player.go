package player

import (
	"log"
	"sync"
	"time"
)

var (
	buffer        []byte
	offset        int
	isPaused      bool
	isPlaying     bool
	ticker        *time.Ticker
	mu            sync.Mutex
	ChunkSize     = 256
	IntervalMs    = 20
	OnFinishTrack func()
	SendChunk     func([]byte)
)

func SetBuffer(b []byte) {
	mu.Lock()
	defer mu.Unlock()

	buffer = b
	offset = 0
	isPaused = false
	isPlaying = false
	log.Println("🎶 WAV buffer set")
}

func Play() {
	mu.Lock()
	offset = 0

	if offset >= len(buffer) {
		log.Println("✅ Playback finished")
		isPlaying = false
		mu.Unlock()
		if OnFinishTrack != nil {
			go OnFinishTrack()
		}
		return
	}

	isPaused = false
	isPlaying = true
	mu.Unlock()

	startTickerLocked()
	log.Println("▶️ Start playing from offset:", offset)
}

func Pause() {
	mu.Lock()
	defer mu.Unlock()

	if !isPlaying || isPaused {
		log.Println("⚠️ Cannot pause: Not playing or already paused")
		return
	}

	isPaused = true
	log.Println("⏸ Paused at offset:", offset)
}

func Resume() {
	mu.Lock()

	if !isPaused {
		log.Println("⚠️ Not paused")
		mu.Unlock()
		return
	}

	isPaused = false
	isPlaying = true

	if ticker != nil {
		ticker.Stop()
	}
	mu.Unlock()

	startTickerLocked()
	log.Println("⏯ Resume playing from offset:", offset)
}

func Stop() {
	mu.Lock()
	defer mu.Unlock()

	if ticker != nil {
		ticker.Stop()
		ticker = nil
	}

	offset = 0
	isPaused = false
	isPlaying = false
	log.Println("⏹ Playback stopped")
}

func startTickerLocked() {
	ticker = time.NewTicker(time.Duration(IntervalMs) * time.Millisecond)

	go func() {
		for range ticker.C {
			mu.Lock()
			if isPaused || offset >= len(buffer) {
				if offset >= len(buffer) {
					log.Println("✅ Playback finished")
				}
				ticker.Stop()
				ticker = nil
				isPlaying = false
				mu.Unlock()
				if offset >= len(buffer) && OnFinishTrack != nil {
					go OnFinishTrack()
				}
				return
			}

			end := offset + ChunkSize
			if end > len(buffer) {
				end = len(buffer)
			}
			chunk := buffer[offset:end]
			offset = end
			mu.Unlock()

			if SendChunk != nil {
				//log.Printf("📤 Sending chunk [%d:%d]", offset, end)
				SendChunk(chunk)
			}
		}
	}()
}

func GetBuffer() []byte {
	mu.Lock()
	defer mu.Unlock()
	return buffer
}

func GetFullBuffer() []byte {
	mu.Lock()
	defer mu.Unlock()
	return buffer
}

func GetOffsetInfo() (int, int) {
	mu.Lock()
	defer mu.Unlock()
	return offset, len(buffer)
}
