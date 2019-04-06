package urlbase

import (
	"fmt"
	// "io"
	// "icode.baidu.com/baidu/go-lib/log"
	"io/ioutil"
	// "minispider/mini_spider/mytools"
	"net/http"
	"os"
	"regexp"
	// "strconv"
	"strings"
	"time"
	"log"
)

func GetUrl(ch chan *Urls, ch_down chan string) {
	seed := <-ch
	seed1 := seed.Url
	deepth := seed.Depth

	file := "../log/" + time.Now().Format("20180102") + "_deepth.txt"
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if nil != err {
		panic(err)
	}
	loger := log.New(logFile, "前缀", log.Ldate|log.Ltime|log.Lshortfile)
	loger.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	loger.SetPrefix("test_deppth")
	loger.Output(2, seed1)
	loger.Output(2, string(deepth))
	fmt.Printf("url深度为＝＝%d\n",deepth)
	defer logFile.Close()

	// log.Logger.Info("当前种子URL地址为:" + seed1)
	// log.Logger.Info("当前URL深度为:" + string(deepth))
	var num, num_all int

	// urlWritePath := mytools.ChangeUrl(seed1)
	resp, _ := http.Get(strings.Split(seed1, "\"")[1])
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	// reg := regexp.MustCompile(`((ht|f)tps?)://[w]{0,3}.baidu.com/link\?[a-zA-z=0-9-\s]*`)
	reg := regexp.MustCompile(`href="((ht|f)tps?):.*?"`)
	// f, _ := os.OpenFile("../data/url.data", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	// defer f.Close()
	// fmt.Println("baidu.com page html:", body)
	// f, _ := os.Create("../data/" + string(urlWritePath) + ".data")
	// defer f.Close()
	for _, d := range reg.FindAllString(string(body), -1) {
		if deepth <= 0 {
			fmt.Println("这是最终地址，不能再爬取了，深度为0，url为：" + d)
			num_all++
			break
		}
		fmt.Printf("#############这是种子地址:%s\n ", d)
		num_all++
		num++
		var url Urls
		url.Url = d
		url.Depth = deepth - 1
		loger.Output(2, url.Url)
		loger.Output(2, string(url.Depth))
		if url.Depth == 0 {
			// log.Logger.Info("将该url放入下载页面的通道,因为URL深度为0。" + d)
			ch_down <- url.Url
		} else {
			// log.Logger.Info("当前URL放入种子库，作为新的种子URL地址：" + d + "深度为" + string(url.Depth))

			ch <- &url
			ch_down <- url.Url
		}

		/*
			ff, _ := os.OpenFile("../data/"+string(urlWritePath)+".data", os.O_RDWR, 0666)
			file, _ := ioutil.ReadAll(ff)
			dd := strings.Split(d, "")
			dddd := ""
			for _, ddd := range dd {
				if ddd == "?" {
					ddd = `\?`
				}
				dddd += ddd
			}
			if mytools.CheckRegexp(string(file), dddd, 0).(string) == "" {
				var temurl Urls
				temurl.Url = d
				temurl.Depth = deepth - 1
				ch <- &temurl
				io.WriteString(f, temurl.Url+"\n")
				if num_all%10 == 0 {
					fmt.Printf("已经检索地址数量：%d 个\n", num_all)
				}
				num++
			}
			ff.Close()
		*/

	}
	// log.Logger.Info("在种子地址: " + seed1)
	// log.Logger.Info("中总的发现了网络地址个数为：" + strconv.Itoa(len(reg.FindAllString(string(body), -1))))
	// log.Logger.Info("\n去重后网络地址数：" + strconv.Itoa(num))
	// log.Logger.Info("\n储存一级url成功！\n")
	return
}
