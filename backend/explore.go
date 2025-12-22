package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
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
	// 注意：必须在 conf.MustLoad 之前加载，以便 go-zero 的 env 标签能读取到环境变量
	loadEnvFile()

	var c config.Config
	// go-zero 的 conf.MustLoad 会自动从环境变量读取配置（通过 env 标签）
	conf.MustLoad(*configFile, &c)

	// 手动处理一些特殊配置（环境变量优先级高于 YAML）
	loadConfigFromEnv(&c)

	// 创建服务器，显式启用CORS支持
	// 允许所有来源（开发环境），生产环境应限制为特定域名
	// 允许必要的请求头
	server := rest.MustNewServer(c.RestConf,
		rest.WithCors("*"),
		rest.WithCorsHeaders("Content-Type", "Authorization"),
	)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	// 注册静态文件服务（可选，用于 Docker 部署）
	// 如果使用 Nginx，可以跳过此步骤
	registerStaticFileServer(server)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}

// loadEnvFile 加载.env文件
// 从项目根目录查找 .env 文件（相对于当前可执行文件或配置文件所在目录）
func loadEnvFile() {
	// 尝试多个可能的 .env 文件路径
	envPaths := []string{
		".env",       // 当前工作目录
		"../.env",    // 上一级目录（如果从 backend 目录运行）
		"../../.env", // 上两级目录（如果从 backend 子目录运行）
	}

	// 尝试从配置文件路径推断项目根目录
	if configPath := *configFile; configPath != "" {
		if absPath, err := filepath.Abs(configPath); err == nil {
			// 如果配置文件是 etc/explore.yaml，则 .env 应该在项目根目录
			configDir := filepath.Dir(absPath)
			if strings.HasSuffix(configDir, "etc") {
				projectRoot := filepath.Dir(configDir)
				envPaths = append([]string{filepath.Join(projectRoot, ".env")}, envPaths...)
			}
		}
	}

	// 尝试查找并加载 .env 文件
	for _, envFile := range envPaths {
		if absPath, err := filepath.Abs(envFile); err == nil {
			if _, err := os.Stat(absPath); err == nil {
				loadEnvFileFromPath(absPath)
				return
			}
		}
	}
}

// loadEnvFileFromPath 从指定路径加载 .env 文件
func loadEnvFileFromPath(envFile string) {
	data, err := os.ReadFile(envFile)
	if err != nil {
		return
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// 跳过空行和注释
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// 解析 KEY=VALUE 格式
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			// 移除引号（支持 "value" 或 'value' 格式）
			if len(value) > 0 {
				if (value[0] == '"' && value[len(value)-1] == '"') ||
					(value[0] == '\'' && value[len(value)-1] == '\'') {
					if len(value) > 1 {
						value = value[1 : len(value)-1]
					} else {
						value = ""
					}
				}
			}
			// 只设置未存在的环境变量（避免覆盖已设置的环境变量）
			if os.Getenv(key) == "" {
				os.Setenv(key, value)
			}
		}
	}
}

// loadConfigFromEnv 从环境变量加载所有配置（包括服务配置和AI配置）
// 注意：go-zero 的 conf.MustLoad 已经通过 env 标签自动读取了部分配置
// 这里主要处理一些特殊逻辑和数组类型的配置
func loadConfigFromEnv(c *config.Config) {
	// 覆盖后端服务配置（环境变量优先级最高）
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
// 注意：go-zero 的 env 标签已经自动读取了字符串类型的配置
// 这里主要处理数组类型和默认值逻辑
func loadAIConfigFromEnv(c *config.Config) {
	// 处理图片识别模型列表（数组类型，需要手动解析）
	if len(c.AI.ImageRecognitionModels) == 0 {
		modelsStr := os.Getenv("IMAGE_RECOGNITION_MODELS")
		if modelsStr != "" {
			// 解析逗号分隔的模型列表
			models := parseCommaSeparatedList(modelsStr)
			if len(models) > 0 {
				c.AI.ImageRecognitionModels = models
			} else {
				c.AI.ImageRecognitionModels = configpkg.GetDefaultImageRecognitionModels()
			}
		} else {
			c.AI.ImageRecognitionModels = configpkg.GetDefaultImageRecognitionModels()
		}
	}

	// 确保其他模型配置有默认值（如果环境变量和 YAML 都未设置）
	if c.AI.IntentModel == "" {
		c.AI.IntentModel = configpkg.DefaultIntentModel
	}
	if c.AI.ImageGenerationModel == "" {
		c.AI.ImageGenerationModel = configpkg.DefaultImageGenerationModel
	}
	if c.AI.TextGenerationModel == "" {
		c.AI.TextGenerationModel = configpkg.DefaultTextGenerationModel
	}

	// 处理UseAIModel配置（优先级：环境变量 > 配置文件 > 默认值true）
	useAIModelStr := os.Getenv("USE_AI_MODEL")
	if useAIModelStr != "" {
		// 环境变量优先级最高，解析环境变量值（支持 "true"/"false", "1"/"0", "yes"/"no"）
		useAIModelStr = strings.ToLower(strings.TrimSpace(useAIModelStr))
		c.AI.UseAIModel = (useAIModelStr == "true" || useAIModelStr == "1" || useAIModelStr == "yes")
	} else {
		// 如果环境变量未设置，使用配置文件中的值（go-zero 的 conf.MustLoad 已经自动读取）
		// 如果配置文件中也未设置（零值 false），设置为默认值 true
		// 注意：如果配置文件中显式设置为 false，c.AI.UseAIModel 已经是 false，不会被覆盖
		// 如果配置文件中设置为 true 或未设置，这里会确保值为 true
		if !c.AI.UseAIModel {
			// 配置文件未设置（零值 false），设置为默认值 true
			c.AI.UseAIModel = true
		}
	}
}

// parseCommaSeparatedList 解析逗号分隔的字符串列表
func parseCommaSeparatedList(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// registerStaticFileServer 注册静态文件服务
// 用于 Docker 部署时，后端同时提供前端静态文件服务
// 如果使用 Nginx，可以通过环境变量 ENABLE_STATIC_SERVER=false 禁用
func registerStaticFileServer(server *rest.Server) {
	// 检查是否启用静态文件服务（默认启用）
	enableStatic := os.Getenv("ENABLE_STATIC_SERVER")
	if enableStatic == "false" || enableStatic == "0" {
		return
	}

	// 查找前端静态文件目录
	staticDirs := []string{
		"frontend",       // 当前目录下的 frontend
		"../frontend",    // 上一级目录
		"../../frontend", // 上两级目录
		"/app/frontend",  // Docker 容器中的路径
	}

	var staticDir string
	for _, dir := range staticDirs {
		if absPath, err := filepath.Abs(dir); err == nil {
			if info, err := os.Stat(absPath); err == nil && info.IsDir() {
				// 检查是否有 index.html
				if _, err := os.Stat(filepath.Join(absPath, "index.html")); err == nil {
					staticDir = absPath
					break
				}
			}
		}
	}

	if staticDir == "" {
		// 如果找不到静态文件目录，静默跳过（可能使用 Nginx）
		return
	}

	// 创建文件服务器
	fileServer := http.FileServer(http.Dir(staticDir))

	// 注册静态文件路由（优先级低于 API 路由）
	// 使用通配符匹配所有非 API 路径
	server.AddRoute(rest.Route{
		Method: http.MethodGet,
		Path:   "/",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			// 如果是 API 路径，不处理（由 API 路由处理）
			if strings.HasPrefix(r.URL.Path, "/api") {
				http.NotFound(w, r)
				return
			}
			// 提供静态文件服务
			fileServer.ServeHTTP(w, r)
		},
	})

	// 处理前端路由（SPA 路由回退到 index.html）
	server.AddRoute(rest.Route{
		Method: http.MethodGet,
		Path:   "/*",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			// 如果是 API 路径，不处理
			if strings.HasPrefix(r.URL.Path, "/api") {
				http.NotFound(w, r)
				return
			}

			// 检查文件是否存在
			filePath := filepath.Join(staticDir, r.URL.Path)
			if info, err := os.Stat(filePath); err == nil && !info.IsDir() {
				// 文件存在，直接提供
				fileServer.ServeHTTP(w, r)
				return
			}

			// 文件不存在，可能是前端路由，返回 index.html
			indexPath := filepath.Join(staticDir, "index.html")
			if _, err := os.Stat(indexPath); err == nil {
				http.ServeFile(w, r, indexPath)
			} else {
				http.NotFound(w, r)
			}
		},
	})
}
