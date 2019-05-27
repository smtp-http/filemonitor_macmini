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
    "os"
    "io"
    "bufio"
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


func GetLog(file string,start int,lineNum int) ([]string,int){ //
    f, _ := os.Open(file) //日志文件路径
    defer f.Close()
    var resultSlice []string
    buf := bufio.NewReader(f)                                      //读取日志文件里的字符流
    for {                                                          //逐行读取日志文件
        line, err := buf.ReadString('\n')
        resultSlice = append(resultSlice, line)
        if err != nil {
            if err == io.EOF {
                break //表示文件读取完了
            }
        }
    }

    length := len(resultSlice)
    fmt.Println(length) //打印出结果的总条数
    // bubbleSort(resultSlice)       //对结果排序
    if start > length - 1{
        return nil,0
    }

    if lineNum > length - start {
        return resultSlice[0:length - start -1],length - start
    } else {
        return resultSlice[length - start - lineNum - 1:length - start - 1],lineNum
    }
}


/////////////////////////////////////////////////////// file list ////////////////////////////////////////////////////////


type ResGetLogList struct {
    Files []string  `json:"fileList"`
}

func (l *LogFile)GetLogList(w http.ResponseWriter, r *http.Request) {
        
    var res ResGetLogList
    var err error


    res.Files,err = GetSpecificExtensionFiles(l.Path,config.GetLogConfig().Extension)
    if err != nil {
        res.Files = nil;
    }

    ret, _ := json.Marshal(res)
    w.Header().Set("Access-Control-Allow-Origin", "*")
    fmt.Fprint(w, string(ret))

}


////////////////////////////////////////////////////// file content //////////////////////////////////////////////////////
type ReqGetContent struct {
    FileName string     `json:"fileName"`
    Start   int         `json:"start"`
    LineNum int         `json:"lineNum"`
}

type ResGetContent struct {
    Status      string      `json:"status"`
    FileName    string     `json:"fileName"`
    LineAmount  int        `json:"lineAmount"` 
    Content     []string   `json:"Content"`
}

func (l *LogFile)GetLogContent(w http.ResponseWriter, r *http.Request) {
    body, _ := ioutil.ReadAll(r.Body)
    defer r.Body.Close()
    body_str := string(body)
    fmt.Println(body_str)

    var req ReqGetContent
    var res ResGetContent

    if err := json.Unmarshal(body, &req); err == nil {
        fmt.Println(req)

        res.Content,res.LineAmount = GetLog(req.FileName,req.Start,req.LineNum)
        if res.Content == nil {
            res.Status = "fail"
        } else {
            res.Status = "sucess"
        }
        res.FileName = req.FileName
        res.LineAmount = req.LineNum

        ret, _ := json.Marshal(res)
        w.Header().Set("Access-Control-Allow-Origin", "*")
        fmt.Fprint(w, string(ret))
    } else {
        res.Status = "fail"

        fmt.Println(err)
        ret, _ := json.Marshal(res)
        w.Header().Set("Access-Control-Allow-Origin", "*")
        fmt.Fprint(w, string(ret))
    }
}


