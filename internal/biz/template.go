package biz

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lts8989/sd_api/internal/model/db_model"
	"github.com/lts8989/sd_api/utils"
	"github.com/lts8989/sd_sdk/sdk"
	"os"
)

func CreateTask(ctx utils.Context, tempId int32, params any) error {
	byteParams, err := json.Marshal(params)
	if err != nil {
		return err
	}
	paramsMap := make(map[string]string)
	if err = json.Unmarshal(byteParams, &paramsMap); err != nil {
		return err
	}

	//读取数据库，通过id查找到模版文件名
	temp := db_model.Templates{}
	if err := utils.MysqlDB.First(&temp, tempId).Error; err != nil {
		return errors.New("没有找到模版：" + err.Error())
	}

	//读取模版文件内容
	filePath := "prompt_temp/" + temp.FileName
	templateContent, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	//替换模版中的参数
	for k, v := range paramsMap {
		key := []byte(fmt.Sprintf("<%%%s%%>", k))
		templateContent = bytes.Replace(templateContent, key, []byte(v), -1)
	}

	//调用sd接口，下发绘图命令
	resp, err := sdk.Prompt(utils.Conf.SdServCfg.ClientId, templateContent)
	if err != nil {
		return err
	}

	task := db_model.ExecutionTasks{
		TemplateId: int(tempId),
		Parameters: string(byteParams),
	}

	//处理sd的业务错误
	err = nil
	if len(resp.PromptId) == 0 {
		task.PromptId = ""
		errResult, _ := json.Marshal(resp)
		task.Result = string(errResult)
		err = errors.New(resp.Error.Message)
	} else {
		task.PromptId = resp.PromptId
		task.Result = "{}"
	}

	//promptid存入数据库
	if dbErr := utils.MysqlDB.Create(&task).Error; dbErr != nil {
		return errors.New("创建任务失败:" + dbErr.Error())
	}

	return err
}
