package gut

/**
 * This file is part of Gourd.
 *
 * @link     http://gourd.kyour.cn
 * @document http://gourd.kyour.cn/doc
 * @contact  kyour@vip.qq.com
 * @license  https://https://github.com/kyour-cn/gourd/blob/master/LICENSE
 */

import "strconv"

//int类型转字符串
func Int2str(i int) string {
	return strconv.Itoa(i)
}

//字符串装int
func Str2int(s string) (int, error) {
	return strconv.Atoi(s)
}
