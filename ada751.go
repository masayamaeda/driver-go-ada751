package ada751

import (
	//"fmt"
	"io"
	"time"

	"github.com/tarm/goserial"
)

// MySensor is ...
type MySensor struct {
	sensor     io.ReadWriteCloser
	password   uint32
	address    uint32
	fingerID   uint16
	confidence uint16
}

// NewMySensor is ...
func NewMySensor(password uint32, address uint32) *MySensor {
	mySensor := &MySensor{sensor: nil, password: password, address: address}
	mySensor.Begin("/dev/ttyAMA0", 57600)
	return mySensor
}

// GetFingerInfo is ...
func (mySensor *MySensor) GetFingerInfo() (uint16, uint16) {

	return mySensor.fingerID, mySensor.confidence
}

// Begin is ...
func (mySensor *MySensor) Begin(name string, baud int) {

	// センサ起動待ち（１秒）
	time.Sleep(1 * time.Second)

	// シリアルポートオープン
	uartConf := &serial.Config{Name: name, Baud: baud}
	mySensor.sensor, _ = serial.OpenPort(uartConf)
}

// VerifyPassword is ...
func (mySensor *MySensor) VerifyPassword() (result bool) {

	packet := []uint8{MessageVerifyPassword,
		uint8(mySensor.password >> 24), uint8(mySensor.password >> 16), uint8(mySensor.password >> 8), uint8(mySensor.password)}

	//送信処理
	mySensor.writePacket(mySensor.address, CommandPacket, uint16(len(packet)+2), packet)

	//シリアル受信
	len := mySensor.getReply(packet, 0)

	if len == 1 && packet[0] == AckPacket && packet[1] == Ok {
		return true
	}

	return false
}

// GetImage is ...
func (mySensor *MySensor) GetImage() (value uint8, result bool) {

	packet := []uint8{MessageGetImage, 0x00}

	mySensor.writePacket(mySensor.address, CommandPacket, uint16(len(packet)+2), packet)

	len := mySensor.getReply(packet, 0)

	if len != 1 && packet[0] != AckPacket {
		return 0, false
	}

	return packet[1], true
}

// Image2Tz is ...
func (mySensor *MySensor) Image2Tz(slot uint8) (value uint8, result bool) {

	packet := []uint8{MessageImage2tz, slot}

	mySensor.writePacket(mySensor.address, CommandPacket, uint16(len(packet)+2), packet)

	len := mySensor.getReply(packet, 0)

	if len != 1 && packet[0] != AckPacket {
		return 0, false
	}

	return packet[1], true
}

// CreateModel is ...
func (mySensor *MySensor) CreateModel() (value uint8, result bool) {

	packet := []uint8{MessageRegModel, 0x00}

	mySensor.writePacket(mySensor.address, CommandPacket, uint16(len(packet)+2), packet)

	len := mySensor.getReply(packet[:], 0)

	if len != 1 && packet[0] != AckPacket {
		return 0, false
	}

	return packet[1], true
}

// StoreModel is ...
func (mySensor *MySensor) StoreModel(id uint16) (value uint8, result bool) {

	packet := []uint8{MessageStore, 0x01, uint8(id >> 8), uint8(id & 0xFF)}

	mySensor.writePacket(mySensor.address, CommandPacket, uint16(len(packet)+2), packet)

	len := mySensor.getReply(packet, 0)

	if len != 1 && packet[0] != AckPacket {
		return 0, false
	}

	return packet[1], true
}

// LoadModel is ...
func (mySensor *MySensor) LoadModel(id uint16) (value uint8, result bool) {

	packet := []uint8{MessageLoad, 0x01, uint8(id >> 8), uint8(id & 0xFF)}

	mySensor.writePacket(mySensor.address, CommandPacket, uint16(len(packet)+2), packet)

	len := mySensor.getReply(packet, 0)

	if len != 1 && packet[0] != AckPacket {
		return 0, false
	}

	return packet[1], true
}

// GetModel is ...
func (mySensor *MySensor) GetModel() (value uint8, result bool) {

	packet := []uint8{MessageUpload, 0x01}

	mySensor.writePacket(mySensor.address, CommandPacket, uint16(len(packet)+2), packet)

	len := mySensor.getReply(packet, 0)

	if len != 1 && packet[0] != AckPacket {
		return 0, false
	}

	return packet[1], true
}

// DeleteModel is ...
func (mySensor *MySensor) DeleteModel(id uint16) (value uint8, result bool) {

	packet := []uint8{MessageDelete, uint8(id >> 8), uint8(id & 0xFF), 0x00, 0x01}

	mySensor.writePacket(mySensor.address, CommandPacket, uint16(len(packet)+2), packet)

	len := mySensor.getReply(packet, 0)

	if len != 1 && packet[0] != AckPacket {
		return 0, false
	}

	return packet[1], true
}

// EmptyDatabase is ...
func (mySensor *MySensor) EmptyDatabase() (value uint8, result bool) {

	packet := []uint8{MessageEmpty}

	mySensor.writePacket(mySensor.address, CommandPacket, uint16(len(packet)+2), packet)

	len := mySensor.getReply(packet, 0)

	if len != 1 && packet[0] != AckPacket {
		return 0, false
	}

	return packet[1], true
}

// FingerFastSearch is ...
func (mySensor *MySensor) FingerFastSearch() (value uint8, result bool) {

	mySensor.fingerID = 0xFFFF
	mySensor.confidence = 0xFFFF

	//speed search of slot #1 starting at page 0x0000 and page #0x00A3
	packet := []uint8{MessageHiSpeedSearch, 0x01, 0x00, 0x00, 0x00, 0xA3}

	mySensor.writePacket(mySensor.address, CommandPacket, uint16(len(packet)+2), packet)

	len := mySensor.getReply(packet, 0)

	if len != 1 && packet[0] != AckPacket {
		return 0, false
	}

	mySensor.fingerID = uint16(packet[2])
	mySensor.fingerID <<= 8
	mySensor.fingerID |= uint16(packet[3])

	mySensor.confidence = uint16(packet[4])
	mySensor.confidence <<= 8
	mySensor.confidence |= uint16(packet[5])

	return packet[1], true
}

//
func (mySensor *MySensor) GetTemplateCount() (value uint8, result bool, templateCount uint16) {

	// get number of templates in memory
	packet := []uint8{MessageTemplateCount, 0x00, 0x00, 0x00}

	mySensor.writePacket(mySensor.address, CommandPacket, uint16(len(packet)), packet)

	len := mySensor.getReply(packet, 0)

	if len != 1 && packet[0] != AckPacket {
		return 0, false, templateCount
	}

	templateCount = uint16(packet[2])
	templateCount <<= 8
	templateCount |= uint16(packet[3])

	return packet[1], true, templateCount
}

// writePacket is ...
func (mySensor *MySensor) writePacket(addr uint32, packettype uint8,
	len uint16, packet []uint8) {

	mySensor.sensor.Write([]byte{byte(StartCode >> 8)})
	mySensor.sensor.Write([]byte{byte(StartCode & 0xFF)})
	mySensor.sensor.Write([]byte{byte(addr >> 24)})
	mySensor.sensor.Write([]byte{byte(addr >> 16)})
	mySensor.sensor.Write([]byte{byte(addr >> 8)})
	mySensor.sensor.Write([]byte{byte(addr)})
	mySensor.sensor.Write([]byte{byte(packettype)})
	mySensor.sensor.Write([]byte{byte(len >> 8)})
	mySensor.sensor.Write([]byte{byte(len)})

	sum := uint16(uint8(len>>8) + uint8(len&0xFF) + packettype)
	var i uint16
	for i = 0; i < len-2; i++ {
		mySensor.sensor.Write([]byte{packet[i]})
		sum += uint16(packet[i])
	}
	mySensor.sensor.Write([]byte{byte(sum >> 8)})
	mySensor.sensor.Write([]byte{byte(sum & 0xFF)})
}

// getReply is ...
func (mySensor *MySensor) getReply(packet []uint8, timeout uint16) (len uint8) {
	reply := make([]byte, 128)
	var idx uint8
	var timer uint16

	maxtime := DefaultTimeout
	if timeout > 0 {
		maxtime = timeout
	}

	idx = 0

	for {
		// 1バイトづつ受信する
		buf := make([]byte, 1)
		for {
			_, err := mySensor.sensor.Read(buf)
			if err == nil {
				break
			}
			timer++
			if timer >= maxtime {
				return Timeout
			}
		}

		if idx == 0 && buf[0] != uint8(StartCode>>8) {
			continue
		}
		reply[idx] = buf[0]
		idx++

		// check packet!
		if idx >= 9 {
			if reply[0] != uint8(StartCode>>8) || reply[1] != uint8(StartCode&0xFF) {
				return BadPacket
			}
			packettype := reply[6]

			len := uint16(reply[7])
			len <<= 8
			len |= uint16(reply[8])
			len -= 2

			if uint16(idx) <= (len + 10) {
				continue
			}
			packet[0] = packettype
			var i uint16
			for i = 0; i < len; i++ {
				packet[1+i] = reply[9+i]
			}
			return uint8(len)
		}
	}
}
