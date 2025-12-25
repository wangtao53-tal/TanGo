package base

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/tango/explore/internal/tools"
	"github.com/zeromicro/go-zero/core/logx"
)

// PronunciationHintTool pronunciation_hint工具实现
// 用于提供发音提示，主要用于Language Agent
type PronunciationHintTool struct {
	logger logx.Logger
}

// NewPronunciationHintTool 创建pronunciation_hint工具实例
func NewPronunciationHintTool(logger logx.Logger) tools.Tool {
	return &PronunciationHintTool{
		logger: logger,
	}
}

// Name 返回工具名称
func (t *PronunciationHintTool) Name() string {
	return "pronunciation_hint"
}

// Description 返回工具描述
func (t *PronunciationHintTool) Description() string {
	return "提供单词发音提示，帮助孩子学习正确发音。输入单词，返回发音提示和发音技巧。"
}

// Parameters 返回工具参数定义（JSON Schema格式）
func (t *PronunciationHintTool) Parameters() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"word": map[string]interface{}{
				"type":        "string",
				"description": "要查询发音的单词",
			},
		},
		"required": []string{"word"},
	}
}

// Execute 执行工具
func (t *PronunciationHintTool) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	// 设置超时
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	// 提取参数
	word, ok := params["word"].(string)
	if !ok {
		return nil, fmt.Errorf("参数word必须是字符串类型")
	}

	if word == "" {
		return nil, fmt.Errorf("单词不能为空")
	}

	t.logger.Infow("执行pronunciation_hint工具",
		logx.Field("word", word),
	)

	// TODO: 后续可以接入真实的发音API
	// 当前使用Mock实现
	result := t.getPronunciationHint(word)

	return result, nil
}

// getPronunciationHint 获取发音提示（Mock实现）
func (t *PronunciationHintTool) getPronunciationHint(word string) map[string]interface{} {
	// Mock发音数据库
	pronunciations := map[string]map[string]interface{}{
		"apple": {
			"word":        "apple",
			"phonetic":    "/ˈæpl/",
			"hint":        "读作'艾-普-欧'，注意'æ'的发音像'艾'",
			"tip":         "可以分成两个音节：ap-ple",
		},
		"book": {
			"word":        "book",
			"phonetic":    "/bʊk/",
			"hint":        "读作'布-克'，注意'ʊ'的发音",
			"tip":         "只有一个音节，发音简单",
		},
		"happy": {
			"word":        "happy",
			"phonetic":    "/ˈhæpi/",
			"hint":        "读作'哈-皮'，注意重音在第一个音节",
			"tip":         "可以分成两个音节：hap-py",
		},
	}

	// 查找匹配的单词（不区分大小写）
	wordLower := strings.ToLower(word)
	if entry, ok := pronunciations[wordLower]; ok {
		return entry
	}

	// 如果没有找到匹配的单词，返回通用回答
	return map[string]interface{}{
		"word":     word,
		"phonetic": "/.../",
		"hint":     fmt.Sprintf("'%s'的发音需要多练习。", word),
		"tip":      "可以尝试慢慢读，注意每个音节的发音。",
	}
}

