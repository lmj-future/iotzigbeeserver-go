package common

//定义commandType常量
const (
	ReadBasic                   string = "readBasic"
	SwitchCommand               string = "switchCommand"
	GetSwitchState              string = "getSwitchState"
	SocketCommand               string = "socketCommand"
	GetSocketState              string = "getSocketState"
	ChildLock                   string = "childLock"
	OnOffUSB                    string = "onOffUSB"
	BackGroundLight             string = "backGroundLight"
	PowerOffMemory              string = "powerOffMemory"
	HistoryElectricClear        string = "historyElectricClear"
	SensorWriteReq              string = "sensorWriteReq"
	SensorEnrollRsp             string = "sensorEnrollRsp"
	SensorReadReq               string = "sensorReadReq"
	IntervalCommand             string = "intervalCommand"
	WindowCommand               string = "windowCommand"
	GetWindowState              string = "getWindowState"
	SetThreshold                string = "setThreshold"
	StartWarning                string = "startWarning"
	ReadAttribute               string = "readAttribute"
	HonyarSocketRead            string = "honyarSocketRead"
	HeimanScenesDefaultResponse string = "heimanScenesDefaultResponse"
	AddScene                    string = "addScene"
	ViewScene                   string = "viewScene"
	RemoveScene                 string = "removeScene"
	RemoveAllScenes             string = "removeAllScenes"
	StoreScene                  string = "storeScene"
	RecallScene                 string = "recallScene"
	GetSceneMembership          string = "getSceneMembership"
	EnhancedAddScene            string = "enhancedAddScene"
	EnhancedViewScene           string = "enhancedViewScene"
	CopyScene                   string = "copyScene"
	SendKeyCommand              string = "sendKeyCommand"
	StudyKey                    string = "studyKey"
	DeleteKey                   string = "deleteKey"
	CreateID                    string = "createID"
	GetIDAndKeyCodeList         string = "getIDAndKeyCodeList"
	SyncTime                    string = "syncTime"
	SetLanguage                 string = "setLanguage"
	SetUnitOfTemperature        string = "setUnitOfTemperature"
)

//ZclUpMsg 定义上行ZCL数据结构体
type ZclUpMsg struct {
	MsgType     string
	DevEUI      string
	ProfileID   uint16
	GroupID     uint16
	ClusterID   uint16
	SrcEndpoint uint8
	ZclData     string
}

//Command 定义Command
type Command struct {
	Cmd                     uint8
	DstEndpointIndex        int
	AttributeName           string
	AttributeID             uint16
	Value                   interface{}
	Interval                uint16
	Disturb                 uint16
	MaxTemperature          int16
	MinTemperature          int16
	MaxHumidity             uint16
	MinHumidity             uint16
	WarningControl          uint8
	WarningTime             uint16
	TransactionID           uint8
	GroupID                 uint16
	SceneID                 uint8
	TransitionTime          uint16
	SceneName               string
	KeyID                   uint8
	Mode                    uint8
	FromGroupID             uint16
	FromSceneID             uint16
	ToGroupID               uint16
	ToSceneID               uint16
	InfraredRemoteID        uint8
	InfraredRemoteKeyCode   uint8
	InfraredRemoteModelType uint8
}

//ZclDownMsg 定义下行ZCL数据结构体
type ZclDownMsg struct {
	MsgType     string
	DevEUI      string
	T300ID      string
	CommandType string
	Command     Command
	ClusterID   uint16
	MsgID       interface{}
}

// DownMsg 下行交互数据结构体
type DownMsg struct {
	ProfileID        uint16
	ClusterID        uint16
	DstEndpointIndex int
	ZclData          string
	SN               uint8
	MsgType          string
	DevEUI           string
	MsgID            interface{}
}
