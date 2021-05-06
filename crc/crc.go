package crc

import (
	"encoding/hex"
	"strconv"
	"strings"

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
	var builder strings.Builder
	builder.WriteString("0000")
	builder.WriteString(strconv.FormatInt(int64(crc), 16))

	return strings.Repeat(builder.String()[builder.Len()-4:], 1)
}

//Check Check
func Check(msg []byte) bool {
	if CRC(append(msg[:0:0], msg[:len(msg)-2]...)) == hex.EncodeToString(append(msg[:0:0], msg[len(msg)-2:]...)) { //byte转hex字符
		return true
	}
	globallogger.Log.Warnln("CRC check failed ! crc=", CRC(append(msg[:0:0], msg[:len(msg)-2]...)))
	return false
}
