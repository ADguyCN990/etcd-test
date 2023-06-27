package dao

import (
	"etcd-test/interval/model/entity"
	"gorm.io/gorm"
	"time"
)

type Dao struct {
}

var db *gorm.DB

// Init 初始化
func (D *Dao) Init(a *gorm.DB) {
	db = a
}

// AddTask 新增一个任务
func (D *Dao) AddTask(taskName *string) error {
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return err
	}
	task := entity.TbTask{
		TaskName:     taskName,
		CreatedTime:  time.Now().In(location),
		ModifiedTime: time.Now().In(location),
		DeleteFlag:   0,
		SolvedFlag:   0,
	}
	db.Debug().Create(&task)
	return nil
}

// FindTasksAfterTimestamp 查询大于指定时间戳的所有任务
func (D *Dao) FindTasksAfterTimestamp(timestamp time.Time) (*[]int64, error) {
	var ids []int64
	if err := db.Where("modified_time > ? AND delete_flag = ?", timestamp, 0).Model(&entity.TbTask{}).Pluck("id", &ids).Error; err != nil {
		return nil, err
	}
	return &ids, nil
}

// FindTasksNotSolved FindTaskNotSolved 查询未发送消息的任务
func (D *Dao) FindTasksNotSolved() (*[]int64, error) {
	var ids []int64
	result := db.Model(&entity.TbTask{}).Where("solved_flag = ? AND delete_flag = ?", 0, 0).Pluck("id", &ids)
	if result.Error != nil {
		return nil, result.Error
	}
	return &ids, nil
}

// FindTaskById 根据ID查询任务
func (D *Dao) FindTaskById(id int64) (*entity.TbTask, error) {
	var task *entity.TbTask
	result := db.First(&task, id).Where("delete_flag = ?", 0)
	if result.Error != nil {
		return nil, result.Error
	}
	return task, nil
}

func (D *Dao) SolveTask(id int64) error {
	task, err := D.FindTaskById(id)
	if err != nil {
		return err
	}
	task.SolvedFlag = 1
	if err := db.Save(&task).Error; err != nil {
		return err
	}
	return nil
}
