package config

import (
	"time"

	"github.com/lib/pq"
)

// ZoneStatus ZoneStatus
type ZoneStatus struct {
	Data           uint16 `json:"data" bson:"data,omitempty" gorm:"column:zonestatusdata"`
	Alarm1         uint16 `json:"alarm1" bson:"alarm1,omitempty" gorm:"column:alarm1"`
	Alarm2         uint16 `json:"alarm2" bson:"alarm2,omitempty" gorm:"column:alarm2"`
	Tamper         uint16 `json:"tamper" bson:"tamper,omitempty" gorm:"column:tamper"`
	Battery        uint16 `json:"battery" bson:"battery,omitempty" gorm:"column:battery"`
	SupervisionRpt uint16 `json:"supervisionRpt" bson:"supervisionRpt,omitempty" gorm:"column:supervisionrpt"`
	RestoreRpt     uint16 `json:"restoreRpt" bson:"restoreRpt,omitempty" gorm:"column:restorerpt"`
	Trouble        uint16 `json:"trouble" bson:"trouble,omitempty" gorm:"column:trouble"`
	ACMains        uint16 `json:"ACMains" bson:"ACMains,omitempty" gorm:"column:acmains"`
	Test           uint16 `json:"test" bson:"test,omitempty" gorm:"column:test"`
	BatteryDefect  uint16 `json:"batteryDefect" bson:"batteryDefect,omitempty" gorm:"column:batterydefect"`
}

// ZoneStatusChangeNotification ZoneStatusChangeNotification
type ZoneStatusChangeNotification struct {
	ZoneStatus     ZoneStatus `json:"zoneStatus" bson:"zoneStatus,omitempty" gorm:"embedded"`
	ExtendedStatus uint8      `json:"extendedStatus" bson:"extendedStatus,omitempty" gorm:"column:extendedstatus"`
	ZoneID         uint8      `json:"zoneID" bson:"zoneID,omitempty" gorm:"column:zoneid"`
	Delay          uint16     `json:"delay" bson:"delay,omitempty" gorm:"column:delay"`
}

//Attribute Attribute
type Attribute struct {
	Status                       []string                     `json:"status" bson:"status,omitempty" gorm:"-"`
	StatusPG                     pq.StringArray               `json:"statuspg" bson:"-" gorm:"type:text[];column:statuspg"`
	ZoneStatusChangeNotification ZoneStatusChangeNotification `json:"zoneStatusChangeNotification" bson:"zoneStatusChangeNotification,omitempty" gorm:"embedded"`
	LatestCommand                string                       `json:"latestCommand" bson:"latestCommand,omitempty" gorm:"column:latestcommand"`
	Data                         string                       `json:"data" bson:"data,omitempty" gorm:"column:attributedata"`
}

//BindInfo BindInfo
type BindInfo struct {
	ClusterIDIn  []string `json:"clusterIDIn" bson:"clusterIDIn,omitempty"`
	ClusterIDOut []string `json:"clusterIDOut" bson:"clusterIDOut,omitempty"`
	ClusterID    []string `json:"clusterID" bson:"clusterID,omitempty"`
	SrcEndpoint  string   `json:"srcEndpoint" bson:"srcEndpoint,omitempty"`
	DstAddrMode  string   `json:"dstAddrMode" bson:"dstAddrMode,omitempty"`
	DstAddress   string   `json:"dstAddress" bson:"dstAddress,omitempty"`
	DstEndpoint  string   `json:"dstEndpoint" bson:"dstEndpoint,omitempty"`
}

// TableName TableName
func (t TerminalInfo) TableName() string {
	return "terminalinfos"
}

//TerminalInfo TerminalInfo
type TerminalInfo struct {
	ID                         uint           `json:"id" bson:"-" gorm:"primary_key"`
	TmnName                    string         `json:"tmnName" bson:"tmnName,omitempty" gorm:"column:tmnname"`
	DevEUI                     string         `json:"devEUI" bson:"devEUI,omitempty" gorm:"column:deveui;uniqueIndex"`
	UserName                   string         `json:"userName" bson:"userName,omitempty" gorm:"column:username"`
	ScenarioID                 string         `json:"scenarioID" bson:"scenarioID,omitempty" gorm:"column:scenarioid"`
	OIDIndex                   string         `json:"OIDIndex" bson:"OIDIndex,omitempty" gorm:"column:oidindex"`
	ProfileID                  string         `json:"profileID" bson:"profileID,omitempty" gorm:"column:profileid"`
	FirmTopic                  string         `json:"firmTopic" bson:"firmTopic,omitempty" gorm:"column:firmtopic"` //厂商ID
	ManufacturerName           string         `json:"manufacturerName" bson:"manufacturerName,omitempty" gorm:"column:manufacturername"`
	NwkAddr                    string         `json:"nwkAddr" bson:"nwkAddr,omitempty" gorm:"column:nwkaddr"`                                                       //设备网络地址（同设备短地址）
	CapabilityFlags            string         `json:"capabilityFlags" bson:"capabilityFlags,omitempty" gorm:"column:capabilityflags"`                               //设备功能标记 &关系 00:终端设备 01:协调器 02:全功能设备 04:长时间供电 08:有省电功能 40:有加密功能 80:
	JoinType                   string         `json:"joinType" bson:"joinType,omitempty" gorm:"column:jointype"`                                                    //设备入网方式 00:第一次入网 01:不加密重新入网 02:加密重新入网
	TmnType                    string         `json:"tmnType" bson:"tmnType,omitempty" gorm:"column:tmntype"`                                                       //设备类型，single_switch/double_switch/triple_switch/lamp/......
	TmnType2                   string         `json:"tmnType2" bson:"tmnType2,omitempty" gorm:"column:tmntype2"`                                                    //终端管理对应的设备型号
	Attribute                  Attribute      `json:"attribute" bson:"attribute,omitempty" gorm:"embedded"`                                                         //设备的属性列表，依赖于tmnType，如tmnType: 'three_phase_switch'，此字段为attribute: { status: ['ON'/'OFF', 'ON'/'OFF', 'ON'/'OFF'] }
	Temperature                float64        `json:"temperature" bson:"temperature,omitempty" gorm:"column:temperature"`                                           //温度
	Humidity                   float64        `json:"humidity" bson:"humidity,omitempty" gorm:"column:humidity"`                                                    //湿度
	PM25                       uint64         `json:"PM25" bson:"PM25,omitempty" gorm:"column:pm25"`                                                                //PM2.5
	Power                      int            `json:"power" bson:"power,omitempty" gorm:"column:power"`                                                             //电量
	CurrentSummationDelivered  uint64         `json:"currentSummationDelivered" bson:"currentSummationDelivered,omitempty" gorm:"column:currentsummationdelivered"` //CurrentSummationDelivered
	InstantaneousDemand        int64          `json:"instantaneousDemand" bson:"instantaneousDemand,omitempty" gorm:"column:instantaneousdemand"`                   //InstantaneousDemand
	Language                   string         `json:"language" bson:"language,omitempty" gorm:"column:language"`                                                    //语言
	UnitOfTemperature          string         `json:"unitOfTemperature" bson:"unitOfTemperature,omitempty" gorm:"column:unitoftemperature"`                         //温度单位
	BatteryState               string         `json:"batteryState" bson:"batteryState,omitempty" gorm:"column:batterystate"`                                        //电池状态
	PM10                       string         `json:"PM10" bson:"PM10,omitempty" gorm:"column:pm10"`                                                                //PM10
	Formaldehyde               float64        `json:"formaldehyde" bson:"formaldehyde,omitempty" gorm:"column:formaldehyde"`                                        //甲醛
	TVOC                       string         `json:"TVOC" bson:"TVOC,omitempty" gorm:"column:tvoc"`                                                                //TVOC
	AQI                        string         `json:"AQI" bson:"AQI,omitempty" gorm:"column:aqi"`                                                                   //AQI
	Disturb                    string         `json:"disturb" bson:"disturb,omitempty" gorm:"column:disturb"`                                                       //勿扰模式开关
	MaxTemperature             string         `json:"maxTemperature" bson:"maxTemperature,omitempty" gorm:"column:maxtemperature"`                                  //温度阈值上限
	MinTemperature             string         `json:"minTemperature" bson:"minTemperature,omitempty" gorm:"column:mintemperature"`                                  //温度阈值下限
	MaxHumidity                string         `json:"maxHumidity" bson:"maxHumidity,omitempty" gorm:"column:maxhumidity"`                                           //湿度阈值上限
	MinHumidity                string         `json:"minHumidity" bson:"minHumidity,omitempty" gorm:"column:minhumidity"`                                           //湿度阈值下限
	EndpointCount              int            `json:"endpointCount" bson:"endpointCount,omitempty" gorm:"column:endpointcount"`                                     // default 0,设备端口号总数
	Endpoint                   []string       `json:"endpoint" bson:"endpoint,omitempty" gorm:"-"`                                                                  //default: "00",设备端口号
	EndpointPG                 pq.StringArray `json:"endpointPG" bson:"-" gorm:"type:text[];column:endpointpg"`
	BindInfo                   []BindInfo     `json:"bindInfo" bson:"bindInfo,omitempty" gorm:"-"` //{ clusterIDIn: [],//设备支持的输入簇, clusterIDOut: [],//设备支持的输出簇, clusterID: []绑定的簇,// srcEndpoint: 源设备端口号, dstAddrMode: '03', dstAddress: 目的设备Mac地址, dstEndpoint: 目的设备端口号}
	BindInfoPG                 string         `json:"bindInfoPG" bson:"-" gorm:"column:bindinfopg"`
	NeighborTableListRecords   []string       `json:"neighborTableListRecords" bson:"neighborTableListRecords,omitempty" gorm:"-"` //TerminalNetworkInfoUp结果
	NeighborTableListRecordsPG pq.StringArray `json:"neighborTableListRecordsPG" bson:"-" gorm:"type:text[];column:neighbortablelistrecords"`
	Online                     bool           `json:"online" bson:"online,omitempty" gorm:"column:online"`             //在线状态
	LeaveState                 bool           `json:"leaveState" bson:"leaveState,omitempty" gorm:"column:leavestate"` //离网状态
	PermitJoin                 bool           `json:"permitJoin" bson:"permitJoin,omitempty" gorm:"column:permitjoin"` //允许其他设备加入该设备
	ACMac                      string         `json:"ACMac" bson:"ACMac,omitempty" gorm:"column:acmac"`
	APMac                      string         `json:"APMac" bson:"APMac,omitempty" gorm:"column:apmac"`
	T300ID                     string         `json:"T300ID" bson:"T300ID,omitempty" gorm:"column:t300id"`
	FirstAddr                  string         `json:"firstAddr" bson:"firstAddr,omitempty" gorm:"column:firstaddr"`
	SecondAddr                 string         `json:"secondAddr" bson:"secondAddr,omitempty" gorm:"column:secondaddr"`
	ThirdAddr                  string         `json:"thirdAddr" bson:"thirdAddr,omitempty" gorm:"column:thirdaddr"`
	ModuleID                   string         `json:"moduleID" bson:"moduleID,omitempty" gorm:"column:moduleid"`
	ChildNum                   int            `json:"childNum" bson:"childNum,omitempty" gorm:"column:childnum"`
	ChildNwkAddr               []string       `json:"childNwkAddr" bson:"childNwkAddr,omitempty" gorm:"-"` //设备孩子节点列表，包含2个字节的短地址
	ChildNwkAddrPG             pq.StringArray `json:"childNwkAddrPG" bson:"-" gorm:"type:text[];column:childnwkaddr"`
	UpdateTime                 time.Time      `json:"updateTime" bson:"updateTime,omitempty" gorm:"column:updatetime"`
	IsDiscovered               bool           `json:"isDiscovered" bson:"isDiscovered,omitempty" gorm:"column:isdiscovered"`
	IsNeedBind                 bool           `json:"isNeedBind" bson:"isNeedBind,omitempty" gorm:"column:isneedbind"`
	IsReadBasic                bool           `json:"isReadBasic" bson:"isReadBasic,omitempty" gorm:"column:isreadbasic"`
	IsExist                    bool           `json:"isExist" bson:"isExist,omitempty" gorm:"column:isexist"`          //应用是否有添加过
	CreateTime                 time.Time      `json:"createTime" bson:"createTime,omitempty" gorm:"column:createtime"` //终端首次上线时间
	UDPVersion                 string         `json:"UDPVersion" bson:"UDPVersion,omitempty" gorm:"column:udpversion"` //UDP version
	Interval                   int            `json:"interval" bson:"interval,omitempty" gorm:"column:interval"`
}

//TerminalInfo2 TerminalInfo2
type TerminalInfo2 struct {
	DevEUI           string    `json:"devEUI" bson:"devEUI,omitempty"`
	ProfileID        string    `json:"profileID" bson:"profileID,omitempty"`
	ManufacturerName string    `json:"manufacturerName" bson:"manufacturerName,omitempty"`
	TmnType          string    `json:"tmnType" bson:"tmnType,omitempty"`       //设备类型，single_switch/double_switch/triple_switch/lamp/......
	Online           bool      `json:"online" bson:"online,omitempty"`         //在线状态
	LeaveState       bool      `json:"leaveState" bson:"leaveState,omitempty"` //离网状态
	FirstAddr        string    `json:"firstAddr" bson:"firstAddr,omitempty"`
	SecondAddr       string    `json:"secondAddr" bson:"secondAddr,omitempty"`
	ThirdAddr        string    `json:"thirdAddr" bson:"thirdAddr,omitempty"`
	IsExist          bool      `json:"isExist" bson:"isExist,omitempty"`       //应用是否有添加过
	CreateTime       time.Time `json:"createTime" bson:"createTime,omitempty"` //终端首次上线时间
}

// TableName TableName
func (t SocketInfo) TableName() string {
	return "socketinfos"
}

//SocketInfo SocketInfo
type SocketInfo struct {
	ID         uint      `json:"id" bson:"-" gorm:"primary_key"`
	ACMac      string    `json:"ACMac" bson:"ACMac,omitempty" gorm:"column:acmac"`
	APMac      string    `json:"APMac" bson:"APMac,omitempty" gorm:"column:apmac"`
	Family     string    `json:"family" bson:"family,omitempty" gorm:"column:family"` //源socket的IP协议
	IPAddr     string    `json:"IPAddr" bson:"IPAddr,omitempty" gorm:"column:ipaddr"` //源socket的IP
	IPPort     int       `json:"IPPort" bson:"IPPort,omitempty" gorm:"column:ipport"` //源socket的port
	UpdateTime time.Time `json:"updateTime" bson:"updateTime,omitempty" gorm:"column:updatetime"`
	CreateTime time.Time `json:"createTime" bson:"createTime,omitempty" gorm:"column:createtime"`
}

// Interval Interval
type Interval struct {
	Interval string `json:"interval"`
}

//DevMgrInfo DevMgrInfo
type DevMgrInfo struct {
	TmnType     string     `json:"tmnType"`
	TmnName     string     `json:"tmnName"`
	TmnDevSN    string     `json:"tmnDevSN"`
	TmnOIDIndex string     `json:"tmnOIDIndex"`
	SceneID     string     `json:"sceneID"`
	TenantID    string     `json:"tenantID"`
	FirmTopic   string     `json:"firmTopic"`
	LinkType    string     `json:"linkType"`
	Property    []Interval `json:"property"`
}
