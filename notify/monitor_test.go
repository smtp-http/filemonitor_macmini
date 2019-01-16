package notify
import (
	"testing"
	"github.com/smtp-http/filemonitor_macmini/conn"
)



func Test_tcpservice(t *testing.T) {
	monitor := GetFileMonitorInstance()
	tcpserver := conn.GetServerInstance()
	go tcpserver.ServerRun("0.0.0.0","6688")
	monitor.SetTcpserver(tcpserver)
	monitor.Monitor()
}
