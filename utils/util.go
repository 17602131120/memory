package utils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type MMUtil struct {
}

func (this *MMUtil) MatchSpiderName(runtimeName string) string {

	//sungq/spiderDianping/app/spiders.(*AllShop).Run 取出结构体名称
	valid := regexp.MustCompile("\\(\\*.{0,}\\)")
	Name := valid.FindString(runtimeName)
	Name = strings.Replace(Name, "(*", "", 1)
	Name = strings.Replace(Name, ")", "", 1)

	return Name

}
func (this *MMUtil) GetDifferentCode(prefix string) string {

	code := fmt.Sprintf("%s_%s", prefix, time.Now().Format("20060102150405"))

	for i := 0; i < 10; i++ {
		code = fmt.Sprintf("%s%d", code, rand.Intn(10))

	}
	return code
}
func (this *MMUtil) GetgoID() int {

	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)

	return int(n)
}

//解析 map[string]interface{} 数据格式
func Printmap(m map[string]interface{}) {
	for k, v := range m {
		switch value := v.(type) {
		case nil:
			log.Println(k, "is nil", "null")
		case string:
			log.Println(k, "is string", value)
		case int:
			log.Println(k, "is int", value)
		case float64:
			log.Println(k, "is float64", value)
		case []interface{}:
			log.Println(k, "is an array:")
			for i, u := range value {
				log.Println(i, u)
			}
		case map[string]interface{}:
			log.Println(k, "is an map:")
			Printmap(value)
		default:
			log.Println(k, "is unknown type", fmt.Sprintf("%T", v))
		}
	}
}

func (this *MMUtil) substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

func (this *MMUtil) GetParentDirectory(dirctory string) string {
	return this.substr(dirctory, 0, strings.LastIndex(dirctory, "/"))
}

func (this *MMUtil) GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

func (this *MMUtil) FileRead(filename string) []string {

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	var (
		str string
		seq string
	)
	str = string(b)

	goos := runtime.GOOS

	if "windows" == goos {
		seq = "\r\n"
	} else if "darwin" == goos { //macOS
		seq = "\n"
	} else {
		seq = "\n"
	}

	return strings.Split(str, seq)

}

func (this *MMUtil) FileUpdate(filename string, seq string, lineLen int, replaceIndex int, replaceStr string) bool {

	temps := this.FileRead(filename)
	var temps_array []string
	for _, value := range temps {

		_line := strings.Split(value, ",")
		if len(_line) == lineLen {

			if _line[0] == seq {
				_line[replaceIndex] = replaceStr
				value = strings.Join(_line, ",")
				//log.Println(value)
			}

		}
		temps_array = append(temps_array, value)
	}

	d1 := []byte(strings.Join(temps_array, "\n"))
	err := ioutil.WriteFile(filename, d1, 0644)
	if err != nil {
		panic(err)

		return false
	}

	return true

}
