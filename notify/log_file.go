package notify

import (
    "github.com/smtp-http/filemonitor_macmini/config"
    "github.com/smtp-http/filemonitor_macmini/conn"
    "fmt"
    //"path/filepath"
    "io/ioutil"
    "sync"
    "net/http"
    "encoding/json"
)


type LogFile struct {
	Path string
    Server *conn.HttpServer
}

var log_instance *LogFile
var log_once sync.Once
 
func GetLogFileInstance() *LogFile {
    log_once.Do(func() {
    	loader := config.GetLoader()
		loader.Load("./logCfg.json",config.GetLogConfig())
        log_instance = &LogFile{}
        log_instance.Path = config.GetLogConfig().LogPath
        log_instance.Server = conn.GetHttpServerInstance()

        log_instance.Server.AddHandleFunc("/FileList",log_instance.GetLogList)
        log_instance.Server.AddHandleFunc("/LogContent",log_instance.GetLogContent)
    })
    return log_instance
}


/*
func (l *LogFile) GetFileList() ([]string,error) {
	var files []string

	return files,nil
}



func (l *LogFile) GetFileContent(fileName string,pos int,line int) (string,error) {
	var content string


	return content,nil
}
*/

////////////////////////////////////////////////////// file content //////////////////////////////////////////////////////
type ReqGetContent struct {
    FileName string     `json:"fileName"`
    Start   int         `json:"start"`
    LineNum int         `json:"lineNum"`
}

type ResGetContent struct {
    FileName string     `json:"fileName"`
    LineAmount  int     `json:"lineAmount"` 
    Content     string  `json:"Content"`
}

func (l *LogFile)GetLogContent(w http.ResponseWriter, r *http.Request) {
    body, _ := ioutil.ReadAll(r.Body)
    defer r.Body.Close()
    body_str := string(body)
    fmt.Println(body_str)

    var req ReqGetContent
   // var e error

    //logFile := notify.GetLogFileInstance()

    if err := json.Unmarshal(body, &req); err == nil {
        fmt.Println(req)
        
        var res ResGetContent
        //res.Content,e = logFile.GetFileContent(req.FileName,req.Start,req.LineNum)
        //if(e != nil){
        //    res.Content = ""
        //}
        
        res.FileName = req.FileName
        res.LineAmount = req.LineNum

        ret, _ := json.Marshal(res)
        fmt.Fprint(w, string(ret))
    } else {
        fmt.Println(err)
    }
}

/////////////////////////////////////////////////////// file list ////////////////////////////////////////////////////////


type ResGetLogList struct {
    Files []string  `json:"fileList"`
}

func (l *LogFile)GetLogList(w http.ResponseWriter, r *http.Request) {
        
    var res ResGetLogList
    //var err error
    // TODO: get file list
    //logFile := notify.GetLogFileInstance()

    //res.Files,error = logFile.GetFileList()

    ret, _ := json.Marshal(res)
    fmt.Fprint(w, string(ret))

}
