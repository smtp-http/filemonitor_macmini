package notify
import (
	"testing"
	//"time"
	//"github.com/smtp-http/filemonitor_macmini/conn"
)


/*
func Test_tcpservice(t *testing.T) {
	monitor := GetFileMonitorInstance()
	tcpserver := conn.GetServerInstance()
	go tcpserver.ServerRun("0.0.0.0","6688")
	monitor.SetTcpserver(tcpserver)
	monitor.Monitor()


	
}



func Test_DirectoryMonitor(t *testing.T) {
	cancel_sig := make(chan string)

	dm := new(DirectoryMonitor)
	go dm.StartMonitor(cancel_sig,"D:\\tmp")

	time.Sleep(10*time.Second)

    cancel_sig <- "cancel"
}


func Test_Dispatcher(t *testing.T) {
	disp := new(Dispatcher)

	disp.Dispatch()

	//time.Sleep(10*time.Second)
}*/

func Test_FindTargetFolder(t *testing.T) {
	cur_dir := "D:\\tmp"
	FindTargetFolder(cur_dir)
}