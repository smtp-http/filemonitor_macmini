package notify

import (
    "github.com/smtp-http/filemonitor_macmini/config"
    "fmt"
    "path/filepath"
    "io/ioutil"
    "sync"
)


type LogFile struct {
	Path string
}

var log_instance *LogFile
var log_once sync.Once
 
func GetLogFileInstance() *LogFile {
    log_once.Do(func() {
    	loader := config.GetLoader()
		loader.Load("./logCfg.json",config.GetLogConfig())
        log_instance = &LogFile{}
        log_instance.Path = config.GetLogConfig().LogPath
    })
    return log_instance
}



func (l *LogFile) GetFileList() ([]string,error) {
	var files []string

	return files,nil
}

func (l *LogFile) GetFileContent(fileName string,pos int,line int) (string,error) {
	var content string


	return content,nil
}

