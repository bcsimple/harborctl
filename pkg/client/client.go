package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/bcsimple/harborctl/internal/harborctl/app/cmd/root"
	"github.com/bcsimple/harborctl/pkg/config"
	"github.com/bcsimple/harborctl/utils/table"
	"github.com/goharbor/go-client/pkg/harbor"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/artifact"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/search"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

var (
	InsecureTransport http.RoundTripper = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second, // nolint:gomnd
			KeepAlive: 30 * time.Second, // nolint:gomnd
			DualStack: true,
		}).DialContext,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // nolint:gosec
		},
		MaxIdleConns:          100,              // nolint:gomnd
		IdleConnTimeout:       90 * time.Second, // nolint:gomnd
		TLSHandshakeTimeout:   10 * time.Second, // nolint:gomnd
		ExpectContinueTimeout: 1 * time.Second,
	}
)

type harborClient struct {
	//harbor client
	*client.HarborAPI
	//print format
	table.TableInformation
	//harborctl connection config info
	*config.HarborConnectInfo
	//url拼接
	*url.URL
	//http client some api must use origin http request
	HttpClient *http.Client
}

func NewHarborClient(options *root.GlobalOptions) *harborClient {

	clientConfig := config.NewConnectConfiguration()
	info := clientConfig.GetDefaultConnectInfo()
	if options.Context != "" {
		var err error
		if info, err = clientConfig.GetConnectInfo(options.Context); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	Hclient := &harborClient{
		HarborConnectInfo: info,
	}

	Hclient.HarborConnectInfo.Scheme = "http"

	Hclient.URL = &url.URL{
		Scheme: Hclient.HarborConnectInfo.Scheme,
		User:   url.UserPassword(Hclient.HarborConnectInfo.User, Hclient.HarborConnectInfo.Password),
		Host:   Hclient.HarborConnectInfo.Host,
	}
	_ = harbor.Config{
		URL:       Hclient.URL,
		Transport: InsecureTransport,
	}

	clientSet, err := harbor.NewClientSet(&harbor.ClientSetConfig{
		URL:      fmt.Sprintf("%s://%s", Hclient.HarborConnectInfo.Scheme, Hclient.HarborConnectInfo.Host),
		Insecure: true,
		Username: Hclient.HarborConnectInfo.User,
		Password: Hclient.HarborConnectInfo.Password,
	})

	Hclient.HarborAPI = clientSet.V2()
	if err != nil {
		fmt.Printf("init client err:%s\n", err.Error())
		panic(err)
	}

	Hclient.TableInformation = table.NewTableInformation(options.FormatStyle)
	Hclient.HttpClient = &http.Client{}

	return Hclient
}

// 全局search方法
func (h *harborClient) searchAll(name string) (*models.Search, error) {
	searchResults, err := h.Search.Search(context.Background(), &search.SearchParams{
		Q: name,
	})
	if err != nil {
		return nil, err
	}

	return searchResults.GetPayload(), nil
}

// scan use
func (h *harborClient) GetImageTagsByImageName(name string) map[string]string {
	searchResult, err := h.searchAll(name)
	if err != nil {
		panic(err)
	}
	if len(searchResult.Repository) == 0 {
		//return "无"字样
		return map[string]string{"无": "无"}
	}
	imageMap := make(map[string]string)
	for _, repository := range searchResult.Repository {
		imageName, imageTag := h.searchTag(repository.ProjectName, repository.RepositoryName)
		if filepath.Base(imageName) == name {
			imageMap[imageName] = imageTag
		}

		if imageName == name {
			imageMap[imageName] = imageTag
		}
	}
	return imageMap
}

// scan use get tag
func (h *harborClient) searchTag(projectName, repositoryName string) (imageName string, imageTags string) {

	repositoryNames := strings.TrimPrefix(repositoryName, fmt.Sprintf("%s/", projectName))

	if strings.Contains(repositoryNames, "/") {
		repositoryNames = url.PathEscape(repositoryNames)
	}
	res, err := h.Artifact.ListArtifacts(context.Background(), &artifact.ListArtifactsParams{
		ProjectName:    projectName,
		RepositoryName: repositoryNames,
	})
	if err != nil {
		return "", ""
	}
	var tags []string
	for _, result := range res.Payload {
		for _, tag := range result.Tags {
			tags = append(tags, tag.Name)
		}
	}
	// tag 进行排序
	sort.Slice(tags, func(i, j int) bool {
		return tags[i] > tags[j]
	})
	return repositoryName, strings.Join(tags, ",")
}
