package frame

import (
	"errors"

	"github.com/dyrkin/bin"
)

// Configuration Configuration
type Configuration struct {
	transactionIDProvider            func() uint8
	FrameType                        Type
	FrameTypeConfigured              bool
	ManufacturerCode                 uint16
	ManufacturerCodeConfigured       bool
	Direction                        Direction
	DirectionConfigured              bool
	DisableDefaultResponse           bool
	DisableDefaultResponseConfigured bool
	CommandID                        uint8
	CommandIDConfigured              bool
	Command                          interface{}
	CommandConfigured                bool
}

// Builder Builder
type Builder interface {
	BuildFrameType(frameType Type) Builder
	BuildManufacturerCode(manufacturerCode uint16) Builder
	BuildDirection(direction Direction) Builder
	BuildDisableDefaultResponse(disableDefaultResponse bool) Builder
	BuildCommandID(commandID uint8) Builder
	BuildCommand(command interface{}) Builder
	Build() (*Frame, error)
}

var defaultTransactionIDProvider func() uint8

// New New
func New() Builder {
	return &Configuration{transactionIDProvider: defaultTransactionIDProvider}
}

// BuildIDGenerator BuildIDGenerator
func (f *Configuration) BuildIDGenerator(transactionIDProvider func() uint8) Builder {
	f.transactionIDProvider = transactionIDProvider
	return f
}

// BuildFrameType BuildFrameType
func (f *Configuration) BuildFrameType(frameType Type) Builder {
	f.FrameType = frameType
	f.FrameTypeConfigured = true
	return f
}

// BuildManufacturerCode BuildManufacturerCode
func (f *Configuration) BuildManufacturerCode(manufacturerCode uint16) Builder {
	f.ManufacturerCode = manufacturerCode
	f.ManufacturerCodeConfigured = true
	return f
}

// BuildDirection BuildDirection
func (f *Configuration) BuildDirection(direction Direction) Builder {
	f.Direction = direction
	f.DirectionConfigured = true
	return f
}

// BuildDisableDefaultResponse BuildDisableDefaultResponse
func (f *Configuration) BuildDisableDefaultResponse(disableDefaultResponse bool) Builder {
	f.DisableDefaultResponse = disableDefaultResponse
	f.DisableDefaultResponseConfigured = true
	return f
}

// BuildCommandID BuildCommandID
func (f *Configuration) BuildCommandID(commandID uint8) Builder {
	f.CommandID = commandID
	f.CommandIDConfigured = true
	return f
}

// BuildCommand BuildCommand
func (f *Configuration) BuildCommand(command interface{}) Builder {
	f.Command = command
	f.CommandConfigured = true
	return f
}

// Build Build
func (f *Configuration) Build() (*Frame, error) {
	if err := f.validateConfiguration(); err != nil {
		return nil, err
	}
	frame := &Frame{}
	frame.FrameControl = &Control{}
	frame.FrameControl.FrameType = f.FrameType
	frame.FrameControl.ManufacturerSpecific = flag(f.ManufacturerCodeConfigured)
	frame.FrameControl.Direction = f.Direction
	frame.FrameControl.DisableDefaultResponse = flag(f.DisableDefaultResponse)
	frame.ManufacturerCode = f.ManufacturerCode
	frame.TransactionSequenceNumber = f.transactionIDProvider()
	frame.CommandIdentifier = f.CommandID
	if f.CommandConfigured {
		frame.Payload = bin.Encode(f.Command)
	} else {
		frame.Payload = make([]uint8, 0)
	}
	return frame, nil
}

func (f *Configuration) validateConfiguration() error {
	if !f.FrameTypeConfigured {
		return errors.New("frame type must be set")
	}
	if !f.CommandIDConfigured {
		return errors.New("command id must be set")
	}
	if !f.DirectionConfigured {
		return errors.New("direction must be set")
	}
	return nil
}

func flag(flag bool) uint8 {
	if flag {
		return 1
	}
	return 0
}

// MakeDefaultTransactionIDProvider MakeDefaultTransactionIDProvider
func MakeDefaultTransactionIDProvider() func() uint8 {
	transactionID := uint8(1)
	return func() uint8 {
		transactionID = transactionID + 1
		if transactionID > 255 {
			transactionID = 1
		}
		return transactionID
	}
}

func init() {
	defaultTransactionIDProvider = MakeDefaultTransactionIDProvider()
}
