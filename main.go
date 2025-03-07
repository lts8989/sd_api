package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/lts8989/comfyui-go-api/internal/biz"
	"github.com/lts8989/comfyui-go-api/internal/control"
	"github.com/lts8989/comfyui-go-api/utils"
	sdk_log "github.com/lts8989/comfyui-go-sdk/log"
	"github.com/lts8989/comfyui-go-sdk/sdk"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"path"
	"reflect"
	"time"
)

func init() {
	err := utils.InitConfig()
	if err != nil {
		return
	}

	if err = ConfigLocalFileSystemLogger(utils.Conf.LogConfig); err != nil {
		log.Fatal("本地日志初始化失败：" + err.Error())
		os.Exit(1)
	}

	if err = utils.InitDB(utils.Conf.DbCfg); err != nil {
		log.Fatal("数据库加载失败：" + err.Error())
	}

	sdkConfig := utils.Conf.SdServCfg
	sdk.Setup(sdkConfig.Domain, sdkConfig.ClientId, 30, 5)

	sdk_log.InitLogger(log.StandardLogger())
}

// ConfigLocalFileSystemLogger 初始化本地文件系统日志
func ConfigLocalFileSystemLogger(cfg utils.LogCfg) error {
	baseLogPath := path.Join(cfg.Path, cfg.FileName)
	writer, err := rotatelogs.New(
		baseLogPath+"%Y%m%d.log",                                                 // 分割后的文件名称
		rotatelogs.WithLinkName(baseLogPath),                                     // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(time.Duration(cfg.MaxAge)),                         // 设置最大保存时间(7天)
		rotatelogs.WithRotationCount(48),                                         // 最多存365个文件
		rotatelogs.WithRotationTime(time.Second*time.Duration(cfg.RotationTime)), // 设置日志切割时间间隔(1天)
	)
	if err != nil {
		return errors.New(fmt.Sprintf("Conf local file system logger error. %+v", errors.WithStack(err)))
	}
	lfHook := lfshook.NewHook(lfshook.WriterMap{
		log.DebugLevel: writer, // 为不同级别设置不同的输出目的
		log.InfoLevel:  writer,
		log.WarnLevel:  writer,
		log.ErrorLevel: writer,
		log.FatalLevel: writer,
		log.PanicLevel: writer,
	}, &log.TextFormatter{TimestampFormat: "2006-01-02 15:04:05"})
	if err != nil {
		return errors.New(fmt.Sprintf("Conf local file system logger error. %+v", errors.WithStack(err)))
	}
	log.AddHook(lfHook)
	return nil
}

// 初始化web服务
func initWebServer() error {
	log.Info("初始化web服务")
	uni := ut.New(zh.New())
	utils.Trans, _ = uni.GetTranslator("zh")
	router := gin.New()
	router.Use(Recover)
	router.Use(LogMiddleware())
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("simpleMap", simpleMap)
		if err != nil {
			return err
		}
		_ = zh_translations.RegisterDefaultTranslations(v, utils.Trans)
	}

	sdApi := router.Group("/sdapi")
	sdApi.GET("ping", utils.Build(control.Ping))
	sdApi.POST("create_task", utils.Build(control.CreateTask))
	sdApi.POST("history", utils.Build(control.History))
	err := router.Run(fmt.Sprintf(":%d", utils.Conf.ServCfg.Port))

	return err
}
func Recover(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			c.JSON(http.StatusOK, &utils.Response{
				Code:        utils.RespCodeBizError,
				MessageData: errorToString(r),
				Data:        nil,
			})
			c.Abort()
		}
	}()
	c.Next()
}
func errorToString(r interface{}) string {
	switch v := r.(type) {
	case error:
		return v.Error()
	default:
		return r.(string)
	}
}

func LogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log only when path is not being skipped

		// Stop timer
		end := time.Now()
		timeSubNano := end.Sub(start)
		timeSub := timeSubNano.Milliseconds()
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		bodySize := c.Writer.Size()
		if raw != "" {
			path = path + "?" + raw
		}

		log.Infof("[GIN] %dms | %s | %s | %d | %d | %s",
			timeSub, clientIP, method, statusCode, bodySize, path,
		)

	}
}

func simpleMap(fl validator.FieldLevel) bool {
	str := fl.Field().String()
	fl.Field().Interface()
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(str), &result); err != nil {
		return false
	}

	for _, val := range result {
		kind := reflect.TypeOf(val).Kind()
		switch kind {
		case reflect.Bool,
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64, reflect.String:
			continue
		default:

			return false
		}
	}

	return true
}

func main() {
	go sdk.ConnectToWebSocket(biz.ReceivedMsg)
	if err := initWebServer(); err != nil {
		log.Fatal("初始化web服务失败:" + err.Error())
		return
	}
}
