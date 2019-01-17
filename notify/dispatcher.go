package notify

import (
	"fmt"
)

type Dispatcher struct {

}


type DispMsg struct {
	MonitorName string
	NextPath       string
}

var DispMsgCh chan DispMsg = make(chan DispMsg)

func (d *Dispatcher) Dispatch () {
	//cancel_sig := make(chan string)

	root_monitor := GetRootDirectoryMonitorInstance()

	go root_monitor.StartMonitor(nil,"D:\\tmp")

	for {
		fmt.Println("disp begin!")
        select {
            case dm := <- DispMsgCh :
                fmt.Printf("disp msg ch: %s  :  %s  \n",dm.MonitorName,dm.NextPath)
                if dm.MonitorName == "root" {

                } else if dm.MonitorName == "middle" {

                } else if dm.MonitorName == "file" {

                } else {
                	fmt.Printf("wrong MonitorName: %s\n",dm.MonitorName)
                }
        }
    }
}