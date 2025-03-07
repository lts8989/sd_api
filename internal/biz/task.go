package biz

import (
	"encoding/json"
	"fmt"
	"github.com/lts8989/comfyui-go-api/internal/model/db_model"
	"github.com/lts8989/comfyui-go-api/utils"
	sdk_model "github.com/lts8989/comfyui-go-sdk/model"
	"github.com/lts8989/comfyui-go-sdk/sdk"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"os"
)

func ReceivedMsg(data sdk_model.WsReceive) error {
	if t, ok := sdk_model.ReceiveTypeMap[data.Type]; ok {
		promptId := data.Data.PromptID
		task := db_model.ExecutionTasks{}
		if dbErr := utils.MysqlDB.Where("prompt_id=?", promptId).Take(&task).Error; dbErr != nil {
			return fmt.Errorf("没有找到任务,prompt_id:%s,error:%v", promptId, dbErr)
		}
		if int8(task.Status) >= t {
			log.Infof("任务状态不需要更新,prompt_id:%s,db状态:%s,ws状态:%s", promptId, sdk_model.ReceiveTypeDescMap[int8(task.Status)], sdk_model.ReceiveTypeDescMap[t])
			return nil
		}

		resultStr, _ := json.Marshal(data)

		txErr := utils.MysqlDB.Transaction(func(tx *gorm.DB) error {
			upTask := db_model.ExecutionTasks{
				Status: int(t),
				Result: string(resultStr),
			}
			if dbErr := tx.Model(&task).Updates(upTask).Error; dbErr != nil {
				log.Errorf("更新任务失败,prompt_id:%s,error:%v", promptId, dbErr)
				return dbErr
			}

			if len(data.Data.Output.Images) > 0 {
				images := make([]db_model.ResultImages, 0, len(data.Data.Output.Images))
				for _, v := range data.Data.Output.Images {
					images = append(images, db_model.ResultImages{
						Filename:  v.Filename,
						Subfolder: v.Subfolder,
						Type:      v.Type,
						PromptId:  promptId,
					})
				}
				if dbErr := tx.Omit("created_at").CreateInBatches(images, len(images)).Error; dbErr != nil {
					log.Errorf("记录生成结果失败,prompt_id:%s,error:%v", promptId, dbErr)
					return dbErr
				}

				if dbErr := DownLoadImg(images); dbErr != nil {
					log.Errorf("下载图片失败,prompt_id:%s,error:%v", promptId, dbErr)
					return dbErr
				}
			}

			return nil
		})
		return txErr
	} else {
		log.Infof("不是任务相关的type,忽略:%s", data.Type)
		return nil
	}
}

func DownLoadImg(imgs []db_model.ResultImages) error {
	for _, v := range imgs {
		log.Infof("begin download img filename:%s", v.Filename)
		req := sdk_model.ViewReq{
			Filename:  v.Filename,
			Subfolder: v.Subfolder,
			Type:      v.Type,
		}

		filepath := "img/" + req.Filename

		if _, err := os.Stat(filepath); err == nil {
			log.Infof("File already exists: %s", req.Filename)
			continue
		}

		//调用sdapi接口下载图片文件
		fileData, err := sdk.View(req)
		if err != nil {
			return err
		}

		//保存图片文件到硬盘
		err = os.WriteFile(filepath, fileData, 0644)
		if err != nil {
			return err
		}
		log.Infof("end download img filename:%s", v.Filename)
	}
	return nil
}

func FetchTask(promptId string) error {
	task := db_model.ExecutionTasks{}
	if dbErr := utils.MysqlDB.Where("prompt_id=?", promptId).Take(&task).Error; dbErr != nil {
		return fmt.Errorf("没有找到任务,prompt_id:%s,error:%v", promptId, dbErr)
	}

	if int8(task.Status) == sdk_model.ReceiveTypeMap[sdk_model.ReceiveTypeSuccess] {
		return fmt.Errorf("任务已经完成")
	}

	//从sd的api接口中获取执行结果
	imgList, err := sdk.History(promptId)
	if err != nil {
		return err
	}

	modelList := make([]db_model.ResultImages, 0)
	for _, m := range imgList {
		modelList = append(modelList, db_model.ResultImages{
			Filename:  m.Filename,
			Subfolder: m.Subfolder,
			Type:      m.Type,
		})
	}

	if err = DownLoadImg(modelList); err != nil {
		return err
	}

	//todo: 更新任务状态，插入图片表

	return nil
}
