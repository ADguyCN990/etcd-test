package main

import (
	"etcd-test/interval/dao"
	"etcd-test/interval/router"
	"etcd-test/interval/service"
	"etcd-test/start"
	"github.com/sirupsen/logrus"
	"sync"
)

var st = &start.Init{}
var Dao = &dao.Dao{}
var rt = &router.Router{}
var se = &service.Service{}
var wg = sync.WaitGroup{}

func main() {

	// 加载配置文件
	err := st.InitConfig()
	if err != nil {
		logrus.WithError(err).Fatal("无法加载配置文件")
	}
	logrus.Info("加载配置文件成功")

	// 启动数据库
	db, err := st.Database()
	if err != nil {
		logrus.WithError(err).Fatal("启动数据库失败")
	}
	logrus.Info("启动数据库成功")

	// 初始化Dao层
	Dao.Init(db)

	// 启动iris
	app, err := st.Iris()
	if err != nil {
		logrus.WithError(err).Fatal("启动iris失败")
	}
	logrus.Info("启动iris成功")

	// 初始化Router
	err = rt.Init(app)
	if err != nil {
		logrus.WithError(err).Fatal("初始化Router失败")
	}
	logrus.Info("成功初始化Router")

	wg.Add(2)

	// 监听数据库新增数据
	go func() {
		err := se.Listen()
		if err != nil {
			logrus.WithError(err).Fatal("数据库监听失败")
		}
	}()

	// 开始监听端口
	go func() {
		err = st.IrisListen(app)
		if err != nil {
			logrus.WithError(err).Fatal("iris监听端口失败")
		}
		logrus.Info("成功监听端口")
	}()

	wg.Wait()

}
