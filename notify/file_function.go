package notify

import (
	//"os"
	"io/ioutil"
	"fmt"
)

func FindTargetFolder(cur_dir string) {
	
	//f, _ := os.Stat("next folder")

	files, _ := ioutil.ReadDir(cur_dir)
    for _, f := range files {
        fmt.Println(f.Name())
        if f.IsDir() {
        	fmt.Println(f.Name() + " is dir!")
    	}
    }
    
}


func FindTargetFolderAfterReset() {

}