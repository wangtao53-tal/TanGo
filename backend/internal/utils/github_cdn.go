package utils

import (
	"regexp"
	"strings"
)

// IsGitHubRawURL 检测是否是GitHub raw URL
func IsGitHubRawURL(url string) bool {
	return strings.HasPrefix(url, "https://raw.githubusercontent.com/") ||
		strings.HasPrefix(url, "http://raw.githubusercontent.com/")
}

// ConvertToJSDelivrCDN 将GitHub raw URL转换为jsDelivr CDN URL
// 原始格式: https://raw.githubusercontent.com/{owner}/{repo}/{branch}/{path}
// CDN格式: https://cdn.jsdelivr.net/gh/{owner}/{repo}@{branch}/{path}
func ConvertToJSDelivrCDN(rawURL string) (string, error) {
	if !IsGitHubRawURL(rawURL) {
		return rawURL, nil // 不是GitHub raw URL，直接返回
	}

	// 移除协议前缀
	urlWithoutProtocol := strings.TrimPrefix(rawURL, "https://")
	urlWithoutProtocol = strings.TrimPrefix(urlWithoutProtocol, "http://")

	// 移除 raw.githubusercontent.com/ 前缀
	urlWithoutProtocol = strings.TrimPrefix(urlWithoutProtocol, "raw.githubusercontent.com/")

	// 分割路径: owner/repo/branch/path...
	parts := strings.SplitN(urlWithoutProtocol, "/", 4)
	if len(parts) < 4 {
		// URL格式不正确，返回原始URL
		return rawURL, nil
	}

	owner := parts[0]
	repo := parts[1]
	branch := parts[2]
	path := parts[3]

	// 构建jsDelivr CDN URL
	cdnURL := "https://cdn.jsdelivr.net/gh/" + owner + "/" + repo + "@" + branch + "/" + path

	return cdnURL, nil
}

// ExtractGitHubRawURLInfo 提取GitHub raw URL的信息（owner, repo, branch, path）
func ExtractGitHubRawURLInfo(rawURL string) (owner, repo, branch, path string, ok bool) {
	if !IsGitHubRawURL(rawURL) {
		return "", "", "", "", false
	}

	// 使用正则表达式提取
	re := regexp.MustCompile(`(?:https?://)?raw\.githubusercontent\.com/([^/]+)/([^/]+)/([^/]+)/(.+)`)
	matches := re.FindStringSubmatch(rawURL)
	if len(matches) != 5 {
		return "", "", "", "", false
	}

	return matches[1], matches[2], matches[3], matches[4], true
}
