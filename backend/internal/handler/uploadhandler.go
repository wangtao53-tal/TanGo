package handler

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/tango/explore/internal/logic"
	"github.com/tango/explore/internal/svc"
	"github.com/tango/explore/internal/types"
	"github.com/tango/explore/internal/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func UploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 检查Content-Type，支持两种方式：
		// 1. multipart/form-data（推荐，更高效）
		// 2. application/json（base64，向后兼容）
		contentType := r.Header.Get("Content-Type")
		
		var req types.UploadRequest
		var err error

		// 判断是multipart还是JSON
		if contentType != "" && len(contentType) >= 19 && contentType[:19] == "multipart/form-data" {
			// 方式1: multipart/form-data 文件上传
			err = parseMultipartUpload(r, &req)
		} else {
			// 方式2: application/json base64上传（向后兼容）
			err = httpx.Parse(r, &req)
		}

		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewUploadLogic(r.Context(), svcCtx)
		resp, err := l.Upload(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

// parseMultipartUpload 解析multipart/form-data上传
func parseMultipartUpload(r *http.Request, req *types.UploadRequest) error {
	// 解析multipart form，最大10MB
	err := r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		return utils.NewAPIError(400, "解析multipart表单失败", err.Error())
	}

	// 获取文件
	file, header, err := r.FormFile("image")
	if err != nil {
		// 如果没有文件，尝试从form字段获取base64（兼容）
		imageData := r.FormValue("imageData")
		if imageData != "" {
			req.ImageData = imageData
			req.Filename = r.FormValue("filename")
			return nil
		}
		return utils.NewAPIError(400, "未找到图片文件", "请使用字段名'image'上传文件")
	}
	defer file.Close()

	// 读取文件内容
	fileData, err := io.ReadAll(file)
	if err != nil {
		return utils.NewAPIError(400, "读取文件失败", err.Error())
	}

	// 验证文件大小
	maxSize := int64(10 * 1024 * 1024) // 10MB
	if int64(len(fileData)) > maxSize {
		return utils.ErrImageTooLarge
	}

	// 验证文件格式
	if !utils.IsValidImageFormat(fileData) {
		return utils.ErrImageFormatInvalid
	}

	// 将文件数据转换为base64（保持与现有逻辑兼容）
	req.ImageData = base64.StdEncoding.EncodeToString(fileData)
	
	// 获取文件名
	if header.Filename != "" {
		req.Filename = header.Filename
	} else {
		req.Filename = r.FormValue("filename")
	}

	logx.Infow("解析multipart上传成功",
		logx.Field("filename", req.Filename),
		logx.Field("size", len(fileData)),
	)

	return nil
}

// UploadStreamHandler 图片上传并流式返回Handler
func UploadStreamHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 设置SSE响应头
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")

		// 处理OPTIONS预检请求
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// 从请求体解析JSON参数（POST请求）
		var req types.UploadRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			// 发送错误事件
			errorEvent := types.StreamEvent{
				Type:    "error",
				Content: map[string]interface{}{"message": "请求参数解析失败: " + err.Error()},
			}
			errorJSON, _ := json.Marshal(errorEvent)
			fmt.Fprintf(w, "event: error\ndata: %s\n\n", string(errorJSON))
			w.(http.Flusher).Flush()
			return
		}

		// 参数验证
		if req.ImageData == "" {
			errorEvent := types.StreamEvent{
				Type:    "error",
				Content: map[string]interface{}{"message": "图片数据不能为空"},
			}
			errorJSON, _ := json.Marshal(errorEvent)
			fmt.Fprintf(w, "event: error\ndata: %s\n\n", string(errorJSON))
			w.(http.Flusher).Flush()
			return
		}

		// 图片上传
		uploadLogic := logic.NewUploadLogic(r.Context(), svcCtx)
		uploadResp, err := uploadLogic.Upload(&req)
		if err != nil {
			// 发送错误事件
			errorEvent := types.StreamEvent{
				Type:    "error",
				Content: map[string]interface{}{"message": "图片上传失败: " + err.Error()},
			}
			errorJSON, _ := json.Marshal(errorEvent)
			fmt.Fprintf(w, "event: error\ndata: %s\n\n", string(errorJSON))
			w.(http.Flusher).Flush()
			return
		}

		// 发送图片上传完成事件
		uploadedEvent := types.StreamEvent{
			Type:    "image_uploaded",
			Content: map[string]interface{}{"url": uploadResp.Url, "filename": uploadResp.Filename},
		}
		uploadedJSON, _ := json.Marshal(uploadedEvent)
		fmt.Fprintf(w, "event: image_uploaded\ndata: %s\n\n", string(uploadedJSON))
		w.(http.Flusher).Flush()

		// 注意：此handler已废弃，请使用统一接口 /api/conversation/stream
		// 调用统一流式接口（传入图片URL）
		streamLogic := logic.NewStreamLogic(r.Context(), svcCtx)
		streamReq := types.UnifiedStreamConversationRequest{
			MessageType: "image",
			Image:       uploadResp.Url,
			// 注意：UploadRequest中没有这些字段，需要从查询参数或其他地方获取
			// SessionId:             "",
			// Message:               "",
			// UserAge:               0,
			// MaxContextRounds:      20,
			// IdentificationContext: nil,
		}
		if err := streamLogic.StreamConversationUnified(w, streamReq); err != nil {
			// 错误已在StreamConversationUnified中处理
			return
		}
	}
}
