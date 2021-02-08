package zcl

import (
	"encoding/hex"
	"fmt"

	"github.com/dyrkin/bin"
	"github.com/dyrkin/znp-go"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globallogger"
	"github.com/h3c/iotzigbeeserver-go/zcl/zcl-go/cluster"
	"github.com/h3c/iotzigbeeserver-go/zcl/zcl-go/frame"
	"github.com/h3c/iotzigbeeserver-go/zcl/zcl-go/reflection"
)

//CommandExtractor ...
type CommandExtractor func(commandDescriptors map[uint8]*cluster.CommandDescriptor) (uint8, *cluster.CommandDescriptor, error)

//ClusterQuery ...
type ClusterQuery func(c map[cluster.ID]*cluster.Cluster) (cluster.ID, *cluster.Cluster, error)

//CommandQuery ...
type CommandQuery func(c *cluster.Cluster) (uint8, *cluster.CommandDescriptor, error)

//FrameControl ...
type FrameControl struct {
	FrameType              frame.Type      //00通用，01为cluster特殊
	ManufacturerSpecific   bool            //特殊厂商标识，1表示厂商自定义，需要有Manufacturer Code；0表示标准的
	Direction              frame.Direction //0：Client to Server；1：Server to Client；承载属性实体为Server，操作影响属性为Client
	DisableDefaultResponse bool            //默认回复标识，0表示默认有回复，1表示出现error才有默认回复
}

//Frame ZCL报文格式
type Frame struct {
	FrameControl              *FrameControl //FrameControl
	ManufacturerCode          uint16        //厂商编码
	TransactionSequenceNumber uint8         //帧序列号
	CommandIdentifier         uint8         //command ID
	CommandName               string        //
	Command                   interface{}   //
}

//IncomingMessage ZNP格式，终端数据报文
type IncomingMessage struct {
	GroupID              uint16
	ClusterID            uint16
	SrcAddr              string
	SrcEndpoint          uint8
	DstEndpoint          uint8
	WasBroadcast         bool
	LinkQuality          uint8
	SecurityUse          bool
	Timestamp            uint32
	TransactionSeqNumber uint8
	Data                 *Frame
}

// Zcl Zcl
type Zcl struct {
	library *cluster.Library
}

// New New
func New() *Zcl {
	return &Zcl{cluster.New()}
}

// ToZclIncomingMessage ToZclIncomingMessage
func (z *Zcl) ToZclIncomingMessage(m *znp.AfIncomingMessage) (*IncomingMessage, error) {
	im := &IncomingMessage{}
	im.GroupID = m.GroupID
	im.ClusterID = m.ClusterID
	im.SrcAddr = m.SrcAddr
	im.SrcEndpoint = m.SrcEndpoint
	im.DstEndpoint = m.DstEndpoint
	im.WasBroadcast = m.WasBroadcast > 0
	im.LinkQuality = m.LinkQuality
	im.SecurityUse = m.SecurityUse > 0
	im.Timestamp = m.Timestamp
	im.TransactionSeqNumber = m.TransSeqNumber
	data, err := z.toZclFrame(m.Data, m.ClusterID)
	im.Data = data
	return im, err
}

func (z *Zcl) toZclFrame(data []uint8, clusterID uint16) (*Frame, error) {
	frame := frame.Decode(data)
	f := &Frame{}
	f.FrameControl = z.toZclFrameControl(frame.FrameControl)
	f.ManufacturerCode = frame.ManufacturerCode
	f.TransactionSequenceNumber = frame.TransactionSequenceNumber
	f.CommandIdentifier = frame.CommandIdentifier
	cmd, name, err := z.toZclCommand(clusterID, frame)
	f.CommandName = name
	f.Command = cmd
	return f, err
}

func (z *Zcl) toZclCommand(clusterID uint16, f *frame.Frame) (interface{}, string, error) {
	var cd *cluster.CommandDescriptor
	var ok bool
	switch f.FrameControl.FrameType {
	case frame.FrameTypeGlobal:
		// globallogger.Log.Infof("[toZclCommand][frame.FrameTypeGlobal] Frame: %+v", *f)
		if cd, ok = z.library.Global()[f.CommandIdentifier]; !ok {
			return nil, "", fmt.Errorf("unsupported global cmd identifier %d", f.CommandIdentifier)
		}
		// globallogger.Log.Infof("[toZclCommand][frame.FrameTypeGlobal] Frame: %+v", cd.Command)
		cmd := cd.Command
		copy := reflection.Copy(cmd)
		// globallogger.Log.Infof("[toZclCommand][frame.FrameTypeGlobal] copy: %+v", copy)
		bin.Decode(f.Payload, copy)
		// globallogger.Log.Infof("[toZclCommand][frame.FrameTypeGlobal] copy: %+v", copy)
		z.patchName(copy, clusterID, f.CommandIdentifier)
		return copy, cd.Name, nil
	case frame.FrameTypeLocal:
		// globallogger.Log.Infof("[toZclCommand][frame.FrameTypeLocal]......")
		var c *cluster.Cluster
		if c, ok = z.library.Clusters()[cluster.ID(clusterID)]; !ok {
			return nil, "", fmt.Errorf("unknown cluster %d", clusterID)
		}
		var commandDescriptors map[uint8]*cluster.CommandDescriptor
		switch f.FrameControl.Direction {
		case frame.DirectionClientServer:
			commandDescriptors = c.CommandDescriptors.Received
		case frame.DirectionServerClient:
			commandDescriptors = c.CommandDescriptors.Generated
		}
		if cd, ok = commandDescriptors[f.CommandIdentifier]; !ok {
			return nil, "", fmt.Errorf("cluster %d doesn't support this cmd %d", clusterID, f.CommandIdentifier)
		}
		cmd := cd.Command
		copy := reflection.Copy(cmd)
		bin.Decode(f.Payload, copy)
		return copy, cd.Name, nil
	}
	return nil, "", fmt.Errorf("unknown frame type")
}

func (z *Zcl) patchName(cmd interface{}, clusterID uint16, commandID uint8) {
	switch cmd := cmd.(type) {
	case *cluster.ReadAttributesResponse:
		for _, v := range cmd.ReadAttributeStatuses {
			v.AttributeName = z.getAttributeName(clusterID, v.AttributeID)
		}
		// globallogger.Log.Infof("llllll%+v", cmd.ReadAttributeStatuses[0])
		// globallogger.Log.Infof("llllll%+v", cmd.ReadAttributeStatuses[0].Attribute)
	case *cluster.WriteAttributesCommand:
		for _, v := range cmd.WriteAttributeRecords {
			v.AttributeName = z.getAttributeName(clusterID, v.AttributeID)
		}
	case *cluster.WriteAttributesUndividedCommand:
		for _, v := range cmd.WriteAttributeRecords {
			v.AttributeName = z.getAttributeName(clusterID, v.AttributeID)
		}
	case *cluster.WriteAttributesNoResponseCommand:
		for _, v := range cmd.WriteAttributeRecords {
			v.AttributeName = z.getAttributeName(clusterID, v.AttributeID)
		}
	case *cluster.WriteAttributesResponse:
		for _, v := range cmd.WriteAttributeStatuses {
			v.AttributeName = z.getAttributeName(clusterID, v.AttributeID)
		}
	case *cluster.ConfigureReportingCommand:
		for _, v := range cmd.AttributeReportingConfigurationRecords {
			v.AttributeName = z.getAttributeName(clusterID, v.AttributeID)
		}
	case *cluster.ConfigureReportingResponse:
		for _, v := range cmd.AttributeStatusRecords {
			v.AttributeName = z.getAttributeName(clusterID, v.AttributeID)
		}
	case *cluster.ReadReportingConfigurationCommand:
		for _, v := range cmd.AttributeRecords {
			v.AttributeName = z.getAttributeName(clusterID, v.AttributeID)
		}
	case *cluster.ReadReportingConfigurationResponse:
		for _, v := range cmd.AttributeReportingConfigurationResponseRecords {
			v.AttributeName = z.getAttributeName(clusterID, v.AttributeID)
		}
	case *cluster.ReportAttributesCommand:
		for _, v := range cmd.AttributeReports {
			v.AttributeName = z.getAttributeName(clusterID, v.AttributeID)
		}
		// globallogger.Log.Infof("llllll%+v", cmd.AttributeReports[0])
		// globallogger.Log.Infof("llllll%+v", cmd.AttributeReports[0].Attribute)
	case *cluster.DiscoverAttributesResponse:
		for _, v := range cmd.AttributeInformations {
			v.AttributeName = z.getAttributeName(clusterID, v.AttributeID)
		}
	case *cluster.ReadAttributesStructuredCommand:
		for _, v := range cmd.AttributeSelectors {
			v.AttributeName = z.getAttributeName(clusterID, v.AttributeID)
		}
	case *cluster.WriteAttributesStructuredCommand:
		for _, v := range cmd.WriteAttributeStructuredRecords {
			v.AttributeName = z.getAttributeName(clusterID, v.AttributeID)
		}
	case *cluster.WriteAttributesStructuredResponse:
		for _, v := range cmd.WriteAttributeStatusRecords {
			v.AttributeName = z.getAttributeName(clusterID, v.AttributeID)
		}
	case *cluster.DiscoverAttributesExtendedResponse:
		for _, v := range cmd.ExtendedAttributeInformations {
			v.AttributeName = z.getAttributeName(clusterID, v.AttributeID)
		}
	}
}

func (z *Zcl) getAttributeName(clusterID uint16, attributeID uint16) string {
	if cluster, ok := z.library.Clusters()[cluster.ID(clusterID)]; ok {
		if attributeDescriptor, ok := cluster.AttributeDescriptors[attributeID]; ok {
			return attributeDescriptor.Name
		}
	}
	return ""
}

func (z *Zcl) toZclFrameControl(frameControl *frame.Control) *FrameControl {
	fc := &FrameControl{}
	fc.FrameType = frameControl.FrameType
	fc.ManufacturerSpecific = frameControl.ManufacturerSpecific > 0
	fc.Direction = frameControl.Direction
	fc.DisableDefaultResponse = frameControl.DisableDefaultResponse > 0
	return fc
}

// EncFrameConfigurationToHexString EncFrameConfigurationToHexString
func (z *Zcl) EncFrameConfigurationToHexString(frameConfiguration frame.Configuration, transactionID uint8) (string, uint8) {
	frameObj := frame.New()
	if frameConfiguration.FrameTypeConfigured {
		frameObj.BuildFrameType(frameConfiguration.FrameType)
	}
	if frameConfiguration.ManufacturerCodeConfigured {
		frameObj.BuildManufacturerCode(frameConfiguration.ManufacturerCode)
	}
	if frameConfiguration.DirectionConfigured {
		frameObj.BuildDirection(frameConfiguration.Direction)
	}
	if frameConfiguration.DisableDefaultResponseConfigured {
		frameObj.BuildDisableDefaultResponse(frameConfiguration.DisableDefaultResponse)
	}
	if frameConfiguration.CommandIDConfigured {
		frameObj.BuildCommandID(frameConfiguration.CommandID)
	}
	if frameConfiguration.CommandConfigured {
		frameObj.BuildCommand(frameConfiguration.Command)
	}
	frameConfig, err := frameObj.Build()
	if err != nil {
		globallogger.Log.Errorf("[EncFrameConfigurationToHexString]: err: %+v", err)
	}
	if transactionID != 0 {
		frameConfig.TransactionSequenceNumber = transactionID
	}
	// globallogger.Log.Infof("[EncFrameConfigurationToHexString]: frameConfig.FrameControl: %+v", frameConfig.FrameControl)
	// globallogger.Log.Infof("[EncFrameConfigurationToHexString]: frameConfig: %+v", frameConfig)
	frameBuf := frame.Encode(frameConfig)
	// globallogger.Log.Infof("[EncFrameConfigurationToHexString]: frameBuf: %+v", frameBuf)
	return hex.EncodeToString(frameBuf), frameConfig.TransactionSequenceNumber
}

// ClusterLibrary ClusterLibrary
func (z *Zcl) ClusterLibrary() *cluster.Library {
	return z.library
}
