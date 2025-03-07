package db_model

import "github.com/lts8989/comfyui-go-api/utils"

type ResultImages struct {
	Id        int          `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"` // Primary Key
	Filename  string       `gorm:"column:filename;NOT NULL" json:"filename"`       // 文件名
	Subfolder string       `gorm:"column:subfolder;NOT NULL" json:"subfolder"`     // 子文件夹
	Type      string       `gorm:"column:type;NOT NULL" json:"type"`               // 类型
	PromptId  string       `gorm:"column:prompt_id;NOT NULL" json:"prompt_id"`
	CreatedAt utils.MyTime `gorm:"column:created_at;default:CURRENT_TIMESTAMP;NOT NULL" json:"created_at"` // 记录创建时间，默认为当前时间
}

func (m *ResultImages) TableName() string {
	return "result_images"
}
