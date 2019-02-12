package notify

import (
    "sync"
    //"os"
    "time"
    "fmt"
)

type UpdateScanner struct {
	UpdateTime     int64
}

var scanner_instance *UpdateScanner
var scanner_once sync.Once
 
func GetUpdateScannerInstance() *UpdateScanner {
    scanner_once.Do(func() {
        scanner_instance = &UpdateScanner{}
        scanner_instance.UpdateTime = 0
    })
    return scanner_instance
}


func (u *UpdateScanner) Scan(root string) {
	ticker := time.NewTicker(time.Second * 5)
	
	var err error
	u.UpdateTime,err = GetLastProcessTime()
	if err != nil {
		fmt.Printf("GetLastProcessTime error : %v\n",err)
		return
	}

    for _ = range ticker.C {
        fmt.Printf("ticked at %v\n", time.Now())
        
    }
}