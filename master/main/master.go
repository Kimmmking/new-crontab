package main

import (
	"flag"
	"fmt"
	"new-crontab/master"
	"runtime"
)

var (
	confFile string // 配置文件路径
)

// 解析命令行参数
func initArgs() {
	// master -config ./master.json
	flag.StringVar(&confFile, "config", "./master.json", "指定master.json的路径")
	flag.Parse()
}

// 初始化线程数量
func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	var (
		err error
	)

	// 初始化命令行参数，获取配置文件
	initArgs()

	// 初始化线程
	initEnv()

	// 启动 Api http 服务
	if err = master.InitApiServer(); err != nil {
		goto ERR
	}

	// 初始化任务管理模块
	if err = master.InitJobMrg(); err != nil {
		goto ERR
	}

	// 加载配置
	if err = master.InitConfig(confFile); err != nil {
		goto ERR
	}
	fmt.Println(master.G_config)

ERR:
	fmt.Println(err)
}
