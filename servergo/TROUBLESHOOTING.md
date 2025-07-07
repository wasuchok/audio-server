# การแก้ไขปัญหาเพลงไม่เล่น

## สาเหตุหลัก
เพลงไม่เล่นเพราะ **ESP32 ไม่ได้เชื่อมต่อ** กับเซิร์ฟเวอร์

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

### 2. ตรวจสอบ Network
- ตรวจสอบว่า ESP32 และเซิร์ฟเวอร์อยู่ในเครือข่ายเดียวกัน
- ตรวจสอบ Firewall ว่าไม่ได้บล็อก port 5555

### 3. ทดสอบการเชื่อมต่อ
```bash
# ทดสอบจาก ESP32 ไปยังเซิร์ฟเวอร์
telnet [SERVER_IP] 5555
```

### 4. ตรวจสอบ ESP32 Code
ตรวจสอบว่า ESP32 code มีการเชื่อมต่อ TCP ไปยังเซิร์ฟเวอร์:

```cpp
// ตัวอย่าง ESP32 code
WiFiClient client;
if (client.connect(serverIP, 5555)) {
    Serial.println("Connected to server");
    // รับข้อมูล audio
    while (client.available()) {
        byte audioData = client.read();
        // ส่งไปยัง DAC หรือ speaker
    }
}
```

## การ Debug เพิ่มเติม

### 1. เปิด Debug Mode
ใน `player/player.go` เปิด comment บรรทัดนี้:
```go
log.Printf("📤 Sending chunk [%d:%d]", offset, end)
```

### 2. ตรวจสอบ Audio Buffer
ตรวจสอบว่าไฟล์ WAV ถูกโหลดสำเร็จ:
```
📦 Loaded output/Always.wav (1234567 bytes)
```

### 3. ตรวจสอบ WebSocket
ตรวจสอบว่า WebSocket control ทำงาน:
```
🎧 Control WebSocket connected
📥 Command received: play
▶️ Start playing from offset: 0
```

## สรุป
ปัญหาหลักคือ ESP32 ไม่ได้เชื่อมต่อ ทำให้ไม่มีอุปกรณ์สำหรับเล่นเสียง ต้องแก้ไขการเชื่อมต่อ ESP32 ก่อน 