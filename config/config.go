package config

import (
	"os"

	"github.com/jinzhu/configor"
	"github.com/mrlightwood/golang-products-api/helpers"
)

type Config struct {
	ConfigFile string
	LogLevel   uint32 `default:"4"`
	Api        struct {
		HttpPort int  `default:"8080"`
		Logging  bool `default:"false"`
	}
	Store struct {
		Dbpath string `required:"true"`
	}
}

func NewConfig(configFile string) (*Config, error) {
	config := &Config{ConfigFile: configFile}
	if err := configor.Load(config, configFile); err != nil {
		return nil, err
	}
	// Create database file if not exists
	config.Store.Dbpath = helpers.RootDir() + config.Store.Dbpath
	if _, err := os.Stat(config.Store.Dbpath); os.IsNotExist(err) {
		os.Create(config.Store.Dbpath)
	}
	return config, nil
}
