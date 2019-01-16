package main

import (
	"github.com/smtp-http/filenotify_go/notify"
	"github.com/smtp-http/filenotify_go/conn"
	"github.com/smtp-http/filenotify_go/config"
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
