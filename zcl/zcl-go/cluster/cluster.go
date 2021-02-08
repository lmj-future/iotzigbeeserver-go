package cluster

// AttributeDescriptor AttributeDescriptor
type AttributeDescriptor struct {
	Name   string
	Type   ZclDataType
	Access Access
}

// CommandDescriptor CommandDescriptor
type CommandDescriptor struct {
	Name    string
	Command interface{}
}

// CommandDescriptors CommandDescriptors
type CommandDescriptors struct {
	Received  map[uint8]*CommandDescriptor
	Generated map[uint8]*CommandDescriptor
}

// Cluster Cluster
type Cluster struct {
	Name                 string
	AttributeDescriptors map[uint16]*AttributeDescriptor
	CommandDescriptors   *CommandDescriptors
}

// Library Library
type Library struct {
	global   map[uint8]*CommandDescriptor
	clusters map[ID]*Cluster
}

// Access Access
type Access uint8

// Read Write Reportable Scene
const (
	Read       Access = 0x01
	Write      Access = 0x02
	Reportable Access = 0x04
	Scene      Access = 0x08
)

// ZHADeviceID ZHADeviceID
type ZHADeviceID struct {
	DeviceName string
	DeviceID   uint16
}

// ZHADeviceID List
var (
	OnOffSwitchDevice ZHADeviceID = ZHADeviceID{
		DeviceName: "onOffSwitch",
		DeviceID:   0,
	}
	// ZHADeviceID = ZHADeviceID{"levelControlSwitch": 1}
	OnOffOutputDevice ZHADeviceID = ZHADeviceID{
		DeviceName: "onOffOutput",
		DeviceID:   2,
	}
	// ZHADeviceID = ZHADeviceID{"levelControllableOutput": 3}
	SceneSelectorDevice ZHADeviceID = ZHADeviceID{
		DeviceName: "sceneSelector",
		DeviceID:   4,
	}
	// ZHADeviceID = ZHADeviceID{"configurationTool": 5}
	RemoteControlDevice ZHADeviceID = ZHADeviceID{
		DeviceName: "remoteControl",
		DeviceID:   6,
	}
	// ZHADeviceID = ZHADeviceID{"combinedInterface": 7}
	// ZHADeviceID = ZHADeviceID{"rangeExtender": 8}
	MainsPowerOutletDevice ZHADeviceID = ZHADeviceID{
		DeviceName: "mainsPowerOutlet",
		DeviceID:   9,
	}
	// ZHADeviceID = ZHADeviceID{"doorLock": 10}
	// ZHADeviceID = ZHADeviceID{"doorLockController": 11}
	// ZHADeviceID = ZHADeviceID{"simpleSensor": 12}
	// ZHADeviceID = ZHADeviceID{"consumptionAwarenessDevice": 13}
	// ZHADeviceID = ZHADeviceID{"homeGateway": 80}
	SmartPlugDevice ZHADeviceID = ZHADeviceID{
		DeviceName: "smartPlug",
		DeviceID:   81,
	}
	// ZHADeviceID = ZHADeviceID{"whiteGoods": 82}
	// ZHADeviceID = ZHADeviceID{"meterInterface": 83}
	// ZHADeviceID = ZHADeviceID{"testDevice": 255}
	OnOffLightDevice ZHADeviceID = ZHADeviceID{
		DeviceName: "onOffLight",
		DeviceID:   256,
	}
	// ZHADeviceID = ZHADeviceID{"dimmableLight": 257}
	// ZHADeviceID = ZHADeviceID{"coloredDimmableLight": 258}
	OnOffLightSwitchDevice ZHADeviceID = ZHADeviceID{
		DeviceName: "onOffLightSwitch",
		DeviceID:   259,
	}
	// ZHADeviceID = ZHADeviceID{"dimmerSwitch": 260}
	// ZHADeviceID = ZHADeviceID{"colorDimmerSwitch": 261}
	// ZHADeviceID = ZHADeviceID{"lightSensor": 262}
	OccupancySensorDevice ZHADeviceID = ZHADeviceID{
		DeviceName: "occupancySensor",
		DeviceID:   263,
	}
	// ZHADeviceID = ZHADeviceID{"shade": 512}
	// ZHADeviceID = ZHADeviceID{"shadeController": 513}
	WindowCoveringDeviceDevice ZHADeviceID = ZHADeviceID{
		DeviceName: "windowCoveringDevice",
		DeviceID:   514,
	}
	WindowCoveringControllerDevice ZHADeviceID = ZHADeviceID{
		DeviceName: "windowCoveringController",
		DeviceID:   515,
	}
	// ZHADeviceID = ZHADeviceID{"heatingCoolingUnit": 768}
	// ZHADeviceID = ZHADeviceID{"thermostat": 769}
	TemperatureSensorDevice ZHADeviceID = ZHADeviceID{
		DeviceName: "temperatureSensor",
		DeviceID:   770,
	}
	// ZHADeviceID = ZHADeviceID{"pump": 771}
	// ZHADeviceID = ZHADeviceID{"pumpController": 772}
	// ZHADeviceID = ZHADeviceID{"pressureSensor": 773}
	// ZHADeviceID = ZHADeviceID{"flowSensor": 774}
	// ZHADeviceID = ZHADeviceID{"miniSplitAc": 775}
	// ZHADeviceID = ZHADeviceID{"iasControlIndicatingEquipment": 1024}
	// ZHADeviceID = ZHADeviceID{"iasAncillaryControlEquipment": 1025}
	IasZoneDevice ZHADeviceID = ZHADeviceID{
		DeviceName: "iasZone",
		DeviceID:   1026,
	}
	IasWarningDevice ZHADeviceID = ZHADeviceID{
		DeviceName: "iasWarningDevice",
		DeviceID:   1027,
	}
)

// ID ClusterID
type ID uint16

// ClusterID List
const (
	Basic                          ID = 0x0000
	PowerConfiguration             ID = 0x0001
	DeviceTemperatureConfiguration ID = 0x0002
	Identify                       ID = 0x0003
	Scenes                         ID = 0x0005
	OnOff                          ID = 0x0006
	LevelControl                   ID = 0x0008
	Alarms                         ID = 0x0009
	Time                           ID = 0x000a
	MultistateInput                ID = 0x0012
	OTA                            ID = 0x0019
	TemperatureMeasurement         ID = 0x0402
	RelativeHumidityMeasurement    ID = 0x0405
	ElectricalMeasurement          ID = 0x0b04
	IASZone                        ID = 0x0500
	IASWarningDevice               ID = 0x0502
	SmartEnergyMetering            ID = 0x0702
)

// MAILEKE私有属性
const (
	MAILEKEMeasurement ID = 0x0415
)

// HEIMAN私有属性
const (
	HEIMANPM25Measurement         ID = 0x042a
	HEIMANFormaldehydeMeasurement ID = 0x042b
	HEIMANScenes                  ID = 0xfc80
	HEIMANAirQualityMeasurement   ID = 0xfc81
	HEIMANAirQualityPM10          ID = 0xfc82
	HEIMANAirQualityTVOC          ID = 0xfc83
	HEIMANInfraredRemote          ID = 0xfc82
)

// HONYAR私有属性
const (
	HONYARScenes ID = 0xfe05
)

// New New
func New() *Library {
	return &Library{
		global: map[uint8]*CommandDescriptor{
			0x00: {"ReadAttributes", &ReadAttributesCommand{}},
			0x01: {"ReadAttributesResponse", &ReadAttributesResponse{}},
			0x02: {"WriteAttributes", &WriteAttributesCommand{}},
			0x03: {"WriteAttributesUndivided", &WriteAttributesUndividedCommand{}},
			0x04: {"WriteAttributesResponse", &WriteAttributesResponse{}},
			0x05: {"WriteAttributesNoResponse", &WriteAttributesNoResponseCommand{}},
			0x06: {"ConfigureReporting", &ConfigureReportingCommand{}},
			0x07: {"ConfigureReportingResponse", &ConfigureReportingResponse{}},
			0x08: {"ReadReportingConfiguration", &ReadReportingConfigurationCommand{}},
			0x09: {"ReadReportingConfigurationResponse", &ReadReportingConfigurationResponse{}},
			0x0a: {"ReportAttributes", &ReportAttributesCommand{}},
			0x0b: {"DefaultResponse", &DefaultResponseCommand{}},
			0x0c: {"DiscoverAttributes", &DiscoverAttributesCommand{}},
			0x0d: {"DiscoverAttributesResponse", &DiscoverAttributesResponse{}},
			0x0e: {"ReadAttributesStructured", &ReadAttributesStructuredCommand{}},
			0x0f: {"WriteAttributesStructured", &WriteAttributesStructuredCommand{}},
			0x10: {"WriteAttributesStructuredResponse", &WriteAttributesStructuredResponse{}},
			0x11: {"DiscoverCommandsReceived", &DiscoverCommandsReceivedCommand{}},
			0x12: {"DiscoverCommandsReceivedResponse", &DiscoverCommandsReceivedResponse{}},
			0x13: {"DiscoverCommandsGenerated", &DiscoverCommandsGeneratedCommand{}},
			0x14: {"DiscoverCommandsGeneratedResponse", &DiscoverCommandsGeneratedResponse{}},
			0x15: {"DiscoverAttributesExtended", &DiscoverAttributesExtendedCommand{}},
			0x16: {"DiscoverAttributesExtendedResponse", &DiscoverAttributesExtendedResponse{}},
		},
		clusters: map[ID]*Cluster{
			Basic: {
				Name: "Basic",
				AttributeDescriptors: map[uint16]*AttributeDescriptor{
					0x0000: {"ZLibraryVersion", ZclDataTypeUint8, Read},
					0x0001: {"ApplicationVersion", ZclDataTypeUint8, Read},
					0x0002: {"StackVersion", ZclDataTypeUint8, Read},
					0x0003: {"HWVersion", ZclDataTypeUint8, Read},
					0x0004: {"ManufacturerName", ZclDataTypeCharStr, Read},
					0x0005: {"ModelIdentifier", ZclDataTypeCharStr, Read},
					0x0006: {"DateCode", ZclDataTypeCharStr, Read},
					0x0007: {"PowerSource", ZclDataTypeEnum8, Read},
					0x0010: {"LocationDescription", ZclDataTypeCharStr, Read | Write},
					0x0011: {"PhysicalEnvironment", ZclDataTypeEnum8, Read | Write},
					0x0012: {"DeviceEnabled", ZclDataTypeBoolean, Read | Write},
					0x0013: {"AlarmMask", ZclDataTypeBitmap8, Read | Write},
					0x0014: {"DisableLocalConfig", ZclDataTypeBitmap8, Read | Write},
					0x4000: {"SWBuildID", ZclDataTypeCharStr, Read},
				},
				CommandDescriptors: &CommandDescriptors{
					Received: map[uint8]*CommandDescriptor{
						0x00: {"ResetToFactoryDefaults", &ResetToFactoryDefaultsCommand{}},
					},
				},
			},
			PowerConfiguration: {
				Name: "PowerConfiguration",
				AttributeDescriptors: map[uint16]*AttributeDescriptor{
					0x0000: {"MainsVoltage", ZclDataTypeUint16, Read},
					0x0001: {"MainsFrequency", ZclDataTypeUint8, Read},
					0x0010: {"MainsAlarmMask", ZclDataTypeBitmap8, Read | Write},
					0x0011: {"MainsVoltageMinThreshold", ZclDataTypeUint16, Read | Write},
					0x0012: {"MainsVoltageMaxThreshold", ZclDataTypeUint16, Read | Write},
					0x0013: {"MainsVoltageDwellTripPoint", ZclDataTypeUint16, Read | Write},
					0x0020: {"BatteryVoltage", ZclDataTypeUint8, Read},
					0x0021: {"BatteryPercentageRemaining", ZclDataTypeUint8, Read | Reportable},
					0x0030: {"BatteryManufacturer", ZclDataTypeCharStr, Read | Write},
					0x0031: {"BatterySize", ZclDataTypeEnum8, Read | Write},
					0x0032: {"BatteryAHrRating", ZclDataTypeUint16, Read | Write},
					0x0033: {"BatteryQuantity", ZclDataTypeUint8, Read | Write},
					0x0034: {"BatteryRatedVoltage", ZclDataTypeUint8, Read | Write},
					0x0035: {"BatteryAlarmMask", ZclDataTypeBitmap8, Read | Write},
					0x0036: {"BatteryVoltageMinThreshold", ZclDataTypeUint8, Read | Write},
					0x0037: {"BatteryVoltageThreshold1", ZclDataTypeUint8, Read | Write},
					0x0038: {"BatteryVoltageThreshold2", ZclDataTypeUint8, Read | Write},
					0x0039: {"BatteryVoltageThreshold3", ZclDataTypeUint8, Read | Write},
					0x003a: {"BatteryPercentageMinThreshold", ZclDataTypeUint8, Read | Write},
					0x003b: {"BatteryPercentageThreshold1", ZclDataTypeUint8, Read | Write},
					0x003c: {"BatteryPercentageThreshold2", ZclDataTypeUint8, Read | Write},
					0x003d: {"BatteryPercentageThreshold3", ZclDataTypeUint8, Read | Write},
					0x003e: {"BatteryAlarmState", ZclDataTypeBitmap32, Read},
				},
			},
			DeviceTemperatureConfiguration: {
				Name: "DeviceTemperatureConfiguration",
				AttributeDescriptors: map[uint16]*AttributeDescriptor{
					0x0000: {"CurrentTemperature", ZclDataTypeInt16, Read},
					0x0001: {"MinTempExperienced", ZclDataTypeInt16, Read},
					0x0002: {"MaxTempExperienced", ZclDataTypeInt16, Read},
					0x0003: {"OverTempTotalDwell", ZclDataTypeInt16, Read},
					0x0010: {"DeviceTempAlarmMask", ZclDataTypeBitmap16, Read | Write},
					0x0011: {"LowTempThreshold", ZclDataTypeInt16, Read | Write},
					0x0012: {"HighTempThreshold", ZclDataTypeInt16, Read | Write},
					0x0013: {"LowTempDwellTripPoint", ZclDataTypeUint24, Read | Write},
					0x0014: {"HighTempDwellTripPoint", ZclDataTypeUint24, Read | Write},
				},
			},
			Identify: {
				Name: "Identify",
				AttributeDescriptors: map[uint16]*AttributeDescriptor{
					0x0000: {"IdentifyTime", ZclDataTypeInt16, Read | Write},
				},
				CommandDescriptors: &CommandDescriptors{
					Received: map[uint8]*CommandDescriptor{
						0x00: {"Identify", &IdentifyCommand{}},
						0x01: {"IdentifyQuery", &IdentifyQueryCommand{}},
						0x40: {"TriggerEffect ", &TriggerEffectCommand{}},
					},
					Generated: map[uint8]*CommandDescriptor{
						0x00: {"IdentifyQueryResponse ", &IdentifyQueryResponse{}},
					},
				},
			},
			Scenes: {
				Name: "Scenes",
				AttributeDescriptors: map[uint16]*AttributeDescriptor{
					0x0000: {"SceneCount", ZclDataTypeUint8, Read},
					0x0001: {"CurrentScene", ZclDataTypeUint8, Read},
					0x0002: {"CurrentGroup", ZclDataTypeUint16, Read},
					0x0003: {"SceneValid", ZclDataTypeBoolean, Read},
					0x0004: {"NameSupport", ZclDataTypeBitmap8, Read},
					0x0005: {"LastConfiguredBy", ZclDataTypeIeeeAddr, Read},
				},
				CommandDescriptors: &CommandDescriptors{
					Received: map[uint8]*CommandDescriptor{
						0x00: {"AddScene", &AddSceneCommand{}},
						0x01: {"ViewScene", &ViewSceneCommand{}},
						0x02: {"RemoveScene", &RemoveSceneCommand{}},
						0x03: {"RemoveAllScenes", &RemoveAllScenesCommand{}},
						0x04: {"StoreScene", &StoreSceneCommand{}},
						0x05: {"RecallScene", &RecallSceneCommand{}},
						0x06: {"GetSceneMembership", &GetSceneMembership{}},
						0x40: {"EnhancedAddScene", &EnhancedAddSceneCommand{}},
						0x41: {"EnhancedViewScene", &EnhancedViewSceneCommand{}},
						0x42: {"CopyScene", &CopySceneCommand{}},
					},
					Generated: map[uint8]*CommandDescriptor{
						0x00: {"AddSceneResponse", &AddSceneResponse{}},
						0x01: {"ViewSceneResponse", &ViewSceneResponse{}},
						0x02: {"RemoveSceneResponse", &RemoveSceneResponse{}},
						0x03: {"RemoveAllScenesResponse", &RemoveAllScenesResponse{}},
						0x04: {"StoreSceneResponse", &StoreSceneResponse{}},
						0x06: {"GetSceneMembershipResponse", &GetSceneMembershipResponse{}},
						0x40: {"EnhancedAddSceneResponse", &EnhancedAddSceneResponse{}},
						0x41: {"EnhancedViewSceneResponse", &EnhancedViewSceneResponse{}},
						0x42: {"CopySceneResponse", &CopySceneResponse{}},
					},
				},
			},
			OnOff: {
				Name: "OnOff",
				AttributeDescriptors: map[uint16]*AttributeDescriptor{
					0x0000: {"OnOff", ZclDataTypeBoolean, Read | Reportable | Scene},
					0x4000: {"GlobalSceneControl", ZclDataTypeBoolean, Read},
					0x4001: {"OnTime", ZclDataTypeUint16, Read | Write},
					0x4002: {"OffWaitTime", ZclDataTypeUint16, Read | Write},
					0x8000: {"ChildLock", ZclDataTypeBoolean, Read | Write},       //鸿雁童锁开关
					0x8005: {"BackGroundLight", ZclDataTypeBoolean, Read | Write}, //鸿雁背景光开关
					0x8006: {"PowerOffMemory", ZclDataTypeBoolean, Read | Write},  //鸿雁断电记忆开关
				},
				CommandDescriptors: &CommandDescriptors{
					Received: map[uint8]*CommandDescriptor{
						0x00: {"Off", &OffCommand{}},
						0x01: {"On", &OnCommand{}},
						0x02: {"Toggle ", &ToggleCommand{}},
						0x40: {"OffWithEffect ", &OffWithEffectCommand{}},
						0x41: {"OnWithRecallGlobalScene ", &OnWithRecallGlobalSceneCommand{}},
						0x42: {"OnWithTimedOff ", &OnWithTimedOffCommand{}},
					},
				},
			},
			LevelControl: {
				Name: "LevelControl",
				AttributeDescriptors: map[uint16]*AttributeDescriptor{
					0x0000: {"CurrentLevel", ZclDataTypeUint8, Read | Reportable},
					0x0001: {"RemainingTime", ZclDataTypeUint16, Read},
					0x0010: {"OnOffTransitionTime", ZclDataTypeUint16, Read | Write},
					0x0011: {"OnLevel", ZclDataTypeUint8, Read | Write},
					0x0012: {"OnTransitionTime", ZclDataTypeUint16, Read | Write},
					0x0013: {"OffTransitionTime", ZclDataTypeUint16, Read | Write},
					0x0014: {"DefaultMoveRate", ZclDataTypeUint16, Read | Write},
				},
				CommandDescriptors: &CommandDescriptors{
					Received: map[uint8]*CommandDescriptor{
						0x00: {"MoveToLevel ", &MoveToLevelCommand{}},
						0x01: {"Move", &MoveCommand{}},
						0x02: {"Step ", &StepCommand{}},
						0x03: {"Stop ", &StopCommand{}},
						0x04: {"MoveToLevel/OnOff", &MoveToLevelOnOffCommand{}},
						0x05: {"Move/OnOff", &MoveOnOffCommand{}},
						0x06: {"Step/OnOff", &StepOnOffCommand{}},
						0x07: {"Stop/OnOff", &StopOnOffCommand{}},
					},
				},
			},
			Alarms: {
				Name: "Alarms",
				AttributeDescriptors: map[uint16]*AttributeDescriptor{
					0x0000: {"AlarmCount", ZclDataTypeUint16, Read},
				},
				CommandDescriptors: &CommandDescriptors{
					Received: map[uint8]*CommandDescriptor{
						0x00: {"ResetAlarm", &ResetAlarmCommand{}},
						0x01: {"ResetAllAlarms", &ResetAllAlarmsCommand{}},
						0x02: {"GetAlarm", &GetAlarmCommand{}},
						0x03: {"ResetAlarmLog", &ResetAlarmLogCommand{}},
					},
					Generated: map[uint8]*CommandDescriptor{
						0x00: {"Alarm", &AlarmCommand{}},
						0x01: {"GetAlarmResponse", &GetAlarmResponse{}},
					},
				},
			},
			Time: {
				Name: "Time",
				AttributeDescriptors: map[uint16]*AttributeDescriptor{
					0x0000: {"Time", ZclDataTypeUtc, Read | Write},
					0x0001: {"TimeStatus", ZclDataTypeBitmap8, Read | Write},
					0x0002: {"TimeZone", ZclDataTypeInt32, Read | Write},
					0x0003: {"DstStart", ZclDataTypeUint32, Read | Write},
					0x0004: {"DstEnd", ZclDataTypeUint32, Read | Write},
					0x0005: {"DstShift", ZclDataTypeInt32, Read | Write},
					0x0006: {"StandardTime", ZclDataTypeUint32, Read},
					0x0007: {"LocalTime", ZclDataTypeUint32, Read},
					0x0008: {"LastSetTime", ZclDataTypeUtc, Read},
					0x0009: {"ValidUntilTime", ZclDataTypeUtc, Read | Write},
				},
			},
			MultistateInput: {
				Name: "MultistateInput",
				AttributeDescriptors: map[uint16]*AttributeDescriptor{
					0x000E: {"StateText", ZclDataTypeArray, Read | Write},
					0x001C: {"Description", ZclDataTypeCharStr, Read | Write},
					0x004A: {"NumberOfStates", ZclDataTypeUint16, Read | Write},
					0x0051: {"OutOfService", ZclDataTypeBoolean, Read | Write},
					0x0055: {"PresentValue", ZclDataTypeUint16, Read | Write},
					0x0067: {"Reliability", ZclDataTypeEnum8, Read | Write},
					0x006F: {"StatusFlags", ZclDataTypeBitmap8, Read},
					0x0100: {"ApplicationType", ZclDataTypeUint32, Read},
				},
			},
			OTA: {
				Name: "OTA",
				AttributeDescriptors: map[uint16]*AttributeDescriptor{
					0x0000: {"UpgradeServerID", ZclDataTypeIeeeAddr, Read},
					0x0001: {"FileOffset", ZclDataTypeUint32, Read},
					0x0002: {"CurrentFileVersion", ZclDataTypeUint32, Read},
					0x0003: {"CurrentZigBeeStackVersion", ZclDataTypeUint16, Read},
					0x0004: {"DownloadedFileVersion", ZclDataTypeUint32, Read},
					0x0005: {"DownloadedZigBeeStackVersion", ZclDataTypeUint16, Read},
					0x0006: {"ImageUpgradeStatus", ZclDataTypeEnum8, Read},
					0x0007: {"ManufacturerID", ZclDataTypeUint16, Read},
					0x0008: {"ImageTypeID ", ZclDataTypeUint16, Read},
					0x0009: {"MinimumBlockPeriod ", ZclDataTypeUint16, Read},
					0x000a: {"ImageStamp ", ZclDataTypeUint32, Read},
				},
			},
			TemperatureMeasurement: {
				Name: "TemperatureMeasurement",
				AttributeDescriptors: map[uint16]*AttributeDescriptor{
					0x0000: {"MeasuredValue", ZclDataTypeInt16, Read},
					0x0001: {"MinMeasuredValue", ZclDataTypeInt16, Read},
					0x0002: {"MaxMeasuredValue", ZclDataTypeInt16, Read},
					0x0003: {"Tolerance", ZclDataTypeUint16, Read},
				},
			},
			RelativeHumidityMeasurement: {
				Name: "RelativeHumidityMeasurement",
				AttributeDescriptors: map[uint16]*AttributeDescriptor{
					0x0000: {"MeasuredValue", ZclDataTypeUint16, Read},
					0x0001: {"MinMeasuredValue", ZclDataTypeUint16, Read},
					0x0002: {"MaxMeasuredValue", ZclDataTypeUint16, Read},
					0x0003: {"Tolerance", ZclDataTypeUint16, Read},
				},
			},
			ElectricalMeasurement: {
				Name: "ElectricalMeasurement",
				AttributeDescriptors: map[uint16]*AttributeDescriptor{
					0x0000: {"MeasurementType", ZclDataTypeBitmap32, Read},

					0x0100: {"DCVoltage", ZclDataTypeInt16, Read},
					0x0101: {"DCVoltageMin", ZclDataTypeInt16, Read},
					0x0102: {"DCVoltageMax", ZclDataTypeInt16, Read},
					0x0103: {"DCCurrent", ZclDataTypeInt16, Read},
					0x0104: {"DCCurrentMin", ZclDataTypeInt16, Read},
					0x0105: {"DCCurrentMax", ZclDataTypeInt16, Read},
					0x0106: {"DCPower", ZclDataTypeInt16, Read},
					0x0107: {"DCPowerMin", ZclDataTypeInt16, Read},
					0x0108: {"DCPowerMax", ZclDataTypeInt16, Read},

					0x0200: {"DCVoltageMultiplier", ZclDataTypeUint16, Read},
					0x0201: {"DCVoltageDivisor", ZclDataTypeUint16, Read},
					0x0202: {"DCCurrentMultiplier", ZclDataTypeUint16, Read},
					0x0203: {"DCCurrentDivisor", ZclDataTypeUint16, Read},
					0x0204: {"DCPowerMultiplier", ZclDataTypeUint16, Read},
					0x0205: {"DCPowerDivisor", ZclDataTypeUint16, Read},

					0x0300: {"ACFrequency", ZclDataTypeUint16, Read},
					0x0301: {"ACFrequencyMin", ZclDataTypeUint16, Read},
					0x0302: {"ACFrequencyMax", ZclDataTypeUint16, Read},
					0x0303: {"NeutralCurrent", ZclDataTypeUint16, Read},
					0x0304: {"TotalActivePower", ZclDataTypeInt32, Read},
					0x0305: {"TotalReactivePower", ZclDataTypeInt32, Read},
					0x0306: {"ApparentPower", ZclDataTypeUint32, Read},
					0x0307: {"Measured1stHarmonicCurrent", ZclDataTypeInt16, Read},
					0x0308: {"Measured3rdHarmonicCurrent", ZclDataTypeInt16, Read},
					0x0309: {"Measured5thHarmonicCurrent", ZclDataTypeInt16, Read},
					0x030a: {"Measured7thHarmonicCurrent", ZclDataTypeInt16, Read},
					0x030b: {"Measured9thHarmonicCurrent", ZclDataTypeInt16, Read},
					0x030c: {"Measured11thHarmonicCurrent", ZclDataTypeInt16, Read},
					0x030d: {"MeasuredPhase1stHarmonicCurrent", ZclDataTypeInt16, Read},
					0x030e: {"MeasuredPhase3rdHarmonicCurrent", ZclDataTypeInt16, Read},
					0x030f: {"MeasuredPhase5thHarmonicCurrent", ZclDataTypeInt16, Read},
					0x0310: {"MeasuredPhase7thHarmonicCurrent", ZclDataTypeInt16, Read},
					0x0311: {"MeasuredPhase9thHarmonicCurrent", ZclDataTypeInt16, Read},
					0x0312: {"MeasuredPhase11thHarmonicCurrent", ZclDataTypeInt16, Read},

					0x0400: {"ACFrequencyMultiplier", ZclDataTypeUint16, Read},
					0x0401: {"ACFrequencyDivisor", ZclDataTypeUint16, Read},
					0x0402: {"PowerMultiplier", ZclDataTypeUint32, Read},
					0x0403: {"PowerDivisor", ZclDataTypeUint32, Read},
					0x0404: {"HarmonicCurrentMultiplier", ZclDataTypeInt8, Read},
					0x0405: {"PhaseHarmonicCurrentMultiplier", ZclDataTypeInt8, Read},

					0x0500: {"Reserved", ZclDataTypeInt16, Read},
					0x0501: {"LineCurrent", ZclDataTypeUint16, Read},
					0x0502: {"ActiveCurrent", ZclDataTypeInt16, Read},
					0x0503: {"ReactiveCurrent", ZclDataTypeInt16, Read},
					0x0504: {"Reserved", ZclDataTypeInt8, Read},
					0x0505: {"RMSVoltage", ZclDataTypeUint16, Read},
					0x0506: {"RMSVoltageMin", ZclDataTypeUint16, Read},
					0x0507: {"RMSVoltageMax", ZclDataTypeUint16, Read},
					0x0508: {"RMSCurrent", ZclDataTypeUint16, Read},
					0x0509: {"RMSCurrentMin", ZclDataTypeUint16, Read},
					0x050a: {"RMSCurrentMax", ZclDataTypeUint16, Read},
					0x050b: {"ActivePower", ZclDataTypeInt16, Read},
					0x050c: {"ActivePowerMin", ZclDataTypeInt16, Read},
					0x050d: {"ActivePowerMax", ZclDataTypeInt16, Read},
					0x050e: {"ReactivePower", ZclDataTypeInt16, Read},
					0x050f: {"ApparentPower", ZclDataTypeUint16, Read},
					0x0510: {"PowerFactor", ZclDataTypeInt8, Read},
					0x0511: {"AverageRMSVoltageMeasurementPeriod", ZclDataTypeUint16, Read | Write},
					0x0512: {"AverageRMSOverVoltageCounter", ZclDataTypeUint16, Read | Write},
					0x0513: {"AverageRMSUnderVoltageCounter", ZclDataTypeUint16, Read | Write},
					0x0514: {"RMSExtremeOverVoltagePeriod", ZclDataTypeUint16, Read | Write},
					0x0515: {"RMSExtremeUnderVoltagePeriod", ZclDataTypeUint16, Read | Write},
					0x0516: {"RMSVoltageSagPeriod", ZclDataTypeUint16, Read | Write},
					0x0517: {"RMSVoltageSwellPeriod", ZclDataTypeUint16, Read | Write},

					0x0600: {"ACVoltageMultiplier", ZclDataTypeUint16, Read},
					0x0601: {"ACVoltageDivisor", ZclDataTypeUint16, Read},
					0x0602: {"ACCurrentMultiplier", ZclDataTypeUint16, Read},
					0x0603: {"ACCurrentDivisor", ZclDataTypeUint16, Read},
					0x0604: {"ACPowerMultiplier", ZclDataTypeUint16, Read},
					0x0605: {"ACPowerDivisor", ZclDataTypeUint16, Read},

					0x0700: {"DCOverloadAlarmsMask", ZclDataTypeBitmap8, Read | Write},
					0x0701: {"DCVoltageOverload", ZclDataTypeInt16, Read},
					0x0702: {"DCVoltageOverload", ZclDataTypeInt16, Read},

					0x0800: {"ACAlarmsMask", ZclDataTypeBitmap16, Read | Write},
					0x0801: {"ACVoltageOverload", ZclDataTypeInt16, Read},
					0x0802: {"ACCurrentOverload", ZclDataTypeInt16, Read},
					0x0803: {"ACActivePowerOverload", ZclDataTypeInt16, Read},
					0x0804: {"ACReactivePowerOverload", ZclDataTypeInt16, Read},
					0x0805: {"AverageRMSOverVoltage", ZclDataTypeInt16, Read},
					0x0806: {"AverageRMSUnderVoltage", ZclDataTypeInt16, Read},
					0x0807: {"RMSExtremeOverVoltage", ZclDataTypeInt16, Read | Write},
					0x0808: {"RMSExtremeUnderVoltage", ZclDataTypeInt16, Read | Write},
					0x0809: {"RMSVoltageSag", ZclDataTypeInt16, Read | Write},
					0x080a: {"RMSVoltageSwell", ZclDataTypeInt16, Read | Write},

					0x0900: {"ReservedPhB", ZclDataTypeInt16, Read},
					0x0901: {"LineCurrentPhB", ZclDataTypeUint16, Read},
					0x0902: {"ActiveCurrentPhB", ZclDataTypeInt16, Read},
					0x0903: {"ReactiveCurrentPhB", ZclDataTypeInt16, Read},
					0x0904: {"ReservedPhB", ZclDataTypeInt8, Read},
					0x0905: {"RMSVoltagePhB", ZclDataTypeUint16, Read},
					0x0906: {"RMSVoltageMinPhB", ZclDataTypeUint16, Read},
					0x0907: {"RMSVoltageMaxPhB", ZclDataTypeUint16, Read},
					0x0908: {"RMSCurrentPhB", ZclDataTypeUint16, Read},
					0x0909: {"RMSCurrentMinPhB", ZclDataTypeUint16, Read},
					0x090a: {"RMSCurrentMaxPhB", ZclDataTypeUint16, Read},
					0x090b: {"ActivePowerPhB", ZclDataTypeInt16, Read},
					0x090c: {"ActivePowerMinPhB", ZclDataTypeInt16, Read},
					0x090d: {"ActivePowerMaxPhB", ZclDataTypeInt16, Read},
					0x090e: {"ReactivePowerPhB", ZclDataTypeInt16, Read},
					0x090f: {"ApparentPowerPhB", ZclDataTypeUint16, Read},
					0x0910: {"PowerFactorPhB", ZclDataTypeInt8, Read},
					0x0911: {"AverageRMSVoltageMeasurementPeriodPhB", ZclDataTypeUint16, Read | Write},
					0x0912: {"AverageRMSOverVoltageCounterPhB", ZclDataTypeUint16, Read | Write},
					0x0913: {"AverageRMSUnderVoltageCounterPhB", ZclDataTypeUint16, Read | Write},
					0x0914: {"RMSExtremeOverVoltagePeriodPhB", ZclDataTypeUint16, Read | Write},
					0x0915: {"RMSExtremeUnderVoltagePeriodPhB", ZclDataTypeUint16, Read | Write},
					0x0916: {"RMSVoltageSagPeriodPhB", ZclDataTypeUint16, Read | Write},
					0x0917: {"RMSVoltageSwellPeriodPhB", ZclDataTypeUint16, Read | Write},

					0x0a00: {"ReservedPhC", ZclDataTypeInt16, Read},
					0x0a01: {"LineCurrentPhC", ZclDataTypeUint16, Read},
					0x0a02: {"ActiveCurrentPhC", ZclDataTypeInt16, Read},
					0x0a03: {"ReactiveCurrentPhC", ZclDataTypeInt16, Read},
					0x0a04: {"ReservedPhC", ZclDataTypeInt8, Read},
					0x0a05: {"RMSVoltagePhC", ZclDataTypeUint16, Read},
					0x0a06: {"RMSVoltageMinPhC", ZclDataTypeUint16, Read},
					0x0a07: {"RMSVoltageMaxPhC", ZclDataTypeUint16, Read},
					0x0a08: {"RMSCurrentPhC", ZclDataTypeUint16, Read},
					0x0a09: {"RMSCurrentMinPhC", ZclDataTypeUint16, Read},
					0x0a0a: {"RMSCurrentMaxPhC", ZclDataTypeUint16, Read},
					0x0a0b: {"ActivePowerPhC", ZclDataTypeInt16, Read},
					0x0a0c: {"ActivePowerMinPhC", ZclDataTypeInt16, Read},
					0x0a0d: {"ActivePowerMaxPhC", ZclDataTypeInt16, Read},
					0x0a0e: {"ReactivePowerPhC", ZclDataTypeInt16, Read},
					0x0a0f: {"ApparentPowerPhC", ZclDataTypeUint16, Read},
					0x0a10: {"PowerFactorPhC", ZclDataTypeInt8, Read},
					0x0a11: {"AverageRMSVoltageMeasurementPeriodPhC", ZclDataTypeUint16, Read | Write},
					0x0a12: {"AverageRMSOverVoltageCounterPhC", ZclDataTypeUint16, Read | Write},
					0x0a13: {"AverageRMSUnderVoltageCounterPhC", ZclDataTypeUint16, Read | Write},
					0x0a14: {"RMSExtremeOverVoltagePeriodPhC", ZclDataTypeUint16, Read | Write},
					0x0a15: {"RMSExtremeUnderVoltagePeriodPhC", ZclDataTypeUint16, Read | Write},
					0x0a16: {"RMSVoltageSagPeriodPhC", ZclDataTypeUint16, Read | Write},
					0x0a17: {"RMSVoltageSwellPeriodPhC", ZclDataTypeUint16, Read | Write},
				},
				CommandDescriptors: &CommandDescriptors{
					Received: map[uint8]*CommandDescriptor{
						0x00: {"GetProfileInfoResponse", &GetProfileInfoResponse{}},
						0x01: {"GetMeasurementProfileResponse", &GetMeasurementProfileResponse{}},
					},
					Generated: map[uint8]*CommandDescriptor{
						0x00: {"GetProfileInfoCommand", &GetProfileInfoCommand{}},
						0x01: {"GetMeasurementProfileCommand", &GetMeasurementProfileCommand{}},
					},
				},
			},
			IASZone: {
				Name: "IASZone",
				AttributeDescriptors: map[uint16]*AttributeDescriptor{
					0x0000: {"ZoneState", ZclDataTypeEnum8, Read},
					0x0001: {"ZoneType", ZclDataTypeEnum16, Read},
					0x0002: {"ZoneStatus", ZclDataTypeBitmap16, Read},
					0x0010: {"IAS_CIE_Address", ZclDataTypeIeeeAddr, Read | Write},
					0x0011: {"ZoneID", ZclDataTypeUint8, Read},
					0x0012: {"NumberOfZoneSensitivityLevelsSupported", ZclDataTypeUint8, Read},
					0x0013: {"CurrentZoneSensitivityLevel", ZclDataTypeUint8, Read | Write},
				},
				CommandDescriptors: &CommandDescriptors{
					Received: map[uint8]*CommandDescriptor{
						0x00: {"ZoneEnrollResponse", &ZoneEnrollResponse{}},
						0x01: {"InitiateNormalOperationMode", &InitiateNormalOperationModeCommand{}},
						0x02: {"InitiateTestMode", &InitiateTestModeCommand{}},
					},
					Generated: map[uint8]*CommandDescriptor{
						0x00: {"ZoneStatusChangeNotification", &ZoneStatusChangeNotificationCommand{}},
						0x01: {"ZoneEnrollRequest", &ZoneEnrollCommand{}},
					},
				},
			},
			IASWarningDevice: {
				Name: "IASWarningDevice",
				AttributeDescriptors: map[uint16]*AttributeDescriptor{
					0x0000: {"MaxDuration", ZclDataTypeUint16, Read | Write},
				},
				CommandDescriptors: &CommandDescriptors{
					Received: map[uint8]*CommandDescriptor{
						0x00: {"StartWarning", &StartWarning{}},
						0x01: {"Squawk", &Squark{}},
					},
				},
			},
			SmartEnergyMetering: {
				Name: "SmartEnergyMetering",
				AttributeDescriptors: map[uint16]*AttributeDescriptor{
					0x0000: {"CurrentSummationDelivered", ZclDataTypeUint48, Read},
					0x0001: {"CurrentSummationReceived", ZclDataTypeUint48, Read},
					0x0002: {"CurrentMaxDemandDelivered", ZclDataTypeUint48, Read},
					0x0003: {"CurrentMaxDemandReceived", ZclDataTypeUint48, Read},
					0x0004: {"DFTSummation", ZclDataTypeUint48, Read},
					0x0005: {"DailyFreezeTime", ZclDataTypeUint16, Read},
					0x0006: {"PowerFactor", ZclDataTypeInt8, Read},
					0x0007: {"ReadingSnapShotTime", ZclDataTypeUtc, Read},
					0x0008: {"CurrentMaxDemandDeliveredTime", ZclDataTypeUtc, Read},
					0x0009: {"CurrentMaxDemandReceivedTime", ZclDataTypeUtc, Read},
					0x000a: {"DefaultUpdatePeriod", ZclDataTypeUint8, Read},
					0x000b: {"FastPollUpdatePeriod", ZclDataTypeUint8, Read},
					0x000c: {"CurrentBlockPeriodConsumptionDelivered", ZclDataTypeUint48, Read},
					0x000d: {"DailyConsumptionTarget", ZclDataTypeUint24, Read},
					0x000e: {"CurrentBlock", ZclDataTypeEnum8, Read},
					0x000f: {"ProfileIntervalPeriod", ZclDataTypeEnum8, Read},
					0x0010: {"IntervalReadReportingPeriod", ZclDataTypeUint16, Read},
					0x0011: {"PresetReadingTime", ZclDataTypeUint16, Read},
					0x0012: {"VolumePerReport", ZclDataTypeUint16, Read},
					0x0013: {"FlowRestriction", ZclDataTypeUint8, Read},
					0x0014: {"SupplyStatus", ZclDataTypeEnum8, Read},
					0x0015: {"CurrentInletEnergyCarrierSummation", ZclDataTypeUint48, Read},
					0x0016: {"CurrentOutletEnergyCarrierSummation", ZclDataTypeUint48, Read},
					0x0017: {"InletTemperature", ZclDataTypeInt24, Read},
					0x0018: {"OutletTemperature", ZclDataTypeInt24, Read},
					0x0019: {"ControlTemperature", ZclDataTypeInt24, Read},
					0x001a: {"CurrentInletEnergyCarrierDemand", ZclDataTypeInt24, Read},
					0x001b: {"CurrentOutletEnergyCarrierDemand", ZclDataTypeInt24, Read},
					0x001c: {"PreviousBlockPeriodConsumptionDelivered", ZclDataTypeUint48, Read},
					0x0100: {"CurrentTier1SummationDelivered", ZclDataTypeUint48, Read},
					0x0101: {"CurrentTier1SummationReceived", ZclDataTypeUint48, Read},
					0x0102: {"CurrentTier2SummationDelivered", ZclDataTypeUint48, Read},
					0x0103: {"CurrentTier2SummationReceived", ZclDataTypeUint48, Read},
					0x0104: {"CurrentTier3SummationDelivered", ZclDataTypeUint48, Read},
					0x0105: {"CurrentTier3SummationReceived", ZclDataTypeUint48, Read},
					0x0106: {"CurrentTier4SummationDelivered", ZclDataTypeUint48, Read},
					0x0107: {"CurrentTier4SummationReceived", ZclDataTypeUint48, Read},
					0x0108: {"CurrentTier5SummationDelivered", ZclDataTypeUint48, Read},
					0x0109: {"CurrentTier5SummationReceived", ZclDataTypeUint48, Read},
					0x010a: {"CurrentTier6SummationDelivered", ZclDataTypeUint48, Read},
					0x010b: {"CurrentTier6SummationReceived", ZclDataTypeUint48, Read},
					0x010c: {"CurrentTier7SummationDelivered", ZclDataTypeUint48, Read},
					0x010d: {"CurrentTier7SummationReceived", ZclDataTypeUint48, Read},
					0x010e: {"CurrentTier8SummationDelivered", ZclDataTypeUint48, Read},
					0x010f: {"CurrentTier8SummationReceived", ZclDataTypeUint48, Read},
					0x0110: {"CurrentTier9SummationDelivered", ZclDataTypeUint48, Read},
					0x0111: {"CurrentTier9SummationReceived", ZclDataTypeUint48, Read},
					0x0112: {"CurrentTier10SummationDelivered", ZclDataTypeUint48, Read},
					0x0113: {"CurrentTier10SummationReceived", ZclDataTypeUint48, Read},
					0x0114: {"CurrentTier11SummationDelivered", ZclDataTypeUint48, Read},
					0x0115: {"CurrentTier11SummationReceived", ZclDataTypeUint48, Read},
					0x0116: {"CurrentTier12SummationDelivered", ZclDataTypeUint48, Read},
					0x0117: {"CurrentTier12SummationReceived", ZclDataTypeUint48, Read},
					0x0118: {"CurrentTier13SummationDelivered", ZclDataTypeUint48, Read},
					0x0119: {"CurrentTier13SummationReceived", ZclDataTypeUint48, Read},
					0x011a: {"CurrentTier14SummationDelivered", ZclDataTypeUint48, Read},
					0x011b: {"CurrentTier14SummationReceived", ZclDataTypeUint48, Read},
					0x011c: {"CurrentTier15SummationDelivered", ZclDataTypeUint48, Read},
					0x011d: {"CurrentTier15SummationReceived", ZclDataTypeUint48, Read},
					0x0200: {"Status", ZclDataTypeBitmap8, Read},
					0x0201: {"RemainingBatteryLife", ZclDataTypeUint8, Read},
					0x0202: {"HoursInOperation", ZclDataTypeUint24, Read},
					0x0203: {"HoursInFault", ZclDataTypeUint24, Read},
					0x0300: {"UnitofMeasure", ZclDataTypeEnum8, Read},
					0x0301: {"Multiplier", ZclDataTypeUint24, Read},
					0x0302: {"Divisor", ZclDataTypeUint24, Read},
					0x0303: {"SummationFormatting", ZclDataTypeBitmap8, Read},
					0x0304: {"DemandFormatting", ZclDataTypeBitmap8, Read},
					0x0305: {"HistoricalConsumptionFormatting", ZclDataTypeBitmap8, Read},
					0x0306: {"MeteringDeviceType", ZclDataTypeBitmap8, Read},
					0x0307: {"SiteID", ZclDataTypeOctetStr, Read},
					0x0308: {"MeterSerialNumber", ZclDataTypeOctetStr, Read},
					0x0309: {"EnergyCarrierUnitOfMeasure", ZclDataTypeEnum8, Read},
					0x030a: {"EnergyCarrierSummationFormatting", ZclDataTypeBitmap8, Read},
					0x030b: {"EnergyCarrierDemandFormatting", ZclDataTypeBitmap8, Read},
					0x030c: {"TemperatureUnitOfMeasure", ZclDataTypeEnum8, Read},
					0x030d: {"TemperatureFormatting", ZclDataTypeBitmap8, Read},
					0x0400: {"InstantaneousDemand", ZclDataTypeInt24, Read},
					0x0401: {"CurrentDayConsumptionDelivered", ZclDataTypeUint24, Read},
					0x0402: {"CurrentDayConsumptionReceived", ZclDataTypeUint24, Read},
					0x0403: {"PreviousDayConsumptionDelivered", ZclDataTypeUint24, Read},
					0x0404: {"PreviousDayConsumptionReceived", ZclDataTypeUint24, Read},
					0x0405: {"CurrentPartialProfileIntervalStartTimeDelivered", ZclDataTypeUtc, Read},
					0x0406: {"CurrentPartialProfileIntervalStartTimeReceived", ZclDataTypeUtc, Read},
					0x0407: {"CurrentPartialProfileIntervalValueDelivered", ZclDataTypeUint24, Read},
					0x0408: {"CurrentPartialProfileIntervalValueReceived", ZclDataTypeUint24, Read},
					0x0409: {"CurrentDayMaxPressure", ZclDataTypeUint48, Read},
					0x040a: {"CurrentDayMinPressure", ZclDataTypeUint48, Read},
					0x040b: {"PreviousDayMaxPressure", ZclDataTypeUint48, Read},
					0x040c: {"PreviousDayMinPressure", ZclDataTypeUint48, Read},
					0x040d: {"CurrentDayMaxDemand", ZclDataTypeInt24, Read},
					0x040e: {"PreviousDayMaxDemand", ZclDataTypeInt24, Read},
					0x040f: {"CurrentMonthMaxDemand", ZclDataTypeInt24, Read},
					0x0410: {"CurrentYearMaxDemand", ZclDataTypeInt24, Read},
					0x0411: {"CurrentDayMaxEnergyCarrierDemand", ZclDataTypeInt24, Read},
					0x0412: {"PreviousDayMaxEnergyCarrierDemand", ZclDataTypeInt24, Read},
					0x0413: {"CurrentMonthMaxEnergyCarrierDemand", ZclDataTypeInt24, Read},
					0x0414: {"CurrentMonthMinEnergyCarrierDemand", ZclDataTypeInt24, Read},
					0x0415: {"CurrentYearMaxEnergyCarrierDemand", ZclDataTypeInt24, Read},
					0x0416: {"CurrentYearMinEnergyCarrierDemand", ZclDataTypeInt24, Read},
					0x0500: {"MaxNumberOfPeriodsDelivered", ZclDataTypeUint8, Read},
					0x0600: {"CurrentDemandDelivered", ZclDataTypeInt24, Read},
					0x0601: {"DemandLimit", ZclDataTypeUint24, Read},
					0x0602: {"DemandIntegrationPeriod", ZclDataTypeUint8, Read},
					0x0603: {"NumberOfDemandSubintervals", ZclDataTypeUint8, Read},
					0x0800: {"GenericAlarmMask", ZclDataTypeBitmap16, Read | Write},
					0x0801: {"ElectricityAlarmMask", ZclDataTypeBitmap32, Read | Write},
					0x0802: {"GenericFlow/PressureAlarmMask", ZclDataTypeBitmap16, Read | Write},
					0x0803: {"WaterSpecificAlarmMask", ZclDataTypeBitmap16, Read | Write},
					0x0804: {"HeatandCoolingSpecificAlarmMask", ZclDataTypeBitmap16, Read | Write},
					0x0805: {"GasSpecificAlarmMask", ZclDataTypeBitmap16, Read | Write},
					0x8000: {"HistoryElectricClear", ZclDataTypeBoolean, Read | Write}, //鸿雁清空历史电量
				},
				CommandDescriptors: &CommandDescriptors{
					Received: map[uint8]*CommandDescriptor{
						0x00: {"GetProfile", &GetProfile{}},
						0x01: {"RequestMirrorResponse", &RequestMirrorResponse{}},
						0x02: {"MirrorRemoved", &MirrorRemoved{}},
						0x03: {"RequestFastPollMode", &RequestFastPollMode{}},
					},
					Generated: map[uint8]*CommandDescriptor{
						0x00: {"GetProfileResponse", &GetProfileResponse{}},
						0x01: {"RequestMirror", &RequestMirror{}},
						0x02: {"RemoveMirror", &RemoveMirror{}},
						0x03: {"RequestFastPollModeResponse", &RequestFastPollModeResponse{}},
					},
				},
			},
			HEIMANScenes: {
				Name:                 "HEIMANScenes",
				AttributeDescriptors: map[uint16]*AttributeDescriptor{},
				CommandDescriptors: &CommandDescriptors{
					Received: map[uint8]*CommandDescriptor{
						0xf0: {"Cinema", nil},
						0xf1: {"AtHome", nil},
						0xf2: {"Sleep", nil},
						0xf3: {"GoOut", nil},
						0xf4: {"Repast", nil},
					},
					Generated: map[uint8]*CommandDescriptor{},
				},
			},
			HEIMANAirQualityMeasurement: {
				Name: "HEIMANAirQualityMeasurement",
				AttributeDescriptors: map[uint16]*AttributeDescriptor{
					0xf000: {"Language", ZclDataTypeUint8, Read},
					0xf001: {"UnitOfTemperature", ZclDataTypeUint8, Read},
					0xf002: {"BatteryState", ZclDataTypeUint8, Read},
					0xf003: {"PM10MeasuredValue", ZclDataTypeUint16, Read},
					0xf004: {"TVOCMeasuredValue", ZclDataTypeUint8, Read},
					0xf005: {"AQIMeasuredValue", ZclDataTypeUint16, Read},
					0xf006: {"MaxTemperature", ZclDataTypeInt16, Read},
					0xf007: {"MinTemperature", ZclDataTypeInt16, Read},
					0xf008: {"MaxHumidity", ZclDataTypeUint16, Read},
					0xf009: {"MinHumidity", ZclDataTypeUint16, Read},
					0xf00a: {"Disturb", ZclDataTypeUint16, Read},
				},
				CommandDescriptors: &CommandDescriptors{
					Received: map[uint8]*CommandDescriptor{
						0x00: {"SetLanguage", &HEIMANSetLanguage{}},
						0x01: {"SetUnitOfTemperature", &HEIMANSetUnitOfTemperature{}},
						0x02: {"GetTime", nil},
					},
					Generated: map[uint8]*CommandDescriptor{
						0x00: {"SetLanguageResponse", nil},
						0x01: {"SetUnitOfTemperatureResponse", nil},
						0x02: {"GetTimeResponse", nil},
					},
				},
			},
			HEIMANPM25Measurement: {
				Name: "HEIMANPM25Measurement",
				AttributeDescriptors: map[uint16]*AttributeDescriptor{
					0x0000: {"MeasuredValue", ZclDataTypeInt16, Read},
					0x0001: {"MinMeasuredValue", ZclDataTypeInt16, Read},
					0x0002: {"MaxMeasuredValue", ZclDataTypeInt16, Read},
					0x0003: {"Tolerance", ZclDataTypeUint16, Read},
				},
			},
			HEIMANFormaldehydeMeasurement: {
				Name: "HEIMANFormaldehydeMeasurement",
				AttributeDescriptors: map[uint16]*AttributeDescriptor{
					0x0000: {"MeasuredValue", ZclDataTypeUint16, Read},
					0x0001: {"MinMeasuredValue", ZclDataTypeUint16, Read},
					0x0002: {"MaxMeasuredValue", ZclDataTypeUint16, Read},
					0x0003: {"Tolerance", ZclDataTypeUint16, Read},
				},
			},
			HEIMANInfraredRemote: {
				Name:                 "HEIMANInfraredRemote",
				AttributeDescriptors: map[uint16]*AttributeDescriptor{},
				CommandDescriptors: &CommandDescriptors{
					Received: map[uint8]*CommandDescriptor{
						0xf0: {"SendKeyCommand", &HEIMANInfraredRemoteSendKeyCommand{}},
						0xf1: {"StudyKey", &HEIMANInfraredRemoteStudyKey{}},
						0xf3: {"DeleteKey", &HEIMANInfraredRemoteDeleteKey{}},
						0xf4: {"CreateID", &HEIMANInfraredRemoteCreateID{}},
						0xf6: {"GetIDAndKeyCodeList", &HEIMANInfraredRemoteGetIDAndKeyCodeList{}},
					},
					Generated: map[uint8]*CommandDescriptor{
						0xf2: {"StudyKeyResponse", &HEIMANInfraredRemoteStudyKeyResponse{}},
						0xf5: {"CreateIDResponse", &HEIMANInfraredRemoteCreateIDResponse{}},
						0xf7: {"GetIDAndKeyCodeListResponse", &HEIMANInfraredRemoteGetIDAndKeyCodeListResponse{}},
					},
				},
			},
			HONYARScenes: {
				Name: "HONYARScenes",
				AttributeDescriptors: map[uint16]*AttributeDescriptor{
					0x0000: {"ScenesCmd", ZclDataTypeUint8, Read},
				},
			},
			MAILEKEMeasurement: {
				Name: "MAILEKEMeasurement",
				AttributeDescriptors: map[uint16]*AttributeDescriptor{
					0x0000: {"PM1MeasuredValue", ZclDataTypeUint16, Read},
					0x0001: {"PM25MeasuredValue", ZclDataTypeUint16, Read},
					0x0002: {"PM10MeasuredValue", ZclDataTypeUint16, Read},
				},
			},
		},
	}
}

// Clusters Clusters
func (cl *Library) Clusters() map[ID]*Cluster {
	return cl.clusters
}

// Global Global
func (cl *Library) Global() map[uint8]*CommandDescriptor {
	return cl.global
}
