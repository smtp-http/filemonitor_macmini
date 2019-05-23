package conn

import (
	"fmt"
	"net/http"
	"sync"
	"encoding/json"
	"io/ioutil"
	"github.com/smtp-http/filemonitor_macmini/notify"
)


////////////////////////////////////////////////// http server /////////////////////////////////////////////////////////


type HttpServer struct {
	Port string
}


var instance *HttpServer
var once sync.Once
 
func GetHttpServerInstance() *HttpServer {
    once.Do(func() {
        instance = &HttpServer{}
    })
    return instance
}


func (s *HttpServer) StartServer(port string) {

	http.HandleFunc("/LogContent",GetLogContent)
	http.HandleFunc("/FileList",GetLogList)

	http.ListenAndServe(":3001", nil)
}



////////////////////////////////////////////////////// file content //////////////////////////////////////////////////////
type ReqGetContent struct {
	FileName string 	`json:"fileName"`
	Start	int 		`json:"start"`
	LineNum int 		`json:"lineNum"`
}

type ResGetContent struct {
	FileName string 	`json:"fileName"`
	LineAmount 	int 	`json:"lineAmount"`	
	Content 	string 	`json:"Content"`
}
func GetLogContent(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
    defer r.Body.Close()
	body_str := string(body)
	fmt.Println(body_str)

	var req ReqGetContent
	var e error

	logFile := notify.GetLogFileInstance()

	if err := json.Unmarshal(body, &req); err == nil {
		fmt.Println(req)
		
		var res ResGetContent
		res.Content,e = logFile.GetFileContent(req.FileName,req.Start,req.LineNum)
		if(e != nil){
			res.Content = ""
		}
		
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
	[]Files	string 	`json:"fileList"`
}

func GetLogList(w http.ResponseWriter, r *http.Request) {
		
	var res ResGetLogList
	var err error
	// TODO: get file list
	logFile := notify.GetLogFileInstance()

	res.Files,error = logFile.GetFileList()

	ret, _ := json.Marshal(res)
	fmt.Fprint(w, string(ret))

}