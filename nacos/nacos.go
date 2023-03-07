package nacos

import (
	"github.com/syklinux/golib/log"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

const (
	defaultCacheDir            = "/tmp/cache"
	defaultLogLevel            = "info"
	defaultLogDir              = "/tmp/"
	defaultContextPath         = "/nacos"
	defaultNotLoadCacheAtStart = true
	defaultPort                = 8848
	defaultTimeoutMs           = 5000
)

// Config nacos Config
type Config struct {
	Servers     string `json:"servers"`
	NameSpaceID string `json:"nameSpaceID"`
	Port        uint64 `json:"port"`
	ContextPath string `json:"contextPath"`
	LogLevel    string `json:"logLevel"`
}

type Nacos struct {
	ServerConfig []constant.ServerConfig
	ClientConfig *constant.ClientConfig
}

var ConfClient config_client.IConfigClient

// Init Init
func Init(conf Config) {
	na := newNacos()
	na.NewServerConfig(conf)
	na.NewClientConfig(conf)
	na.NewConfClient()
}

// newNacos newNacos
func newNacos() *Nacos {
	return new(Nacos)
}

// NewServerConfig 生成server配置
func (na *Nacos) NewServerConfig(conf Config) {
	var (
		port        uint64
		contextPath string
	)
	if conf.Port == 0 {
		port = defaultPort
	} else {
		port = conf.Port
	}

	if conf.ContextPath == "" {
		contextPath = defaultContextPath
	} else {
		contextPath = conf.ContextPath
	}

	sc := []constant.ServerConfig{
		*constant.NewServerConfig(conf.Servers, port, constant.WithContextPath(contextPath)),
	}
	na.ServerConfig = sc
}

// NewClientConfig NewClientConfig
func (na *Nacos) NewClientConfig(conf Config) {
	var (
		logLevel string
	)

	if conf.LogLevel == "" {
		logLevel = defaultLogLevel
	} else {
		logLevel = conf.LogLevel
	}

	cc := *constant.NewClientConfig(
		constant.WithNamespaceId(conf.NameSpaceID),
		constant.WithTimeoutMs(defaultTimeoutMs),
		constant.WithNotLoadCacheAtStart(defaultNotLoadCacheAtStart),
		constant.WithLogDir(defaultLogDir),
		constant.WithCacheDir(defaultCacheDir),
		constant.WithLogLevel(logLevel),
	)

	na.ClientConfig = &cc
}

// NewConfClient NewConfClient
func (na *Nacos) NewConfClient() {
	client, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  na.ClientConfig,
			ServerConfigs: na.ServerConfig,
		},
	)
	if err != nil {
		panic(err)
	}

	log.Info("new client success")
	ConfClient = client
}

// GetConfig 获取配置
func GetConfig(dataID string, group string) (string, error) {
	content, err := ConfClient.GetConfig(vo.ConfigParam{
		DataId: dataID,
		Group:  group})
	return content, err
}

// Close CloseClient
func Close() {
	ConfClient.CloseClient()
}
