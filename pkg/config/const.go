package config

import (
	"github.com/bcsimple/harborctl/utils/file"
	"path/filepath"
)

const (
	configDir = ".harborctl"
	//configDir  = "config"
	configFile = "config.yaml"
)

var (
	pathDir = filepath.Join(file.HomeDir(), configDir)
	//pathDir  = filepath.Join("./", configDir)
	pathFile = filepath.Join(pathDir, configFile)
)
