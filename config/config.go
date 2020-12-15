package config

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"sync"
)

type Config struct {
	Domain struct {
		DomainId string `yaml:"domainId"`
		Ak       string `yaml:"ak"`
		Sk       string `yaml:"sk"`
	} `yaml:"domain"`
	Rms struct {
		Endpoint      string   `yaml:"endpoint"`
		ResourceTypes []string `yaml:"resourceTypes"`
		Limit         int32    `yaml:"limit"`
	} `yaml:"rms"`
	Smn struct {
		Items []struct {
			RegionId string `yaml:"regionId"`
			Endpoint string `yaml:"endpoint"`
		} `yaml:"items"`
	} `yaml:"smn"`
	Mysql struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Network  string `yaml:"network"`
		Server   string `yaml:"server"`
		Port     int32  `yaml:"port"`
		Database string `yaml:"database"`
	}
	Sync struct {
		Spec string `yaml:"spec"`
	} `yaml:"sync"`
}

var singletonConfig *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(loadConfig)
	return singletonConfig
}

func loadConfig() {
	f, err := os.Open("conf/properties.yml")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	singletonConfig = &cfg
}
