package gut

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
)

//返回字符串ms5值
func Md5(s string) string {

	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))

}

//json序列化
func JsonEn(obj interface{}) (string, error) {

	jsonStr, err := json.Marshal(obj)

	return string(jsonStr), err
}

//json反序列化 -map
func JsonDe(jsonStr string) (interface{}, error) {

	var mapResult map[string]interface{}

	err := json.Unmarshal([]byte(jsonStr), &mapResult)

	return mapResult, err

}
