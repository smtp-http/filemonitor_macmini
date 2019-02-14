package main

import (
	"github.com/smtp-http/filemonitor_macmini/notify"
	"github.com/smtp-http/filemonitor_macmini/conn"
	"github.com/smtp-http/filemonitor_macmini/config"
)


func main() {

	loader := config.GetLoader()
	loader.Load("./config.json",config.GetConfig())


	conn.GetHttpClientInstance().HttpSetUrl(config.GetConfig().Url)
	tcpserver := conn.GetServerInstance()
	go tcpserver.ServerRun(config.GetConfig().Ip,config.GetConfig().Port)

	tcpclient := new(conn.TcpClient)
	tcpclient.Init(config.GetConfig().Host)
	//tcpclient.SendTest()

	//================== disp ====================
	//disp := new(notify.Dispatcher)

	//disp.Dispatch(config.GetConfig().RootDirectory)

	// ================= updater ==================
	//updater := notify.GetUpdateScannerInstance()

	//updater.ScanFile(config.GetConfig().RootDirectory)

	//================ finder =============


	finder := notify.GetFinderInstance()
	finder.Monitor()

}
