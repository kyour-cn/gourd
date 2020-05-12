package gut

import (
	"strings"
	"time"
)

//获取当前时间戳
func Time() int64 {
	return time.Now().Unix()
}

//time时间戳转字符串
func Date(format string, i int64) string {

	if i == -1 {
		i = time.Now().Unix()
	}

	if format == "" {
		format = "Y-m-d H:i:s"
	}

	format = strings.Replace(format, "Y", "2006", -1)
	format = strings.Replace(format, "m", "01", -1)
	format = strings.Replace(format, "d", "02", -1)
	format = strings.Replace(format, "H", "15", -1)
	format = strings.Replace(format, "i", "04", -1)
	format = strings.Replace(format, "s", "05", -1)

	return time.Unix(i, 0).Format(format)
}
