package nimble

import (
	"fmt"
	"os"
	"path"
)

func OpenRecordFile(filename string) (*os.File, error) {
	return os.Open(recordDir + filename)
}

func CreateRecordFile(filename string) (file *os.File, err error) {
	filepath := recordDir + filename
	dirpath := path.Dir(filepath)
	err = os.MkdirAll(dirpath, os.ModeDir|os.ModePerm)
	if err != nil {
		fmt.Printf("MkdirAll %s: err=%v", dirpath, err)
		return
	}
	return os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0666)
}
