package cluster

// ResetToFactoryDefaultsCommand ResetToFactoryDefaultsCommand
type ResetToFactoryDefaultsCommand struct {
}

// IdentifyCommand IdentifyCommand
type IdentifyCommand struct {
	IdentifyTime uint16
}

// IdentifyQueryCommand IdentifyQueryCommand
type IdentifyQueryCommand struct{}

// TriggerEffectCommand TriggerEffectCommand
type TriggerEffectCommand struct {
	EffectIdentifier uint8
	EffectVariant    uint8
}

// IdentifyQueryResponse IdentifyQueryResponse
type IdentifyQueryResponse struct {
	Timeout uint16
}

// AddSceneCommand AddSceneCommand
type AddSceneCommand struct {
	GroupID        uint16
	SceneID        uint8
	TransitionTime uint16
	SceneName      string
	KeyID          uint8
}

// ViewSceneCommand ViewSceneCommand
type ViewSceneCommand struct {
	GroupID uint16
	SceneID uint8
}

// RemoveSceneCommand RemoveSceneCommand
type RemoveSceneCommand struct {
	GroupID uint16
	SceneID uint8
}

// RemoveAllScenesCommand RemoveAllScenesCommand
type RemoveAllScenesCommand struct {
	GroupID uint16
}

// StoreSceneCommand StoreSceneCommand
type StoreSceneCommand struct {
	GroupID uint16
	SceneID uint8
}

// RecallSceneCommand RecallSceneCommand
type RecallSceneCommand struct {
	GroupID uint16
	SceneID uint8
}

// GetSceneMembership GetSceneMembership
type GetSceneMembership struct {
	GroupID uint16
}

// EnhancedAddSceneCommand EnhancedAddSceneCommand
type EnhancedAddSceneCommand struct {
	GroupID        uint16
	SceneID        uint8
	TransitionTime uint16
	SceneName      string
}

// EnhancedViewSceneCommand EnhancedViewSceneCommand
type EnhancedViewSceneCommand struct {
	GroupID uint16
	SceneID uint8
}

// CopySceneCommand CopySceneCommand
type CopySceneCommand struct {
	Mode        uint8
	FromGroupID uint16
	FromSceneID uint16
	ToGroupID   uint16
	ToSceneID   uint16
}

// AddSceneResponse AddSceneResponse
type AddSceneResponse struct {
	Status  uint8
	GroupID uint16
	SceneID uint8
}

// ViewSceneResponse ViewSceneResponse
type ViewSceneResponse struct {
	Status         uint8
	GroupID        uint16
	SceneID        uint8
	TransitionTime uint16
	SceneName      string
}

// RemoveSceneResponse RemoveSceneResponse
type RemoveSceneResponse struct {
	Status  uint8
	GroupID uint16
	SceneID uint8
}

// RemoveAllScenesResponse RemoveAllScenesResponse
type RemoveAllScenesResponse struct {
	Status  uint8
	GroupID uint16
}

// StoreSceneResponse StoreSceneResponse
type StoreSceneResponse struct {
	Status  uint8
	GroupID uint16
	SceneID uint8
}

// GetSceneMembershipResponse GetSceneMembershipResponse
type GetSceneMembershipResponse struct {
	Status     uint8
	Capacity   uint8
	GroupID    uint16
	SceneCount uint8
	SceneList  []uint8
}

// EnhancedAddSceneResponse EnhancedAddSceneResponse
type EnhancedAddSceneResponse struct {
	Status  uint8
	GroupID uint16
	SceneID uint8
}

// EnhancedViewSceneResponse EnhancedViewSceneResponse
type EnhancedViewSceneResponse struct {
	Status         uint8
	GroupID        uint16
	SceneID        uint8
	TransitionTime uint16
	SceneName      string
}

// CopySceneResponse CopySceneResponse
type CopySceneResponse struct {
	Status      uint8
	FromGroupID uint16
	FromSceneID uint8
}

// OffCommand OffCommand
type OffCommand struct{}

// OnCommand OnCommand
type OnCommand struct{}

// ToggleCommand ToggleCommand
type ToggleCommand struct{}

// OffWithEffectCommand OffWithEffectCommand
type OffWithEffectCommand struct {
	EffectIdentifier uint8
	EffectVariant    uint8
}

// OnWithRecallGlobalSceneCommand OnWithRecallGlobalSceneCommand
type OnWithRecallGlobalSceneCommand struct{}

// OnWithTimedOffCommand OnWithTimedOffCommand
type OnWithTimedOffCommand struct {
	OnOffControl uint8
	OnTime       uint16
	OffWaitTime  uint16
}

// MoveToLevelCommand MoveToLevelCommand
type MoveToLevelCommand struct {
	Level          uint8
	TransitionTime uint16
}

// MoveCommand MoveCommand
type MoveCommand struct {
	MoveMode uint8
	Rate     uint8
}

// StepCommand StepCommand
type StepCommand struct {
	StepMode       uint8
	StepSize       uint8
	TransitionTime uint16
}

// StopCommand StopCommand
type StopCommand struct{}

// MoveToLevelOnOffCommand MoveToLevelOnOffCommand
type MoveToLevelOnOffCommand struct {
	Level          uint8
	TransitionTime uint16
}

// MoveOnOffCommand MoveOnOffCommand
type MoveOnOffCommand struct {
	MoveMode uint8
	Rate     uint8
}

// StepOnOffCommand StepOnOffCommand
type StepOnOffCommand struct {
	StepMode       uint8
	StepSize       uint8
	TransitionTime uint16
}

// StopOnOffCommand StopOnOffCommand
type StopOnOffCommand struct{}

// ResetAlarmCommand ResetAlarmCommand
type ResetAlarmCommand struct {
	AlarmCode         uint8
	ClusterIdentifier uint16
}

// ResetAllAlarmsCommand ResetAllAlarmsCommand
type ResetAllAlarmsCommand struct{}

// GetAlarmCommand GetAlarmCommand
type GetAlarmCommand struct{}

// ResetAlarmLogCommand ResetAlarmLogCommand
type ResetAlarmLogCommand struct{}

// AlarmCommand AlarmCommand
type AlarmCommand struct {
	AlarmCode         uint8
	ClusterIdentifier uint16
}

// GetAlarmResponse GetAlarmResponse
type GetAlarmResponse struct {
	Status            uint8
	AlarmCode         uint8
	ClusterIdentifier uint16
	TimeStamp         uint32
}

// GetProfileInfoResponse GetProfileInfoResponse
type GetProfileInfoResponse struct {
	ProfileCount          uint8
	ProfileIntervalPeriod uint8
	MaxNumberOfIntervals  uint8
	ListOfAttributes      []uint16
}

// GetProfileInfoCommand GetProfileInfoCommand
type GetProfileInfoCommand struct {
}

// GetMeasurementProfileResponse GetMeasurementProfileResponse
type GetMeasurementProfileResponse struct {
	StartTime                  uint32
	Status                     uint8
	ProfileIntervalPeriod      uint8
	NumberOfIntervalsDelivered uint8
	AttributeID                uint8
	AttributeValues            []uint16
}

// GetMeasurementProfileCommand GetMeasurementProfileCommand
type GetMeasurementProfileCommand struct {
	AttributeID       uint16
	StartTime         uint32
	NumberOfIntervals uint8
}

// ZoneEnrollResponse ZoneEnrollResponse
type ZoneEnrollResponse struct {
	ResponseCode uint8
	ZoneID       uint8
}

// InitiateNormalOperationModeCommand InitiateNormalOperationModeCommand
type InitiateNormalOperationModeCommand struct {
}

// InitiateTestModeCommand InitiateTestModeCommand
type InitiateTestModeCommand struct {
	TestModeDuration            uint8
	CurrentZoneSensitivityLevel uint8
}

// ZoneStatusChangeNotificationCommand ZoneStatusChangeNotificationCommand
type ZoneStatusChangeNotificationCommand struct {
	ZoneStatus     uint16
	ExtendedStatus uint8
	ZoneID         uint8
	Delay          uint16
}

// ZoneEnrollCommand ZoneEnrollCommand
type ZoneEnrollCommand struct {
	ZoneType         uint16
	ManufacturerCode uint16
}

// StartWarning StartWarning
type StartWarning struct {
	WarningControl  uint8
	WarningDuration uint16
	StrobeDutyCycle uint8
	StrobeLevel     uint8
}

// Squark Squark
type Squark struct {
	SquarkControl uint8
}

// GetProfile GetProfile
type GetProfile struct{}

// RequestMirrorResponse RequestMirrorResponse
type RequestMirrorResponse struct{}

// MirrorRemoved MirrorRemoved
type MirrorRemoved struct{}

// RequestFastPollMode RequestFastPollMode
type RequestFastPollMode struct{}

// GetProfileResponse GetProfileResponse
type GetProfileResponse struct{}

// RequestMirror RequestMirror
type RequestMirror struct{}

// RemoveMirror RemoveMirror
type RemoveMirror struct{}

// RequestFastPollModeResponse RequestFastPollModeResponse
type RequestFastPollModeResponse struct{}

// HEIMANSetLanguage HEIMANSetLanguage
type HEIMANSetLanguage struct {
	Language uint8
}

// HEIMANSetUnitOfTemperature HEIMANSetUnitOfTemperature
type HEIMANSetUnitOfTemperature struct {
	UnitOfTemperature uint8
}

// HEIMANInfraredRemoteSendKeyCommand HEIMANInfraredRemoteSendKeyCommand
type HEIMANInfraredRemoteSendKeyCommand struct {
	ID      uint8 `json:"ID"`
	KeyCode uint8 `json:"KeyCode"`
}

// HEIMANInfraredRemoteStudyKey HEIMANInfraredRemoteStudyKey
type HEIMANInfraredRemoteStudyKey struct {
	ID      uint8 `json:"ID"`
	KeyCode uint8 `json:"KeyCode"`
}

// HEIMANInfraredRemoteDeleteKey HEIMANInfraredRemoteDeleteKey
type HEIMANInfraredRemoteDeleteKey struct {
	ID      uint8 `json:"ID"`
	KeyCode uint8 `json:"KeyCode"`
}

// HEIMANInfraredRemoteCreateID HEIMANInfraredRemoteCreateID
type HEIMANInfraredRemoteCreateID struct {
	ModelType uint8 `json:"ModelType"`
}

// HEIMANInfraredRemoteGetIDAndKeyCodeList HEIMANInfraredRemoteGetIDAndKeyCodeList
type HEIMANInfraredRemoteGetIDAndKeyCodeList struct{}

// HEIMANInfraredRemoteStudyKeyResponse HEIMANInfraredRemoteStudyKeyResponse
type HEIMANInfraredRemoteStudyKeyResponse struct {
	ID           uint8 `json:"ID"`
	KeyCode      uint8 `json:"KeyCode"`
	ResultStatus uint8 `json:"ResultStatus"`
}

// HEIMANInfraredRemoteCreateIDResponse HEIMANInfraredRemoteCreateIDResponse
type HEIMANInfraredRemoteCreateIDResponse struct {
	ID        uint8 `json:"ID"`
	ModelType uint8 `json:"ModelType"`
}

// IDKeyCode IDKeyCode
// type IDKeyCode struct {
// 	ID        uint8   `json:"ID"`
// 	ModelType uint8   `json:"ModelType"`
// 	KeyNum    uint8   `json:"KeyNum"`
// 	KeyCode   []uint8 `json:"KeyCode"`
// }

// HEIMANInfraredRemoteGetIDAndKeyCodeListResponse HEIMANInfraredRemoteGetIDAndKeyCodeListResponse
type HEIMANInfraredRemoteGetIDAndKeyCodeListResponse struct {
	PacketNumSum        uint8   `json:"ID"`
	CurrentPacketNum    uint8   `json:"CurrentPacketNum"`
	CurrentPacketLength uint8   `json:"CurrentPacketLength"`
	IDKeyCodeList       []uint8 `json:"IDKeyCodeList"`
}
