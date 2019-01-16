package conn

import (
	"fmt"
	"github.com/smtp-http/zero"
	"container/list"
	"sync"
)


type TcpServer struct {
	SessionList *list.List
}

func (s *TcpServer) Notify (msg []byte) bool {

	for e := s.SessionList.Front(); e != nil; e = e.Next() {
		session := e.Value.(*zero.Session)
		conn := session.GetConn()
		
		err := conn.SendBytes(msg)
		if err != nil {
			fmt.Printf("send err: %x, remove\n",e)

			s.SessionList.Remove(e)
		}
	}

	return true
}

 
var tcpserver_instance *TcpServer
var tcp_once sync.Once
 
func GetServerInstance() *TcpServer {
    tcp_once.Do(func() {
        tcpserver_instance = &TcpServer{}
    })
    return tcpserver_instance
}





func (s *TcpServer) ServerRun(ip string,port string) {
	
	s.SessionList = list.New()

	ss, err := zero.NewSocketService(ip + ":" + port)
	if err != nil {
		return
	}

	ss.RegMessageHandler(HandleMessage)
	ss.RegConnectHandler(HandleConnect)
	ss.RegDisconnectHandler(HandleDisconnect)

	ss.Serv()
}


func HandleMessage(s *zero.Session, msg *zero.Message) {
	fmt.Println("receive msgID:", msg)
	fmt.Println("receive data:", string(msg.GetData()))
}

func HandleDisconnect(s *zero.Session, err error) {
	fmt.Println(s.GetConn().GetName() + " lost.")
	tcpserver := GetServerInstance()
	for e := tcpserver.SessionList.Front(); e != nil; e = e.Next() {
		session := e.Value.(*zero.Session)
		if s == session {
			tcpserver.SessionList.Remove(e)
		}
	}
	
}

func HandleConnect(s *zero.Session) {
	fmt.Println(s.GetConn().GetName() + " connected." )
	tcpserver := GetServerInstance()
	tcpserver.SessionList.PushBack(s)
}