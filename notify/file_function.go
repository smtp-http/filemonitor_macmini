package notify

import (
	"os"
	"io/ioutil"
	"fmt"
	"time"
	"path/filepath"
	"errors"
	"github.com/smtp-http/filemonitor_macmini/config"
)


func FindTargetFolder(cur_dir string) {
	var nex_dir string = cur_dir
	var try_times int = 0
	for {
		err,target_path := findNewestSubFolder(nex_dir)
		if err != nil {
			if err.Error() == "No files or folder!" {
				fmt.Println(err.Error())
				
				try_times ++
				if try_times > config.GetConfig().TryTimes {
					fmt.Printf("Find target folder timeout!")
					dm := DispMsg{}
					dm.MonitorName = "file"
					dm.Action = "timeout"
					dm.NextPath = target_path
					DispMsgCh <- dm
					return
				}
				time.Sleep(10 * time.Millisecond)
				continue
			}

			dm := DispMsg{}
			dm.MonitorName = "file"
			dm.Action = "find_error"
			dm.NextPath = target_path
			DispMsgCh <- dm
			fmt.Printf("find newest sub folder err: %v \n",err)
			return
		}

		
		folder := filepath.Base(target_path)

		fmt.Printf("folder: %v    des folder: %v\n",folder,config.GetConfig().DestinationFolder)

		if folder == config.GetConfig().DestinationFolder {
			dm := DispMsg{}
			dm.MonitorName = "file"
			dm.Action = "create_file_monitor"
			dm.NextPath = target_path
			DispMsgCh <- dm
			fmt.Printf("newest sub folder: %v\n",target_path)
			break

		}

		nex_dir = filepath.Join(nex_dir,folder)

	}
	
}



func findNewestSubFolder(cur_dir string) (error,string){
	//f, _ := os.Stat("next folder")
	var creat_time_unix int64 = 0
	var dir_name string = ""
	files, _ := ioutil.ReadDir(cur_dir)
	if len(files) <= 0 {
		return errors.New("No files or folder!"),""
	}

    for _, f := range files {
        fmt.Println(f.Name())
        if f.IsDir() {
        	err,creat_time := getFileModTime(filepath.Join(cur_dir,f.Name())) 
        	if err != nil {
        		fmt.Printf("get file: %v creatime err!\n",f.Name())
        		return err,""
        	}
        	fmt.Printf("unix time: %v\n",creat_time.Unix())
        	if creat_time.Unix() > creat_time_unix {
        		creat_time_unix = creat_time.Unix()
        		fmt.Printf("time: %v \n",creat_time_unix)
        		dir_name = f.Name()
        	}
    	}
    }
    
    return nil,dir_name
}


func getFileModTime(path string) (error,time.Time) { 
	f, err := os.Open(path) 
	if err != nil { 
		fmt.Println("open file error") 
		return err,time.Now()
	} 
	defer f.Close() 
	fi, err := f.Stat() 
	if err != nil { 
		fmt.Println("stat fileinfo error") 
		return err,time.Now()
	} 
	return nil,fi.ModTime()
}
