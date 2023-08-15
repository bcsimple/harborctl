package client

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/bcsimple/harborctl/internal/harborctl/app/cmd/root"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/replication"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"io/ioutil"
	"net/http"
	neturl "net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type Replication struct {
	*harborClient
	replicationTitle []string
}

func NewReplication(options *root.GlobalOptions) *Replication {
	return &Replication{
		harborClient: NewHarborClient(options),
		replicationTitle: []string{
			"序号ID",
			"名称",
			"目标目标url",
			"目标名字空间",
			"触发类型",
			"过滤名称",
			"过滤资源",
			"是否开启",
			"是否覆盖",
			"复制描述",
		},
	}
}

func (r *Replication) ModifyReplication(id string) {
	replcationInfo := r.GetReplicationInfo(id)
	var (
		name      string
		namespace string
		isStart   string
		reader    = bufio.NewReader(os.Stdin)
	)
	fmt.Printf("请输入需要更新的过滤名称:")
	name = tripeSpace(reader.ReadString('\n'))
	fmt.Printf("请输入需要更新的目标名字空间:")
	namespace = tripeSpace(reader.ReadString('\n'))
	if strings.TrimSpace(namespace) != "" {
		replcationInfo.DestNamespace = namespace
	}
	if strings.TrimSpace(name) != "" {
		for _, filter := range replcationInfo.Filters {
			if filter.Type == "name" {
				filter.Value = name
			}
		}
	}
	result, err := r.Replication.UpdateReplicationPolicy(context.Background(), &replication.UpdateReplicationPolicyParams{
		ID:     parseStringToInt64(id),
		Policy: replcationInfo,
	})
	if err != nil {
		fmt.Println("更新复制规则失败! err:", err.Error())
		os.Exit(1)
	}
	fmt.Println(result.XRequestID, "更新成功!")

	fmt.Printf("请确认是否开始复制:y/n:")
	fmt.Scan(&isStart)
	switch isStart {
	case "Y", "yes", "y", "Yes":
		r.StartExecution(id)
	case "no", "No", "n", "N":
		os.Exit(1)
	}
}

func (r *Replication) SearchReplication(name string) {
	if name == "" {
		fmt.Println("name is empty. please input again!")
		return
	}

	url := &neturl.URL{
		Host: r.HarborConnectInfo.Host,
		Path: "/api/v2.0/replication/policies",
		//User:   neturl.UserPassword("admin", "123qqq...A"),
		Scheme: r.HarborConnectInfo.Scheme,
	}
	request, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		os.Exit(1)
	}

	request.Header.Set("accept", "application/json")
	client := &http.Client{}

	request.SetBasicAuth(r.HarborConnectInfo.User, r.HarborConnectInfo.Password)
	params := make(neturl.Values)
	params.Add("name", name)

	request.URL.RawQuery = params.Encode()
	response, err := client.Do(request)

	if err != nil {
		panic(err)
	}
	registryInfos := make([]*ReplicationInfo, 0, 100)
	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&registryInfos); err != nil {
		panic(err)
	}
	r.printReplicationTable1(registryInfos)
}

func (r *Replication) GetReplicationInfo(id string) *models.ReplicationPolicy {
	replcationInfo, err := r.SearchReplicationByID(id, false)
	if err != nil {
		fmt.Println("获取复制规则失败,id:", id)
	}
	return replcationInfo
}
func (r *Replication) SearchReplicationByID(id string, isPrint bool) (*models.ReplicationPolicy, error) {
	replicationInfo, err := r.Replication.GetReplicationPolicy(context.Background(), &replication.GetReplicationPolicyParams{
		ID: parseStringToInt64(id),
	})
	if err != nil {
		return nil, err
	}
	if isPrint {
		r.printReplicationTable([]*models.ReplicationPolicy{replicationInfo.Payload})
	}
	return replicationInfo.Payload, nil
}

func (r *Replication) printReplicationTable1(replicationInfos []*ReplicationInfo) {
	allData := make([][]string, 0)
	for _, replicationInfo := range replicationInfos {
		data := make([]string, 0, 0)
		data = append(data, strconv.FormatInt(replicationInfo.ID, 10))
		data = append(data, replicationInfo.Name)
		data = append(data, replicationInfo.DestRegistry.URL)
		data = append(data, replicationInfo.DestNamespace)
		data = append(data, replicationInfo.Trigger.Type)
		filterNameValue := ""
		filterNameResource := ""
		for _, filter := range replicationInfo.Filters {
			if filter.Type == "name" {
				filterNameValue = filter.Value.(string)
			}
			if filter.Type == "resource" {
				filterNameResource = filter.Value.(string)
			}
		}
		data = append(data, filterNameValue)
		data = append(data, filterNameResource)
		data = append(data, strconv.FormatBool(replicationInfo.Enabled))
		data = append(data, strconv.FormatBool(replicationInfo.Override))
		data = append(data, replicationInfo.Description)
		allData = append(allData, data)
	}
	r.TableInformation.SetTitles(r.replicationTitle).SetData(allData).Output()
}

func (r *Replication) printReplicationTable(replicationInfos []*models.ReplicationPolicy) {
	allData := make([][]string, 0)
	for _, replicationInfo := range replicationInfos {
		data := make([]string, 0, 0)
		data = append(data, strconv.FormatInt(replicationInfo.ID, 10))
		data = append(data, replicationInfo.Name)
		data = append(data, replicationInfo.DestRegistry.URL)
		data = append(data, replicationInfo.DestNamespace)
		data = append(data, replicationInfo.Trigger.Type)
		filterNameValue := ""
		filterNameResource := ""
		for _, filter := range replicationInfo.Filters {
			if filter.Type == "name" {
				filterNameValue = filter.Value.(string)
			}
			if filter.Type == "resource" {
				filterNameResource = filter.Value.(string)
			}
		}
		data = append(data, filterNameValue)
		data = append(data, filterNameResource)
		data = append(data, strconv.FormatBool(replicationInfo.Enabled))
		data = append(data, strconv.FormatBool(replicationInfo.Override))
		data = append(data, replicationInfo.Description)
		allData = append(allData, data)
	}
	r.TableInformation.SetTitles(r.replicationTitle).SetData(allData).Output()
}

// 根据replicationID 直接开始推送
func (r *Replication) StartExecution(id string) {
	_, err := r.Replication.StartReplication(context.Background(), &replication.StartReplicationParams{
		Execution: &models.StartReplicationExecution{
			PolicyID: parseStringToInt64(id),
		},
	})
	if err != nil {
		fmt.Println("开始复制失败! err:", err.Error())
		os.Exit(1)
	}
	fmt.Printf("复制规则id: %s 开始复制发送成功!!\n", id)

}

//从配置文件读取需要推送的项目 并开始遍历推送

type ReplicationRules struct {
	Name         string `json:"name"`
	Tag          string `json:"tag"`
	Resource     string `json:"resource"`
	DstNamespace string `json:"dst_namespace"`
}

func (re *Replication) StartExecutionFromConfig(id string, path string) {
	if _, isExist := exists(path); !isExist {
		fmt.Println("配置文件不存在 路径:", path)
		os.Exit(1)
	}
	rs := make([]*ReplicationRules, 0)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("read replication rules failed !! path:", path)
		panic(err)
	}

	err = json.Unmarshal(data, &rs)
	if err != nil {
		fmt.Println("json unmarshal replication rules failed !! path:", path)
		panic(err)
	}

	//构建推送策略 更新并推送
	for _, r := range rs {
		replcationInfo := re.GetReplicationInfo(id)
		//replcationInfo := GetReplicationInfo(id)
		filters := make([]*models.ReplicationFilter, 0)
		filters = append(filters, &models.ReplicationFilter{
			Type:  "name",
			Value: r.Name,
		})

		filters = append(filters, &models.ReplicationFilter{
			Type:  "tag",
			Value: r.Tag,
		})

		filters = append(filters, &models.ReplicationFilter{
			Type:  "resource",
			Value: r.Resource,
		})
		replcationInfo.Filters = filters
		replcationInfo.DestNamespace = r.DstNamespace

		result, err := re.Replication.UpdateReplicationPolicy(context.Background(), &replication.UpdateReplicationPolicyParams{
			ID:     parseStringToInt64(id),
			Policy: replcationInfo,
		})
		if err != nil {
			fmt.Println("更新复制规则失败! err:", err.Error())
			os.Exit(1)
		} else {
			fmt.Println(result.XRequestID, "更新成功!")
			re.StartExecution(id)
		}

	}

}

func (r *Replication) CreateReplication(id string, isPull bool) error {
	registryInstance := &Registry{
		harborClient: r.harborClient,
	}
	registry, err := registryInstance.SearchRegistryByID(id, false)
	if err != nil {
		return fmt.Errorf("查询仓库失败! err:%s\n", err.Error())
	}
	var speed int32 = 0
	var count int8 = 1
	policy := &models.ReplicationPolicy{
		DestRegistry: registry,
		Override:     true,
		Name:         fmt.Sprintf("%s_%s", registry.Name, "replication"),
		Enabled:      true,
		Description:  fmt.Sprintf("%s 仓库的复制规则", registry.Name),
		Speed:        &speed,
		Trigger: &models.ReplicationTrigger{
			Type: "manual",
		},
		DestNamespaceReplaceCount: &count,
		Filters: []*models.ReplicationFilter{
			{
				Type:  "resource",
				Value: "image",
			},
			{
				Type:  "name",
				Value: "kube_system/**",
			},
		},
	}
	//change pulled module
	if isPull {
		policy.DestRegistry = nil
		policy.SrcRegistry = registry
		policy.Name = fmt.Sprintf("%s_%s_%s", registry.Name, "pulled", "replication")
	}

	params := &replication.CreateReplicationPolicyParams{
		Policy: policy,
	}
	result, err := r.Replication.CreateReplicationPolicy(context.Background(), params)
	if err != nil {
		return fmt.Errorf("常见复制策略失败,err: %s\n", err)
	}

	fmt.Println(result.XRequestID, "创建成功!")
	return nil
}

func (r *Replication) SearchExecutionByID(id string) {
	idNumber := parseStringToInt64(id)
	sort := "sort=-start_time"
	results, err := r.Replication.ListReplicationExecutions(context.Background(), &replication.ListReplicationExecutionsParams{
		PolicyID: &idNumber,
		Sort:     &sort,
	})
	if err != nil {
		fmt.Println("查找执行器失败! 复制规则id:", id)
		os.Exit(1)
	}
	var (
		executionTitle = []string{
			"执行器ID",
			"执行器所所属规则ID",
			"开始时间",
			"结束时间",
			"执行结果",
		}
	)
	allData := make([][]string, 0)
	for _, execution := range results.Payload {
		data := make([]string, 0, 0)
		data = append(data, strconv.FormatInt(execution.ID, 10))
		data = append(data, strconv.FormatInt(execution.PolicyID, 10))
		data = append(data, time.Time(execution.StartTime).Format("2006/01/02 15:04:05"))
		data = append(data, time.Time(execution.EndTime).Format("2006/01/02 15:04:05"))
		data = append(data, execution.Status)
		allData = append(allData, data)
	}
	r.TableInformation.SetTitles(executionTitle).SetData(allData).Output()
}

func (r *Replication) SearchTasksByID(id string, size string) {
	idNumber := parseStringToInt64(id)
	sizeNumber := parseStringToInt64(size)
	sort := "sort=-start_time"
	results, err := r.Replication.ListReplicationTasks(context.Background(), &replication.ListReplicationTasksParams{
		ID:       idNumber,
		Sort:     &sort,
		PageSize: &sizeNumber,
	})
	if err != nil {
		fmt.Println("查找执行器对应的任务失败! 执行器id:", id)
		os.Exit(1)
	}
	var (
		executionTitle = []string{
			"任务ID",
			"任务所属的执行器ID",
			"源资源名称",
			"目标资源名称",
			"开始时间",
			"结束时间",
			"执行结果",
		}
	)
	allData := make([][]string, 0)
	for _, task := range results.Payload {
		data := make([]string, 0, 0)
		data = append(data, strconv.FormatInt(task.ID, 10))
		data = append(data, strconv.FormatInt(task.ExecutionID, 10))
		data = append(data, task.SrcResource)
		data = append(data, task.DstResource)
		data = append(data, time.Time(task.StartTime).Format("2006/01/02 15:04:05"))
		data = append(data, time.Time(task.EndTime).Format("2006/01/02 15:04:05"))
		data = append(data, task.Status)
		allData = append(allData, data)
	}
	r.TableInformation.SetTitles(executionTitle).SetData(allData).Output()
}

func parseStringToInt64(id string) int64 {
	idNumber, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		fmt.Println("转化id为整数类型失败! err:", err)
		os.Exit(1)
	}
	return idNumber
}
