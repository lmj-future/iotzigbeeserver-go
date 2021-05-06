package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

// 初始化 udpReceiveTotal counter类型指标， 表示接收udp总流量
var udpReceiveTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "udp_receive",
		Help: "the total number of udp frame received",
	},
	[]string{},
)

// 初始化 udpReceiveByLabel counter类型指标， 表示以某一个label为统计单位的接收udp总流量
var udpReceiveByLabel = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "udp_receive_by_label",
		Help: "the total number of udp frame received by label",
	},
	[]string{"label"},
)

// 初始化 udpReceiveByAddress counter类型指标， 表示以源地址为统计单位的接收udp总流量
var udpReceiveByAddress = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "udp_receive_by_address",
		Help: "the total number of udp frame received by address",
	},
	[]string{"address"},
)

// 初始化 udpReceiveByGWSN counter类型指标， 表示以网关序列号为统计单位的接收udp总流量
var udpReceiveByGwSN = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "udp_receive_by_gwSN",
		Help: "the total number of udp frame received by gwSN",
	},
	[]string{"gwSN"},
)

// 初始化 udpReceiveByGWSN counter类型指标， 表示以网关序列号和模块ID为统计单位的接收udp总流量
var udpReceiveByGwSNAndModuleID = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "udp_receive_by_gwSN_and_moduleID",
		Help: "the total number of udp frame received by gwSN and moduleID",
	},
	[]string{"gwSN", "moduleID"},
)

// 初始化 udpReceiveByDevSN counter类型指标， 表示以设备序列号为统计单位的接收udp总流量
var udpReceiveByDevSN = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "udp_receive_by_devSN",
		Help: "the total number of udp frame received by devSN",
	},
	[]string{"devSN"},
)

// 初始化 udpSendTotal counter类型指标， 表示发送udp总流量
var udpSendTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "udp_send",
		Help: "the total number of udp frame sended",
	},
	[]string{},
)

func init() {
	// 注册监控指标
	prometheus.MustRegister(udpReceiveTotal)
	prometheus.MustRegister(udpReceiveByLabel)
	prometheus.MustRegister(udpReceiveByAddress)
	prometheus.MustRegister(udpReceiveByGwSN)
	prometheus.MustRegister(udpReceiveByGwSNAndModuleID)
	prometheus.MustRegister(udpReceiveByDevSN)
	prometheus.MustRegister(udpSendTotal)
}

// 接收udp总流量统计计数
func CountUdpReceiveTotal() {
	udpReceiveTotal.With(prometheus.Labels{}).Inc()
}

// 以源地址为统计单位的接收udp总流量统计计数
func CountUdpReceiveByLabel(label string) {
	udpReceiveByLabel.With(prometheus.Labels{"label": label}).Inc()
}

// 以源地址为统计单位的接收udp总流量统计计数
func CountUdpReceiveByAddress(address string) {
	udpReceiveByAddress.With(prometheus.Labels{"address": address}).Inc()
}

// 以网关序列号为统计单位的接收udp总流量统计计数
func CountUdpReceiveByGwSN(gwSN string) {
	udpReceiveByGwSN.With(prometheus.Labels{"gwSN": gwSN}).Inc()
}

// 以网关序列号和模块ID为统计单位的接收udp总流量统计计数
func CountUdpReceiveByGwSNAndModuleID(gwSN string, moduleID string) {
	udpReceiveByGwSNAndModuleID.With(prometheus.Labels{"gwSN": gwSN, "moduleID": moduleID}).Inc()
}

// 以设备序列号为统计单位的接收udp总流量统计计数
func CountUdpReceiveByDevSN(devSN string) {
	udpReceiveByDevSN.With(prometheus.Labels{"devSN": devSN}).Inc()
}

// 发送udp总流量统计计数
func CountUdpSendTotal() {
	udpSendTotal.With(prometheus.Labels{}).Inc()
}

// prometheus Start
// func PrometheusStart() {
// 	http.Handle("/metrics", promhttp.Handler())
// 	// expose prometheus metrics接口
// 	http.ListenAndServe(":8081", nil)
// }
