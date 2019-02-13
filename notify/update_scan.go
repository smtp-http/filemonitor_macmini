package notify

import (
    "sync"
    //"os"
    "time"
    "fmt"
    "encoding/csv"
    "io/ioutil"
    "io"
    "strings"
    "path/filepath"
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


func (u *UpdateScanner) ScanFile(root string) {
	ticker := time.NewTicker(time.Second * 2)
	
	var err error
	u.UpdateTime,err = GetLastProcessTime()
	if err != nil {
		fmt.Printf("GetLastProcessTime error : %v\n",err)
		return
	}

    for _ = range ticker.C {
    	/*
    	tm := time.Now()
        fmt.Printf("ticked at %v\n", tm)
        err := UpdateLastProcessTime(tm.Unix())
		if err != nil {
			fmt.Printf("GetLastProcessTime error : %v\n",err)
			return
		}*/

		files,err := GetAllFiles(root)
		if err != nil {
			fmt.Printf("GetAllFiles error : %v\n",err)
			return
		}

		for i, file := range files {
			err,tm := getFileModTime(file)
			if err != nil {
				fmt.Printf("Get file create time err %d : %v\n",i,err)
			} else {
				fileName := filepath.Base(file)
				if fileName == "Records.csv" {
					
				} else if fileName == "Attributes.csv" {
					
				} else {
					continue
				}

				fmt.Printf("%d  ---   %v\n",i,tm.Unix())
				cntb,err := ioutil.ReadFile(file)
				if err != nil {
					fmt.Printf("Read file err %v\n",err)
					continue
				}
/*
				r2 := csv.NewReader(strings.NewReader(string(cntb)))
				ss,_ := r2.ReadAll()
				fmt.Println(ss)
				//fmt.Println(ss)
				sz := len(ss)
				for i:=0;i<sz;i++{
					fmt.Println(ss[i])
				}*/
				records,e := ReadCSV(strings.NewReader(string(cntb)))
				if e != nil {
					fmt.Printf("++ Read csv error %v \n",e)
					continue
				}

				for i,record := range records {
					fmt.Printf("--- r[%d] = %v\n",i,record)
				}
			}
			
		}

		return
    }
}


type Result struct{
	Ip     string
	Sn     string
	Status string
}

// ReadCSV 展示了如何处理CSV
// 接收的参数通过io.Reader传入
func ReadCSV(b io.Reader) ([][]string, error) {

	//返回的是csv.Reader
	r := csv.NewReader(b)

	// 分隔符和注释是csv.Reader结构体中的字段
	r.Comma = ';'
	r.Comment = '-'

	var records [][]string

	// 读取并返回一个字符串切片和错误信息
	// 我们也可以将其用于字典键或其他形式的查找
	// 此处忽略了返回的切片 目的是跳过csv首行标题
	_, err := r.Read()
	if err != nil && err != io.EOF {
		return nil, err
	}

	// 循环直到全部处理完毕
	for {
		record, err := r.Read() 
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		//fmt.Printf("record: %v\n",record)
		//_, err := strconv.ParseInt(record[2], 10, 64)
		//if err != nil {
		//	return nil, err
		//}

		//m := Records{/*record[0], record[1], int(year)*/}
		records = append(records, record)
	}
	return records, nil
}
