package util

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestGetFileType(t *testing.T) {
	//f, err := os.Open("C:\\Users\\Administrator\\Desktop\\api.html")
	f, err := os.Open("C:\\Users\\WenTe\\Desktop\\video\\664bc4e86cfae46338056e7ec016555e.mp4")
	if err != nil {
		t.Logf("open error: %v", err)
	}
	fSrc, err := ioutil.ReadAll(f)
	fileType := GetFileType(fSrc[:10])
	log.Println(fileType)
}
