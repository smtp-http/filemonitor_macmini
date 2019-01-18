package notify

import (
	"fmt"
)

type Dispatcher struct {

}


type DispMsg struct {
	MonitorName string
	Action      string
	NextPath    string
}

var DispMsgCh chan DispMsg = make(chan DispMsg)

func (d *Dispatcher) Dispatch () {
	md_cancel_sig := make(chan string)
	file_cancel_sig := make(chan string)

	root_monitor := GetRootDirectoryMonitorInstance()

	go root_monitor.StartMonitor(nil,"D:\\tmp")

	for {
		fmt.Println("disp begin!")
        select {
            case dm := <- DispMsgCh :
                fmt.Printf("disp msg ch: %s  :  %s  \n",dm.MonitorName,dm.NextPath)
                if dm.MonitorName == "root" {
                	if dm.Action == "create_dir" {                                         // 停止对 middle dir 的监听
	                	if GetMiddleDirectoryMonitorInstance().IsRunning {
	                		GetMiddleDirectoryMonitorInstance().MonitorPath = dm.NextPath
	                		md_cancel_sig <- "cancel"
	                	} else {
	                		GetMiddleDirectoryMonitorInstance().IsRunning = true
	                		GetMiddleDirectoryMonitorInstance().StartMonitor(md_cancel_sig,dm.NextPath)

	                		go FindTargetFolder(dm.NextPath)
	                	}


                	} else if dm.Action == "stop_ok" {
                		GetMiddleDirectoryMonitorInstance().IsRunning = true
	                	GetMiddleDirectoryMonitorInstance().StartMonitor(md_cancel_sig,GetMiddleDirectoryMonitorInstance().MonitorPath)

	                	go FindTargetFolder(GetMiddleDirectoryMonitorInstance().MonitorPath)
                	}

                } else if dm.MonitorName == "middle" {
                	if dm.Action == "create_dir" {
                		go FindTargetFolder(dm.NextPath)
                	}
                } else if dm.MonitorName == "file" {
                	if dm.Action == "create_file_monitor" {   // 这条消息来自 FindTargetFolder（）
                		if GetFileMonitorInstance().IsRunning {
                			GetFileMonitorInstance().MonitorPath = dm.NextPath
                			md_cancel_sig <- "cancel"
                		} else {
                			GetFileMonitorInstance().IsRunning = true
                			GetFileMonitorInstance().StartMonitor(file_cancel_sig,dm.NextPath)
                		}
                	} else if dm.Action == "stop_ok" {
                		GetFileMonitorInstance().IsRunning = true
                		GetFileMonitorInstance().StartMonitor(file_cancel_sig,GetFileMonitorInstance().MonitorPath)
                	}

                } else {
                	fmt.Printf("wrong MonitorName: %s\n",dm.MonitorName)
                }
        }
    }
}