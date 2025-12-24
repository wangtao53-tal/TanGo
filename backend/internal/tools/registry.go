package tools

import (
	"sync"

	"github.com/zeromicro/go-zero/core/logx"
)

// ToolRegistry 工具注册表
// 负责管理所有可用工具，支持注册、获取、列表等功能
type ToolRegistry struct {
	tools map[string]Tool
	mu    sync.RWMutex
	logger logx.Logger
}

// NewToolRegistry 创建新的工具注册表
func NewToolRegistry(logger logx.Logger) *ToolRegistry {
	return &ToolRegistry{
		tools:  make(map[string]Tool),
		logger: logger,
	}
}

// Register 注册工具到注册表
// 如果工具已存在，会覆盖原有工具
func (r *ToolRegistry) Register(tool Tool) {
	if tool == nil {
		r.logger.Errorw("尝试注册nil工具")
		return
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	name := tool.Name()
	if name == "" {
		r.logger.Errorw("工具名称为空，无法注册")
		return
	}

	r.tools[name] = tool
	r.logger.Infow("工具已注册", logx.Field("tool", name))
}

// GetTool 根据名称获取工具
// 返回工具实例和是否存在
func (r *ToolRegistry) GetTool(name string) (Tool, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tool, ok := r.tools[name]
	return tool, ok
}

// ListTools 列出所有已注册的工具
func (r *ToolRegistry) ListTools() []Tool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tools := make([]Tool, 0, len(r.tools))
	for _, tool := range r.tools {
		tools = append(tools, tool)
	}
	return tools
}

// GetToolsForAgent 根据Agent类型获取工具列表
// agentType: "Science", "Language", "Humanities"
func (r *ToolRegistry) GetToolsForAgent(agentType string) []Tool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// 定义Agent类型到工具的映射
	agentToolMap := map[string][]string{
		"Science": {
			"simple_fact_lookup",
			"get_current_time",
			"image_generate_simple",
		},
		"Language": {
			"simple_dictionary",
			"pronunciation_hint",
			"get_current_time", // 添加时间工具，支持"几点了"这类问题
		},
		"Humanities": {
			// Humanities Agent可以不使用工具
		},
	}

	toolNames, ok := agentToolMap[agentType]
	if !ok {
		return []Tool{}
	}

	tools := make([]Tool, 0)
	for _, name := range toolNames {
		if tool, exists := r.tools[name]; exists {
			tools = append(tools, tool)
		}
	}

	return tools
}

// DefaultRegistry 全局默认工具注册表
var DefaultRegistry *ToolRegistry

// InitDefaultRegistry 初始化全局默认工具注册表
func InitDefaultRegistry(logger logx.Logger) {
	DefaultRegistry = NewToolRegistry(logger)
}

// GetDefaultRegistry 获取全局默认工具注册表
// 如果未初始化，会创建一个新的注册表
func GetDefaultRegistry(logger logx.Logger) *ToolRegistry {
	if DefaultRegistry == nil {
		InitDefaultRegistry(logger)
	}
	return DefaultRegistry
}

