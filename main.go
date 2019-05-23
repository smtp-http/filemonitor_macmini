package main

import (
	"github.com/smtp-http/filemonitor_macmini/notify"
	"github.com/smtp-http/filemonitor_macmini/config"
)


func main() {

	loader := config.GetLoader()
	loader.Load("./logCfg.json",config.GetLogConfig())


	logFile := notify.GetLogFileInstance()
	logFile.Server.StartServer(config.GetLogConfig().HttpPort)
}
