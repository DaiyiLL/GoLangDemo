package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Stalk struct {
	Title     string // title
	ConUrl    string // 内容的url
	ConDesc   string // 段子内容
	LikeCnt   int64  // 喜欢的点击次数
	UnLikeCnt int64  // 不喜欢的点击次数
	CommentCnt int64  // 评论次数
	CommentUrl string // 评论详情

	AuthorName string  // 用户的name
	AuthorUrl  string  // 用户url
	AuthorImg  string  // 用户头像
}

func main01() {
	// https://m.pengfue.com/xiaohua_1.html

	//  段子title

	/*
	<div class="humorListTitle">
		<h1 class="f18"><a href="https://m.pengfue.com/content/1857807/" title="会放坏">会放坏</a></h1>
	</div>
	 */

	// 段子主要内容
	/*
	<div class="con-img">
		女：“我购物车里的那些水果牛奶饮料零食你赶紧给我买！” 男：“着什么急呀？” 女：“天越来越热了，会放坏的。”
	</div>
	 */

	// 段子
	/*
	<li class="ding"><i class="iconfont icon-xiao1"></i><em>0</em></li>
	<li class="cai"><i class="iconfont icon-ku1"></i><em>0</em></li>
	<li title="评论"><a href="https://m.pengfue.com/content/1857805/#commentWrap"><i class="iconfont icon-pinglun1"></i> <span>0</span></a></li>
	 */

	var start, end int
	fmt.Println("请输入爬取的起始页(>=1):")
	fmt.Scan(&start)
	fmt.Println("请输入爬取的终止也(>=start):")
	fmt.Scan(&end)

	toWork(start, end)
}

func toWork(start, end int)  {
	fmt.Printf("正在爬取%d到%d页....\n", start, end)
	//page := make(chan int)

	page := make(chan int)

	for i:=start; i<=end; i++ {
		go SpiderPageDB(i, page)
	}



	for i:=start; i<=end; i++ {
		index := <- page
		fmt.Println("第" + strconv.Itoa(index) + "页爬取完成")
	}
}

func SpiderPageDB(index int, page chan int) {
	// Url
	url := "https://m.pengfue.com/xiaohua_" + strconv.Itoa(index) + ".html"
	result, err := HttpGetDB(url)
	if err != nil {
		LogError("HttpGetDB", err)
		return
	}

	//fmt.Println(result)
	//ret := regexp.MustCompile(`<h1 class="f18"><a href="(?s:(.*?))"`)
	ret := regexp.MustCompile(`<div class="humorListTitle">(?s:(.*?))</div>`)
	humorListTitles := ret.FindAllStringSubmatch(result, -1)

	ret1 := regexp.MustCompile(`<div class="con-img">(?s:(.*?))</div>`)
	contentList := ret1.FindAllStringSubmatch(result, -1)

	ret2 := regexp.MustCompile(`<li class="ding"><i class="iconfont icon-xiao1"></i><em>(?s:(.*?))</em></li>`)
	likeCntList := ret2.FindAllStringSubmatch(result, -1)

	ret3 := regexp.MustCompile(`<li class="cai"><i class="iconfont icon-ku1"></i><em>(?s:(.*?))</em></li>`)
	unlikeCntList := ret3.FindAllStringSubmatch(result, -1)

	// <a href="https://m.pengfue.com/content/1857805/#commentWrap"><i class="iconfont icon-pinglun1"></i> <span>0</span></a>
	ret4 := regexp.MustCompile(`<li title="评论">(?s:(.*?))</li>`)
	commentList := ret4.FindAllStringSubmatch(result, -1)

	ret5 := regexp.MustCompile(`<div class="head-name">(?s:(.*?))</div>`)
	authorList := ret5.FindAllStringSubmatch(result, -1)

	// 以title开始
	lenght := len(humorListTitles)
	if len(likeCntList) < lenght && len(contentList) < lenght && len(commentList) < lenght && len(authorList) < lenght {
		return
	}

	// 分析小项
	urlRet := regexp.MustCompile(`<h1 class="f18"><a href="(?s:(.*?))"`)
	titleRet := regexp.MustCompile(`title="(?s:(.*?))">`)

	cmtUrlRet := regexp.MustCompile(`<a href="(?s:(.*?))">`)
	cmtCntRet := regexp.MustCompile(`<span>(?s:(.*?))</span>`)

	authorNameRet := regexp.MustCompile(`<a class="dp-b" href="(?s:(.*?))</a>`)
	authorImgRet  := regexp.MustCompile(`<img src="(?s:(.*?))"`)

	dataList := make([]Stalk, 0)
	for i, v := range humorListTitles {
		titleDesc := v[1]
		cmtDesc := commentList[i][1]
		authorDesc := authorList[i][1]

		//fmt.Println(titleDesc)
		titleResult := titleRet.FindAllStringSubmatch(titleDesc, -1)
		urlResult := urlRet.FindAllStringSubmatch(titleDesc, -1)
		likeCntResult := likeCntList[i][1]
		unlikeCntResult := unlikeCntList[i][1]
		conResult := strings.Trim(contentList[i][1], "")


		cmtUrlResult := cmtUrlRet.FindAllStringSubmatch(cmtDesc, -1)
		cmtCntResult := cmtCntRet.FindAllStringSubmatch(cmtDesc, -1)

		// 作者的昵称和url
		authorNameResult := authorNameRet.FindAllStringSubmatch(authorDesc, -1)[0][1]
		//fmt.Println(authorNameResult)
		authorInfoList := strings.Split(authorNameResult, `">`)
		authorName, authorUrl := "", ""
		if len(authorInfoList) >= 2 {
			authorName = authorInfoList[0]
			authorUrl  = authorInfoList[1]
		}
		authorImgResult := authorImgRet.FindAllStringSubmatch(authorDesc, -1)

		stalk := Stalk{
			DYToString(titleResult[0][1]),
			DYToString(urlResult[0][1]),
			DYToString(conResult),
			DYToInt64(likeCntResult),
			DYToInt64(unlikeCntResult),
			DYToInt64(cmtCntResult[0][1]),
			DYToString(cmtUrlResult[0][1]),
			authorName,
			authorUrl,
			DYToString(authorImgResult[0][1]),
		}
		dataList = append(dataList, stalk)
	}

	// 保存到文件中
	saveDataToFile(index, dataList)

	// 防止住go程提前结束
	page <- index
}

func saveDataToFile(idx int, dataList []Stalk)  {
	path := "Stalk_" + strconv.Itoa(idx) + ".txt"
	file, err := os.Create(path)
	if err != nil {
		LogError("os.Create", err)
		return
	}
	defer file.Close()

	// 遍历数组写入文件中
	for i, v := range dataList {
		file.WriteString(strconv.Itoa(i) + "行: " + fmt.Sprintf("%v\n", v))
	}
}

func HttpGetDB(url string) (result string, err error)  {
	result = ""
	resp, err := http.Get(url)
	if err != nil {
		LogError("http.Client.Do", err)
		return result, err
	}
	defer resp.Body.Close()

	buf := make([]byte, 4096)
	var build strings.Builder
	for {
		n, err := resp.Body.Read(buf)
		if n == 0 {
			break
		}
		if err != nil && err != io.EOF {
			return "",err
		}
		build.WriteString(string(buf[:n]))
	}
	result = build.String()
	return
}

