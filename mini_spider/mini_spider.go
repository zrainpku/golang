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

var hash map[string]int
var lock sync.Mutex
var idx_url int

func saveUrl(){
	// timer := time.NewTimer(1 * time.Second)
 //    <-timer.C

	f,_ :=os.OpenFile("../data/ans.txt",os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	defer f.Close()
	ff, _ := os.OpenFile("../data/ans.txt", os.O_RDWR, 0666)
	file, _ := ioutil.ReadAll(ff)
	ff.Close()
	lock.Lock()

	// for idx:=idx_url;idx<len(hash);idx++{
	// 	key := hash[idx]
	// 	check := regexp.MustCompile(key)
	// 	check_ans :=check.FindAllString(string(file), -1)
	// 	if len(check_ans)==0{
	// 		fmt.Printf("找到了一个新的url: %s\n",key)
	// 		key+="\n"
	//         _,_=f.Write([]byte(key))
	// 	}
	// }
	for key,val :=range hash{
		if val<idx_url{
			continue
		}
		check := regexp.MustCompile(key)
		check_ans :=check.FindAllString(string(file), -1)
		if len(check_ans)==0{
			fmt.Printf("找到了一个新的url: %s\n",key)
			key+="\n"
	        _,_=f.Write([]byte(key))
		}
	}
	idx_url=len(hash)-1
	// hash = nil
	// hash =make(map[string]byte)
	lock.Unlock()
	// url+="\n"
	// _,_=f.Write([]byte(url))
	
	// saveUrl(ans)
}





func tempgetUrl(ch chan *urlbase.Urls, ch_down chan string) {
	seed := <-ch
	seed1 := seed.Url
	deepth := seed.Depth

	resp, err := http.Get(strings.Split(seed1, "\"")[1])
	if err != nil{
		return
		// fmt.Println("ERROR!")
	}

	//  porsible...
	ff, _ := os.OpenFile("../data/ans.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	file1, _ := ioutil.ReadAll(ff)
	check := regexp.MustCompile(seed1)
	check_ans :=check.FindAllString(string(file1), -1)
	if len(check_ans)==0{
		fmt.Printf("新的网页，把种子网页准备加入ans: %s\n" ,seed1)
		key := seed1 + "\n"	
		_,_=ff.Write([]byte(key))
		ff.Close()
	}else{
		ff.Close()
		return
	}

	// if _, ok := hash[seed1]; ok {
	// 	fmt.Printf("直接跳过: %s\n" ,seed1)
	// 	return
	// }

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
	fmt.Printf("当前的url深度为: %d\n",deepth)
	defer logFile.Close()


	// resp, err := http.Get(strings.Split(seed1, "\"")[1])
	// if err != nil{

	// }
	// defer resp.Body.Close()
	body, err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil{
		return
		// fmt.Println("ERROR!")
	}
	resp.Body.Close()
	reg := regexp.MustCompile(`href="((ht|f)tps?):.*?"`)

	reg_ans :=reg.FindAllString(string(body), -1)
	// fmt.Printf("当前的种子URL: %s\n" , strings.Split(seed1, "\"")[1])
	// fmt.Printf("当前的种子页面新的链接数为: %d\n" ,len(reg_ans))
	
	for _, d := range reg_ans {

		check_pku := regexp.MustCompile(`.*?pku.*?`)
	    check_ans_pku :=check_pku.FindAllString(d, -1)
	    if len(check_ans_pku)==0{
		    fmt.Printf("非官方网页 跳转\n")
		    continue
	    }

	    check_pdf := regexp.MustCompile(`.*?.pdf"`)
	    check_ans_pdf :=check_pdf.FindAllString(d, -1)
	    if len(check_ans_pdf)>0{
		    fmt.Printf("pdf网页 跳转\n")
		    continue
	    }

	    var idx_str int=strings.Index(d,"?")
	    fmt.Printf("当前的URL: %s\n" ,d)
	    fmt.Printf("带问号的 index %d \n",idx_str)
	    var dd string =d
	    if idx_str>=0{
	    	//save d    cmp dd
	    	ddd :=  strings.Split(d,"?")
	    	fmt.Printf("带问号的 left %s \n",ddd[0])
	    	fmt.Printf("带问号的 right %s \n",ddd[1])
	    	dd= ddd[0]
	    }

		f2,_ := os.OpenFile("../data/seed.txt",os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
		f3, _ := os.OpenFile("../data/ans.txt",os.O_APPEND|os.O_RDWR, 0666)
		file2, _ := ioutil.ReadAll(f2)
	    file3, _ := ioutil.ReadAll(f3)
	    check2 := regexp.MustCompile(dd)
	    check3 := regexp.MustCompile(dd)
	    check_ans2 := check2.FindAllString(string(file2), -1)
	    check_ans3 := check3.FindAllString(string(file3), -1)

	    if len(check_ans3)==0 && len(check_ans2)==0{
	    	// fmt.Printf("新的种子URL: %s\n" ,d)
	    	key2 := d + "\n"	
		    _,_=f2.Write([]byte(key2))
		    var url urlbase.Urls
		    url.Url = d
		    url.Depth = deepth - 1
		    ch <- &url
	    }else if len(check_ans3)==0{
	    	// continue
	    	// fmt.Printf("新的网页URL: %s\n" ,d)
	    	key3 := d + "\n"
	    	_,_=f3.Write([]byte(key3))
	    }else{
	    	fmt.Printf("next: \n")
	    }
	    f2.Close()
	    f3.Close()




		// if _, ok := hash[d]; ok {
		// 	fmt.Printf("直接跳过: %s\n" ,d)
		// 	continue
		// }else{
		// 	id :=len(hash)
		// 	hash[d]=id
		// 	fmt.Printf("++++111: %s,--------%d\n" ,d,len(d))
		// } 


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
	hash =make(map[string]int)
	idx_url=0
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

	ch := make(chan *urlbase.Urls, 1000)
	ch_down := make(chan string, 1000)
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
	for len(ch) != 0  {
		t := time.NewTimer(time.Second * 1)
		<-t.C

		fmt.Printf("len( ch) is :%d \n", len(ch))
		// log.Logger.Info("种子通道含有的URL数量为（最多50）:" + string(len(ch)))

		// if len(ch) >0{
			// go func() {
				// urlbase.GetUrl(ch, ch_down)
				tempgetUrl(ch,ch_down)
			// }()
		// }


		// if (len(hash)-idx_url)>50 || len(ch)<2{
		// 	fmt.Printf("len(ch) is :%d ,len(hash) is : %d, idx_url=%d, \n", len(ch),len(hash),idx_url)
		// 	saveUrl()
		// }else if len(hash)%10==0{
		// 	fmt.Printf("length of **hash** is :%d\n", len(hash))
		// }
		// if len(hash)>600 {
		// 	time.Sleep(time.Duration(10)*time.Second)
		// }
		// fmt.Printf("length of ch_down is :%d\n", len(ch_down))
		timer := time.NewTimer(4 * time.Second)
        <-timer.C
	}
	close(ch)
	close(ch_down)

	waitGroup.Wait()

}