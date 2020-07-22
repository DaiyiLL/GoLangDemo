package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

func main() {
	// https://tieba.baidu.com/f?kw=绝对求生&ie=utf-8&pn=50

	var start, end int
	fmt.Println("请输入爬取的起始页(>=1):")
	fmt.Scan(&start)
	fmt.Println("请输入爬取的终止也(>=start):")
	fmt.Scan(&end)

	working(start, end)
}

// 爬取页面操作
func working(start int, end int)  {
	fmt.Printf("正在爬取第%d页到%d页...\n", start, end)

	// 循环爬取每一页的数据
	for i:=start; i<=end; i++ {
		url := "https://tieba.baidu.com/f?kw=绝地求生&ie=utf-8&pn=" + strconv.Itoa((i - 1) * 50)
		result, err := HttpGet(url)
		if err != nil {
			fmt.Println("HttpGet err", err)
			continue
		}
		//fmt.Println("result = \n", result)
		f, err := os.Create("第" + strconv.Itoa(i) + "页.html")
		if err != nil {
			fmt.Println("os.Create err:", err)
			continue
		}

		f.WriteString(result)
		f.Close()
	}
}

func HttpGet(url string) (result string, err error)  {
	resp, err1 := http.Get(url)
	if err1 != nil {
		err = err1
		result = ""
		return
	}
	defer resp.Body.Close()

	// 循环读取网页数据，传出给调用者
	buf := make([]byte, 4096)
	for {
		n, err2 := resp.Body.Read(buf)
		if n == 0 {
			fmt.Println("读取网页完成")
			break
		}
		if err2 != nil && err2 != io.EOF {
			err = err2
			break
		}
		// 累加每一次循环读到的 buf数据，存入到result一次性返回
		result += string(buf[:n])
	}
	return
}
