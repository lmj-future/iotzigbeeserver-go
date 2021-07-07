package dgram

import (
	"net"
	"strconv"
	"strings"

	"github.com/h3c/iotzigbeeserver-go/constant"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globallogger"
)

var udpServer = UDPServer{}

//ServiceSocket ServiceSocket
type ServiceSocket interface {
	On(event string, cb interface{})
	Receive(data []byte) ([]byte, RInfo, error)
	Send(sendBuf []byte, IPPort int, IPAddr string) error
	Close() error
}

//UDPServer UDPServer
type UDPServer struct {
	UDPAddr *net.UDPAddr
	UDPConn *net.UDPConn
}

//RInfo RInfo
type RInfo struct {
	Family  string
	Address string
	Port    int
}

//CreateUDPSocket CreateUDPSocket
func CreateUDPSocket(bindPort int) ServiceSocket {
	//address := "127.0.0.1:8080"
	//udpAddr, err := net.ResolveUDPAddr("udp4", address)
	udpAddr := &net.UDPAddr{
		IP:   net.ParseIP(constant.Constant.LocalHost), //net.IP{0, 0, 0, 0},
		Port: bindPort,
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	//defer conn.Close()
	if err != nil {
		globallogger.Log.Errorln("connect udpserver failed, err:", err.Error())
		//os.Exit(1)
	}

	udpServer.UDPAddr = udpAddr
	udpServer.UDPConn = conn
	globallogger.Log.Errorf("CreateUDPSocket success: %+v", udpServer)
	return udpServer
}

//Send Send
func (udps UDPServer) Send(sendBuf []byte, IPPort int, IPAddr string) error {
	var ipBuilder strings.Builder
	ipBuilder.WriteString(IPAddr)
	ipBuilder.WriteString(":")
	ipBuilder.WriteString(strconv.FormatInt(int64(IPPort), 10))
	udpAddr, _ := net.ResolveUDPAddr("udp4", ipBuilder.String())
	_, err := udps.UDPConn.WriteToUDP(sendBuf, udpAddr)
	if err != nil {
		globallogger.Log.Errorln("send msg err:", err.Error())
	}
	return err
	// conn, err := net.DialUDP("udp",
	// 	&net.UDPAddr{IP: net.ParseIP(constant.Constant.LocalHost), Port: udps.UDPAddr.Port},
	// 	&net.UDPAddr{IP: net.ParseIP(IPAddr), Port: IPPort})
	// if err != nil {
	// 	globallogger.Log.Errorln("send msg err:", err.Error())
	// 	return err
	// } else {
	// 	defer conn.Close()
	// 	conn.Write(sendBuf)
	// 	return nil
	// }
}

//On On
func (udps UDPServer) On(event string, cb interface{}) {
	switch event {
	case "message":
	case "error":
	case "listening":
	}
}

//Receive Receive
func (udps UDPServer) Receive(data []byte) ([]byte, RInfo, error) {
	n, addr, err := udps.UDPConn.ReadFromUDP(data)
	if err != nil {
		globallogger.Log.Errorln("failed read udp msg, error:", err.Error())
		return nil, RInfo{}, err
	}
	return append(data[:0:0], data[:n]...), RInfo{Address: addr.IP.String(), Port: addr.Port, Family: addr.Network()}, err
}

//Close Close
func (udps UDPServer) Close() error {
	return udps.UDPConn.Close()
}
