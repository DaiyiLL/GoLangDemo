package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

//https://movie.douban.com/top250?start=0&filter=  //1
//https://movie.douban.com/top250?start=25&filter= //2


// 电影名称   <img width="100" alt="肖申克的救赎" `<img width="100" alt="(?s:(.*?))"`
// 评分      <span class="rating_num" property="v:average">9.7</span>  `<span class="rating_num" property="v:average">(?s:(.*?))</span>`
// 评分人数   <span>2046145人评价</span> `<span>(?s:(.*?))人评价</span> `

func main()  {
	var start, end int
	fmt.Println("请输入爬取的起始页(>=1):")
	fmt.Scan(&start)
	fmt.Println("请输入爬取的终止也(>=start):")
	fmt.Scan(&end)

	toWork(start, end)
}

func toWork(start int, end int)  {
	fmt.Printf("正在爬取%d到%d页...\n", start, end)
	page := make(chan int)

	for i:=start; i<=end; i++ {
		go SpiderPageDB(i, page)
	}

	for i:= start; i<=end; i++ {
		index := <- page
		fmt.Println("第" + strconv.Itoa(index) + "页爬取完成")
	}
}

// 爬取一个豆瓣页面数据信息
func SpiderPageDB(index int, page chan  int)  {
	// 获取url地址
	url := "https://movie.douban.com/top250?filter=&start=" + strconv.Itoa((index - 1) * 25)
	// 爬取url对应页面的数据
	result, err := HttpGetDB(url)
	if err != nil {
		fmt.Println("HttpGetDB err:", err)
		return
	}

	//fmt.Println("result = ", result)
	// 解析、编译正则表达式
	ret := regexp.MustCompile(`<img width="100" alt="(?s:(.*?))"`)
	// 提取需要的信息
	filmNames := ret.FindAllStringSubmatch(result, -1)

	ret1 := regexp.MustCompile(`<span class="rating_num" property="v:average">(?s:(.*?))</span>`)
	// 提取需要的信息
	filmScores := ret1.FindAllStringSubmatch(result, -1)

	ret2 := regexp.MustCompile(`<span>(?s:(\d*?))人评价</span>`)
	// 提取需要的信息
	filmCritics := ret2.FindAllStringSubmatch(result, -1)

	//fmt.Println(filmNames)
	//fmt.Println(filmScores)
	//fmt.Println(filmCritics)

	Save2File(index, filmNames, filmScores, filmCritics)

	page <- index
}

func Save2File(idx int, filmNames, filmScores, filmCritics [][]string)  {
	file, err := os.Create("douban/豆瓣电影第" + strconv.Itoa(idx) + "页.txt")
	if err != nil {
		fmt.Println("os.Create err: ", err)
		return
	}
	defer file.Close()

	n := len(filmNames)
	file.WriteString("电影名称" + "\t\t\t\t" + "评分" + "\t\t" + "评分人数" + "\n")
	//fmt.Println(filmCritics)
	for i:=0; i<n; i++ {
		file.WriteString(filmNames[i][1] + "\t\t\t\t" + filmScores[i][1] + "\t\t" + filmCritics[0][1] + "\n")
	}
}

// 爬取result
func HttpGetDB(url string) (result string, err error) {
	req, _ := http.NewRequest("GET", url, nil)
	// 比如说设置个token
	req.Header.Set("Referer", "https://book.douban.com/tag/%E5%B0%8F%E8%AF%B4")
	// 再设置个json
	req.Header.Set("User-Agent","Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/50.0.2661.94 Safari/537.36")


	resp, err1 := (&http.Client{}).Do(req)

	//resp, err1 := http.Get(url)
	if err1 != nil {
		err = err1
		return
	}
	defer resp.Body.Close()

	// 循环爬取整夜的数据
	buf := make([]byte, 4096)
	for {
		n, err2 := resp.Body.Read(buf)
		if n == 0 {
			break
		}
		if err2 != nil && err2 != io.EOF {
			err = err2
			return
		}
		result += string(buf[:n])
	}
	return
}
