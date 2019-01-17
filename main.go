package main

import (
	"github.com/smtp-http/filemonitor_macmini/notify"
	"github.com/smtp-http/filemonitor_macmini/conn"
	"github.com/smtp-http/filemonitor_macmini/config"
	"time"
	"fmt"
)


func main() {
	cancel_sig := make(chan string)

	loader := config.GetLoader()
	loader.Load("./config.json",config.GetConfig())


	monitor := notify.GetFileMonitorInstance()
	conn.GetHttpClientInstance().HttpSetUrl(config.GetConfig().Url)
	tcpserver := conn.GetServerInstance()
	go tcpserver.ServerRun(config.GetConfig().Ip,config.GetConfig().Port)

	monitor.SetTcpserver(tcpserver)
	go monitor.Monitor(cancel_sig,config.GetConfig().Path)

	time.Sleep(1000*time.Second)

    cancel_sig <- "cancel"

    for {
        time.Sleep(10*time.Second)
        fmt.Println("go away!")
        break
    }
}
