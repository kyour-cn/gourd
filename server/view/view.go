package view

import (
	"github.com/kyour-cn/gourd/common"
	"github.com/kyour-cn/gourd/utils/gut"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
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
	file  []string
	Data  map[string]interface{}
	funcs []template.FuncMap
}

var config Config

type Config struct {
	ViewPath   string `toml:"view_dir"`
	ViewSuffix string `toml:"view_suffix"`
	init       bool
}

//新建一个模板
func New(tnames ...string) (t Template) {

	if !config.init {
		err := common.ReadConfig("view", &config)
		if err != nil {
			log.Printf("View 配置错误：%v\n", err)
		}
		config.init = true
	}

	t.Data = make(map[string]interface{})

	//参数可传可不传
	for _, n := range tnames {

		t.file = append(t.file, config.ViewPath+n+config.ViewSuffix)

		//t.file = config.ViewPath + n + config.ViewSuffix
	}

	//定义row函数，禁止转义
	t.AddFunc("row", func(s string) template.HTML {
		return template.HTML(s)
	})

	//添加模板函数，时间格式化
	t.AddFunc("date", func(f string, t uint) string {
		return gut.Date(f, int64(t))
	})

	t.AddFunc("default", func(f interface{}, s interface{}) interface{} {
		if f == nil || f == "" {
			return s
		} else {
			return f
		}
	})

	return t
}

//添加模板函数
func (t *Template) AddFunc(n string, f interface{}) {
	t.funcs = append(t.funcs, template.FuncMap{n: f})
}

//设定模板路径
func (t *Template) AddFile(n string) {
	t.file = append(t.file, config.ViewPath+n+config.ViewSuffix)
}

//批量设定模板路径
func (t *Template) AddFiles(n ...string) {
	for _, v := range n {
		t.file = append(t.file, config.ViewPath+v+config.ViewSuffix)
	}
}

//设定全路径模板位置，不会自动增加前缀、后缀
func (t *Template) AddFileFull(file string) {
	t.file = append(t.file, file)
}

//模板变量赋值
func (t *Template) Assign(name string, data interface{}) {
	t.Data[name] = data
}

//设置全部模板变量
func (t *Template) SetAssign(d map[string]interface{}) {
	t.Data = d
}

//渲染输出
func (t *Template) Fetch(w http.ResponseWriter) (err error) {

	var temp *template.Template

	for _, f := range t.file {
		name := filepath.Base(f)
		temp = template.New(name)
		break
	}

	//遍历添加函数
	for _, f := range t.funcs {
		temp = temp.Funcs(f)
	}

	temp, err = temp.ParseFiles(t.file...)

	if err == nil {

		err = temp.Execute(w, t.Data)
		if err != nil {
			log.Printf("Template Execute Err:%v", err)

		}
	} else {
		log.Printf("Template ParseFiles Err:%v", err)
	}

	return err

}
