package router

import (
	"etcd-test/interval/controller"
	"github.com/kataras/iris/v12"
)

type Router struct {
}

func (router *Router) Init(r *iris.Application) error {
	ct := &controller.Controller{}
	r.Post("/api/ds_bpmn/v1/task/add", ct.AddTask)
	r.Get("/ping", ct.Ping)

	return nil
}
