package env


import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"runtime"

	"git.intra.weibo.com/adx/logging"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/labstack/echo"
	echopprof "github.com/sevenNt/echo-pprof"
)

// 灰度机列表，以便决定是否建立监控指标
var GrayScale = []string{
	"10.85.101.122",
	"10.41.12.60",
	"10.75.26.159",
	"10.77.96.168",
}

// HTTPServer配置定义
type HTTPServerConfig struct {
	Disable bool   `toml:"disable" json:"disable"`
	Address string `toml:"address" json:"address"`
	Host    string `toml:"host"` // 本机IP地址
	HostIP  net.IP // 转换后的本机IP地址
	Gray    bool   // 是否为灰度机，根据本地IP对比GrayScale判断得出
}

func (c *HTTPServerConfig) Validate() error {
	if len(c.Address) == 0 {
		return fmt.Errorf("redis_server.address is invalid")
	}
	// 解析本机IP地址 如果为空则自动获取
	c.HostIP = net.ParseIP(c.Host)

	hostIPStr := c.HostIP.String()

	for _, gip := range GrayScale {
		if gip == hostIPStr {
			c.Gray = true
			break
		}
	}

	fmt.Println("local IP Addr:", c.HostIP.String())

	return nil
}

//motan-go conf
type MotanGo struct {
	Open      bool   `json:"open" toml:"open"`
	EngineRoom string `json:"engine_room" toml:"engine_room"`
	Path      string `json:"path" toml:"path"`
}

// 运行环境配置
type ServiceCoreEnvConfig struct {
	LogConfig         logging.LogConfig   `toml:"logging" json:"logging"`
	WebServerConfig   HTTPServerConfig    `toml:"http_server" json:"http_server"`

}

func (acec *ServiceCoreEnvConfig) Validate() error {
	if err := acec.LogConfig.Validate(); err != nil {
		return err
	}
	if err := acec.WebServerConfig.Validate(); err != nil {
		return err
	}

	return nil
}

type ServiceCoreEnv struct {
	Cfg       *ServiceCoreEnvConfig // 配置
	Logger    logging.Logger      // 日志
	WebServer *echo.Echo          // http web server
}



// 初始话运行环境
func NewServiceCoreEnv(cfg *ServiceCoreEnvConfig) (*ServiceCoreEnv, error) {
	//初始化logger
	logger := logging.NewLoggerWithConfig(&cfg.LogConfig)
	//初始化graphite

	//初始化http web server
	eco := echo.New()
	eco.Logger = logger
	eco.HideBanner = true



	//调试信息
	echopprof.Wrap(eco)

	return &ServiceCoreEnv{
		Cfg:       cfg,
		Logger:    logger,
		WebServer: eco,
	}, nil
}

// 开启运行环境
func (ace *ServiceCoreEnv) Run() error {
	go func() {
		ace.WebServer.Server.Addr = ace.Cfg.WebServerConfig.Address
		if err := gracehttp.Serve(ace.WebServer.Server); err != nil {
			if err.Error() == "http: Server closed" {
				ace.Logger.Warn("http server has closed")
			} else {
				ace.Logger.Errorf("fail to startup http web server, err: %s", err)
				panic(err)
			}
		}
	}()

	pid := fmt.Sprintf("%d\n", os.Getpid())

	pidFile := "/var/run/Service-core.pid"
	// 当操作系统为Mac OS，创建pid到/tmp目录下
	if runtime.GOOS == "darwin" {
		pidFile = "/tmp/Service-core.pid"
	}
	err := ioutil.WriteFile(pidFile, []byte(pid), 0644)
	if err != nil {
		panic(err)
	}
	ace.Logger.Info("http web server has been startup")
	return nil
}

// 关闭运行环境(释放相关资源)
func (ace *ServiceCoreEnv) Close() error {
	if err := ace.WebServer.Close(); err != nil {
		return fmt.Errorf("fail to close web server in Service core env, err: %s", err)
	}
	return nil
}
