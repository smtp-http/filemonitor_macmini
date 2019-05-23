package conn

import (
	"fmt"
	"net/http"
	"sync"
	//"encoding/json"
	//"io/ioutil"
	//"github.com/smtp-http/filemonitor_macmini/notify"
)


////////////////////////////////////////////////// http server /////////////////////////////////////////////////////////


type HttpServer struct {
	Port string
}


var server_instance *HttpServer
var server_once sync.Once
 
func GetHttpServerInstance() *HttpServer {
    server_once.Do(func() {
        server_instance = &HttpServer{}
    })
    return server_instance
}


func (s *HttpServer) StartServer(port string) {

	p := ":" + port
	fmt.Println("======== port: ",p)
	http.ListenAndServe(":3001", nil)
}

func (s *HttpServer)AddHandleFunc(subUrl string,f func( http.ResponseWriter,  *http.Request)) {
	http.HandleFunc(subUrl,f)
}


