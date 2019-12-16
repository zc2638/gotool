package utilx

import (
	"os"
	"strings"
)

/**
 * Created by zc on 2019/12/16.
 */
//自动创建文件夹
func PathCreate(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		mkErr := os.MkdirAll(path, os.ModePerm)
		if mkErr != nil {
			return mkErr
		}
	}
	return nil
}

//自动生成全路径文件
func CreateFile(path string) (*os.File, error) {

	var pathArr = strings.Split(path, "/")
	var pathLen = len(pathArr)
	if pathLen > 1 {
		dir := strings.Join(pathArr[:pathLen-1], "/")
		if err := PathCreate(dir); err != nil {
			return nil, err
		}
	}
	return os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0766)
}