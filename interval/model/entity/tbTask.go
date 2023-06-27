package entity

import "time"

// TbTask  任务表
type TbTask struct {
	ID           int64     `gorm:"column:id" db:"id" json:"id" form:"id"` //  主键id
	TaskName     *string   `gorm:"column:task_name" db:"task_name" json:"task_name" form:"task_name"`
	CreatedTime  time.Time `gorm:"column:created_time" db:"created_time" json:"created_time" form:"created_time"`     //  创建时间
	ModifiedTime time.Time `gorm:"column:modified_time" db:"modified_time" json:"modified_time" form:"modified_time"` //  修改时间
	DeleteFlag   int64     `gorm:"column:delete_flag" db:"delete_flag" json:"delete_flag" form:"delete_flag"`         //  删除标志 1.删除 0.未删除
	SolvedFlag   int64     `gorm:"column:solved_flag" db:"solved_flag" json:"solved_flag" form:"solved_flag"`         //  是否推送过消息 1.已推送 0.未推送
}

func (TbTask) TableName() string {
	return "tb_task"
}
