package player

import (
	"bytes"
	"encoding/binary"
)

// CreateWavHeader generates a 44-byte WAV header for PCM 16bit, mono, 16kHz
func CreateWavHeader(dataLen int) []byte {
	var buf bytes.Buffer

	// WAV Header values
	sampleRate := 16000
	bitsPerSample := 16
	numChannels := 1

	byteRate := sampleRate * numChannels * bitsPerSample / 8
	blockAlign := numChannels * bitsPerSample / 8
	subchunk2Size := dataLen
	chunkSize := 36 + subchunk2Size

	// Write the WAV header
	buf.WriteString("RIFF")
	binary.Write(&buf, binary.LittleEndian, uint32(chunkSize))
	buf.WriteString("WAVE")
	buf.WriteString("fmt ")
	binary.Write(&buf, binary.LittleEndian, uint32(16))            // Subchunk1Size
	binary.Write(&buf, binary.LittleEndian, uint16(1))             // AudioFormat (1 = PCM)
	binary.Write(&buf, binary.LittleEndian, uint16(numChannels))   // NumChannels
	binary.Write(&buf, binary.LittleEndian, uint32(sampleRate))    // SampleRate
	binary.Write(&buf, binary.LittleEndian, uint32(byteRate))      // ByteRate
	binary.Write(&buf, binary.LittleEndian, uint16(blockAlign))    // BlockAlign
	binary.Write(&buf, binary.LittleEndian, uint16(bitsPerSample)) // BitsPerSample
	buf.WriteString("data")
	binary.Write(&buf, binary.LittleEndian, uint32(subchunk2Size)) // Subchunk2Size

	return buf.Bytes()
}

// MakeWavChunk creates a full WAV chunk (header + raw PCM data)
func MakeWavChunk(data []byte) []byte {
	header := CreateWavHeader(len(data))
	return append(header, data...)
}
