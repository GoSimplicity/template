package errs

import (
	"errors"
)

// 定义统一的错误类型
var (
	// 业务错误
	ErrInvalidParameter     = errors.New("参数错误")
	ErrInvalidOperation     = errors.New("无效的操作")
	ErrDatabaseError        = errors.New("数据库错误")
	ErrExternalServiceError = errors.New("外部服务调用错误")
)
