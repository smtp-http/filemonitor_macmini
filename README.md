# filemonitor_macmini
File monitor running on macmini

## Usage

####  更多usage参看 [其他例子](https://github.com/smtp-http/filemonitor_macmini/wiki)

```
go get -u github.com/smtp-http/filemonitor_macmini
```


```go
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

	disp := new(notify.Dispatcher)

	disp.Dispatch()
}
```


#### 需要处理的目录结构：

![image](https://github.com/smtp-http/filemonitor_macmini/blob/master/images/path.png)
