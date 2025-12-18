package utils

import (
	"fmt"
	"net/http"
)

// APIError 表示API错误
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail,optional"`
}

// Error 实现error接口
func (e *APIError) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("code: %d, message: %s, detail: %s", e.Code, e.Message, e.Detail)
	}
	return fmt.Sprintf("code: %d, message: %s", e.Code, e.Message)
}

// NewAPIError 创建新的API错误
func NewAPIError(code int, message string, detail ...string) *APIError {
	err := &APIError{
		Code:    code,
		Message: message,
	}
	if len(detail) > 0 {
		err.Detail = detail[0]
	}
	return err
}

// 预定义错误
var (
	ErrInvalidRequest    = NewAPIError(http.StatusBadRequest, "请求参数无效")
	ErrImageRequired     = NewAPIError(http.StatusBadRequest, "图片数据不能为空")
	ErrInvalidImage      = NewAPIError(http.StatusBadRequest, "图片格式无效")
	ErrAgeRequired       = NewAPIError(http.StatusBadRequest, "年龄参数必填")
	ErrInvalidAge        = NewAPIError(http.StatusBadRequest, "年龄必须在3-18之间")
	ErrObjectNameRequired = NewAPIError(http.StatusBadRequest, "对象名称不能为空")
	ErrCategoryRequired   = NewAPIError(http.StatusBadRequest, "对象类别不能为空")
	ErrShareNotFound      = NewAPIError(http.StatusNotFound, "分享链接不存在或已过期")
	ErrInternalServer     = NewAPIError(http.StatusInternalServerError, "服务器内部错误")
)

