package sonny

import (
	"fmt"
	"github.com/labstack/echo"
	"net/http"
	_ "sync"

	"service_core/env"
)

// 服务名
const (
	SONNY_SERVICE_NAME = "Sonny"
	SONNY_SERVICE_API  = "/v1/service/core/sonny"
)

// sonny服务定义
type SonnyService struct {
	name           string                     //服务名
	env            *env.ServiceCoreEnv          // 运行环境
}

func NewSonnyService(env *env.ServiceCoreEnv) (*SonnyService, error) {
	return &SonnyService{
		name: SONNY_SERVICE_NAME,
		env:  env,
	}, nil
}

// 初始化&启动
func (ss *SonnyService) Start() error {
	// 初始化依赖模块
	fmt.Printf("初始化依赖模块")

	// 注册API
	//true->false
	//ss.registerAPI()
	return nil
}

// 关闭&释放资源
func (ss *SonnyService) Stop() error {
	// 关闭依赖模块
	fmt.Printf("关闭依赖模块")
	return nil
}

// 重启服务
func (ss *SonnyService) Reload() error {
	// TODO 暂未实现
	return nil
}

// 模块自检
func (ss *SonnyService) Try() error {
	// TODO 暂未实现
	return nil
}

// 模块命名
func (ss *SonnyService) Name() string {
	return ss.name
}
// 处理HTTP请求
func (ss *SonnyService) processHttpPostRequest(c echo.Context) error {

	return c.JSONBlob(http.StatusOK, []byte("bt"))
}
// 注册API
func (ss *SonnyService) registerAPI() {
	// 注册sonny服务api
	ss.env.WebServer.POST(SONNY_SERVICE_API, ss.processHttpPostRequest)
}
