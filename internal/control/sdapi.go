package control

import (
	"github.com/lts8989/comfyui-go-api/internal/biz"
	"github.com/lts8989/comfyui-go-api/internal/model"
	"github.com/lts8989/comfyui-go-api/utils"
	"github.com/lts8989/comfyui-go-sdk/sdk"
)

func CreateTask(c *utils.Context) {
	var u model.CreateTaskForm
	if err := c.ShouldBind(&u); err != nil {
		c.Error(err)
		return
	}

	if err := biz.CreateTask(*c, u.TempId, u.Params); err != nil {
		c.Error(err)
		return
	}
	c.Success(nil)
}

func Ping(c *utils.Context) {
	if a, err := sdk.GetSystemStats(); err != nil {
		c.Error(err)
	} else {
		c.Success(a)
	}
}

func History(c *utils.Context) {
	var u model.HistoryForm
	if err := c.ShouldBind(&u); err != nil {
		c.Error(err)
		return
	}

	err := biz.FetchTask(u.PromptId)
	if err != nil {
		c.Error(err)
		return
	}

	c.Success(nil)
}
