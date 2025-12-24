package base

import (
	"context"
	"fmt"
	"time"

	"github.com/tango/explore/internal/tools"
	"github.com/zeromicro/go-zero/core/logx"
)

// SimpleDictionaryTool simple_dictionary工具实现
// 用于查找单词释义，主要用于Language Agent
type SimpleDictionaryTool struct {
	logger logx.Logger
}

// NewSimpleDictionaryTool 创建simple_dictionary工具实例
func NewSimpleDictionaryTool(logger logx.Logger) tools.Tool {
	return &SimpleDictionaryTool{
		logger: logger,
	}
}

// Name 返回工具名称
func (t *SimpleDictionaryTool) Name() string {
	return "simple_dictionary"
}

// Description 返回工具描述
func (t *SimpleDictionaryTool) Description() string {
	return "查找单词释义和例句，用于语言学习。输入单词，返回释义、例句和用法。"
}

// Parameters 返回工具参数定义（JSON Schema格式）
func (t *SimpleDictionaryTool) Parameters() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"word": map[string]interface{}{
				"type":        "string",
				"description": "要查询的单词，例如：'apple'、'book'、'happy'",
			},
		},
		"required": []string{"word"},
	}
}

// Execute 执行工具
func (t *SimpleDictionaryTool) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
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

	t.logger.Infow("执行simple_dictionary工具",
		logx.Field("word", word),
	)

	// TODO: 后续可以接入真实的词典API
	// 当前使用Mock实现
	result := t.lookupWord(word)

	return result, nil
}

// lookupWord 查找单词（Mock实现）
func (t *SimpleDictionaryTool) lookupWord(word string) map[string]interface{} {
	// Mock词典数据库
	dictionary := map[string]map[string]interface{}{
		"apple": {
			"word":        "apple",
			"pronunciation": "/ˈæpl/",
			"meaning":     "苹果",
			"example":     "I like to eat an apple every day.",
			"example_cn":  "我每天喜欢吃一个苹果。",
		},
		"book": {
			"word":        "book",
			"pronunciation": "/bʊk/",
			"meaning":     "书",
			"example":     "I read a book before bed.",
			"example_cn":  "我睡前读一本书。",
		},
		"happy": {
			"word":        "happy",
			"pronunciation": "/ˈhæpi/",
			"meaning":     "快乐的",
			"example":     "I am happy to see you.",
			"example_cn":  "我很高兴见到你。",
		},
	}

	// 查找匹配的单词（不区分大小写）
	wordLower := toLower(word)
	if entry, ok := dictionary[wordLower]; ok {
		return entry
	}

	// 如果没有找到匹配的单词，返回通用回答
	return map[string]interface{}{
		"word":        word,
		"pronunciation": "/.../",
		"meaning":     fmt.Sprintf("'%s'是一个英语单词。", word),
		"example":     fmt.Sprintf("Can you use '%s' in a sentence?", word),
		"example_cn":  fmt.Sprintf("你能用'%s'造一个句子吗？", word),
	}
}

// toLower 简单的转小写函数
func toLower(s string) string {
	result := ""
	for _, r := range s {
		if r >= 'A' && r <= 'Z' {
			result += string(r + 32)
		} else {
			result += string(r)
		}
	}
	return result
}

