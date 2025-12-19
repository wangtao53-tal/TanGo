package logic

import (
	"context"
	"encoding/base64"

	"github.com/google/uuid"
	"github.com/tango/explore/internal/svc"
	"github.com/tango/explore/internal/types"
)

type VoiceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVoiceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VoiceLogic {
	return &VoiceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// RecognizeVoice 识别语音
func (l *VoiceLogic) RecognizeVoice(req *types.VoiceRequest) (*types.VoiceResponse, error) {
	// 解码音频数据
	_, err := base64.StdEncoding.DecodeString(req.Audio)
	if err != nil {
		return nil, err
	}

	// TODO: 调用语音识别模型（通过eino框架）
	// 当前使用Mock数据
	recognizedText := "这是Mock语音识别结果。待接入真实语音识别模型后，将实现真实的语音转文本功能。"

	// 生成或使用现有会话ID
	sessionId := req.SessionId
	if sessionId == "" {
		sessionId = uuid.New().String()
	}

	// 可选：调用意图识别（当前未使用结果，待后续实现）
	// intentLogic := NewIntentLogic(l.ctx, l.svcCtx)
	// intentReq := &types.IntentRequest{
	// 	Message:   recognizedText,
	// 	SessionId: sessionId,
	// }
	// _, _ = intentLogic.RecognizeIntent(intentReq)

	return &types.VoiceResponse{
		Text:      recognizedText,
		SessionId: sessionId,
	}, nil
}
