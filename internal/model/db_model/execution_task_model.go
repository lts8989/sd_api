package db_model

import "github.com/lts8989/comfyui-go-api/utils"

type ExecutionTasks struct {
	Id         int          `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"` // 自增主键，唯一标识每个任务
	PromptId   string       `gorm:"column:prompt_id" json:"prompt_id"`
	TemplateId int          `gorm:"column:template_id;NOT NULL" json:"template_id"`                // 模版 ID
	Parameters string       `gorm:"column:parameters;NOT NULL" json:"parameters"`                  // 参数对，使用 JSON 格式存储任务的参数
	Status     int          `gorm:"column:status;NOT NULL" json:"status"`                          // 任务状态。1、待处理；2、进行中；3、已完成；4、失败
	Result     string       `gorm:"column:result;NOT NULL" json:"result"`                          // 执行结果，存储任务执行后的结果信息
	CreatedAt  utils.MyTime `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"` // 记录创建时间，默认为当前时间
	UpdatedAt  utils.MyTime `gorm:"column:updated_at;default:CURRENT_TIMESTAMP" json:"updated_at"` // 记录更新时间，自动更新为当前时间
}

func (m *ExecutionTasks) TableName() string {
	return "execution_tasks"
}
