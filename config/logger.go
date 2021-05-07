package config

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/op/go-logging"
)

var (
	//A global logging tool that needs to be initialized with `InitLogger()` before use
	Logger           *logging.Logger
	defaultFormatter = `%{time:2006/01/02-15:04:05.000} %{shortfile} %{color:bold}[%{level:.4s}]%{color:reset} %{message}`
)

//Initialize the global logging tool
func InitLogger() {
	conf := GlobalConfig.LogConfig
	if conf.Prefix == "" {
		_ = fmt.Errorf("logger prefix not found")
	}
	if conf.ModuleName == "" {
		conf.ModuleName = "default"
	}
	logger := logging.MustGetLogger(conf.ModuleName)
	var backends []logging.Backend
	backends = registerStdout(conf, backends)
	backends = registerFile(conf, backends)

	logging.SetBackend(backends...)
	Logger = logger
}

func registerStdout(c LogConfig, backends []logging.Backend) []logging.Backend {
	if c.Stdout != "" {
		level, err := logging.LogLevel(c.Stdout)
		if err != nil {
			fmt.Println(err)
		}
		backends = append(backends, createBackend(os.Stdout, c, level))
	}

	return backends
}

func registerFile(c LogConfig, backends []logging.Backend) []logging.Backend {
	var logDir string
	if c.LogPath == "" {
		logDir = "log"
	} else {
		logDir = c.LogPath
	}
	if c.File != "" {
		if ok := pathExists(logDir); !ok {
			// directory not exist
			fmt.Println("create log directory")
			_ = os.Mkdir(logDir, os.ModePerm)
		}
		fileWriter, err := rotatelogs.New(
			logDir+string(os.PathSeparator)+c.ModuleName+" %Y-%m-%d-%H-%M.log",
			// maximum time to save log files
			rotatelogs.WithMaxAge(7*24*time.Hour),
			// time period of log file switching
			rotatelogs.WithRotationTime(24*time.Hour),
		)
		if err != nil {
			fmt.Println(err)
			return backends
		}
		level, err := logging.LogLevel(c.File)
		if err != nil {
			fmt.Println(err)
		}
		backends = append(backends, createBackend(fileWriter, c, level))
	}

	return backends
}

func createBackend(w io.Writer, c LogConfig, level logging.Level) logging.Backend {
	backend := logging.NewLogBackend(w, c.Prefix, 0)
	stdoutWriter := false
	if w == os.Stdout {
		stdoutWriter = true
	}
	format := getLogFormatter(c, stdoutWriter)
	backendLeveled := logging.AddModuleLevel(logging.NewBackendFormatter(backend, format))
	backendLeveled.SetLevel(level, c.ModuleName)
	return backendLeveled
}

func getLogFormatter(c LogConfig, stdoutWriter bool) logging.Formatter {
	pattern := defaultFormatter
	if !stdoutWriter {
		// Color is only required for console output
		// Other writers don't need %{color} tag
		pattern = strings.Replace(pattern, "%{color:bold}", "", -1)
		pattern = strings.Replace(pattern, "%{color:reset}", "", -1)
	}
	if !c.LogFile {
		// Remove %{logfile} tag
		pattern = strings.Replace(pattern, "%{longfile}", "", -1)
	}
	return logging.MustStringFormatter(pattern)
}

func pathExists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		return os.IsExist(err)
	}
	return true
}
