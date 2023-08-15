package client

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/bcsimple/harborctl/internal/harborctl/app/cmd/root"
	"github.com/bcsimple/harborctl/pkg/client/scan"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

type ScanImage struct {
	*harborClient
	arch                  string
	FilePathAndImagesName map[string][]*scan.Image `json:"file_path_and_images_name"`
	FilePathAndImages     map[string][]*scan.Image `json:"file_path_and_images"`
	WithFile              bool
	WithCompare           bool
	WithCompareOnlyFalse  bool
	path                  string
	release               string
	titleFile             []string
	titleCompare          []string
}

func NewScanImage(options *root.GlobalOptions, withFile, withCompare, withCompareOnlyFalse bool, path, release string) *ScanImage {

	scanInstance := &ScanImage{
		harborClient:          NewHarborClient(options),
		FilePathAndImagesName: make(map[string][]*scan.Image, 0),
		FilePathAndImages:     make(map[string][]*scan.Image, 0),
		WithFile:              withFile,
		WithCompare:           withCompare,
		WithCompareOnlyFalse:  withCompareOnlyFalse,
		path:                  path,
		release:               release,
		titleFile:             titleFile,
		titleCompare:          titleCompare,
	}

	scanInstance.arch = "x86_64"
	if runtime.GOARCH != "amd64" {
		scanInstance.arch = "aarch64"
	}
	if path != "" {
		scanInstance.findImagesFromFile()
	}
	return scanInstance
}

func (s *ScanImage) findImagesFromFile() {
	s.FilePathAndImagesName[s.path] = make([]*scan.Image, 0)
	file, err := os.Open(s.path)
	if err != nil {
		panic(err)
	}

	var buffer = bufio.NewReader(file)

	tmpDataMap := make(map[string]interface{}, 0)
	counter := 1
	for {
		line, _, err := buffer.ReadLine()
		counter++
		if err != nil || err == io.EOF {
			break
		}
		//镜像文件中 若是以#号开头的全部忽略掉
		if bytes.HasPrefix(bytes.TrimSpace(line), []byte("#")) || bytes.Contains(bytes.TrimSpace(line), []byte("%")) {
			continue
		}
		//镜像文件中 要是包含了#号 则直接删除掉 取切片第一个
		if bytes.Contains(line, []byte("#")) {
			line = bytes.Split(line, []byte("#"))[0]
		}
		//切分冒号的字符串
		slice := bytes.Split(line, []byte(": "))
		if len(slice) != 2 {
			continue
		}
		//收集所有的镜像k v的值 若map中k相同 则使用较长的v的值做最后的使用镜像
		if v, ok := tmpDataMap[string(slice[0])]; ok {
			if len(string(slice[1])) > len(v.(string)) {
				tmpDataMap[string(slice[0])] = string(slice[1])
			}
		} else {
			tmpDataMap[string(slice[0])] = string(slice[1])
		}

	}
	imagesMap := make([]*scan.Image, 0)

	for k, v := range tmpDataMap {
		if strings.HasPrefix(strings.ToLower(k), "image_") {
			value := v.(string)
			//替换掉镜像文件中的paas_release
			value = strings.ReplaceAll(value, "{{ paas_release }}", s.release)
			if strings.Contains(value, "}}/") {
				prefixLength := strings.Split(value, "}}/")
				if len(prefixLength) != 2 {
					continue
				}
				imageProjectName := prefixLength[1]

				if strings.Contains(imageProjectName, ":") {
					valueSlice := strings.Split(imageProjectName, ":")
					//处理新版本的image yaml数据
					imageName := valueSlice[0]
					if strings.Contains(imageName, "/") {
						imageName = filepath.Base(imageName)
					}
					//去掉分号双引号
					tag := strings.TrimRight(valueSlice[1], "'\"")
					imagesMap = append(imagesMap, scan.NewImage(imageName, valueSlice[0], strings.ReplaceAll(tag, "\"", "")))
					continue
				}
			}
		}
	}
	s.FilePathAndImagesName[s.path] = imagesMap
	s.WithImagesFromHarbor()
}

func (s *ScanImage) WithImagesFromHarbor() {

	// 如果参数 -d / -c 都为true 则直接退出  不搜索harbor
	// 不管怎么处理都要合并所有数据到统一的字段 FilePathAndImages 中

	if s.WithCompare || s.WithCompareOnlyFalse {
		// 获取harbor 客户端
		s.FilePathAndImagesName = s.withHarbor(s.FilePathAndImagesName)
	}

	for path, image := range s.FilePathAndImagesName {
		s.FilePathAndImages[path] = image
	}

}

func (s *ScanImage) withHarbor(m map[string][]*scan.Image) map[string][]*scan.Image {
	if len(m) == 0 {
		return m
	}
	for path, value := range m {

		if len(value) == 0 {
			continue
		}
		m[path] = make([]*scan.Image, 0)
		for _, image := range value {
			if image.ImageName == "" {
				image.ImageNameInHarbor = map[string]string{"空": "空"}
				continue
			}
			harborImage := s.GetImageTagsByImageName(image.ImageNameWithProjectName)

			if harborImage["无"] == "无" {
				harborImage = s.GetImageTagsByImageName(image.ImageName)
			}
			//搜索出来的镜像如果小于等于1 直接继续
			if len(harborImage) <= 1 {
				image.ImageNameInHarbor = harborImage
				continue
			}

			if ok, imageName, imageTag := s.FilterHarborImageTagEqual(image, harborImage); ok {
				image.ImageNameInHarbor = map[string]string{
					imageName: imageTag,
				}
				continue
			}

			if ok, imageName, imageTag := s.FilterHarborImageTagContains(image, harborImage); ok {
				image.ImageNameInHarbor = map[string]string{
					imageName: imageTag,
				}
				continue
			}

			if ok, imageName, imageTag := s.FilterHarborImageTagNoContains(image, harborImage); ok {
				image.ImageNameInHarbor = map[string]string{
					imageName: imageTag,
				}
				continue
			}
			if ok, imageName, imageTag := s.FilterHarborImageTagNoEqual(image, harborImage); ok {
				image.ImageNameInHarbor = map[string]string{
					imageName: imageTag,
				}
			}
		}
		m[path] = value
	}
	return m
}

//  是否tag等于
//  是否tag包含
//  是否tag不包含
//  是否tag不等于

func (s *ScanImage) FilterHarborImageTagEqual(image *scan.Image, harborImage map[string]string) (ok bool, imageName string, imageTag string) {
	for imageName, imageTag = range harborImage {
		if strings.Contains(imageName, image.ImageName) && imageTag == image.ImageTag {
			return true, imageName, imageTag
		}
	}
	return false, imageName, imageTag
}

func (s *ScanImage) FilterHarborImageTagNoEqual(image *scan.Image, harborImage map[string]string) (ok bool, imageName string, imageTag string) {
	for imageName, imageTag = range harborImage {
		if strings.Contains(imageName, image.ImageName) && imageTag != image.ImageTag {
			return true, imageName, imageTag
		}
	}
	return false, imageName, imageTag
}

func (s *ScanImage) FilterHarborImageTagContains(image *scan.Image, harborImage map[string]string) (ok bool, imageName string, imageTag string) {
	for imageName, imageTag = range harborImage {
		if strings.Contains(imageName, image.ImageName) && strings.Contains(imageTag, image.ImageTag) {
			return true, imageName, imageTag
		}
	}
	return false, imageName, imageTag
}

func (s *ScanImage) FilterHarborImageTagNoContains(image *scan.Image, harborImage map[string]string) (ok bool, imageName string, imageTag string) {
	for imageName, imageTag = range harborImage {
		if strings.Contains(imageName, image.ImageName) && !strings.Contains(imageTag, image.ImageTag) {
			return true, imageName, imageTag
		}
	}
	return false, imageName, imageTag
}

func (s *ScanImage) ReadDataFromFile(path string) []byte {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("读取文件失败! 文件名称:", path)
		os.Exit(1)
	}
	return data
}

func (s *ScanImage) PrintFile() {
	datas := s.getImagesFromFile()
	s.TableInformation.SetTitles(titleFile).SetData(datas).Output()
}

func (s *ScanImage) getImagesFromFile() [][]string {
	counter := 1
	datas := make([][]string, 0)

	for path, images := range s.FilePathAndImages {
		for _, image := range images {
			datas = append(datas, []string{strconv.Itoa(counter), image.ImageNameWithProjectName, image.ImageTag, path})
			counter++
		}
	}
	return datas
}

func (s *ScanImage) PrintCompare() *ScanImage {
	s.print(false)
	return s
}

func (s *ScanImage) PrintDiff() *ScanImage {
	s.print(true)
	return s
}

func (s *ScanImage) print(isDiff bool) {
	counter := 1
	datas := make([][]string, 0)
	for _, images := range s.FilePathAndImages {
		for _, image := range images {
			isTrue := "false"
			for harborImageName, harborImageTag := range image.ImageNameInHarbor {
				//给与临时变量harbor Tag
				hTag := strings.TrimSpace(harborImageTag)
				if hTag == "" {
					hTag = "latest"
				}
				//给与临时变量file Tag
				iTag := strings.TrimSpace(image.ImageTag)
				if iTag == "" {
					iTag = "latest"
				}

				if strings.TrimSpace(image.ImageNameWithProjectName) == strings.TrimSpace(harborImageName) {
					if strings.Contains(hTag, ",") {
						for _, tag := range strings.Split(hTag, ",") {
							if strings.TrimSpace(tag) == iTag {
								isTrue = "True"
							}
						}
					} else {
						if hTag == iTag {
							isTrue = "True"
						}
					}
				}

				if isDiff {
					if isTrue == "false" {
						datas = append(datas, []string{strconv.Itoa(counter), image.ImageNameWithProjectName, image.ImageTag, harborImageName, harborImageTag, isTrue})
						counter++
					}
				} else {
					datas = append(datas, []string{strconv.Itoa(counter), image.ImageNameWithProjectName, image.ImageTag, harborImageName, harborImageTag, isTrue})
					counter++
				}
			}
		}
	}
	s.TableInformation.SetTitles(titleCompare).SetData(datas).Output()
}

var (
	titleFile = []string{
		"序列",
		"镜像名",
		"文件中的镜像tag",
		"文件路径",
	}

	titleCompare = []string{
		"序列",
		"镜像名称",
		"文件中的镜像tag",
		"harbor中镜像名称",
		"harbor中镜像tag",
		"是否一致(true:一致,false:不一致)",
	}
)
