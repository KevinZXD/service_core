package config

import (
	"encoding/json"
	"fmt"

	"github.com/BurntSushi/toml"

	"service_core/env"
)


type ServiceCoreConfig struct {
	Env          *env.ServiceCoreEnvConfig `toml:"env" json:"env"`     // 基础环境配置
}

// 配置加载
func NewServiceCoreConfig(fname string) (*ServiceCoreConfig, error) {
	config := &ServiceCoreConfig{}
	if _, err := toml.DecodeFile(fname, config); err != nil {
		return nil, err
	}
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}

// 校验配置
func (acc *ServiceCoreConfig) Validate() error {
	if err := acc.Env.Validate(); err != nil {
		return err
	}
	return nil
}

// 配置打印输出
func (acc *ServiceCoreConfig) String() string {
	b, err := json.MarshalIndent(*acc, "", "\t")
	if err != nil {
		fmt.Printf("json err: %s\n", err)
	}
	return string(b)
}
