# การแก้ไขปัญหาเพลงไม่เล่น (MP3)

## สาเหตุหลัก
เพลงไม่เล่นเพราะ **ESP32 ไม่ได้เชื่อมต่อ** กับเซิร์ฟเวอร์ หรือ **ESP32 ไม่รองรับ MP3 decoder**

## วิธีตรวจสอบ

### 1. ตรวจสอบ Log
เมื่อรันโปรแกรม ให้ดู log ต่อไปนี้:

```
⚠️ ESP32 not connected - waiting for connection on port 5555
⚠️ ESP32 not connected - audio will not play
⚠️ ESP32 not connected - audio chunks not being sent
```

### 2. ตรวจสอบการเชื่อมต่อ ESP32
- ESP32 ต้องเชื่อมต่อกับ TCP port 5555
- เมื่อเชื่อมต่อสำเร็จ จะเห็น log:
```
✅ ESP32 connected via TCP
🎵 Audio output ready - ESP32 can now receive audio data
✅ ESP32 connected - audio ready
```

## วิธีแก้ไข

### 1. ตรวจสอบ ESP32
- ตรวจสอบว่า ESP32 กำลังรันและเชื่อมต่อกับ WiFi
- ตรวจสอบว่า ESP32 เชื่อมต่อกับ IP address และ port 5555 ของเซิร์ฟเวอร์
- **สำคัญ**: ESP32 ต้องมี MP3 decoder library (เช่น ESP32-audioI2S)

### 2. ตรวจสอบ Network
- ตรวจสอบว่า ESP32 และเซิร์ฟเวอร์อยู่ในเครือข่ายเดียวกัน
- ตรวจสอบ Firewall ว่าไม่ได้บล็อก port 5555

### 3. ทดสอบการเชื่อมต่อ
```bash
# ทดสอบจาก ESP32 ไปยังเซิร์ฟเวอร์
telnet [SERVER_IP] 5555
```

### 4. ตรวจสอบ ESP32 Code
ตรวจสอบว่า ESP32 code มีการเชื่อมต่อ TCP ไปยังเซิร์ฟเวอร์และรองรับ MP3:

```cpp
// ตัวอย่าง ESP32 code สำหรับ MP3
#include "AudioFileSourceTCPStream.h"
#include "AudioGeneratorMP3.h"
#include "AudioOutputI2S.h"

AudioGeneratorMP3 *mp3;
AudioFileSourceTCPStream *file;
AudioOutputI2S *out;

void setup() {
    // เชื่อมต่อ WiFi และ TCP
    WiFiClient client;
    if (client.connect(serverIP, 5555)) {
        Serial.println("Connected to server");
        
        file = new AudioFileSourceTCPStream(client);
        out = new AudioOutputI2S();
        mp3 = new AudioGeneratorMP3();
        mp3->begin(file, out);
    }
}

void loop() {
    if (mp3->isRunning()) {
        mp3->loop();
    }
}
```

## การ Debug เพิ่มเติม

### 1. เปิด Debug Mode
ใน `player/player.go` เปิด comment บรรทัดนี้:
```go
log.Printf("📤 Sending chunk [%d:%d] size: %d", offset-ChunkSize, offset, len(chunk))
```

### 2. ตรวจสอบ MP3 Buffer
ตรวจสอบว่าไฟล์ MP3 ถูกโหลดสำเร็จ:
```
📦 Loaded output/song1.mp3 (1234567 bytes)
```

### 3. ตรวจสอบ WebSocket
ตรวจสอบว่า WebSocket control ทำงาน:
```
🎧 Control WebSocket connected
📥 Command received: play
▶️ Start playing from offset: 0
```

### 4. ตรวจสอบไมค์
ตรวจสอบว่าไมค์ส่ง MP3 ได้:
```
🌐 Mic WebSocket connected
📤 Sending 256 bytes to ESP32
✅ Successfully sent 256 bytes to ESP32
```

## สรุป
ปัญหาหลักคือ ESP32 ไม่ได้เชื่อมต่อ หรือ ESP32 ไม่มี MP3 decoder library ต้องแก้ไขการเชื่อมต่อ ESP32 และติดตั้ง MP3 library ก่อน 