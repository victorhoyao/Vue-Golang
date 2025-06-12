package controller

import (
	"BTaskServer/util/response"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type IConfigController interface {
	GetApplicationConfig(c *gin.Context)
}

type ConfigController struct {
}

func (config ConfigController) GetApplicationConfig(c *gin.Context) {
	// 获取当前工作目录
	workDir, err := os.Getwd()
	if err != nil {
		response.ServerBad(c, nil, "获取工作目录失败")
		return
	}

	// 拼接配置文件路径
	configPath := filepath.Join(workDir, "config", "application.yml")

	// 读取配置文件内容
	content, err := os.ReadFile(configPath)
	if err != nil {
		response.ServerBad(c, nil, "读取配置文件失败")
		return
	}

	response.Success(c, gin.H{"content": string(content)}, "获取配置文件成功")
}

func NewConfigController() IConfigController {
	return ConfigController{}
}
