package config

import (
	"errors"
	"strings"
)

type HarborConnectInfos []HarborConnectInfo

func (infos HarborConnectInfos) Lens() int {
	return len(infos)
}

func (infos *HarborConnectInfos) Update(info *HarborConnectInfo) {
	for index, i := range *infos {
		if i.Name == info.Name {
			(*infos)[index] = *info
		}
	}
}

func (infos *HarborConnectInfos) Append(info *HarborConnectInfo) error {
	if err := infos.check(*info); err != nil {
		return err
	}

	*infos = append(*infos, *encryptPassword(info))
	return nil
}

func (infos HarborConnectInfos) IsHas(info *HarborConnectInfo) bool {
	if infos.Lens() == 0 {
		return false
	}

	for _, i := range infos {
		if i.Name == info.Name {
			return true
		}
	}
	return false
}

func (infos *HarborConnectInfos) Del(info *HarborConnectInfo) {
	if infos.Lens() == 0 {
		return
	}
	tmp := make([]HarborConnectInfo, 0)
	for _, i := range *infos {
		if i.Name == info.Name {
			continue
		}
		tmp = append(tmp, i)
	}
	*infos = tmp
}

func (infos HarborConnectInfos) check(info HarborConnectInfo) error {
	//校验配置文件中的格式是否正确
	if strings.TrimSpace(info.Host) == "" {
		return errors.New("harbor server is null")
	}
	if strings.TrimSpace(info.Name) == "" {
		return errors.New("harbor name is null")
	}

	if strings.TrimSpace(info.User) == "" {
		return errors.New("harbor user is null")
	}
	if strings.TrimSpace(info.Password) == "" {
		return errors.New("harbor password is null")
	}

	if !strings.Contains(info.Host, ":") {
		return errors.New("server illegal format")
	}
	return nil
}

type HarborConnectInfoOpt func(info *HarborConnectInfo)

type HarborConnectInfo struct {
	Name     string   `yaml:"name" json:"name,omitempty"`
	Host     string   `yaml:"server" json:"host"`
	User     string   `yaml:"user" json:"user"`
	Alias    []string `yaml:"alias" json:"alias,omitempty"`
	Password string   `yaml:"password" json:"password"`
	Scheme   string   `yaml:"scheme,omitempty" json:"scheme,omitempty"`
}

func NewHarborConnectInfo(opts ...HarborConnectInfoOpt) *HarborConnectInfo {
	info := &HarborConnectInfo{
		Scheme: "http",
	}
	for _, opt := range opts {
		opt(info)
	}

	return info
}

func WithName(name string) HarborConnectInfoOpt {
	return func(info *HarborConnectInfo) {
		info.Name = name
	}
}

func WithUser(user string) HarborConnectInfoOpt {
	return func(info *HarborConnectInfo) {
		info.User = user
	}
}

func WithServer(host string) HarborConnectInfoOpt {
	return func(info *HarborConnectInfo) {
		info.Host = host
	}
}

func WithPassword(password string) HarborConnectInfoOpt {
	return func(info *HarborConnectInfo) {
		info.Password = password
	}
}

func WithAlias(alias []string) HarborConnectInfoOpt {
	return func(info *HarborConnectInfo) {
		if alias == nil {
			alias = make([]string, 0)
		}
		info.Alias = alias
	}
}
