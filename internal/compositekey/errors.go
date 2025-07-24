package compositekey

import "fmt"

// 에러 코드 정의
const (
	ErrInvalidFormat     = "COMPOSITE_KEY_INVALID_FORMAT"
	ErrInvalidToolName   = "COMPOSITE_KEY_INVALID_TOOL_NAME"
	ErrInvalidDomainName = "COMPOSITE_KEY_INVALID_DOMAIN_NAME"
	ErrInvalidID         = "COMPOSITE_KEY_INVALID_ID"
)

// CompositeKeyError 는 합성키 관련 에러를 나타냅니다.
type CompositeKeyError struct {
	Code    string
	Message string
}

func (e CompositeKeyError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// 에러 생성 함수들
func NewInvalidFormatError(message string) error {
	return CompositeKeyError{
		Code:    ErrInvalidFormat,
		Message: message,
	}
}

func NewInvalidToolNameError(message string) error {
	return CompositeKeyError{
		Code:    ErrInvalidToolName,
		Message: message,
	}
}

func NewInvalidDomainNameError(message string) error {
	return CompositeKeyError{
		Code:    ErrInvalidDomainName,
		Message: message,
	}
}

func NewInvalidIDError(message string) error {
	return CompositeKeyError{
		Code:    ErrInvalidID,
		Message: message,
	}
}
