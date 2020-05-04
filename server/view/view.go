package view

import (
	"github.com/kyour-cn/gourd/common"
	"html/template"
	"log"
	"net/http"
)

/**
 * This file is part of Gourd.
 *
 * @link     http://gourd.kyour.cn
 * @document http://gourd.kyour.cn/doc
 * @contact  kyour@vip.qq.com
 * @license  https://https://github.com/kyour-cn/gourd/blob/master/LICENSE
 */

type Template struct {
	file string
	Data map[string]interface{}
}

var config Config

type Config struct {
	ViewPath   string `toml:"view_dir"`
	ViewSuffix string `toml:"view_suffix"`
	init       bool
}

//新建一个模板
func New(tname ...string) (t Template) {

	if !config.init {
		err := common.ReadConfig("view", &config)
		if err != nil {
			log.Printf("View 配置错误：%v\n", err)
		}
		config.init = true
	}

	t.Data = make(map[string]interface{})

	//参数可传可不传
	for _, n := range tname {
		t.file = config.ViewPath + n + config.ViewSuffix
		return t
	}

	return t
}

//设定模板路径
func (t *Template) SetFile(n string) {
	t.file = config.ViewPath + n + config.ViewSuffix
}

//设定全路径模板位置，不会自动增加前缀、后缀
func (t *Template) SetFileFull(file string) {
	t.file = file
}

//添加模板数据
func (t *Template) AddData(name string, data interface{}) {
	t.Data[name] = data
}

//设置模板数据
func (t *Template) SetData(d map[string]interface{}) {
	t.Data = d
}

//渲染输出
func (t *Template) Fetch(w http.ResponseWriter) error {

	temp, err := template.ParseFiles(t.file)

	if err == nil {
		err = temp.Execute(w, t.Data)
	}

	return err

}
