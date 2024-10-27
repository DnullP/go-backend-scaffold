package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

var config struct {
	Consul ConsulConfig `yaml:"consul"`
	Jaeger JaegerConfig `yaml:"jaeger"`
}
var Consul ConsulConfig
var Jaeger JaegerConfig

func LoadConfig() {
	// 读取 YAML 文件
	data, err := os.ReadFile("/home/go-backend-scaffold/config/common-config.yaml")
	if err != nil {
		log.Fatalf("无法读取配置文件: %v", err)
	}

	// 定义配置变量

	// 解析 YAML 内容到结构体
	err = yaml.Unmarshal(data, &config)

	Consul = config.Consul
	Jaeger = config.Jaeger

	if err != nil {
		log.Fatalf("无法解析 YAML: %v", err)
	}
}
