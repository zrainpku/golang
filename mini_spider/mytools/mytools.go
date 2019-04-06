package mytools

import (
	"fmt"
	// "icode.baidu.com/baidu/go-lib/log"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func GetSeedUrl(path string) []string {
	f, _ := os.OpenFile("../data/url.data", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	allUrl, _ := ioutil.ReadAll(f)
	ret := []string(strings.Split(string(allUrl), "\n"))
	return ret
}

func ChangeUrl(url string) string {
	ret := []byte(url)
	for i, b := range ret {
		if (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') {
			continue
		} else {
			ret[i] = '_'
		}
	}
	fmt.Println("filaname:", string(ret))
	return string(ret)
}

func CheckErr(e error) {
	if e != nil {
		panic(e)
	}
}

func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func CheckFile(dir string, file string) os.FileInfo {
	list, _ := ioutil.ReadDir(dir)
	for _, info := range list {
		if info.Name() == file {
			return info
		}
	}
	return list[0]
}

func SaveFile(file string, cont string) {
	f, _ := os.OpenFile(file, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	defer f.Close()
	io.WriteString(f, cont)
}

func CheckRegexp(cont string, reg string, style int) (result interface{}) {
	check := regexp.MustCompile(reg)
	switch style {
	case 0:
		result = check.FindString(cont)
	case 1:
		result = check.FindAllString(cont, -1)
	default:
		result = check.FindAll([]byte(cont), -1)
	}
	return
}

func DelFirstText(file string) {
	var text = ""
	f, _ := os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0666)
	files, _ := ioutil.ReadAll(f)
	var ss = strings.Split(string(files), "\n")
	for i := 1; i < len(ss)-1; i++ {
		text += ss[i] + "\n"
	}
	defer f.Close()
	ioutil.WriteFile(file, []byte(text), 0666)
	fmt.Print("\n\n删除该地址 => " + ss[0])
}

func PageVisit(url string) (page string, body []byte) {
	resp, err := http.Get(url)
	if err != nil {
		// handle error by zr
		// log.Logger.Info(err)
		return
	}
	defer resp.Body.Close()
	body, _ = ioutil.ReadAll(resp.Body)
	page = string(body)
	return
}
