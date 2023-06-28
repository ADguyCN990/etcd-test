package main

import (
	"etcd-test/interval/dao"
	"etcd-test/interval/service"
	"etcd-test/interval/webServer"
	"etcd-test/pkg"
	"etcd-test/start"
	"github.com/kataras/iris/v12/core/router"
	"github.com/sirupsen/logrus"
	"sync"
)

var st = &start.Init{}
var Dao = &dao.Dao{}
var rt = &router.Router{}
var se = &service.Service{}
var ir = &webServer.Server{}

func main() {

	// 加载配置文件
	err := st.InitConfig()
	if err != nil {
		logrus.WithError(err).Fatal("无法加载配置文件")
	}
	logrus.Info("加载配置文件成功")

	// 启动数据库和gorm
	db, err := st.Database()
	if err != nil {
		logrus.WithError(err).Fatal("启动数据库失败")
	}
	logrus.Info("启动数据库成功")
	Dao.Init(db)

	// 启动iris
	app, err := ir.NewIris()
	if err != nil {
		logrus.WithError(err).Fatal("启动iris失败")
	}
	logrus.Info("启动iris成功")
	err = ir.Init(app)
	if err != nil {
		logrus.WithError(err).Fatal("启动iris失败")
		return
	}

	// 监听数据库新增数据
	go func() {
		err := se.Listen()
		if err != nil {
			logrus.WithError(err).Fatal("数据库监听失败")
		}
		logrus.Info("成功监听数据库的新增数据")
	}()

	// 开始监听端口
	go func() {
		err = ir.IrisListen(start.ConfigData.Server.Port)
	}()

	// 监听中断信号，优雅退出程序
	var wg sync.WaitGroup
	wg.Add(1)
	pkg.WaitExit(&wg, st.Exit)

	wg.Wait()

}
