package client

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/bcsimple/harborctl/internal/harborctl/app/cmd/root"
	"github.com/bcsimple/harborctl/pkg/config"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/registry"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"net/http"
	neturl "net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type Registry struct {
	*harborClient
	registryTitle []string
	ctx           context.Context
}

func NewRegistry(options *root.GlobalOptions) *Registry {
	return &Registry{
		ctx:          context.Background(),
		harborClient: NewHarborClient(options),
		registryTitle: []string{
			"序号ID",
			"仓库名称",
			"仓库url",
			"仓库状态",
			"仓库用户名",
			"仓库类型",
			"仓库描述",
		},
	}
}

func (r *Registry) DeleteRegistry(id int64) error {
	info, err := r.Registry.DeleteRegistry(r.ctx, &registry.DeleteRegistryParams{
		ID: id,
	})
	if err != nil {
		return err
	}
	fmt.Printf("registry: [ %s ] delete success\n", info.XRequestID)
	return nil
}

func (r *Registry) CreateRegistry() error {
	var verifyCreate string
	for {
		fmt.Printf("确定创建新的仓库(y/n):")
		fmt.Scanln(&verifyCreate)
		switch verifyCreate {
		case "Y", "y", "yes", "Yes":
			fmt.Printf("您的输入是:%s\n", verifyCreate)
			params := r.inputParams()
			if params.accessKey != "" && params.accessSecret != "" && params.registryDescription != "" && params.registryUrl != "" && params.registryName != "" {
				if err := r.createRegistry(params); err != nil {
					return fmt.Errorf("仓库: %s 创建失败! err:%s\n", params.registryName, err.Error())
				} else {
					fmt.Printf("仓库: %s 创建成功!\n", params.registryName)
					return nil
				}
			} else {
				fmt.Println("有字段为空! 禁止创建仓库!")
				return nil
			}
		case "N", "n", "No", "no":
			fmt.Printf("您的输入是:%s,\n正在退出....\n", verifyCreate)
			return nil
		default:
			fmt.Println("输入错误!请重新输入")
			continue
		}
	}
}

func (r *Registry) CreateRegistryByConfigInfo(name string, replaceName string) error {
	info, err := config.NewConnectConfiguration().GetConnectInfo(name)
	if err != nil {
		return err
	}
	registryName := info.Name
	if replaceName != "" {
		registryName = replaceName
	}
	return r.createRegistry(&RegistryInputParams{
		registryName:        registryName,
		registryUrl:         info.Host,
		registryDescription: strings.Join(info.Alias, ","),
		accessKey:           info.User,
		accessSecret:        info.Password,
	})
}

func (r *Registry) createRegistry(params *RegistryInputParams) error {
	registryNew := &models.Registry{
		Description: strings.TrimSpace(params.registryDescription),
		Name:        strings.TrimSpace(params.registryName),
		URL:         strings.TrimSpace(params.registryUrl),
		Insecure:    true,
		Type:        "harbor",
		Credential: &models.RegistryCredential{
			AccessKey:    strings.TrimSpace(params.accessKey),
			AccessSecret: strings.TrimSpace(params.accessSecret),
		},
	}
	registryN := registry.NewCreateRegistryParams()
	registryN.SetRegistry(registryNew)
	registryN.SetContext(context.Background())
	result, err := r.Registry.CreateRegistry(context.Background(), registryN)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Printf("registry: %s 创建成功! 仓库的id为: %s\n", params.registryName, result.XRequestID)
	return nil
}

func (r *Registry) SearchRegistry(name string) error {
	params := make(neturl.Values)
	params.Add("q", fmt.Sprintf("name=~%s", name))
	return r.searchRegistry(params.Encode())
}

func (r *Registry) searchRegistry(params string) error {
	r.URL.Path = RegistryURLPath
	request, err := http.NewRequest(http.MethodGet, r.URL.String(), nil)
	if err != nil {
		fmt.Println(request.RequestURI)
		os.Exit(1)
	}
	request.Header.Set("accept", "application/json")
	client := &http.Client{}
	request.SetBasicAuth(r.HarborConnectInfo.User, r.HarborConnectInfo.Password)

	request.URL.RawQuery = params
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	registryInfos := make([]*RegistryInfo, 0, 100)
	decoder := json.NewDecoder(response.Body)
	if err = decoder.Decode(&registryInfos); err != nil {
		return err
	}
	r.print(registryInfos)
	return nil
}

func (r *Registry) SearchRegistryList() error {
	return r.searchRegistry("")
}

func (r *Registry) SearchRegistryByID(id string, isPrint bool) (*models.Registry, error) {

	registryInfo, err := r.Registry.GetRegistry(context.Background(), &registry.GetRegistryParams{
		ID: parseStringToInt64(id),
	})
	if err != nil {
		return nil, err
	}
	registryInfos := []*models.Registry{
		registryInfo.Payload,
	}
	if isPrint {
		r.printRegistryTable(registryInfos)
	}
	return registryInfo.Payload, nil
}

// 打印
func (r *Registry) print(registryInfos []*RegistryInfo) {
	allData := make([][]string, 0)
	for _, registryInfo := range registryInfos {
		data := make([]string, 0, 0)
		data = append(data, strconv.FormatInt(registryInfo.ID, 10))
		data = append(data, registryInfo.Name)
		data = append(data, registryInfo.URL)
		data = append(data, registryInfo.Status)
		data = append(data, registryInfo.Credential.AccessKey)
		data = append(data, registryInfo.Type)
		data = append(data, registryInfo.Description)
		allData = append(allData, data)
	}
	r.TableInformation.SetTitles(r.registryTitle).SetData(allData).Output()
}

func (r *Registry) printRegistryTable(registryInfos []*models.Registry) {
	allData := make([][]string, 0)
	for _, registryInfo := range registryInfos {
		data := make([]string, 0, 0)
		data = append(data, strconv.FormatInt(registryInfo.ID, 10))
		data = append(data, registryInfo.Name)
		data = append(data, registryInfo.URL)
		data = append(data, registryInfo.Status)
		data = append(data, registryInfo.Credential.AccessKey)
		data = append(data, registryInfo.Type)
		data = append(data, registryInfo.Description)
		allData = append(allData, data)
	}
	r.TableInformation.SetTitles(r.registryTitle).SetData(allData).Output()
}

func (r *Registry) inputDefaultInfo(params *RegistryInputParams) {
	data := [][]string{
		{"1", "仓库名称:", params.registryName},
		{"2", "仓库地址:", params.registryUrl},
		{"3", "仓库描述:", params.registryDescription},
		{"4", "仓库用户名:", params.accessKey},
		{"5", "仓库密码:", params.accessSecret},
	}
	r.TableInformation.SetData(data).Output()
}

type RegistryInfo struct {
	ID              int64  `json:"id"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	Type            string `json:"type"`
	URL             string `json:"url"`
	TokenServiceURL string `json:"token_service_url"`
	Credential      struct {
		Type         string `json:"type"`
		AccessKey    string `json:"access_key"`
		AccessSecret string `json:"access_secret"`
	} `json:"credential"`
	Insecure     bool      `json:"insecure"`
	Status       string    `json:"status"`
	CreationTime time.Time `json:"creation_time"`
	UpdateTime   time.Time `json:"update_time"`
}

type RegistryInputParams struct {
	registryName        string
	registryUrl         string
	registryDescription string
	accessKey           string
	accessSecret        string
}

func (r *Registry) inputParams() *RegistryInputParams {
	var registry RegistryInputParams
	for {
		var op string
		fmt.Printf("操作符号(u:修改,d:删除,a:添加,p:查看,e:退出,s:保存):")
		fmt.Scanln(&op)
		switch strings.TrimSpace(op) {
		case "u", "U", "update", "up", "upda":
			r.update(&registry)
		case "a", "add":
			r.add(&registry)
		case "d", "delete":
			r.delete(&registry)
		case "p", "play":
			r.inputDefaultInfo(&registry)
		case "e", "exit", "quit", "q":
			os.Exit(1)
		case "s", "save":
			return &registry
		}
	}
}

func (r *Registry) add(registry *RegistryInputParams) {
	var reader = bufio.NewReader(os.Stdin)
	fmt.Printf("请输入仓库名称:")
	registry.registryName = tripeSpace(reader.ReadString('\n'))
	fmt.Printf("请输入仓库地址:")
	registry.registryUrl = tripeSpace(reader.ReadString('\n'))
	fmt.Printf("请输入仓库描述:")
	registry.registryDescription = tripeSpace(reader.ReadString('\n'))
	fmt.Printf("请输入仓库用户名:")
	registry.accessKey = tripeSpace(reader.ReadString('\n'))
	fmt.Printf("请输入仓库密码:")
	registry.accessSecret = tripeSpace(reader.ReadString('\n'))
	r.inputDefaultInfo(registry)
}

func (r *Registry) update(registry *RegistryInputParams) {
	number := prompt("更新")
	var reader = bufio.NewReader(os.Stdin)
	switch number {
	case "1":
		fmt.Printf("请输入仓库名称:")
		registry.registryName = tripeSpace(reader.ReadString('\n'))
	case "2":
		fmt.Printf("请输入仓库地址:")
		registry.registryUrl = tripeSpace(reader.ReadString('\n'))
	case "3":
		fmt.Printf("请输入仓库描述:")
		registry.registryDescription = tripeSpace(reader.ReadString('\n'))
	case "4":
		fmt.Printf("请输入仓库用户名:")
		registry.accessKey = tripeSpace(reader.ReadString('\n'))
	case "5":
		fmt.Printf("请输入仓库密码:")
		registry.accessSecret = tripeSpace(reader.ReadString('\n'))
	default:
		return
	}
}

func (r *Registry) delete(params *RegistryInputParams) {
	number := prompt("删除")
	switch number {
	case "1":
		params.registryName = ""
	case "2":
		params.registryUrl = ""
	case "3":
		params.registryDescription = ""
	case "4":
		params.accessKey = ""
	case "5":
		params.accessSecret = ""
	default:
		return
	}
}

func tripeSpace(msg string, err error) string {
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(msg)
}

func prompt(opAlias string) string {
	var number string
	fmt.Printf("请输入需要%s的id(1-5):", opAlias)
	fmt.Scan(&number)
	return number
}
