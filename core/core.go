package core


import (
	"fmt"
	_ "io/ioutil"
	"os"
	"os/signal"
	_ "strings"
	"syscall"
	"service_core/config"
	"service_core/env"
	"service_core/sonny"
)

// 子服务接口
type ServiceCoreSubService interface {
	Start() error  // 启动服务
	Stop() error   // 关闭服务
	Reload() error // 重启服务
	Try() error    // 服务自检
	Name() string  // 模块名称
}

// service_core服务定义
type ServiceCore struct {
	cfgFile     string                         // 配置文件
	cfg         *config.ServiceCoreConfig        // 已加载配置
	env         *env.ServiceCoreEnv              // 运行环境
	subServices map[string]ServiceCoreSubService // 子服务列表
}


func NewServiceCore(cfgFile string) (*ServiceCore, error) {
	return &ServiceCore{
		cfgFile: cfgFile,
	}, nil
}

// 加载配置
func (ac *ServiceCore) loadConfig() error {
	config, err := config.NewServiceCoreConfig(ac.cfgFile)
	if err != nil {
		return fmt.Errorf("fail to load config file, err: %s", err)
	}
	ac.cfg = config
	// 初始化运行环境
	env, err := env.NewServiceCoreEnv(ac.cfg.Env)
	if err != nil {
		return fmt.Errorf("fail to init env, err: %s", err)
	}
	ac.env = env
	return nil
}


// 启动Service-core服务
func (ac *ServiceCore) Start() error {
	// 加载配置
	if err := ac.loadConfig(); err != nil {
		return err
	}
	// 注册所有子服务
	if err := ac.registerSubServices(); err != nil {
		return err
	}
	// 启动所有子服务
	if err := ac.startSubServices(); err != nil {
		return err
	}
	// 开启运行环境
	if err := ac.env.Run(); err != nil {
		return err
	}
	// 打印提示信息
	fmt.Println("Service-core has been startup")
	ac.env.Logger.Infof("Service-core has been startup")
	//监听信号, 服务优雅stop和reload
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)
	for {
		select {
		case sig, _ := <-sigCh:
			ac.env.Logger.Errorf("Recevie signal(%s)", sig)
			switch sig {
			case syscall.SIGTERM, syscall.SIGINT: // Stop
				ac.env.Logger.Errorf("Recevie signal(%s), and stoping...", sig)
				return ac.Stop()
			case syscall.SIGHUP: //Reload
				ac.env.Logger.Errorf("Recevie signal(%s), and reloading...", sig)
				ac.Reload()
			}
		}
	}
	return nil
}

// 关闭Service-core服务
func (ac *ServiceCore) Stop() error {
	// 关闭子服务
	if err := ac.stopSubServices(); err != nil {
		err = fmt.Errorf("fail to stop Service-core, err: %s", err)
		ac.env.Logger.Errorf(err.Error())
		return err
	}
	// 关闭运行环境
	if err := ac.env.Close(); err != nil {
		err = fmt.Errorf("fail to stop Service-core, err: %s", err)
		return err
	}
	// 打印提示信息
	ac.env.Logger.Infof("Service-core has been stopped")
	return nil
}

// 重启Service-core服务
func (ac *ServiceCore) Reload() error {
	// TODO 暂未实现
	return nil
}

// 测试Service-core服务
func (ac *ServiceCore) Try() error {
	if err := ac.loadConfig(); err != nil {
		fmt.Println(err.Error())
	}
	return nil
}

// 注册所有自服务
func (ac *ServiceCore) registerSubServices() error {
	subServices := make(map[string]ServiceCoreSubService)
	// 注册服务
	sonnyService, err := sonny.NewSonnyService(ac.env)
	if err != nil {
		return fmt.Errorf("fail to init sonny service, err: %s", err)
	}
	subServices[sonnyService.Name()] = sonnyService

	ac.subServices = subServices
	return nil
}

// 开启所有子服务
func (ac *ServiceCore) startSubServices() error {
	for sname, svc := range ac.subServices {
		if err := svc.Start(); err != nil {
			return fmt.Errorf("fail to start sub service: %s, err: %s", sname, err)
		}
	}
	return nil
}

// 关闭所有子服务
func (ac *ServiceCore) stopSubServices() error {
	for sname, svc := range ac.subServices {
		if err := svc.Stop(); err != nil {
			return fmt.Errorf("fail to stop sub service: %s, err: %s", sname, err)
		}
	}
	return nil
}

// 重启所有子服务
func (ac *ServiceCore) reloadSubServices() error {
	// TODO 暂未实现
	return nil
}
