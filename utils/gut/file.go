package gut

/**
 * This file is part of Gourd.
 *
 * @link     http://gourd.kyour.cn
 * @document http://gourd.kyour.cn/doc
 * @contact  kyour@vip.qq.com
 * @license  https://https://github.com/kyour-cn/gourd/blob/master/LICENSE
 */

import (
	"os"
)

/**
 * 判断路径是否存在
 * @param path 路径
 * @return bool 是否已存在
 */
func PathExists(path string) (exist bool, err error) {

	_, err = os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

/**
 * 检查并创建文件夹
 * @param dir 路径
 * @return code 0=出错，1=成功，2=已存在
 */
func Mkdir(dir string) (code int, err error) {

	exist, err := PathExists(dir)
	if err != nil {
		return
	}

	if exist {
		//2=已存在
		code = 2
	} else {
		// 创建文件夹
		err = os.Mkdir(dir, os.ModePerm)
		if err == nil {
			code = 1
		}
	}
	return
}
