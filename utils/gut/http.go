package gut

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

//http请求-Get
func HttpGet(url string) string {

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer resp.Body.Close()

	var buffer [512]byte
	result := bytes.NewBuffer(nil)
	for {
		n, err := resp.Body.Read(buffer[0:])
		result.Write(buffer[0:n])
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
	}

	return result.String()
}

//响应json数据
func Rjosn(w http.ResponseWriter, res interface{}) {

	rets, _ := JsonEn(res)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, _ = w.Write([]byte(rets))
}
