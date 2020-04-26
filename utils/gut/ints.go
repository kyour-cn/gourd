package gut

import "strconv"

//int类型转字符串
func Int2str(i int) string {
	return strconv.Itoa(i)
}

//字符串装int
func Str2int(s string) (int, error) {
	return strconv.Atoi(s)
}
