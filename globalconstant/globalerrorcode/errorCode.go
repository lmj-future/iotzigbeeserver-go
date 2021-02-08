package globalerrorcode

import "github.com/h3c/iotzigbeeserver-go/publicstruct"

var errorCode publicstruct.ErrCode

func init() {
	errorCode.ZigbeeErrorSuccess = "00000000"
	errorCode.ZigbeeErrorFailed = "00000001"
	errorCode.ZigbeeErrorModuleNotExist = "00000002"
	errorCode.ZigbeeErrorModuleOffline = "00000003"
	errorCode.ZigbeeErrorAPNotExist = "00000004"
	errorCode.ZigbeeErrorAPOffline = "00000005"
	errorCode.ZigbeeErrorTerminalNotExist = "00000006"
	errorCode.ZigbeeErrorChipNotMatch = "00000007"
	errorCode.ZigbeeErrorAPBusy = "00000010"
	errorCode.ZigbeeErrorTimeout = "00000011"
}

// ErrorCode ErrorCode
var ErrorCode = &errorCode
