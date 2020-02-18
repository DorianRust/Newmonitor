package function

import (
	"time"
	"fmt"
	"net/http"
	"io/ioutil"
)


func Get(httpUrl string) float64 {
	start := time.Now()
	result, err := http.Get(httpUrl)

	if err != nil {
		fmt.Print("error", err)
		return -1
	}
	defer func(result  *http.Response) {
		if (result != nil){
			result.Body.Close()
		}

	}(result)
	elapsed := time.Since(start).Seconds()


	return elapsed
}

func CheckHeathcheck(httpUrl string) string {
	result, err := http.Get(httpUrl)

	if err != nil {
		fmt.Print("error", err)
		return "Down"
	}
	defer func(result  *http.Response) {
		if (result != nil){
			result.Body.Close()
		}

	}(result)
	s, _ := ioutil.ReadAll(result.Body) //把  body 内容读入字符串 s
	var content string
	content = fmt.Sprintf("%s", s)     //在返回页面中显示内容。
	if (len(content) != 0 && content == "OK"){
		return "OK"
	}
	return  "Down"
}
