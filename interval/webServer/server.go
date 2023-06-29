package webServer

import (
	"context"
	"etcd-test/interval/controller"
	"github.com/kataras/iris/v12"
	"time"
)

var app *iris.Application

type Server struct {
}

// Init 初始化
func (s *Server) Init(a *iris.Application) error {
	app = a
	ct := controller.Controller{}
	app.Post("/api/ds_bpmn/v1/task/add", ct.AddTask)
	app.Get("/ping", ct.Ping)
	return nil
}

// NewIris 新建一个Iris
func (s *Server) NewIris() (*iris.Application, error) {
	app := iris.New()
	return app, nil
}

// IrisListen Iris开启对端口的监听
func (s *Server) IrisListen(serverPort string) error {
	err := app.Run(iris.Addr(":"+serverPort), iris.WithoutInterruptHandler)
	if err != nil {
		return err
	}
	return nil
}

// ShutdownServer 关闭Iris服务
func (s *Server) ShutdownServer() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := app.Shutdown(ctx)
	if err != nil {
		return err
	}
	return nil
}
