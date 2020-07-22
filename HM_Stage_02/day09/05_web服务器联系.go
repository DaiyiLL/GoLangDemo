package main

import (
	"fmt"
	"net/http"
	"os"
)

func OpenSendFile(fileName string, w http.ResponseWriter)  {
	pathFileName := "/Users/szdjy/Desktop/GO/Demo/HM_Stage_02/day09/resources" + fileName
	f, err := os.Open(pathFileName)
	if err != nil {
		fmt.Println("Open err: ", err)
		w.Write([]byte("No such file or directory"))
		return
	}
	defer f.Close()

	buf := make([]byte, 4096)
	for {
		n, _ := f.Read(buf)   // 从本地将文件内容读取
		if n == 0 {
			return
		}
		w.Write(buf[:n])
	}

}

func handler05(w http.ResponseWriter, r *http.Request)  {
	fmt.Println("客户端请求:", r.URL)

	OpenSendFile(r.URL.String(), w)

}

func main()  {
	http.HandleFunc("/", handler05)
	http.ListenAndServe("127.0.0.1:8000", nil)
}
