package client

import (
	"context"
	"fmt"
	"github.com/bcsimple/harborctl/internal/harborctl/app/cmd/root"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/project"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"os"
	"strconv"
	"strings"
)

type Project struct {
	*harborClient
	searchProjectTitle []string
}

func NewProject(options *root.GlobalOptions) *Project {
	return &Project{
		harborClient: NewHarborClient(options),
		searchProjectTitle: []string{
			"序号",
			"项目名",
			"镜像个数",
			"Chart个数",
			"类型",
		},
	}
}

func (p *Project) SearchProjects(name string) error {
	results := models.Search{}
	if strings.Contains(name, ",") {
		names := strings.Split(name, ",")
		for _, name := range names {
			result, err := p.searchAll(name)
			if err != nil {
				return err
			}
			results.Project = append(results.Project, result.Project...)
		}
	} else {
		result, err := p.searchAll(name)
		if err != nil {
			return err
		}
		results.Project = append(results.Project, result.Project...)
	}
	p.searchProjectOuter(&results)
	return nil
}

func (p *Project) CreateProject(name string) error {
	isPublic := true
	_, err := p.Project.CreateProject(context.Background(), &project.CreateProjectParams{
		Project: &models.ProjectReq{
			Public:      &isPublic,
			ProjectName: name,
		},
	})
	if err != nil {
		return err
	}
	fmt.Printf("project: [ %s ] create success\n", name)
	return nil
}

func (p *Project) DeleteProject(name string) error {
	_, err := p.Project.DeleteProject(context.Background(), &project.DeleteProjectParams{
		ProjectNameOrID: name,
	})
	if err != nil {
		return err
	}
	fmt.Printf("project: [ %s ] delete success\n", name)
	return nil
}

func (p *Project) SearchProjectsList() error {
	projects, err := p.Project.ListProjects(context.Background(), &project.ListProjectsParams{})
	if err != nil {
		return err
	}
	allData := make([][]string, 0)
	for index, projectInfo := range projects.GetPayload() {
		data := make([]string, 0, 0)
		data = append(data, strconv.Itoa(index+1))
		data = append(data, projectInfo.Name)
		data = append(data, strconv.FormatInt(projectInfo.RepoCount, 10))
		data = append(data, strconv.FormatInt(projectInfo.ChartCount, 10))
		data = append(data, "project")
		allData = append(allData, data)
	}
	p.TableInformation.SetTitles(p.searchProjectTitle).SetData(allData).Output()
	return nil
}

func (p *Project) searchProjectBy(id string) *models.Project {
	project1, err := p.Project.GetProject(context.Background(), &project.GetProjectParams{
		ProjectNameOrID: id,
	})
	if err != nil {
		fmt.Println("get project failed! err:", err)
		os.Exit(1)
	}
	return project1.GetPayload()
}

func (p *Project) searchProjectOuter(results *models.Search) {
	allData := make([][]string, 0)
	for index, project1 := range results.Project {
		project1 = p.searchProjectBy(fmt.Sprintf("%d", project1.ProjectID))
		data := make([]string, 0, 0)

		data = append(data, strconv.Itoa(index+1))
		data = append(data, project1.Name)
		data = append(data, strconv.FormatInt(project1.RepoCount, 10))
		data = append(data, strconv.FormatInt(project1.ChartCount, 10))
		data = append(data, "project")
		allData = append(allData, data)
	}
	p.TableInformation.SetTitles(p.searchProjectTitle).SetData(allData).Output()
}
