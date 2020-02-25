package cloud_provider

import (
	"context"
	"fmt"

	uuid "github.com/satori/go.uuid"
)

// ContextKeyType for context.WithValue(
type ContextKeyType string

const (
	// RequestID is Context Key for trace Alert
	RequestID ContextKeyType = "RequestID"
)

// GetRandom 返回 64 位随机字符
func GetRandom() string {
	uuid := uuid.NewV4()
	return uuid.String()
}

// Message 返回打标的信息
func Message(ctx context.Context, msg string) string {
	requestID := ""

	if ctx.Value(RequestID) != nil {
		requestID = ctx.Value(RequestID).(string)
	}

	return fmt.Sprintf("[ReqID:%s] %s", requestID, msg)
}

