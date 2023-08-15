package config

import (
	"encoding/json"
	"fmt"
	"github.com/bcsimple/harborctl/utils/encrypt"
	"github.com/bcsimple/harborctl/utils/file"
	"io/ioutil"
	"os"
)

type SimpleConnectInfo struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Scheme   string `json:"scheme,omitempty"`
	Host     string `json:"host"`
}

func NewSimpleConnectInfo() *SimpleConnectInfo {

	s := &SimpleConnectInfo{}
	if err := json.Unmarshal(readConnectInfo(), s); err != nil {
		fmt.Println(err)
		panic(err)
	}

	if s.User == "" || s.Password == "" || s.Host == "" {
		panic("json unmarshal failed!")
	}
	s.Scheme = "http"
	p, err := encrypt.Decrypt(s.Password)
	if err != nil {
		fmt.Println(`为了安全起见! 请找张顺进行密码加密 不然使用不了任何操作 谢谢配合!`, err)
		os.Exit(1)
	}
	s.Password = p
	return s
}

func (s *SimpleConnectInfo) GetConnectInfo(string) (SimpleConnectInfo, error) {
	return *s, nil
}

func (s *SimpleConnectInfo) SetConnectInfo(info SimpleConnectInfo) error {
	return nil
}

func readConnectInfo() []byte {

	if _, ok := file.Exists(pathDir); !ok {
		fmt.Printf("%s 目录不存在 无法获取harbor的链接地址!\n", pathDir)
		os.Exit(1)
	}

	if _, ok := file.Exists(pathFile); !ok {
		fmt.Printf("%s 文件不存在 无法获取harbor的链接地址!\n", pathFile)
		os.Exit(1)
	}

	data, err := ioutil.ReadFile(pathFile)
	if err != nil {
		fmt.Println("read connect file failed !! path:", pathFile)
		panic(err)
	}
	return data
}
