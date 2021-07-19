package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/asim/go-micro/plugins/config/encoder/yaml/v3"
	"github.com/asim/go-micro/plugins/config/source/etcd/v3"
	logrusPlugin "github.com/asim/go-micro/plugins/logger/logrus/v3"
	"github.com/asim/go-micro/v3/config"
	"github.com/asim/go-micro/v3/config/source"
	"github.com/asim/go-micro/v3/config/source/file"
	"github.com/asim/go-micro/v3/config/source/memory"
	"github.com/asim/go-micro/v3/logger"
	"github.com/sirupsen/logrus"
	goYAML "gopkg.in/yaml.v2"
)

type ConfigDefine struct {
	Source  string `json:source`
	Prefix  string `json:prefix`
	Key     string `json:key`
	Address string
}

var configDefine ConfigDefine

var Schema ConfigSchema_

func setupEnvironment() {
	//registry plugin
	registryPlugin := os.Getenv("MSA_REGISTRY_PLUGIN")
	if "" == registryPlugin {
		registryPlugin = "etcd"
	}
	logger.Infof("MSA_REGISTRY_PLUGIN is %v", registryPlugin)
	os.Setenv("MICRO_REGISTRY", registryPlugin)

	//registry address
	registryAddress := os.Getenv("MSA_REGISTRY_ADDRESS")
	if "" == registryAddress {
		registryAddress = "localhost:2379"
	}
	logger.Infof("MSA_REGISTRY_ADDRESS is %v", registryAddress)
	os.Setenv("MICRO_REGISTRY_ADDRESS", registryAddress)

	//config
	envConfigDefine := os.Getenv("MSA_CONFIG_DEFINE")

	if "" == envConfigDefine {
		logger.Warn("MSA_CONFIG_DEFINE is empty")
		return
	}

	logger.Infof("MSA_CONFIG_DEFINE is %v", envConfigDefine)
	err := json.Unmarshal([]byte(envConfigDefine), &configDefine)
	if err != nil {
		logger.Error(err)
	}
	configDefine.Address = registryAddress
}

func mergeFile(_config config.Config) {
	filepath := configDefine.Prefix + configDefine.Key
	fileSource := file.NewSource(
		file.WithPath(filepath),
	)
	err := _config.Load(fileSource)
	if nil == err {
		logger.Infof("load config %v success", filepath)
		_config.Scan(&Schema)
	} else {
		panic(fmt.Sprintf("load config %v failed: %v", filepath, err))
	}
}

func mergeEtcd(_config config.Config) {
	etcdKey := configDefine.Prefix + configDefine.Key
	etcdSource := etcd.NewSource(
		etcd.WithAddress(configDefine.Address),
		etcd.WithPrefix(configDefine.Prefix),
		etcd.StripPrefix(true),
		source.WithEncoder(yaml.NewEncoder()),
	)
	err := _config.Load(etcdSource)
	if nil == err {
		logger.Infof("load config %v success", etcdKey)
	} else {
		panic(fmt.Sprintf("load config %v failed: %v", etcdKey, err))
	}
	_config.Get(configDefine.Key).Scan(&Schema)
}

func mergeDefault(_config config.Config) {
	memorySource := memory.NewSource(
		memory.WithYAML([]byte(defaultYAML)),
		source.WithEncoder(yaml.NewEncoder()),
	)
	err := _config.Load(memorySource)
	if nil == err {
		logger.Infof("load config default success")
	} else {
		panic(fmt.Sprintf("load config default failed: %v", err))
	}
	_config.Scan(&Schema)
}

func Setup() {
	mode := os.Getenv("MSA_MODE")
	if "" == mode {
		mode = "debug"
	}

	// initialize logger
	if "debug" == mode {
		logger.DefaultLogger = logrusPlugin.NewLogger(
			logger.WithOutput(os.Stdout),
			logger.WithLevel(logger.TraceLevel),
			logrusPlugin.WithTextTextFormatter(new(logrus.TextFormatter)),
		)
		logger.Info("-------------------------------------------------------------")
		logger.Info("- Micro Service Agent -> Setup")
		logger.Info("-------------------------------------------------------------")
		logger.Warn("Running in \"debug\" mode. Switch to \"release\" mode in production.")
		logger.Warn("- using env:	export MSA_MODE=release")
	} else {
		logger.DefaultLogger = logrusPlugin.NewLogger(
			logger.WithOutput(os.Stdout),
			logger.WithLevel(logger.TraceLevel),
			logrusPlugin.WithJSONFormatter(new(logrus.JSONFormatter)),
		)
		logger.Info("-------------------------------------------------------------")
		logger.Info("- Micro Service Agent -> Setup")
		logger.Info("-------------------------------------------------------------")
	}

	conf, err := config.NewConfig()
	if nil != err {
		panic(err)
	}

	setupEnvironment()

	// load default config
	logger.Tracef("default config is: \n\r%v", defaultYAML)

	// merge others
	if "file" == configDefine.Source {
		mergeFile(conf)
	} else if "etcd" == configDefine.Source {
		mergeEtcd(conf)
	} else {
		mergeDefault(conf)
	}

	ycd, err := goYAML.Marshal(&Schema)
	if nil != err {
		logger.Error(err)
	} else {
		logger.Tracef("current config is: \n\r%v", string(ycd))
	}

	level, err := logger.GetLevel(Schema.Logger.Level)
	if nil != err {
		logger.Warnf("the level %v is invalid, just use info level", Schema.Logger.Level)
		level = logger.InfoLevel
	}

	if "debug" == mode {
		logger.Warn("Using \"MSA_DEBUG_LOG_LEVEL\" to switch log's level in \"debug\" mode.")
		logger.Warn("- using env:	export MSA_DEBUG_LOG_LEVEL=debug")
		debugLogLevel := os.Getenv("MSA_DEBUG_LOG_LEVEL")
		if "" == debugLogLevel {
			debugLogLevel = "trace"
		}
		level, _ = logger.GetLevel(debugLogLevel)
	}
	logger.Infof("level is %v now", level)
	logger.Init(
		logger.WithLevel(level),
	)

}
