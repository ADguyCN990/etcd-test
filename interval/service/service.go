package service

import (
	"etcd-test/interval/dao"
	"etcd-test/interval/model/response"
	"github.com/sirupsen/logrus"
)

type Service struct {
}

var d = &dao.Dao{}

// AddTask 新增一个任务
func (S *Service) AddTask(taskName *string) (*response.AddTask, error) {
	err := d.AddTask(taskName)
	if err != nil {
		logrus.WithError(err).Error("添加数据失败")
		message := "Dao层新建任务失败"
		resErr := err.Error()
		return &response.AddTask{
			Message: &message,
			Error:   &resErr,
		}, err
	}
	logrus.Info("成功新建一个任务")
	message := "新建任务成功"
	return &response.AddTask{
		Message: &message,
		Error:   nil,
	}, nil

}
