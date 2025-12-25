package utils

import "strings"

// DetectMarkdown 检测文本是否包含Markdown格式
func DetectMarkdown(text string) bool {
	// 检测常见的Markdown标记
	markdownPatterns := []string{
		"```", // 代码块
		"##",  // 标题
		"###", // 标题
		"- ",  // 无序列表
		"* ",  // 无序列表
		"+ ",  // 无序列表
		"1. ", // 有序列表
		"[",   // 链接开始
		"](",  // 链接
		"**",  // 粗体
		"__",  // 粗体
		"*",   // 斜体
		"_",   // 斜体
		"> ",  // 引用
		"|",   // 表格
		"---", // 分隔线
		"===", // 分隔线
	}

	textLower := strings.ToLower(text)
	for _, pattern := range markdownPatterns {
		if strings.Contains(textLower, pattern) {
			return true
		}
	}

	return false
}
