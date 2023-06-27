package service

import (
	"errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func (S *Service) Listen() error {
	for {
		ids, err := d.FindTasksNotSolved()
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				// 如果不是因为没有数据，那么返回err；否则不管
				logrus.WithError(err).Error("查询未更新solve状态的数据失败")
				return err
			}
		}
		for _, id := range *ids {
			// 处理未发送的任务
			// 1. 向消息队列发送消息
			// 2. 更新数据库字段
			task, err := d.FindTaskById(id)
			if err != nil {
				logrus.WithError(err).Error("根据ID查询数据失败")
				return err
			}
			logrus.Info("向消息队列发送任务：", task)
			// TODO 向消息队列发送任务
			err = d.SolveTask(id)
			if err != nil {
				logrus.WithError(err).Error("更改solved字段失败")
				return err
			}
			logrus.Info("成功处理任务：", task)
		}
	}
}
