package client

type HarborClientInterface interface {
	HarborRegistryInterface
	HarborReplicationInterface
}

type HarborRegistryInterface interface {
	CreateRegistry() error
	CreateRegistryByConfigInfo(string) error
}

type HarborReplicationInterface interface {
}

type HarborChartInterface interface {
	Download()
}

type HarborProjectInterface interface {
	SearchProjects(name string) error
	SearchProjectsList(name string) error
	CreateProject(name string) error
	DeleteProject(name string) error
}

type HarborImageInterface interface {
	SearchProjects(name string) error
}
