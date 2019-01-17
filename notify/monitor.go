package notify;
 
import (
    "github.com/fsnotify/fsnotify"
    "github.com/smtp-http/filemonitor_macmini/conn"
    "github.com/smtp-http/filemonitor_macmini/config"
    "log"
    "fmt"
    "path/filepath"
    "path"
    "io/ioutil"
    "sync"
    "os"
)

type FileMonitor struct {
    m_tcpserver *conn.TcpServer
}

var instance *FileMonitor
var once sync.Once
 
func GetFileMonitorInstance() *FileMonitor {
    once.Do(func() {
        instance = &FileMonitor{}

    })
    return instance
}

func (f *FileMonitor) SetTcpserver(server *conn.TcpServer) {
    f.m_tcpserver = server
}
 
func (f *FileMonitor) StartMonitor(cancel chan string,path_monitor string) {

    var stop bool = false

    //创建一个监控对象
    watch, err := fsnotify.NewWatcher();
    if err != nil {
        log.Fatal(err);
    }
    defer watch.Close();
    //
    err = watch.Add(path_monitor);
    if err != nil {
        log.Fatal(err);
    }
    //处理监控对象的事件

    for {
        if stop {
            break
        }
        select {
            case ev := <-watch.Events:
                {
                    //判断事件发生的类型，如下5种
                    // Create 创建
                    // Write 写入
                    // Remove 删除
                    // Rename 重命名
                    // Chmod 修改权限
                    if ev.Op&fsnotify.Create == fsnotify.Create {
                        fmt.Println("创建文件 : ", ev.Name);
                    }
                    if ev.Op&fsnotify.Write == fsnotify.Write {
                        fmt.Println("写入文件 : ", ev.Name);
                        paths, fileName := filepath.Split(ev.Name) 
                        fmt.Println(paths, fileName) //获取路径中的目录及文件名 E:\data\ test.txt 
                        fmt.Println(filepath.Base(ev.Name)) //获取路径中的文件名test.txt 
                        

                        if path.Ext(ev.Name) == config.GetConfig().FileExtension {
                            fmt.Println(path.Ext(ev.Name)) //获取路径中的文件的后缀 .txt
                            b, err := ioutil.ReadFile(ev.Name) 
                            if err != nil { 
                                fmt.Print(err) 
                            } 
                            fmt.Println(b) 
                            //str := string(b) 
                        
                            f.m_tcpserver.Notify(b)
                            configure := config.GetConfig()
                            if configure.HttpEnable == true {
                                httpClient := conn.GetHttpClientInstance()
                                httpClient.HttpPost(b)
                            }

                        }

                        
                    }

                    if ev.Op&fsnotify.Remove == fsnotify.Remove {
                        fmt.Println("删除文件 : ", ev.Name);
                    }

                    if ev.Op&fsnotify.Rename == fsnotify.Rename {
                        fmt.Println("重命名文件 : ", ev.Name);
                    }
                    
                    if ev.Op&fsnotify.Chmod == fsnotify.Chmod {
                        fmt.Println("修改权限 : ", ev.Name);
                    }
                }

            case can := <- cancel:
                if can == "cancel" {
                    fmt.Println("break!")
                    stop = true
                    break
                }
            case err := <-watch.Errors:
                {
                    log.Println("error : ", err);
                    return;
                }
        }
    }

    fmt.Println("good bye!")
}


////////////////////////////////////////////// directory monitor ///////////////////////////////////////////////

type DirectoryMonitor struct {
    MonitorName string
}

func (d *DirectoryMonitor) StopCurrentMonit(cancel chan string) {

}

var md_instance *DirectoryMonitor
var md_once sync.Once

func GetMiddleDirectoryMonitorInstance() *DirectoryMonitor {
    md_once.Do(func() {
        md_instance = &DirectoryMonitor{}
        md_instance.MonitorName = "middle"
    })
    return md_instance
}

var root_instance *DirectoryMonitor
var root_once sync.Once

func GetRootDirectoryMonitorInstance() *DirectoryMonitor {
    root_once.Do(func() {
        root_instance = &DirectoryMonitor{}
        root_instance.MonitorName = "root"
    })
    return root_instance
}

func (d *DirectoryMonitor) StartMonitor(cancel chan string,path_monitor string) {
    var stop bool = false

    //创建一个监控对象
    watch, err := fsnotify.NewWatcher();
    if err != nil {
        log.Fatal(err);
    }
    defer watch.Close();
    //
    err = watch.Add(path_monitor);
    if err != nil {
        log.Fatal(err);
    }

    //处理监控对象的事件

    for {
        if stop {
            break
        }
        select {
            case ev := <-watch.Events:
                {
                    //判断事件发生的类型，如下5种,判断文件夹只看创建。
                    // Create 创建
  
                    if ev.Op&fsnotify.Create == fsnotify.Create {
                        fmt.Println("创建文件 : ", ev.Name);
                        f, _ := os.Stat(ev.Name)
                        if f.IsDir() {
                            fmt.Println(ev.Name + " is dir.")
                            dm := DispMsg{}
                            dm.MonitorName = d.MonitorName
                            dm.NextPath = ev.Name
                            DispMsgCh <- dm
                        } 
                    }
                }

            case can := <- cancel:
                if can == "cancel" {
                    fmt.Println("break!")
                    stop = true
                    break
                }
            case err := <-watch.Errors:
                {
                    log.Println("error : ", err);
                    return;
                }
        }
    }
}