package request

// AddTask 新增一个任务的request
type AddTask struct {
	TaskName *string `json:"task_name"`
}
