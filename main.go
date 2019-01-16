package main

import (
	"github.com/smtp-http/filemonitor_macmini/notify"
	"github.com/smtp-http/filemonitor_macmini/conn"
	"github.com/smtp-http/filemonitor_macmini/config"
)


func main() {
	loader := config.GetLoader()
	loader.Load("./config.json",config.GetConfig())


	monitor := notify.GetFileMonitorInstance()
	conn.GetHttpClientInstance().HttpSetUrl(config.GetConfig().Url)
	tcpserver := conn.GetServerInstance()
	go tcpserver.ServerRun(config.GetConfig().Ip,config.GetConfig().Port)

	monitor.SetTcpserver(tcpserver)
	monitor.Monitor()
}
