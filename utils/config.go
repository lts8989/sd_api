package utils

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

var Conf *Config

type Config struct {
	DbCfg     DatabaseCfg `yaml:"database"`
	ServCfg   ServiceCfg  `yaml:"service"`
	LogConfig LogCfg      `yaml:"log"`
	SdServCfg SdService   `yaml:"sd_service"`
}

type DatabaseCfg struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Dbname   string `yaml:"dbname"`
}

type ServiceCfg struct {
	Port int `yaml:"port"`
}

type LogCfg struct {
	Path         string `yaml:"path"`
	FileName     string `yaml:"file_name"`
	MaxAge       int    `yaml:"max_age"`
	RotationTime int32  `yaml:"rotation_time"`
}
type SdService struct {
	Domain   string `yaml:"domain"`
	ClientId string `yaml:"client_id"`
}

func InitConfig() error {
	yamlFile, err := os.ReadFile("conf/config.yaml")
	if err != nil {
		log.Fatalf("Failed to read YAML file: %v", err)
		return err
	}
	// 解析YAML数据
	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatalf("Failed to parse YAML data: %v", err)
		return err
	}
	Conf = &config
	return nil
}
