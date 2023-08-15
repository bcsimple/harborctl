package config

import (
	"github.com/bcsimple/harborctl/utils/file"
	"path/filepath"
)

const (
	configDir  = ".harborctl"
	configFile = "config.yaml"
)

var (
	pathDir  = filepath.Join(file.HomeDir(), configDir)
	pathFile = filepath.Join(pathDir, configFile)
)
