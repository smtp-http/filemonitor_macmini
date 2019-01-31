package notify

import (
	"os"
	"io/ioutil"
	"fmt"
	"time"
	"path/filepath"
	"errors"
	"github.com/smtp-http/filemonitor_macmini/config"
	"strconv"
)


func FindTargetFolderWhenReset(root string) {
	var next_dir string = root //config.GetConfig().RootDirectory
	fmt.Printf("next_dir:%v\n",next_dir)
	var try_times int = 0
	for {
		err,target_path := findNewestSubFolder(next_dir)
		if err != nil {
			if err.Error() == "No files or folder!" {
				fmt.Println(err.Error())
				
				try_times ++
				if try_times > config.GetConfig().TryTimes {
					fmt.Printf("Find target folder when reset timeout!")
					return
				}
				time.Sleep(10 * time.Millisecond)
				continue
			}

			fmt.Printf("find newest sub folder err when reset: %v \n",err)
			return
		}

		folder := filepath.Base(target_path)

		if folder == config.GetConfig().DestinationFolder {
			GetFileMonitorInstance().IsRunning = true
            GetFileMonitorInstance().StartMonitor(FileCancelSignal,target_path)
			fmt.Printf("target path: %v\n",target_path)
			break

		}

		next_dir = filepath.Join(next_dir,folder)

	}
}


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


		if folder == config.GetConfig().DestinationFolder {
			dm := DispMsg{}
			dm.MonitorName = "file"
			dm.Action = "create_file_monitor"
			dm.NextPath = target_path
			DispMsgCh <- dm
			fmt.Printf("target path: %v\n",target_path)
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
        
        if f.IsDir() {
        	err,creat_time := getFileModTime(filepath.Join(cur_dir,f.Name())) 
        	if err != nil {
        		fmt.Printf("get file: %v creatime err!\n",f.Name())
        		return err,""
        	}
        	
        	if creat_time.Unix() > creat_time_unix {
        		creat_time_unix = creat_time.Unix()
        		dir_name = filepath.Join(cur_dir,f.Name())
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


func Exists(path string) bool { 
		_, err := os.Stat(path) //os.Stat获取文件信息 
		if err != nil { 
			if os.IsExist(err) { 
				return true 
			} 
			return false 
		} 
		return true 
}


func GetLastProcessTime() (int64,error){
	file_name := config.GetConfig().LastTimeStamp
	if !Exists(file_name) {
		f,err := os.Create(file_name)
		defer f.Close()
		if err !=nil {
		    fmt.Println(err.Error())
		    return 0,err
		} 

		time_stamp := time.Now().Unix()
		time_str := strconv.FormatInt(time_stamp,10)
		_,err = f.Write([]byte(time_str))
		if err != nil {
			fmt.Println(err)
			return 0,err
		}
	}


	f, err := os.OpenFile(file_name, os.O_RDONLY,0600)
    defer f.Close()
    if err !=nil {
        fmt.Println(err.Error())
        return 0,err
    } else {
		contentByte,err := ioutil.ReadAll(f)
		if err != nil {
			fmt.Println(err)
			return 0,err
		}

		fmt.Println(string(contentByte))
		time_rec, er := strconv.ParseInt(string(contentByte), 10, 64)
		if er != nil {
			fmt.Println(er)
			return 0,er
		}
		return time_rec,nil
    }
}



func UpdateLastProcessTime(time_stamp int64) error {

	file_name := config.GetConfig().LastTimeStamp

	if Exists(file_name) {
		err := os.Remove(file_name)               //
		if err != nil {
			fmt.Println("--file remove Error!")
			fmt.Printf("%s", err)
			return err
		}
	}

	f, err := os.OpenFile(file_name, os.O_WRONLY|os.O_CREATE, 0644)
    if err != nil {
        return err
    }

    defer f.Close()
    
    time_str := strconv.FormatInt(time_stamp,10)
	_,err = f.Write([]byte(time_str))
	if err != nil {
		fmt.Println(err)
		return err
	}
    
	return nil
}