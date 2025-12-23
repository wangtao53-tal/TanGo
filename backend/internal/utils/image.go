package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// CleanBase64String 清理 base64 字符串，移除所有空白字符
// 这可以解决传输过程中可能引入的空格、换行符等问题
func CleanBase64String(s string) string {
	// 移除所有空白字符（空格、换行符、制表符等）
	// 使用 strings.ReplaceAll 比 strings.Replace 更高效
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\r", "")
	s = strings.ReplaceAll(s, "\t", "")
	return s
}

// ValidateBase64Image 验证 base64 图片数据
func ValidateBase64Image(imageData string, maxSize int64) error {
	// 检查是否为空
	if imageData == "" {
		return ErrImageDataRequired
	}

	// 清理 base64 字符串，移除可能存在的空白字符
	imageData = CleanBase64String(imageData)
	
	// 清理后再次检查是否为空（可能清理后变成空字符串）
	if imageData == "" {
		return ErrImageDataRequired
	}

	// 检查 base64 字符串长度（base64 编码后大小约为原文件的 4/3）
	if maxSize > 0 && int64(len(imageData)) > maxSize*4/3 {
		return ErrImageTooLarge
	}

	// 尝试解码 base64
	decoded, err := base64.StdEncoding.DecodeString(imageData)
	if err != nil {
		return NewAPIError(400, "Base64 解码失败", err.Error())
	}

	// 检查解码后的大小
	if maxSize > 0 && int64(len(decoded)) > maxSize {
		return ErrImageTooLarge
	}

	// 验证图片格式（通过文件头 Magic Number）
	if !isValidImageFormat(decoded) {
		return ErrImageFormatInvalid
	}

	return nil
}

// isValidImageFormat 通过文件头验证图片格式
func isValidImageFormat(data []byte) bool {
	if len(data) < 4 {
		return false
	}

	// JPEG: FF D8 FF
	if len(data) >= 3 && data[0] == 0xFF && data[1] == 0xD8 && data[2] == 0xFF {
		return true
	}

	// PNG: 89 50 4E 47
	if len(data) >= 4 && data[0] == 0x89 && data[1] == 0x50 && data[2] == 0x4E && data[3] == 0x47 {
		return true
	}

	// WebP: 需要检查 RIFF 头和 WEBP
	if len(data) >= 12 {
		// RIFF header: 52 49 46 46
		if data[0] == 0x52 && data[1] == 0x49 && data[2] == 0x46 && data[3] == 0x46 {
			// WEBP: 57 45 42 50 (在偏移 8 的位置)
			if len(data) >= 12 && data[8] == 0x57 && data[9] == 0x45 && data[10] == 0x42 && data[11] == 0x50 {
				return true
			}
		}
	}

	// GIF: 47 49 46 38 (GIF8)
	if len(data) >= 6 && data[0] == 0x47 && data[1] == 0x49 && data[2] == 0x46 && data[3] == 0x38 {
		return true
	}

	return false
}

// ValidateFilename 验证文件名安全性
func ValidateFilename(filename string) error {
	if filename == "" {
		return nil // 空文件名是允许的（后端会自动生成）
	}

	// 检查路径遍历攻击
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") || strings.Contains(filename, "\\") {
		return ErrFilenameInvalid
	}

	// 只允许字母、数字、连字符、下划线和点
	matched, err := regexp.MatchString(`^[a-zA-Z0-9._-]+$`, filename)
	if err != nil {
		return NewAPIError(500, "文件名验证失败", err.Error())
	}

	if !matched {
		return ErrFilenameInvalid
	}

	// 检查文件扩展名
	ext := strings.ToLower(getFileExtension(filename))
	allowedExts := []string{".jpg", ".jpeg", ".png", ".webp", ".gif"}
	for _, allowedExt := range allowedExts {
		if ext == allowedExt {
			return nil
		}
	}

	return ErrFilenameInvalid
}

// getFileExtension 获取文件扩展名
func getFileExtension(filename string) string {
	idx := strings.LastIndex(filename, ".")
	if idx == -1 || idx == len(filename)-1 {
		return ""
	}
	return filename[idx:]
}

// GenerateFilename 生成唯一文件名
func GenerateFilename(originalExt string) string {
	// 使用时间戳 + 随机字符串
	// 简化版本：使用时间戳 + 简单随机数
	timestamp := fmt.Sprintf("%d", getCurrentTimestamp())
	random := generateRandomString(8)

	// 确保扩展名有效
	ext := strings.ToLower(originalExt)
	if ext == "" {
		ext = ".jpg" // 默认 JPEG
	}
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}

	return fmt.Sprintf("%s-%s%s", timestamp, random, ext)
}

// getCurrentTimestamp 获取当前时间戳（毫秒）
func getCurrentTimestamp() int64 {
	return time.Now().UnixMilli()
}

// generateRandomString 生成随机字符串
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	randomBytes := make([]byte, length)
	if _, err := rand.Read(randomBytes); err != nil {
		// 如果随机数生成失败，使用时间戳作为后备
		timestamp := time.Now().UnixNano()
		for i := range result {
			result[i] = charset[int(timestamp)%len(charset)]
			timestamp = timestamp / int64(len(charset))
		}
		return string(result)
	}
	for i := range result {
		result[i] = charset[randomBytes[i]%byte(len(charset))]
	}
	return string(result)
}
