package base

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/tango/explore/internal/tools"
	"github.com/zeromicro/go-zero/core/logx"
)

// SimpleFactLookupTool simple_fact_lookup工具实现
// 用于查找简单事实，主要用于Science Agent
type SimpleFactLookupTool struct {
	logger logx.Logger
}

// NewSimpleFactLookupTool 创建simple_fact_lookup工具实例
func NewSimpleFactLookupTool(logger logx.Logger) tools.Tool {
	return &SimpleFactLookupTool{
		logger: logger,
	}
}

// Name 返回工具名称
func (t *SimpleFactLookupTool) Name() string {
	return "simple_fact_lookup"
}

// Description 返回工具描述
func (t *SimpleFactLookupTool) Description() string {
	return "查找简单事实，用于科学知识查询。输入查询关键词，返回相关的事实信息。"
}

// Parameters 返回工具参数定义（JSON Schema格式）
func (t *SimpleFactLookupTool) Parameters() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"query": map[string]interface{}{
				"type":        "string",
				"description": "查询关键词，例如：'太阳'、'水循环'、'光合作用'",
			},
		},
		"required": []string{"query"},
	}
}

// Execute 执行工具
func (t *SimpleFactLookupTool) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	// 设置超时
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	// 提取参数
	query, ok := params["query"].(string)
	if !ok {
		return nil, fmt.Errorf("参数query必须是字符串类型")
	}

	if query == "" {
		return nil, fmt.Errorf("查询关键词不能为空")
	}

	t.logger.Infow("执行simple_fact_lookup工具",
		logx.Field("query", query),
	)

	// TODO: 后续可以接入真实的事实查询API
	// 当前使用Mock实现
	result := t.lookupFact(query)

	return result, nil
}

// lookupFact 查找事实（Mock实现）
func (t *SimpleFactLookupTool) lookupFact(query string) map[string]interface{} {
	// Mock事实数据库
	facts := map[string]string{
		"太阳":     "太阳是太阳系的中心恒星，是一颗黄矮星，距离地球约1.5亿公里。",
		"水循环":   "水循环是地球上水在大气、海洋和陆地之间循环的过程，包括蒸发、凝结、降水等环节。",
		"光合作用": "光合作用是植物利用阳光、水和二氧化碳制造氧气和葡萄糖的过程。",
		"重力":     "重力是地球对物体的吸引力，使物体向下落。",
		"磁铁":     "磁铁有南北两极，同极相斥，异极相吸。",
	}

	// 查找匹配的事实
	for key, fact := range facts {
		if key == query || contains(query, key) {
			return map[string]interface{}{
				"query": query,
				"fact":  fact,
				"source": "知识库",
			}
		}
	}

	// 如果没有找到匹配的事实，返回通用回答
	return map[string]interface{}{
		"query": query,
		"fact":  fmt.Sprintf("关于'%s'，这是一个很有趣的话题。", query),
		"source": "知识库",
	}
}

// contains 检查字符串是否包含子串（简单实现）
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0)
}

// ToJSON 将结果转换为JSON字符串（用于工具结果返回）
func (t *SimpleFactLookupTool) ToJSON(result interface{}) string {
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return fmt.Sprintf("%v", result)
	}
	return string(jsonBytes)
}

