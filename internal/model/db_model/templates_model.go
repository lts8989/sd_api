package db_model

import "github.com/lts8989/comfyui-go-api/utils"

type Templates struct {
	Id           int          `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	TemplateName string       `gorm:"column:template_name;NOT NULL" json:"template_name"`            // 模板的名称
	FileName     string       `gorm:"column:file_name;NOT NULL" json:"file_name"`                    // 与模板相关联的文件名
	Parameters   string       `gorm:"column:parameters;NOT NULL" json:"parameters"`                  // 存储与模板相关的参数值
	CreatedAt    utils.MyTime `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"` // 记录的创建时间，默认为当前时间戳
}

func (m *Templates) TableName() string {
	return "templates"
}
