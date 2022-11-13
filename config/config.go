package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

var GlobalConfig *Config
var GlobalMassage *Massage

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
	Port        int32  `yaml:"port"`
	MetricsPort int32  `yaml:"metrics_port"`
	UploadPath  string `yaml:"uploadpath"`
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

type Massage struct {
	Nav           Nav    `yaml:"nav"`
	Index         Index  `yaml:"index"`
	Blog          Blog   `yaml:"blog"`
	About         About  `yaml:"about"`
	CommentAvatar string `yaml:"commentAvatar"`
}

type Nav struct {
	Title string `yaml:"title"`
}

type Index struct {
	Email       string `yaml:"email"`
	Githubname  string `yaml:"githubName"`
	Githublink  string `yaml:"githubLink"`
	Copyright   string `yaml:"copyright"`
	Num         string `yaml:"num"`
	FooterTitle string `yaml:"footerTitle"`
	Profile     string `yaml:"profile"`
}

type Blog struct {
	Copyright   string `yaml:"copyright"`
	DefaultPath string `yaml:"defaultPath"`
	Wechat      string `yaml:"wechat"`
	Alipay      string `yaml:"alipay"`
}

type About struct {
	Name    string   `yaml:"name"`
	Profile string   `yaml:"profile"`
	Avatar  string   `yaml:"avatar"`
	ParaOne string   `yaml:"paraOne"`
	ParaTwo string   `yaml:"paraTwo"`
	Stack   []string `yaml:"stack"`
	Hobby   []string `yaml:"hobby"`
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

func LoadMassage(path string) error {
	massage := &Massage{}
	if path == "" {
		path = "messages.yaml"
	}
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err.Error())
	}
	err = yaml.Unmarshal(yamlFile, massage)
	if err != nil {
		fmt.Println(err.Error())
	}
	GlobalMassage = massage
	return err
}

func init() {
	err := LoadConfig("./config/config.yaml")
	if err != nil {
		fmt.Println("panic load config file failed: ")
		panic(err)
	}
	err = LoadMassage("./config/messages.yaml")
	if err != nil {
		fmt.Println("panic load messages file failed: ")
		panic(err)
	}
	InitLogger()
}
