package main

import (
	"flag"
	"fmt"
	// "icode.baidu.com/baidu/go-lib/log"
	"pkuspider/mini_spider/mytools"
	// "pkuspider/mini_spider/pagebase"
	"pkuspider/mini_spider/urlbase"
	"time"
	"log"
	// "gcfg"
	// "gopkg.in/gcfg.v1"
	"gopkg.in/gcfg.v1"
	"io/ioutil"
	"os"
	"sync"
	"net/http"
	"strings"
	"regexp"
)

var (
	v bool
	h bool
	l string
	c string
)

func init() {
	flag.BoolVar(&v, "v", false, "show version and exit")
	flag.BoolVar(&h, "h", false, "this help")
	flag.StringVar(&l, "l", "../log", "log file path")
	flag.StringVar(&c, "c", "../conf", "config file path")
}

func usage() {
	fmt.Fprintf(os.Stderr, `nginx version: nginx/1.10.0
Usage: nginx [-hvVtTq] [-s signal] [-c filename] [-p prefix] [-g directives]

Options:
`)
	flag.PrintDefaults()
}

func useVersion() {
	fmt.Fprintf(os.Stderr, `the version is mini_spider_1.10`)
}
func saveUrl(url string){
	timer := time.NewTimer(10 * time.Second)
    <-timer.C

	f,_ :=os.OpenFile("../data/ans.txt",os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	defer f.Close()
	url+="\n"
	_,_=f.Write([]byte(url))
	
	// saveUrl(ans)
}
var hash map[string]byte

func tempgetUrl(ch chan *urlbase.Urls, ch_down chan string) {
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

	var num, num_all int

	resp, _ := http.Get(strings.Split(seed1, "\"")[1])
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	reg := regexp.MustCompile(`href="((ht|f)tps?):.*?pku.*?"`)
	for _, d := range reg.FindAllString(string(body), -1) {
		if _, ok := hash[d]; ok {
			continue
		} 
		if deepth <= 0 {
			fmt.Println("这是最终地址，不能再爬取了，深度为0，url为：" + d)
			num_all++
			break
		}
		fmt.Printf("#############这是种子地址:%s\n ", d)
		num_all++
		num++
		var url urlbase.Urls
		url.Url = d
		url.Depth = deepth - 1
		loger.Output(2, url.Url)
		loger.Output(2, string(url.Depth))
		if url.Depth == 0 {
			ch_down <- url.Url
		} else {

			ch <- &url
			ch_down <- url.Url
		}

	}
	return
}

func main() {

	cfig := struct {
		Section struct {
			UrlListFile     string
			OutputDirectory string
			MaxDepth        int
			CrawlInterval   int
			CrawlTimeout    int
			TargetUrl       string
			ThreadCount     int
		}
	}{}
	hash =make(map[string]byte)
	cfgFile, _ := os.OpenFile("../conf/spider.conf", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	cfgStr, _ := ioutil.ReadAll(cfgFile)
	err := gcfg.ReadStringInto(&cfig, string(cfgStr))
	if err != nil {
		fmt.Printf("Failed to parse gcfg data: %s", err)
	}
	fmt.Println("###################", cfig.Section.UrlListFile)
	fmt.Println("###################", cfig.Section.OutputDirectory)

	file := "../log/" + time.Now().Format("20180102") + ".txt"
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if nil != err {
		panic(err)
	}
	loger := log.New(logFile, "前缀", log.Ldate|log.Ltime|log.Lshortfile)
	loger.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	loger.SetPrefix("test_")
	loger.Output(2, "打印一条日志信息")

	// c := flag.String("c", "../conf", "config file path")
	// l := flag.String("l", "../log", "log file path")
	flag.Parse()
	if h {
		flag.Usage()
		return
	}
	if v {
		useVersion()
		return
	}
	// log.Logger.Warn("warn msg")
	// log.Logger.Info("info msg")

	// log.Logger.Info("开始执行程序，并读取种子文件url.data")
	seedUrl := []string{}
	seedUrl = mytools.GetSeedUrl(cfig.Section.UrlListFile)
	// log.Logger.Info("读取种子文件成功文件地址和名字:" + cfig.Section.UrlListFile)
	fmt.Println(seedUrl)
	maxDepth := cfig.Section.MaxDepth

	ch := make(chan *urlbase.Urls, 1)
	ch_down := make(chan string, 1)
	var url urlbase.Urls
	var waitGroup sync.WaitGroup
	for _, ur := range seedUrl {
		// log.Logger.Info(" 起始种子Url为:" + ur)
		url.Url = ur
		url.Depth = maxDepth
		ch <- &url
	}
	// hash :=make(map[string]byte)
	// go saveUrl(hash)
	for len(ch) != 0 || len(ch_down) != 0 {
		t := time.NewTimer(time.Second * 1)
		<-t.C

		fmt.Printf("length of ch is :%d \n", len(ch))
		// log.Logger.Info("种子通道含有的URL数量为（最多50）:" + string(len(ch)))
		if len(ch) != 0 {
			go func() {
				// urlbase.GetUrl(ch, ch_down)
				tempgetUrl(ch,ch_down)
			}()
		}

		fmt.Printf("length of ch_down is :%d\n", len(ch_down))
		if len(ch_down) != 0 {
			go func() {
				down_url := <-ch_down
				loger.Output(2,down_url)
				hash[down_url]=1
				fmt.Printf("length of hash all url is%d\n",len(hash))
				saveUrl(down_url)
				// log.Logger.Info("开始下载web:%s \n", string(down_url))
				// pagebase.SaveHtml(down_url)
				// log.Logger.Info("下载web完毕.%s \n", string(down_url))
			}()

		}
	}
	close(ch)
	close(ch_down)

	waitGroup.Wait()




}
