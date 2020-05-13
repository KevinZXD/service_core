package main

import (
	"errors"
	"fmt"
	"os"
	"service_core/core"
	"github.com/docopt/docopt-go"
)

// 命令行帮助信息
var usage = `
Usage:
    service_core --config=<x.toml> [--test]

Options:
	-h --help     Show this screen.
	--version     Show version.
`

// 命令行选项输入
type CliOption struct {
	CfgFname string `docopt:"--config"`
	IsTest   bool   `docopt:"--test"`
}

func (co *CliOption) Validate() error {
	if len(co.CfgFname) <= 0 {
		return errors.New("config file can not be empty")
	}
	return nil
}

// 执行命令
func doCommand(opt *CliOption) int {
	service, err := core.NewServiceCore(opt.CfgFname)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	if opt.IsTest {
		if err := service.Try(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
		return 0
	}
	if err := service.Start(); err != nil {
		fmt.Println("stdout:", err)
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}
func Run() int {
	var cliOpt = CliOption{}
	opts, err := docopt.ParseArgs(usage, nil, Version)
	if err != nil {
		fmt.Println(usage)
		return 1
	}
	// 数据绑定
	if err := opts.Bind(&cliOpt); err != nil {
		fmt.Println(err.Error())
		return 2
	}
	// 检查命令输入
	if err := cliOpt.Validate(); err != nil {
		fmt.Println(err.Error())
		return 3
	}
	// 执行命令
	return doCommand(&cliOpt)
}

// 入口
func main() {
	os.Exit(Run())
}
