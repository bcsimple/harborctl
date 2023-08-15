package client

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/bcsimple/harborctl/internal/harborctl/app/cmd/root"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/artifact"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type Image struct {
	*harborClient
	imageTitle []string
}

func NewImage(options *root.GlobalOptions) *Image {
	return &Image{
		harborClient: NewHarborClient(options),
		imageTitle: []string{
			"序号",
			"项目名",
			"镜像名",
			"Tag名",
			"Tag数",
			"类型",
		},
	}
}

func (i *Image) SearchAll(name string) error {
	results := models.Search{}
	if strings.Contains(name, ",") {
		names := strings.Split(name, ",")
		for _, name := range names {
			result, err := i.searchAll(name)
			if err != nil {
				return err
			}
			results.Repository = append(results.Repository, result.Repository...)
		}
	} else {
		result, err := i.searchAll(name)
		if err != nil {
			return err
		}
		results.Repository = append(results.Repository, result.Repository...)
	}
	i.searchOuter(&results)
	return nil
}

func (i *Image) searchOuter(results *models.Search) {
	deleteRepeat := make(map[string]int, 0)
	allData := make([][]string, 0)
	for index, repository := range results.Repository {
		data := make([]string, 0, 0)
		//去重复!
		if _, ok := deleteRepeat[repository.RepositoryName]; ok {
			continue
		} else {
			deleteRepeat[repository.RepositoryName] = 0
		}
		data = append(data, strconv.Itoa(index+1))
		data = append(data, repository.ProjectName)
		data = append(data, repository.RepositoryName)
		data = append(data, i.searchTag(repository.ProjectName, repository.RepositoryName))
		data = append(data, strconv.FormatInt(repository.ArtifactCount, 10))
		data = append(data, "image")
		allData = append(allData, data)
	}
	i.TableInformation.SetTitles(i.imageTitle).SetData(allData).Output()
}

func (i *Image) searchTag(projectName, repositoryName string) string {

	repositoryNames := strings.TrimPrefix(repositoryName, fmt.Sprintf("%s/", projectName))

	if strings.Contains(repositoryNames, "/") {
		repositoryNames = url.PathEscape(repositoryNames)
	}
	res, err := i.Artifact.ListArtifacts(context.Background(), &artifact.ListArtifactsParams{
		ProjectName:    projectName,
		RepositoryName: repositoryNames,
	})

	if err != nil {
		return ""
	}
	var tags []string
	for _, result := range res.Payload {
		for _, tag := range result.Tags {
			tags = append(tags, tag.Name)
		}
	}
	return strings.Join(tags, ",")
}

func Test() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	iSum, err := cli.ImageList(ctx, types.ImageListOptions{
		All: true,
	})
	if err != nil {
		panic(err)
	}
	for _, i := range iSum {
		fmt.Printf("%#v\n", i)
	}
}

func (i *Image) DownloadImage(name string, Dpath string) {
	res, err := i.searchAll(name)
	if err != nil {
		panic(err)
	}

	if len(res.Repository) == 0 {
		fmt.Println("image 未找到")
		return
	}

	//拼接所有的tag
	imageRef := make(map[string]string, 0)
	for _, re := range res.Repository {

		tag := i.searchTag(re.ProjectName, re.RepositoryName)
		//未获取tag直接继续
		if tag == "" {
			continue
		}
		nameImage := re.RepositoryName
		if strings.Contains(re.RepositoryName, "/") {
			nameTmp := strings.Split(re.RepositoryName, "/")
			nameImage = nameTmp[len(nameTmp)-1]
		}

		if strings.Contains(tag, ",") {
			tags := strings.Split(tag, ",")
			for _, t := range tags {
				imageName := fmt.Sprintf("%s/%s:%s", i.HarborConnectInfo.Host, re.RepositoryName, t)
				key := fmt.Sprintf("%s-%s.tar.gz", nameImage, t)
				imageRef[key] = imageName
			}
		} else {
			imageName := fmt.Sprintf("%s/%s:%s", i.HarborConnectInfo.Host, re.RepositoryName, tag)
			key := fmt.Sprintf("%s-%s.tar.gz", nameImage, tag)
			imageRef[key] = imageName
		}
	}

	//if directory
	isDir, isExist := exists(Dpath)
	if !isExist {
		if err := os.MkdirAll(Dpath, 0755); err != nil {
			panic(err.Error())
		}
	} else {
		if !isDir {
			Dpath = filepath.Dir(Dpath)
		}
	}
	i.imagePull(imageRef, Dpath)
}

func (i *Image) imagePull(names map[string]string, Dpath string) {

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Println("get docker client failed!")
		panic(err)
	}

	authConfig := types.AuthConfig{
		Username:      i.HarborConnectInfo.User,
		Password:      i.HarborConnectInfo.Password,
		ServerAddress: i.HarborConnectInfo.Host,
	}

	//if _, err := cli.RegistryLogin(context.Background(), authConfig); err != nil {
	//	fmt.Printf("login server:%s failed\n", host)
	//	panic(err)
	//}

	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		panic(err)
	}
	authStr := base64.URLEncoding.EncodeToString(encodedJSON)
	for fileName, imageName := range names {
		out, err := cli.ImagePull(context.Background(), imageName, types.ImagePullOptions{RegistryAuth: authStr})
		if err != nil {
			fmt.Println(err.Error())
			panic(err)
		}
		defer out.Close()
		_, err = io.Copy(os.Stdout, out)
		if err != nil {
			panic(err)
		}
		//imageSha256 := getDigest(string(data))

		save(cli, imageName, filepath.Join(Dpath, fileName))
	}
}

func list(cli *client.Client) {
	iSum, err := cli.ImageList(context.Background(), types.ImageListOptions{
		All: true,
	})
	if err != nil {
		panic(err)
	}
	for _, info := range iSum {
		fmt.Println(info.RepoTags)
		fmt.Println(info.ID)
	}
}

func save(cli *client.Client, imageName string, filename string) {
	fmt.Println(imageName, filename)
	ctx := context.Background()
	fmt.Println(imageName)
	reader, err := cli.ImageSave(ctx, []string{
		imageName,
	})
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	_ = ioutil.WriteFile(filename, data, 0644)
	fmt.Printf("下载成功:%s! 镜像名:%s\n", filename, imageName)
}

func getDigest(name string) string {
	reg := regexp.MustCompile("(?P<name>sha256:.*)\"")
	match := reg.FindStringSubmatch(name)
	groupNames := reg.SubexpNames()
	result := make(map[string]string)

	// 转换为map
	for i, name := range groupNames {
		if i != 0 && name != "" { // 第一个分组为空（也就是整个匹配）
			result[name] = match[i]
		}
	}
	return result["name"]
}
