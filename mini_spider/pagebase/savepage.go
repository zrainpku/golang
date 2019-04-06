package pagebase

import (
	"fmt"
	// "io"
	"io/ioutil"
	"pkuspider/mini_spider/mytools"
	"os"
	// "strconv"
	// "encoding/json"
	// "regexp"
	// "icode.baidu.com/baidu/go-lib/log"
	"strings"
	"time"
)

func SaveHtml(urlFile string) {
	t := time.NewTimer(time.Second * 1)
	<-t.C
	// log.Logger.Info("开始爬取网页:" + urlFile)
	// fUrl, _ := os.OpenFile("../data/"+urlFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	// file, _ := ioutil.ReadAll(fUrl)
	// file1, _ := ioutil.ReadFile("../data/url.data")
	// pageCont, _ := mytools.PageVisit(strings.Split(string(file), "\n")[0])
	// preUrl := []string(strings.Split(string(file), "\n"))
	// for _, tempHrefUrl := range preUrl
	if len(urlFile) < 5 {
		// log.Logger.Info("爬取网页:url非法" + urlFile)
	}
	tempUrl := strings.Split(urlFile, "\"")[1]
	// log.Logger.Info("开始爬取网页:url" + tempUrl)
	fmt.Println(tempUrl)
	pageCont, _ := mytools.PageVisit(tempUrl)
	pageByte := []byte(pageCont)
	filename := mytools.ChangeUrl(tempUrl)
	// var f *os.File
	// var err1 error
	os.Create("../output/" + string(filename))
	fmt.Println(filename)
	// if mytools.CheckFileIsExist("../output/" + string(filename)) { //如果文件存在
	// 	fmt.Println("文件存在")
	// 	// f, err1 = os.OpenFile("../output/zhurui.txt", os.O_APPEND, 0666) //打开文件
	// } else {
	// 	fmt.Println("文件不存在")
	// 	// f, err1 = os.Create("../output/" + string(index) + "1.txt") //创建文件
	// }
	// mytools.CheckErr(err1)
	ioutil.WriteFile("../output/"+string(filename), pageByte, 0644) //写入文件(字符串)

	// pageCont, _ := mytools.PageVisit(strings.Split(preUrl[1], "\"")[1])
	// fmt.Println(pageCont)

	/*
		if mytools.CheckRegexp(mytools.CheckRegexp(pageCont, regTitle, 0).(string), regCheckTitle, 0).(string) != "" {
			fmt.Print(mytools.CheckRegexp(mytools.CheckRegexp(pageCont, regTitle, 0).(string), regCheckTitle, 0).(string))
			fmt.Print("\n有效内容 => " + mytools.CheckRegexp(pageCont, regTitle, 0).(string))
		}
		fmt.Print("\n\n待爬抓网址共" + strconv.Itoa(len(strings.Split(string(file), "\n"))-1) + "个 => " + strings.Split(string(file), "\n")[0] + "\n")
		mytools.DelFirstText("./data/url.txt")
		SaveHtml()
	*/
}
