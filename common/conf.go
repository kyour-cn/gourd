package common

import (
	"github.com/kyour-cn/gourd/utils/toml"
)

var configPath string

//var cache map[string]interface{}

//设定应用配置路径
func SetConfigPath(path string) {

	configPath = path

	//cache = make(map[string]interface{})

}

//取得配置
func ReadConfig(name string, v interface{}) error {

	configPath = "./conf/"

	file := configPath + name + ".conf"

	//log.Printf("Config:%v", file)

	//判断是否存在
	//_, ok := cache[name]
	//if ok {
	//	//直接赋值
	//	v = cache[name]
	//
	//	return nil
	//
	//}

	_, err := toml.DecodeFile(file, v)

	//cache[name] = &v

	return err

}

//取得map类型自定义配置
func ReadConfigMap(name string) (interface{}, error) {

	file := configPath + name + ".conf"

	var v interface{}

	_, err := toml.DecodeFile(file, &v)

	return v, err

}
