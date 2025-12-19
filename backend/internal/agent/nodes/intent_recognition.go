package nodes

import (
	"context"
	"strings"

	"github.com/tango/explore/internal/config"
	"github.com/zeromicro/go-zero/core/logx"
)

// IntentRecognitionNode 意图识别节点
type IntentRecognitionNode struct {
	ctx    context.Context
	config config.AIConfig
	logger logx.Logger
}

// IntentRecognitionResult 意图识别结果
type IntentRecognitionResult struct {
	Intent     string
	Confidence float64
	Reason     string
}

// NewIntentRecognitionNode 创建意图识别节点
func NewIntentRecognitionNode(ctx context.Context, cfg config.AIConfig, logger logx.Logger) (*IntentRecognitionNode, error) {
	return &IntentRecognitionNode{
		ctx:    ctx,
		config: cfg,
		logger: logger,
	}, nil
}

// Execute 执行意图识别
func (n *IntentRecognitionNode) Execute(data *GraphData, context []interface{}) (*IntentRecognitionResult, error) {
	n.logger.Infow("执行意图识别",
		logx.Field("message", data.Text),
		logx.Field("contextLength", len(context)),
	)

	// TODO: 待APP ID提供后，接入真实eino框架调用意图识别模型
	// 当前使用Mock数据（结合规则判断）
	return n.executeMock(data, context)
}

// executeMock Mock实现（待替换为真实eino调用）
func (n *IntentRecognitionNode) executeMock(data *GraphData, context []interface{}) (*IntentRecognitionResult, error) {
	message := strings.ToLower(data.Text)

	// 规则判断：如果包含生成卡片相关关键词，识别为generate_cards意图
	generateCardKeywords := []string{
		"生成", "卡片", "小卡片", "知识卡片",
		"生成卡片", "帮我生成", "我要卡片",
		"create", "card", "generate", "cards",
	}

	for _, keyword := range generateCardKeywords {
		if strings.Contains(message, keyword) {
			result := &IntentRecognitionResult{
				Intent:     "generate_cards",
				Confidence: 0.9,
				Reason:     "检测到生成卡片关键词: " + keyword,
			}
			n.logger.Infow("意图识别完成（Mock-规则）",
				logx.Field("intent", result.Intent),
				logx.Field("confidence", result.Confidence),
			)
			return result, nil
		}
	}

	// 默认返回文本回答意图
	result := &IntentRecognitionResult{
		Intent:     "text_response",
		Confidence: 0.8,
		Reason:     "未检测到生成卡片意图，默认文本回答",
	}

	n.logger.Infow("意图识别完成（Mock-默认）",
		logx.Field("intent", result.Intent),
		logx.Field("confidence", result.Confidence),
	)
	return result, nil
}

// executeReal 真实eino实现（待APP ID提供后实现）
func (n *IntentRecognitionNode) executeReal(data *GraphData, context []interface{}) (*IntentRecognitionResult, error) {
	// TODO: 使用eino框架调用意图识别模型
	// 1. 从config中获取IntentModel配置
	// 2. 使用eino的ChatModel接口
	// 3. 构建prompt，包含用户消息和上下文
	// 4. 调用模型API进行意图分类
	// 5. 解析返回结果，返回意图类型和置信度

	// 示例prompt结构：
	// "请识别以下用户消息的意图：
	//  1. generate_cards: 用户想要生成知识卡片
	//  2. text_response: 用户想要文本回答
	//
	//  用户消息: {message}
	//  上下文: {context}
	//
	//  请返回JSON格式: {\"intent\": \"...\", \"confidence\": 0.9, \"reason\": \"...\"}"

	return nil, nil
}
