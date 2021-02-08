package dgram

import (
	"net"
	"strconv"

	"github.com/h3c/iotzigbeeserver-go/globalconstant/globallogger"
)

var udpServer = UDPServer{}

//ServiceSocket ServiceSocket
type ServiceSocket interface {
	On(event string, cb interface{})
	Receive() ([]byte, RInfo, error)
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
		IP:   net.IP{0, 0, 0, 0},
		Port: bindPort,
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	//defer conn.Close()
	if err != nil {
		globallogger.Log.Errorln("connect udpserver failed, err:" + err.Error())
		//os.Exit(1)
	}

	udpServer.UDPAddr = udpAddr
	udpServer.UDPConn = conn
	globallogger.Log.Infoln("CreateUDPSocket success: ", udpServer)
	return udpServer
}

//Send Send
func (udps UDPServer) Send(sendBuf []byte, IPPort int, IPAddr string) error {
	strconv.FormatInt(int64(IPPort), 10)
	udpAddr, err := net.ResolveUDPAddr("udp4", IPAddr+":"+strconv.FormatInt(int64(IPPort), 10))
	// globallogger.Log.Infoln("send addr", udpAddr)
	_, err = udps.UDPConn.WriteToUDP(sendBuf, udpAddr)
	if err != nil {
		globallogger.Log.Errorln("send msg err:", err.Error())
	}
	return err
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
func (udps UDPServer) Receive() ([]byte, RInfo, error) {
	// 最大读取数据大小
	data := make([]byte, 1024)
	rInfo := RInfo{}
	n, addr, err := udps.UDPConn.ReadFromUDP(data)
	if err != nil {
		globallogger.Log.Errorln("failed read udp msg, error: " + err.Error())
		return nil, rInfo, err
	}
	// globallogger.Log.Infoln("receive addr:", addr)
	//str := string(data[:n])
	msg := data[:n]
	rInfo.Address = addr.IP.String()
	rInfo.Port = addr.Port
	rInfo.Family = addr.Network()
	//fmt.Println("receive from client, msg:", msg)
	//fmt.Printf("receive from client, str:%x", msg)
	//<-limitChan
	return msg, rInfo, err
}

//Close Close
func (udps UDPServer) Close() error {
	return udps.UDPConn.Close()
}
