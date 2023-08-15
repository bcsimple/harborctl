package client

import (
	"encoding/json"
	"fmt"
	"github.com/bcsimple/harborctl/internal/harborctl/app/cmd/root"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Chart struct {
	*harborClient
	searchChartTitle []string
}

func NewChart(options *root.GlobalOptions) *Chart {
	return &Chart{harborClient: NewHarborClient(options),
		searchChartTitle: []string{
			"序号",
			"项目名",
			"Chart名",
			"Chart版本",
			"App版本",
			"类型",
		},
	}
}

// chart 相关操作
func (c *Chart) searchChart(name string) (*models.Search, error) {
	results := &models.Search{}
	if strings.Contains(name, ",") {
		names := strings.Split(name, ",")
		for _, name := range names {
			result, err := c.searchAll(name)
			if err != nil {
				return nil, err
			}
			results.Chart = append(results.Chart, result.Chart...)
		}
	} else {
		result, err := c.searchAll(name)
		if err != nil {
			return nil, err
		}
		results.Chart = append(results.Chart, result.Chart...)
	}
	return results, nil
}

func (c *Chart) SearchChart(name string) error {
	if res, err := c.searchChart(name); err == nil {
		rT := c.SearchTotalVersions(res)
		//return results,nil
		c.searchChartOuter(rT)
		return nil
	} else {
		return err
	}
}

func (c *Chart) searchChartOuter(results *models.Search) {
	deleteRepeat := make(map[string]int, 0)
	allData := make([][]string, 0)
	for index, chart := range results.Chart {
		data := make([]string, 0, 0)
		//去重复!
		if _, ok := deleteRepeat[chart.Name]; ok {
			continue
		} else {
			deleteRepeat[chart.Name] = 0
		}
		data = append(data, strconv.Itoa(index+1))
		data = append(data, strings.Split(chart.Name, "/")[0])
		data = append(data, *chart.Chart.Name)
		data = append(data, *chart.Chart.Version)
		data = append(data, *chart.Chart.AppVersion)
		data = append(data, "chart")
		allData = append(allData, data)
	}
	c.TableInformation.SetTitles(c.searchChartTitle).SetData(allData).Output()

}
func (c *Chart) getChartURL(projectName string, chartName string) string {
	c.URL.Path = fmt.Sprintf("/api/chartrepo/%s/charts/%s", projectName, chartName)
	return c.URL.String()
}

func (c *Chart) downloadChartURL(projectName string, chartName string) string {
	c.URL.Path = fmt.Sprintf("/chartrepo/%s/%s", projectName, chartName)
	return c.URL.String()
}

func (ch *Chart) DownloadChart(name string, Dpath string) {

	res, err := ch.searchChart(name)
	if err != nil {
		panic(err)
	}

	if len(res.Chart) == 0 {
		fmt.Println("chart 未找到")
		return
	}

	downloadUrls := []string{}
	for _, re := range res.Chart {
		charts := ch.getDataByName(re.Name)

		sT := strings.Split(name, "/")
		projectName := sT[0]
		for _, c := range charts {
			for _, u := range c.Urls {
				n := ch.downloadChartURL(projectName, u)
				downloadUrls = append(downloadUrls, n)
			}

		}
	}
	//if directory
	isDir, isExist := exists(Dpath)
	if !isExist {
		if err := os.MkdirAll(Dpath, 0755); err != nil {
			panic(err.Error())
		}
	}

	if !isDir {
		Dpath = filepath.Dir(Dpath)
	}

	for _, url := range downloadUrls {
		name := strings.Split(url, "/")
		//将下载到的文件放入指定的目录中
		filename := filepath.Join(Dpath, name[len(name)-1])
		if err := DownloadFile(url, filename); err != nil {
			fmt.Println("==下载失败 名字是:", filename)
			fmt.Println(err.Error())
			continue
		} else {
			fmt.Println("下载成功!! 名字是:", filename)
		}
	}
}

func exists(path string) (isDir bool, isExist bool) {
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

func DownloadFile(url string, filepath string) error {
	// Create the file
	out, err := os.OpenFile(filepath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func (c *Chart) getDataByName(name string) []*ChartInfo {
	sT := strings.Split(name, "/")
	projectName := sT[0]
	chartName := strings.Join(sT[1:], "/")
	data := request(c.getChartURL(projectName, chartName))
	charts := make([]*ChartInfo, 0)
	if err := json.Unmarshal([]byte(data), &charts); err != nil {
		panic(err)
	}
	return charts
}

func (ch *Chart) SearchTotalVersions(res *models.Search) *models.Search {
	for _, c := range res.Chart {
		charts := ch.getDataByName(c.Name)
		if len(charts) == 1 {
			continue
		}

		vS := make([]string, 0)
		aVS := make([]string, 0)

		for _, chart := range charts {
			vS = append(vS, chart.Version)
			aVS = append(aVS, chart.AppVersion)
		}
		versions := strings.Join(vS, ",")
		appVersions := strings.Join(aVS, ",")
		c.Chart.Version = &versions
		c.Chart.AppVersion = &appVersions
	}
	return res
}

func request(url string) string {
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	return string(data)
}

type ChartInfo struct {
	Name        string   `json:"name"`
	Home        string   `json:"home"`
	Sources     []string `json:"sources"`
	Version     string   `json:"version"`
	Description string   `json:"description"`
	Keywords    []string `json:"keywords"`
	Maintainers []struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"maintainers"`
	Icon        string `json:"icon"`
	APIVersion  string `json:"apiVersion"`
	AppVersion  string `json:"appVersion"`
	Annotations struct {
		ArtifacthubIoLinks    string `json:"artifacthub.io/links"`
		ArtifacthubIoOperator string `json:"artifacthub.io/operator"`
	} `json:"annotations"`
	KubeVersion  string `json:"kubeVersion"`
	Dependencies []struct {
		Name       string `json:"name"`
		Version    string `json:"version"`
		Repository string `json:"repository"`
		Condition  string `json:"condition"`
	} `json:"dependencies"`
	Type    string        `json:"type"`
	Urls    []string      `json:"urls"`
	Created time.Time     `json:"created"`
	Digest  string        `json:"digest"`
	Labels  []interface{} `json:"labels"`
}
