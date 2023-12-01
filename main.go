package main

import (
	"fmt"
	"webapi_erc20/common/config"
	"webapi_erc20/common/logs"
	"webapi_erc20/dao"
	"webapi_erc20/routes"
	"webapi_erc20/service"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {

	// 初始化 config 配置
	cfg := config.NewConfig("./config.yaml")
	// 初始化日志
	logs.Init(cfg.Log)

	logs.Debugf("mode=%+v", cfg.Web)

	//连接数据库
	err := dao.InitMySql(cfg)
	if err != nil {
		panic(err)
	}

	//建立table
	dao.CreateTable()

	//啟動監聽區塊排程
	service.RunListenBlock()

	//注册路由
	r := routes.SetRouter()
	//启动端口为xxxx的项目
	webport := fmt.Sprintf(":%s", cfg.Web.Port)
	r.Run(webport)
}
