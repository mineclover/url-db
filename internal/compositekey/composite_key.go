package compositekey

import (
	"strconv"
	"strings"
)

// CompositeKey 는 합성키 구조를 나타냅니다.
type CompositeKey struct {
	ToolName   string `json:"tool_name"`
	DomainName string `json:"domain_name"`
	ID         int    `json:"id"`
}

// String 은 CompositeKey를 문자열로 변환합니다.
func (ck CompositeKey) String() string {
	return strings.Join([]string{ck.ToolName, ck.DomainName, strconv.Itoa(ck.ID)}, ":")
}

// Create 는 주어진 구성 요소로 합성키를 생성합니다.
func Create(toolName, domainName string, id int) CompositeKey {
	return CompositeKey{
		ToolName:   toolName,
		DomainName: domainName,
		ID:         id,
	}
}

// Parse 는 합성키 문자열을 파싱하여 CompositeKey 구조체로 변환합니다.
func Parse(compositeKey string) (CompositeKey, error) {
	parts := strings.Split(compositeKey, ":")
	if len(parts) != 3 {
		return CompositeKey{}, NewInvalidFormatError("합성키는 정확히 3개의 구성 요소를 가져야 합니다")
	}

	toolName := parts[0]
	domainName := parts[1]
	idStr := parts[2]

	// ID를 정수로 변환
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return CompositeKey{}, NewInvalidIDError("ID는 유효한 정수여야 합니다")
	}

	if id <= 0 {
		return CompositeKey{}, NewInvalidIDError("ID는 양의 정수여야 합니다")
	}

	return CompositeKey{
		ToolName:   toolName,
		DomainName: domainName,
		ID:         id,
	}, nil
}

// IsValid 는 합성키의 유효성을 검사합니다.
func IsValid(compositeKey string) bool {
	_, err := Parse(compositeKey)
	return err == nil
}
