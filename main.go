package main

import (
	"bubble/dao"
	"bubble/models"
	"bubble/routers"
	"bubble/service"
	"bubble/setting"
	"fmt"
	"log"
	"os"
	"time"
)

const defaultConfFile = "./conf/config.ini"

func main() {
	log.Println("连接以太坊客户端")
	//连接以太坊客户端
	blockChain.InitEthClient()

	//开启定时任务
	go startSchedule()

	//监听transfer事件
	go blockChain.QueryTransferInfoFromBlockChain()

	confFile := defaultConfFile
	if len(os.Args) > 2 {
		fmt.Println("use specified conf file: ", os.Args[1])
		confFile = os.Args[1]
	} else {
		fmt.Println("no configuration file was specified, use ./conf/config.ini")
	}
	// 加载配置文件
	if err := setting.Init(confFile); err != nil {
		fmt.Printf("load config from file failed, err:%v\n", err)
		return
	}
	// 创建数据库
	// sql: CREATE DATABASE bubble;
	// 连接数据库
	err := dao.InitMySQL(setting.Conf.MySQLConfig)
	if err != nil {
		fmt.Printf("init mysql failed, err:%v\n", err)
		return
	}
	defer dao.Close() // 程序退出关闭数据库连接
	// 模型绑定
	dao.DB.AutoMigrate(&models.Todo{})
	// 注册路由
	r := routers.SetupRouter()
	if err := r.Run(fmt.Sprintf(":%d", setting.Conf.Port)); err != nil {
		fmt.Printf("server startup failed, err:%v\n", err)
	}
}

func startSchedule() {
	log.Println("开启定时任务")
	ticker := time.NewTicker(15 * time.Second)
	for {
		select {
		case <-ticker.C:
			fmt.Println("执行定时任务")
			blockChain.QueryLatestBlockFromChain()
		}
	}
}
