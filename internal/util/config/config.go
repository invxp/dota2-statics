package config

import (
	"github.com/BurntSushi/toml"
	"log"
	"path/filepath"
	"time"
)

type Config struct {
	Server struct {
		Address string
		APIKey  string
	}
	Log struct {
		Path            string
		MaxAge          time.Duration
		MaxRotationSize int64
	}
}

func Load(currentPath, fileName string) *Config {
	if !filepath.IsAbs(fileName) {
		fileName = filepath.Join(currentPath, fileName)
	}

	config := &Config{}
	_, err := toml.DecodeFile(fileName, config)
	if err != nil {
		log.Panic(err)
	}
	return config
}
