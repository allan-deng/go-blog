package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

var GlobalConfig *Config

//all configuration
type Config struct {
	LogConfig LogConfig `yaml:"log"`
	WebServer WebServer `yaml:"web"`
	WebPage   WebPage   `yaml:"webpage"`
	Mysql     Mysql     `yaml:"mysql"`
}

//log configuration
type LogConfig struct {
	Prefix     string `yaml:"prefix"`
	LogFile    bool   `yaml:"log-file"`
	Stdout     string `yaml:"stdout"`
	File       string `yaml:"file"`
	LogPath    string `yaml:"logfilepath"`
	ModuleName string `yaml:"modulename"`
}

type WebServer struct {
	Port int32 `yaml:"port"`
}

type WebPage struct {
	IcpInfo string `yaml:"icpinfo"`
}

type Mysql struct {
	Host   string `yaml:"host"`
	User   string `yaml:"user"`
	Pwd    string `yaml:"pwd"`
	Dbname string `yaml:"dbname"`
}

//Loading configuration files.If no configuration file location is specified, 'config.yaml' will be loaded by default.
func LoadConfig(path string) error {
	config := &Config{}
	if path == "" {
		path = "config.yaml"
	}
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err.Error())
	}
	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		fmt.Println(err.Error())
	}
	GlobalConfig = config
	return err
}

func init() {
	err := LoadConfig("./config/config.yaml")
	if err != nil {
		fmt.Errorf("Panic load config file failed: ")
		panic(err)

	}
	InitLogger()
}
