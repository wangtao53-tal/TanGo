package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/tango/explore/internal/config"
	"github.com/zeromicro/go-zero/core/logx"
)

// MCPClient MCP客户端
type MCPClient struct {
	config config.MCPServerConfig
	logger logx.Logger
	client *http.Client
}

// NewMCPClient 创建MCP客户端
func NewMCPClient(cfg config.MCPServerConfig, logger logx.Logger) (*MCPClient, error) {
	client := &MCPClient{
		config: cfg,
		logger: logger,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	return client, nil
}

// CallResource 调用MCP资源
func (c *MCPClient) CallResource(ctx context.Context, resourceName string, params map[string]interface{}) (interface{}, error) {
	if c.config.Type != "url" || c.config.URL == "" {
		return nil, fmt.Errorf("不支持的MCP服务器类型: %s", c.config.Type)
	}

	// 构建请求
	reqBody := map[string]interface{}{
		"method": "resources/call",
		"params": map[string]interface{}{
			"resource": resourceName,
			"arguments": params,
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("序列化请求体失败: %w", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequestWithContext(ctx, "POST", c.config.URL, strings.NewReader(string(jsonBody)))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	for k, v := range c.config.Headers {
		req.Header.Set(k, v)
	}

	// TODO: 对于SSE端点，需要实现SSE连接处理
	// 当前实现简单的HTTP POST请求

	// 发送请求
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	// 解析响应
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	c.logger.Infow("MCP资源调用成功",
		logx.Field("resource", resourceName),
		logx.Field("status", resp.StatusCode),
	)

	return result, nil
}

// DiscoverResources 发现MCP服务器可用资源
func (c *MCPClient) DiscoverResources(ctx context.Context) ([]string, error) {
	// TODO: 实现资源发现逻辑
	// 当前返回空列表，后续可以调用MCP的list_resources方法
	return []string{}, nil
}

