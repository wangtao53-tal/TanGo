package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/tango/explore/internal/config"
	configpkg "github.com/tango/explore/internal/config"
	"github.com/tango/explore/internal/handler"
	"github.com/tango/explore/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/explore.yaml", "the config file")

func main() {
	flag.Parse()

	// 加载.env文件（如果存在）
	loadEnvFile()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	// 从环境变量覆盖配置（包括端口、主机等）
	loadConfigFromEnv(&c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}

// loadEnvFile 加载.env文件
func loadEnvFile() {
	envFile := ".env"
	if _, err := os.Stat(envFile); err == nil {
		// 简单的.env文件解析
		data, err := os.ReadFile(envFile)
		if err == nil {
			lines := strings.Split(string(data), "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line == "" || strings.HasPrefix(line, "#") {
					continue
				}
				parts := strings.SplitN(line, "=", 2)
				if len(parts) == 2 {
					key := strings.TrimSpace(parts[0])
					value := strings.TrimSpace(parts[1])
					// 移除引号
					if len(value) > 0 && (value[0] == '"' || value[0] == '\'') {
						if len(value) > 1 {
							value = value[1 : len(value)-1]
						} else {
							value = ""
						}
					}
					os.Setenv(key, value)
				}
			}
		}
	}
}

// loadConfigFromEnv 从环境变量加载所有配置（包括服务配置和AI配置）
func loadConfigFromEnv(c *config.Config) {
	// 覆盖后端服务配置
	if host := os.Getenv("BACKEND_HOST"); host != "" {
		c.Host = host
	}
	if portStr := os.Getenv("BACKEND_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			c.Port = port
		}
	}

	// 加载AI配置
	loadAIConfigFromEnv(c)
}

// loadAIConfigFromEnv 从环境变量加载AI配置
func loadAIConfigFromEnv(c *config.Config) {
	if c.AI.EinoBaseURL == "" {
		c.AI.EinoBaseURL = os.Getenv("EINO_BASE_URL")
	}
	if c.AI.AppID == "" {
		c.AI.AppID = os.Getenv("TAL_MLOPS_APP_ID")
	}
	if c.AI.AppKey == "" {
		c.AI.AppKey = os.Getenv("TAL_MLOPS_APP_KEY")
	}
	if c.AI.IntentModel == "" {
		c.AI.IntentModel = os.Getenv("INTENT_MODEL")
		if c.AI.IntentModel == "" {
			c.AI.IntentModel = configpkg.DefaultIntentModel
		}
	}
	if len(c.AI.ImageRecognitionModels) == 0 {
		modelsStr := os.Getenv("IMAGE_RECOGNITION_MODELS")
		if modelsStr != "" {
			// 解析逗号分隔的模型列表
			models := []string{}
			for _, model := range splitAndTrim(modelsStr, ",") {
				if model != "" {
					models = append(models, model)
				}
			}
			c.AI.ImageRecognitionModels = models
		} else {
			c.AI.ImageRecognitionModels = configpkg.GetDefaultImageRecognitionModels()
		}
	}
	if c.AI.ImageGenerationModel == "" {
		c.AI.ImageGenerationModel = os.Getenv("IMAGE_GENERATION_MODEL")
		if c.AI.ImageGenerationModel == "" {
			c.AI.ImageGenerationModel = configpkg.DefaultImageGenerationModel
		}
	}
	if c.AI.TextGenerationModel == "" {
		c.AI.TextGenerationModel = os.Getenv("TEXT_GENERATION_MODEL")
		if c.AI.TextGenerationModel == "" {
			c.AI.TextGenerationModel = configpkg.DefaultTextGenerationModel
		}
	}
}

// splitAndTrim 分割字符串并去除空白
func splitAndTrim(s, sep string) []string {
	result := []string{}
	for _, part := range splitString(s, sep) {
		trimmed := trimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func splitString(s, sep string) []string {
	result := []string{}
	current := ""
	for i := 0; i < len(s); i++ {
		if i+len(sep) <= len(s) && s[i:i+len(sep)] == sep {
			if current != "" {
				result = append(result, current)
				current = ""
			}
			i += len(sep) - 1
		} else {
			current += string(s[i])
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}

func trimSpace(s string) string {
	start := 0
	end := len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n' || s[start] == '\r') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n' || s[end-1] == '\r') {
		end--
	}
	return s[start:end]
}
