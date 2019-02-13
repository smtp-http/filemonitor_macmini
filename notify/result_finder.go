package notify



import (
	"fmt"
	"io"
	"os"
	"time"
	"sync"
	"path/filepath"
	"github.com/smtp-http/filemonitor_macmini/config"
)


type Finder struct {

}


var finder_instance *Finder
var finder_once sync.Once
 
func GetFinderInstance() *Finder {
    finder_once.Do(func() {
        finder_instance = &Finder{}
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
		n, err := fd.Read(buf)
		if err != nil && err != io.EOF {
			fmt.Printf("read file %s failed: %v\n", file, err)
			return
		}
		if n > 1 {
			//if n != len(buf) {
			//	n--
			//}
			fd.Close()
			time.Sleep(1 * time.Second)
			fmt.Printf("-- n: %v   len(buf):%v\n",n,len(buf))
			fmt.Printf("%s", string(buf[:n]))
			offset += int64(n)
			// TODO: cvs file process
			UpdateTestSummarySeek(offset)
			fd, err = os.Open(file)
			if err != nil {
				fmt.Printf("open file %s failed: %v\n", file, err)
				return
			}
			fd.Seek(offset, 0)
		}// else 

		//if n == len(buf) {
			
		//	fd.Close()
		//	fd, err = os.Open(file)
		//	if err != nil {
		//		fmt.Printf("open file %s failed: %v\n", file, err)
		//		return
		//	}
		//	fmt.Printf("++ offset: %v\n",offset)
		//	fd.Seek(offset, 0)
		//	UpdateTestSummarySeek(offset)
		//}
	}

}