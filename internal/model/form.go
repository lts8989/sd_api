package model

type CreateTaskForm struct {
	TempId int32 `form:"temp_id" binding:"required" json:"temp_id"`
	Params any   `form:"params" binding:"required" json:"params"`
}

type HistoryForm struct {
	PromptId string `form:"prompt_id" binding:"required" json:"prompt_id"`
}
