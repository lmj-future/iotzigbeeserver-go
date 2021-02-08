package crc

import (
	"encoding/hex"
	"strconv"

	"github.com/h3c/iotzigbeeserver-go/globalconstant/globallogger"
)

//CRC CRC
func CRC(msg []byte) string {
	var temp = 0
	var crc = 0xffff
	for i := 0; i < len(msg); i++ {
		crc ^= int(msg[i])
		for j := 0; j < 8; j++ {
			temp = 1 & crc
			crc >>= 1
			if temp == 1 {
				crc ^= 0xa001
			}
		}
	}
	crc ^= 0xffff
	tempcrcstr := "0000" + strconv.FormatInt(int64(crc), 16)
	strcrc := tempcrcstr[len(tempcrcstr)-4:]

	return strcrc
}

//Check Check
func Check(msg []byte) bool {
	var crc = CRC(msg[0 : len(msg)-2])
	if crc == hex.EncodeToString(msg[len(msg)-2:]) { //byte转hex字符
		// globallogger.Log.Infoln("CRC check success !")
		return true
	}
	globallogger.Log.Warnln("CRC check failed ! crc=", crc)
	return false
}
