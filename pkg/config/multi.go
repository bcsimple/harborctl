package config

import (
	"errors"
	"fmt"
	"github.com/bcsimple/harborctl/utils/encrypt"
	"github.com/bcsimple/harborctl/utils/file"
	"github.com/bcsimple/harborctl/utils/table"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

const (
	apiVersion = "v1"
	kind       = "Config"
)

type ConnectConfiguration struct {
	APIVersion     string             `json:"apiVersion" yaml:"apiVersion"`
	Kind           string             `json:"kind" yaml:"kind"`
	Harbors        HarborConnectInfos `json:"harbors" yaml:"harbors"`
	CurrentContext string             `json:"current-context" yaml:"current-context"`
}

func NewConnectConfiguration() HarborConnectInfoInterface {

	c := &ConnectConfiguration{
		APIVersion: apiVersion,
		Kind:       kind,
		Harbors:    make([]HarborConnectInfo, 0),
	}

	if _, ok := file.Exists(pathFile); !ok {
		return c
	}

	d, err := file.ReadFile(pathFile)
	if err != nil {
		fmt.Println("read config file failed error:", err)
		panic(err)
	}
	if err = yaml.Unmarshal(d, c); err != nil {
		fmt.Println("parse config failed error:", err)
		panic(err)
	}
	return c
}

// 核心操作功能
func (c *ConnectConfiguration) GetConnectInfo(name string) (*HarborConnectInfo, error) {
	if _, info, err := c.getConnectInfoByNameOrAlias(name); err == nil {
		return info, nil
	}
	return &HarborConnectInfo{}, errors.New("get connect info failed ")
}

func (c *ConnectConfiguration) GetDefaultConnectInfo() *HarborConnectInfo {
	info, err := c.GetConnectInfo(c.CurrentContext)
	if err != nil {
		fmt.Println("get default connect info failed")
		panic(err)
	}
	return info
}

// 添加
func (c *ConnectConfiguration) SetConnectInfo(info *HarborConnectInfo) {
	if c.Harbors.IsHas(info) {
		fmt.Printf("%s 已存在相同的名字\n", info.Name)
		return
	}

	if err := c.Harbors.Append(info); err != nil {
		fmt.Println("append failed ", err)
		return
	}
	if c.CurrentContext == "" {
		c.setDefaultContext()
	}
	defer func() {
		if err := c.reWrite(); err != nil {
			panic(err)
		}
	}()
}

func (c *ConnectConfiguration) UpdateConnectInfo(info *HarborConnectInfo) {
	if c.Harbors.IsHas(info) {
		c.Harbors.Update(encryptPassword(info))
		return
	}
}

// 删除
func (c *ConnectConfiguration) DelConnectInfo(name string) {
	info, err := c.GetConnectInfo(name)
	if err != nil {
		panic(err)
	}
	c.Harbors.Del(info)

	if c.Harbors.Lens() == 0 {
		c.CurrentContext = ""
	}
	defer func() {
		_ = c.reWrite()
	}()
}

// 设置上下文
func (c *ConnectConfiguration) SetConnectInfoContext(name string) {
	c.setContext(name)
}

func (c *ConnectConfiguration) SetConnectInfoAlias(name string, aliasName string) error {
	if err := c.NameOrAliasNamesMustUnique(aliasName); err != nil {
		return err
	}
	_, info, err := c.getConnectInfoByNameOrAlias(name)
	if err != nil {
		return err
	}
	tmpMap := make(map[string]string, 0)

	for _, a := range info.Alias {
		if _, ok := tmpMap[a]; !ok {
			tmpMap[a] = ""
		}
	}
	tmpMap[aliasName] = ""
	tmpSlice := make([]string, 0)
	for k, _ := range tmpMap {
		tmpSlice = append(tmpSlice, k)
	}
	info.Alias = tmpSlice
	if c.CurrentContext == "" {
		c.setDefaultContext()
	}
	c.UpdateConnectInfo(info)
	defer func() {
		if err = c.reWrite(); err != nil {
			panic(err)
		}
	}()
	return nil
}

// prefix check和名字必须别名必须全局唯一
func (c *ConnectConfiguration) NameOrAliasNamesMustUnique(name string) error {
	for _, harbor := range c.Harbors {
		if name == harbor.Name {

			return fmt.Errorf("%s always exists\n", name)
		}
		for _, a := range harbor.Alias {
			if a == name {

				return fmt.Errorf("aliasname %s always on %s object exists", name, harbor.Name)
			}
		}
	}
	return nil
}

func (c *ConnectConfiguration) List(style string, onlyName, isDecrypt bool) {

	data := make([][]string, 0)
	title := []string{"序号", "名字", "Harbor地址", "用户名", "密码", "别名"}
	for i, info := range c.Harbors {
		pass := info.Password
		if isDecrypt {
			pass, _ = encrypt.Decrypt(info.Password)
		}
		if onlyName {
			data = append(data, []string{strconv.Itoa(i), info.Name})
			continue
		}
		data = append(data, []string{strconv.Itoa(i), info.Name, info.Host, info.User, pass, strings.Join(info.Alias, ",")})
	}
	if onlyName {
		title = []string{"序号", "名字"}
	}
	tableWrite := table.NewTableInformation(style)
	tableWrite.SetTitles(title).SetData(data).Output()

}

func (c *ConnectConfiguration) setDefaultContext() {
	//不存在直接设置为空 存在没一个或者多个直接使用第一个
	if c.Harbors.Lens() == 0 {
		c.CurrentContext = ""
		return
	}
	c.CurrentContext = c.Harbors[0].Name
}

func (c *ConnectConfiguration) setContext(name string) {
	//不存在直接设置为空 存在没一个或者多个直接使用第一个
	if c.Harbors.Lens() == 0 {
		fmt.Println("harbor info list is null")
		return
	}
	n, _, err := c.getConnectInfoByNameOrAlias(name)
	if err != nil || n == "" {
		fmt.Printf("%s harbor entry not found \n", name)
		return
	}
	c.CurrentContext = n
	defer func() {
		if err := c.reWrite(); err != nil {
			panic(err)
		}
	}()
}

// 判断名字或者别名是否存在
func (c *ConnectConfiguration) getConnectInfoByNameOrAlias(name string) (string, *HarborConnectInfo, error) {
	for _, info := range c.Harbors {
		//若名字相同 直接返回名字即可
		if info.Name == name {
			return info.Name, decryptPassword(&info), nil
		}
		//如果名字不相同 直接检查别名中是否包含
		if len(info.Alias) != 0 {
			for _, alias := range info.Alias {
				if strings.TrimSpace(alias) == name {
					return name, decryptPassword(&info), nil
				}
			}
		}
	}
	return "", &HarborConnectInfo{}, errors.New("name not found")
}

func (c *ConnectConfiguration) reWrite() error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	//目录不存在直接创建
	if _, ok := file.Exists(pathDir); !ok {
		if err = os.MkdirAll(pathDir, 0755); err != nil {
			panic(err)
		}
	}
	if err = ioutil.WriteFile(pathFile, data, 0644); err != nil {
		return err
	}
	return nil
}

// 打印所有的信息
func (c *ConnectConfiguration) list() map[string][]string {
	l := make(map[string][]string, 0)
	for _, info := range c.Harbors {
		l[info.Name] = info.Alias
	}
	return l
}

func (c *ConnectConfiguration) View() {
	d, err := ioutil.ReadFile(pathFile)
	if err != nil {
		fmt.Println("read config file error path: ", pathFile)
	}
	fmt.Printf(string(d))
}

func (c *ConnectConfiguration) PrintCurrentContext(completeInfo bool) {
	fmt.Printf("当前上下文: %s\n", c.CurrentContext)
	if completeInfo {
		info := c.GetDefaultConnectInfo()
		fmt.Printf("名称: %s,主机: %s,用户: %s,密码: %s,别名: %s\n", info.Name, info.Host, info.User, info.Password, strings.Join(info.Alias, ","))
	}
}

func encryptPassword(info *HarborConnectInfo) *HarborConnectInfo {
	info.Password, _ = encrypt.Encrypt(info.Password)
	return info
}
func decryptPassword(info *HarborConnectInfo) *HarborConnectInfo {
	info.Password, _ = encrypt.Decrypt(info.Password)
	return info
}
