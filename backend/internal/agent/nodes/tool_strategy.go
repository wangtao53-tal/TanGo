package nodes

import (
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
)

// ToolStrategy 工具使用策略
type ToolStrategy string

const (
	// ToolStrategyDirect 直接使用工具策略：高置信度问题，直接使用工具
	ToolStrategyDirect ToolStrategy = "direct"

	// ToolStrategyEnhance 增强策略：探索性问题，先回答，再提供工具增强
	ToolStrategyEnhance ToolStrategy = "enhance"

	// ToolStrategyNone 不使用工具策略：简单问题，不使用工具
	ToolStrategyNone ToolStrategy = "none"

	// ToolStrategyMultiple 多工具策略：复杂问题，使用多个工具
	ToolStrategyMultiple ToolStrategy = "multiple"
)

// SelectToolStrategy 根据意图和问题内容选择工具使用策略
func SelectToolStrategy(intent string, message string, confidence float64) ToolStrategy {
	// 高置信度问题（置信度≥0.8）：直接使用工具
	if confidence >= 0.8 {
		return ToolStrategyDirect
	}

	// 探索性问题关键词：增强策略
	exploratoryKeywords := []string{"为什么", "怎么", "如何", "是什么", "什么是", "能不能", "会不会"}
	messageLower := strings.ToLower(message)
	for _, keyword := range exploratoryKeywords {
		if strings.Contains(messageLower, keyword) {
			return ToolStrategyEnhance
		}
	}

	// 复杂问题关键词：多工具策略
	complexKeywords := []string{"详细", "深入", "全面", "完整", "所有", "全部"}
	for _, keyword := range complexKeywords {
		if strings.Contains(messageLower, keyword) {
			return ToolStrategyMultiple
		}
	}

	// 简单问题：不使用工具
	return ToolStrategyNone
}

// SelectToolsForIntent 根据意图类型选择推荐工具
func SelectToolsForIntent(intent string, domainAgent string) []string {
	tools := []string{"get_current_time"}

	// 根据意图类型选择工具
	switch intent {
	case "认知型":
		// 认知型问题：需要事实查询
		if domainAgent == "Science" {
			tools = append(tools, "simple_fact_lookup")
		}
	case "探因型":
		// 探因型问题：需要深入查询
		if domainAgent == "Science" {
			tools = append(tools, "simple_fact_lookup")
		}
	case "表达型":
		// 表达型问题：需要语言工具
		if domainAgent == "Language" {
			tools = append(tools, "simple_dictionary", "pronunciation_hint")
		}
	case "游戏型":
		// 游戏型问题：可能需要时间或图片
		if domainAgent == "Science" {
			tools = append(tools, "get_current_time", "image_generate_simple")
		}
	case "情绪型":
		// 情绪型问题：不使用工具，直接回答
		return []string{}
	}

	// 根据领域Agent添加通用工具
	if domainAgent == "Science" && len(tools) == 0 {
		// Science Agent默认工具
		tools = append(tools, "simple_fact_lookup")
	} else if domainAgent == "Language" && len(tools) == 0 {
		// Language Agent默认工具
		tools = append(tools, "simple_dictionary")
	}

	return tools
}

// SelectToolsByKeywords 根据问题关键词选择工具
func SelectToolsByKeywords(message string, domainAgent string) []string {
	tools := []string{}
	messageLower := strings.ToLower(message)

	// 时间相关关键词
	timeKeywords := []string{"时间", "现在", "今天", "几点", "什么时候", "日期"}
	for _, keyword := range timeKeywords {
		if strings.Contains(messageLower, keyword) {
			tools = append(tools, "get_current_time")
			break
		}
	}

	// 图片相关关键词
	imageKeywords := []string{"图片", "图像", "示意图", "画", "图", "看"}
	for _, keyword := range imageKeywords {
		if strings.Contains(messageLower, keyword) {
			if domainAgent == "Science" {
				tools = append(tools, "image_generate_simple")
			}
			break
		}
	}

	// 单词/语言相关关键词
	languageKeywords := []string{"单词", "英语", "怎么说", "发音", "意思", "意思是什么"}
	for _, keyword := range languageKeywords {
		if strings.Contains(messageLower, keyword) {
			if domainAgent == "Language" {
				tools = append(tools, "simple_dictionary", "pronunciation_hint")
			}
			break
		}
	}

	// 事实查询关键词
	factKeywords := []string{"是什么", "什么是", "介绍", "了解", "知道"}
	for _, keyword := range factKeywords {
		if strings.Contains(messageLower, keyword) {
			if domainAgent == "Science" {
				tools = append(tools, "simple_fact_lookup")
			}
			break
		}
	}

	return tools
}

// SelectTools 综合选择工具（根据意图、问题内容和领域Agent）
func SelectTools(intent string, message string, domainAgent string, confidence float64, logger logx.Logger) ([]string, ToolStrategy) {
	// 选择工具使用策略
	strategy := SelectToolStrategy(intent, message, confidence)

	// 如果策略是不使用工具，直接返回
	if strategy == ToolStrategyNone {
		logger.Infow("工具选择：不使用工具",
			logx.Field("intent", intent),
			logx.Field("strategy", strategy),
		)
		return []string{}, strategy
	}

	// 根据意图选择工具
	intentTools := SelectToolsForIntent(intent, domainAgent)

	// 根据关键词选择工具
	keywordTools := SelectToolsByKeywords(message, domainAgent)

	// 合并工具列表（去重）
	toolMap := make(map[string]bool)
	allTools := []string{}

	for _, tool := range intentTools {
		if !toolMap[tool] {
			toolMap[tool] = true
			allTools = append(allTools, tool)
		}
	}

	for _, tool := range keywordTools {
		if !toolMap[tool] {
			toolMap[tool] = true
			allTools = append(allTools, tool)
		}
	}

	logger.Infow("🎯 工具选择完成",
		logx.Field("intent", intent),
		logx.Field("domainAgent", domainAgent),
		logx.Field("strategy", strategy),
		logx.Field("selected_tools", allTools),
		logx.Field("tool_count", len(allTools)),
		logx.Field("intent_tools", intentTools),
		logx.Field("keyword_tools", keywordTools),
	)

	return allTools, strategy
}
