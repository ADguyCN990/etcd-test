package controller

import (
	"etcd-test/interval/model/request"
	"etcd-test/interval/service"
	"github.com/kataras/iris/v12"
	"github.com/sirupsen/logrus"
	"net/http"
)

var sv = &service.Service{}

type Controller struct {
}

// Ping 心跳检测
func (controller *Controller) Ping(c iris.Context) {
	c.StatusCode(200)
	err := c.JSON(map[string]interface{}{
		"message": "Pong",
		"error":   nil,
	})
	if err != nil {
		return
	}
}

// AddTask 新增一个任务
func (controller *Controller) AddTask(c iris.Context) {
	var req request.AddTask
	if err := c.ReadJSON(&req); err != nil {
		c.StatusCode(400)
		logrus.WithError(err).Error("Fail！c.ReadJSON")
		_, err2 := c.WriteString("invalid JSON")
		if err2 != nil {
			logrus.WithError(err2).Error("Fail！c.WriteString")
		}
		return
	}

	taskName := req.TaskName

	// 参数校验
	if taskName == nil {
		c.StatusCode(http.StatusBadRequest)
		message := "taskName为空"
		err2 := c.JSON(map[string]interface{}{
			"message": &message,
			"error":   nil,
		})
		if err2 != nil {
			logrus.WithError(err2).Error("Fail！c.WriteString")
		}
		logrus.Error("新建task时，taskName为空")
		return
	}

	response, err := sv.AddTask(taskName)
	if err != nil {
		c.StatusCode(http.StatusBadRequest)
		err2 := c.JSON(map[string]interface{}{
			"message": *response.Message,
			"error":   *response.Error,
		})
		if err2 != nil {
			logrus.WithError(err2).Error("Fail！c.WriteString")
		}
		return
	}

	c.StatusCode(http.StatusOK)
	err2 := c.JSON(map[string]interface{}{
		"message": *response.Message,
		"error":   nil,
	})
	if err2 != nil {
		logrus.WithError(err2).Error("Fail! c.WriteString")
	}

}
