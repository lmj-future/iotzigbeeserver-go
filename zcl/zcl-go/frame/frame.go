package frame

import "github.com/dyrkin/bin"

// Direction Direction
type Direction uint8

// Direction List
const (
	DirectionClientServer Direction = 0x00
	DirectionServerClient Direction = 0x01
)

// Type Type
type Type uint8

//00表示command是通用的
//01表示command是clusterId特殊的
const (
	FrameTypeGlobal Type = 0x00
	FrameTypeLocal  Type = 0x01
)

// Control Control
type Control struct {
	FrameType              Type      `bits:"0b00000011" bitmask:"start"`
	ManufacturerSpecific   uint8     `bits:"0b00000100"`
	Direction              Direction `bits:"0b00001000"`
	DisableDefaultResponse uint8     `bits:"0b00010000"`
	Reserved               uint8     `bits:"0b11100000" bitmask:"end"`
}

// Frame Frame
type Frame struct {
	FrameControl              *Control
	ManufacturerCode          uint16 `cond:"uint:FrameControl.ManufacturerSpecific==1"`
	TransactionSequenceNumber uint8
	CommandIdentifier         uint8
	Payload                   []uint8
}

// Decode Decode
func Decode(buf []uint8) *Frame {
	frame := &Frame{}
	bin.Decode(buf, frame)
	return frame
}

// Encode Encode
func Encode(frame *Frame) []uint8 {
	return bin.Encode(frame)
}
