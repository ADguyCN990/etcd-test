package main

import (
	"etcd-test/interval/dao"
	"etcd-test/interval/service"
	"etcd-test/interval/start"
	"etcd-test/interval/webServer"
	"etcd-test/pkg"
	"github.com/sirupsen/logrus"
	"sync"
)

var st = &start.Init{}
var Dao = &dao.Dao{}
var se = &service.Service{}
var ir = &webServer.Server{}

func main() {
	// 优雅退出程序用
	var wg sync.WaitGroup

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
	DbRoutine := pkg.NewCloseGoRoutineChannel()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case _, ok := <-DbRoutine:
				if !ok {
					logrus.Info("收到中断信号，关闭监听数据库的协程...")
					return
				}
			default:
				err := se.Listen()
				if err != nil {
					logrus.WithError(err).Fatal("数据库监听失败")
				}
			}
		}

	}()

	// 开始监听端口
	IrisPortRoutine := pkg.NewCloseGoRoutineChannel()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case _, ok := <-IrisPortRoutine:
				logrus.Info("收到关闭协程channel发送的消息了")
				if !ok {
					logrus.Info("收到中断信号，关闭监听Iris端口的协程...")
					return
				}
			default:
				err = ir.IrisListen(start.ConfigData.Server.Port)
				if err != nil {
					logrus.WithError(err).Fatal("IrisServer端口监听失败")
				}
			}
		}
	}()

	// 监听中断信号，优雅退出程序
	wg.Add(1)
	pkg.WaitExit(&wg, st.Exit, &DbRoutine, &IrisPortRoutine)

	wg.Wait()

}
