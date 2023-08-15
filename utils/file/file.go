package file

import (
	"io/ioutil"
	"os"
)

func ReadFile(path string) ([]byte, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func Exists(path string) (isDir bool, isExist bool) {
	info, err := os.Stat(path) //os.Stat获取文件信息

	//判断文件存在与否 不存在直接false isDir 也是false
	if os.IsNotExist(err) {
		return false, false
	}

	if info.IsDir() {
		return true, true
	}
	return false, true
}
