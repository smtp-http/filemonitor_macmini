package notify
import (
	"testing"
	"time"
	//"github.com/smtp-http/filemonitor_macmini/config"
	"fmt"
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
}

func Test_ProcessTargetFolder(t *testing.T) {
	loader := config.GetLoader()
	loader.Load("../config.json",config.GetConfig())

	cur_dir := "D:\\tmp"
	FindTargetFolder(cur_dir)
}
*/

func Test_timestampFile(t *testing.T) {
	tm,err := GetLastProcessTime()
	if err != nil {
		fmt.Println("err:",err)
	}

	fmt.Printf("===== last time: %v\n",tm)

	time_stamp := time.Now().Unix()
	fmt.Printf("+++++ update time: %v \n",time_stamp)
	err = UpdateLastProcessTime(time_stamp)
	if err != nil {
		fmt.Println("------ update err :",err)
	}
}