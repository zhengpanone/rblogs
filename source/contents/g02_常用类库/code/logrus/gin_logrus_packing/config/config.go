package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// LogConfig 结构体定义日志配置
type LogConfig struct {
	WebLogName string `yaml:"web_log_name"`
}

// Conf 结构体定义其他配置
type Conf struct {
	LogConfig LogConfig `yaml:"log_config"`
}

// Config 结构体表示整个配置文件
type Config struct {
	Conf Conf `yaml:"conf"`
}

// 全局配置变量
var Cfg *Config

// LoadConfig 函数从指定的 YAML 文件加载配置
func LoadConfig(filename string) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("无法读取配置文件: %v", err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("无法解析配置文件: %v", err)
	}

	Cfg = &config
}
