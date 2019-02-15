package notify

import (
	"fmt"
	"io"
	"os"
	"time"
	"sync"
	"path/filepath"
	"github.com/smtp-http/filemonitor_macmini/config"
	"github.com/smtp-http/filemonitor_macmini/conn"
	//"encoding/csv"
	"strings"
	"io/ioutil"
)


type Finder struct {
	client *conn.TcpClient
}


var finder_instance *Finder
var finder_once sync.Once
 
func GetFinderInstance() *Finder {
    finder_once.Do(func() {
        finder_instance = &Finder{}
        finder_instance.client = new(conn.TcpClient)
		finder_instance.client.Init(config.GetConfig().Host)
    })
    return finder_instance
}

func(f *Finder)	Monitor() {
	file := filepath.Join(config.GetConfig().RootDirectory,"TestSummary.csv")
	
	buf := make([]byte, 2048)
	offset,e := GetLastTestSummarySeek()
	if e != nil {
		fmt.Printf("Get last test summary seek failed: %v\n",e)
		return
	}

	fd, err := os.Open(file)
	defer fd.Close()
	if err != nil {
		fmt.Printf("open file %s failed: %v\n", file, err)
		return
	}

	fd.Seek(offset,0)

	for {
		n, err := fd.Read(buf[1:])
		if err != nil && err != io.EOF {
			fmt.Printf("read file %s failed: %v\n", file, err)
			return
		}
		if n > 1 {
			fd.Close()
			time.Sleep(100 * time.Millisecond)//(1 * time.Second)
			fmt.Printf("-- n: %v   len(buf):%v\n",n,len(buf))
			str := "\n"
			bstr := []byte(str)
			copy(buf[n+1:n+len(bstr)+1],bstr)
		
			offset += int64(n)

			
			records,e := ReadCSV(strings.NewReader(string(buf[:n + len(bstr)])))
			if e != nil {
				fmt.Printf("++ Read csv error %v \n",e)
				UpdateTestSummarySeek(offset)
				fd.Seek(offset, 0)
				continue
			}

			for i,record := range records {
				// TODO: get result

				fmt.Printf("--- r[%d] = %v\n",i,record)
				fmt.Println("record[0]: ",record[0])
				f.processRecord(record[0])
			}

			fd, err = os.Open(file)
			if err != nil {
				fmt.Printf("open file %s failed: %v\n", file, err)
				return
			}

			UpdateTestSummarySeek(offset)
			fd.Seek(offset, 0)
		}
	}

}

func (f *Finder)processRecord(record string) {
	strs := strings.Split(record,",")
	for i,v := range strs {
		fmt.Printf("[%d]: %s\n",i,v)
	}

	sn := config.GetConfig().SerialNumber
	if len(strs) == 2 {
		ret := sn + "," + "2"
		f.client.Send([]byte(ret))
	} else if len(strs) == 3 {
		ret := sn + "," + "3"
		f.client.Send([]byte(ret))
	} else if len(strs) > 10 {
		ret := sn + "," + strs[0] + "," + get_test_result(strs)
		f.client.Send([]byte(ret))
	} else {
		fmt.Printf("The unknown message!")
	}
}

func get_test_result(record []string) string {
	targetPath := filepath.Join(config.GetConfig().RootDirectory,"Archive")
	targetPath = filepath.Join(targetPath,record[0])

	if !Exists(targetPath) {
		return "NO_RESULT"
	}

	err,subFolder := findNewestSubFolder(targetPath)
	if err != nil {
		fmt.Printf("find newest test folder error: %v\n",err)
		return "NO_NEWEST_RESULT"
	}

	atlaslogs := filepath.Join(subFolder,"AtlasLogs")
	records_csv := filepath.Join(atlaslogs,"Records.csv")

	if !Exists(records_csv) {
		return "NO_RECORDS_CSV_FILE"
	} 

	return read_csv_file(records_csv)
}


func read_csv_file(file string) string{
	cntb,err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Printf("Read file err %v\n",err)
		return "READ_FILE_FAIL"
	}

	records,e := ReadCSV(strings.NewReader(string(cntb)))
	if e != nil {
		fmt.Printf("++ Read csv error %v \n",e)
		return "READ_CSV_FAIL"
	}

	for _,record := range records {
		if processOneLine(record[0]) == "FAIL" {
			return "FAIL"
		}
	}

	return "PASS"
}

func processOneLine(record string) string{
	strs := strings.Split(record,",")
	return strs[4]
}
