package main

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"
)

var downChan chan int
type DouyuGirlDao struct {
	GirlId    int64  `json:"rid"`
	NickName  string `json:"nn"`
	AvatarUrl string `json:"rs16"`
	Title     string `json:"rn"`

	FileName  string
}

type DataList struct {
	RecordCount int64
	RecordList  []*DouyuGirlDao `json:"rl"`
}

type BaseData struct {
	Code     int64    `json:"code"`
	Msg      string   `json:"msg"`
	DataList DataList `json:"data"`
}

const (
	IMG_PREFIX_URL = "https://rpic.douyucdn.cn/live-cover/roomCover"
	RESOURCES_DICTORY_URI = "Resources/"
)

func main()  {
	//url := "https://rpic.douyucdn.cn/live-cover/roomCover/2018/10/28/0a539aaf2e6a086155cad8af9a40b80e_big.png/dy1"
	downChan = make(chan int, 10)

	var start, end int
	fmt.Println("请输入爬取的起始页(>=1):")
	fmt.Scan(&start)
	fmt.Println("请输入爬取的终止页(>=start):")
	fmt.Scan(&end)

	beginSpider(start, end)
}

func beginSpider(start, end int)  {
	fmt.Printf("正在爬取%d到%d页....\n", start, end)

	page := make(chan int)

	for i:=start; i<=end; i++ {
		go SpiderDouyuDB(i, page)
	}

	for i:=start; i<=end; i++ {
		index := <- page
		fmt.Println("第" + strconv.Itoa(index) + "页爬取完成")
	}
}

func SpiderDouyuDB(idx int, page chan int)  {
	url := "https://www.douyu.com/gapi/rknc/directory/yzRec/" + strconv.Itoa(idx)
	result, err := DouyuGet(url)
	if err != nil {
		LogError("HttpGetDB", err)
		page <- idx
		return
	}

	//fmt.Println(string(result))

	data := BaseData{}
	err = jsoniter.Unmarshal(result, &data)
	if err != nil {
		LogError("jsoniter.UnmarshalFromString", err)
		page <- idx
		return
	}
	//fmt.Println("msg:", data)

	fileNameRet := regexp.MustCompile(`([^/]*?).(jpg|png|jpeg|PNG|JPG|JPEG)`)



	for i, girl := range data.DataList.RecordList {
		result := fileNameRet.FindAllString(girl.AvatarUrl, -1)
		if len(result) > 0 {
			girl.FileName = result[0]
		}
		fmt.Println(girl.FileName)
		saveDouyuGirl(idx, i, girl, downChan)
	}

	page <- idx

	//itemRet := regexp.MustCompile(`<li class="layout-Cover-item">(?s:(.*?))</li>`)
	//itemList := itemRet.FindAllStringSubmatch(result, -1)
	//
	//avatarRet := regexp.MustCompile(`src="(?s:(.*?))" class="DyImg-content is-normal" alt=`)
	//titleRet  := regexp.MustCompile(`title="(?s:(.*?))"`)
	//nameRect  := regexp.MustCompile(`<h2 class="DyListCover-user is-template"><svg><use xlink:href="#icon-user_c95acf8"></use></svg>(?s:(.*?))</h2>`)

	//girls := make([]DouyuGirl, 0)
	//for _, v := range itemList {
	//	itemDesc := v[1]
	//
	//	//avatarResult := avatarRet.FindAllStringSubmatch(itemDesc, -1)
	//	//fmt.Println(itemDesc)
	//	avatarUrl := DYToString(avatarRet.FindAllStringSubmatch(itemDesc, -1)[0][1])
	//	//fmt.Println(avatarUrl)
	//	title     := DYToString(titleRet.FindAllStringSubmatch(itemDesc, -1)[0][1])
	//	//fmt.Println(title)
	//	name      := DYToString(nameRect.FindAllStringSubmatch(itemDesc, -1)[0][1])
	//
	//	girls = append(girls, DouyuGirl{name, avatarUrl, title})
	//}
	//
	//for _, girl := range girls {
	//	fmt.Println(girl.AvatarUrl)
	//}
	//page <- idx
	// 下载图片
}

func DouyuGet(url string) ([]byte, error)  {
	resp, err := http.Get(url)
	if err != nil {
		LogError("http.Client.Do", err)
		return []byte(""), err
	}
	defer resp.Body.Close()

	// 循环爬取整夜的数据
	buf := make([]byte, 4096)
	builder := make([]byte, 0)
	//var build strings.Builder
	for {
		n, err := resp.Body.Read(buf)
		if n == 0 {
			break
		}
		if err != nil && err != io.EOF {
			return []byte(""),err
		}
		//build.WriteString(string(buf[:n]))
		builder = append(builder, buf[:n]...)
	}
	return builder, nil
}

func saveDouyuGirl(idx int, index int, dao *DouyuGirlDao, task chan int) (bool, error)  {
	fileDir := RESOURCES_DICTORY_URI + time.Now().Format("2006/01/02/") + strconv.Itoa(idx) + "/"
	fileName := filepath.Join(fileDir, dao.FileName)
	// 常见文件
	err := os.MkdirAll(fileDir, os.ModePerm)
	if err != nil {
		LogError( dao.FileName + " os.MkdirAll:", err)
		return false, err
	}

	file, err := os.Create(fileName)
	if err != nil {
		LogError( dao.FileName + " os.Create:", err)
		return false, err
	}
	defer file.Close()
	resp, err := http.Get(dao.AvatarUrl)
	if err != nil {
		LogError( fileName + " http.Get:", err)
		return false, err
	}
	defer resp.Body.Close()

	buf := make([]byte, 4096)
	for {
		n, err := resp.Body.Read(buf)
		if n == 0 {
			return false, err
		}
		if err != nil && err != io.EOF {
			LogError(dao.FileName + " resp.Body.Read:", err)
			return false, err
		}
		file.Write(buf[:n])
	}
	//task <- index
	return true, nil
}

/*
{
	"rid": 3544169,
	"rn": "早起的大牛想吃肉",
	"uid": 174025350,
	"nn": "程心曲",
	"cid1": 2,
	"cid2": 311,
	"cid3": 1745,
	"iv": 0,
	"av": "avatar_v3/202007/06cda64a915f476eb615c0bdc02633d6",
	"ol": 2088290,
	"url": "/3544169",
	"c2url": "/directory/game/XX",
	"c2name": "颜值",
	"icdata": {
		"549": {
			"url": "",
			"w": 0,
			"h": 0
		}
	},
	"dot": 2103,
	"subrt": 0,
	"topid": 0,
	"oaid": 0,
	"bid": 0,
	"gldid": 0,
	"rs1": "https://rpic.douyucdn.cn/live-cover/roomCover/2020/04/13/60a4ec3d32497b2d2886c995320c7382_small.jpg/dy2",
	"rs16": "https://rpic.douyucdn.cn/live-cover/roomCover/2020/04/13/60a4ec3d32497b2d2886c995320c7382_big.jpg/dy1",
	"utag": [],
	"rpos": 0,
	"rgrpt": 1,
	"rkic": "",
	"rt": 2103,
	"ot": 0,
	"clis": 2,
	"chanid": 0,
	"icv1": [
		[
			{
				"id": 549,
				"url": "https://sta-op.douyucdn.cn/dy-listicon/94e105559334bf35bd197084dc7deece.png",
				"score": 998,
				"w": 0,
				"h": 0
			}
		],
		[],
		[],
		[]
	],
	"ioa": 1,
	"od": "S4粉丝节全站十大主播"
}


 */