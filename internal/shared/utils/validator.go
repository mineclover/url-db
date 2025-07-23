package compositekey

import (
	"regexp"
	"strconv"
	"strings"
)

// 검증에 사용되는 정규표현식
var (
	// 유효한 문자 (영문자, 숫자, 하이픈, 언더스코어)
	validCharsRegex = regexp.MustCompile(`^[a-zA-Z0-9\-_]+$`)
)

// ValidateFormat 은 합성키의 기본 형식을 검증합니다.
func ValidateFormat(compositeKey string) error {
	if compositeKey == "" {
		return NewInvalidFormatError("합성키가 비어있습니다")
	}

	parts := strings.Split(compositeKey, ":")
	if len(parts) != 3 {
		return NewInvalidFormatError("합성키는 정확히 3개의 구성 요소를 가져야 합니다")
	}

	return nil
}

// ValidateToolName 은 도구명을 검증합니다.
func ValidateToolName(toolName string) error {
	if toolName == "" {
		return NewInvalidToolNameError("도구명이 비어있습니다")
	}

	if len(toolName) > MaxToolNameLength {
		return NewTooLongError("도구명이 너무 깁니다")
	}

	if !validCharsRegex.MatchString(toolName) {
		return NewInvalidToolNameError("도구명에 유효하지 않은 문자가 포함되어 있습니다")
	}

	// 하이픈으로 시작하거나 끝나면 안됨
	if strings.HasPrefix(toolName, "-") || strings.HasSuffix(toolName, "-") {
		return NewInvalidToolNameError("도구명은 하이픈으로 시작하거나 끝날 수 없습니다")
	}

	return nil
}

// ValidateDomainName 은 도메인명을 검증합니다.
func ValidateDomainName(domainName string) error {
	if domainName == "" {
		return NewInvalidDomainNameError("도메인명이 비어있습니다")
	}

	if len(domainName) > MaxDomainNameLength {
		return NewTooLongError("도메인명이 너무 깁니다")
	}

	if !validCharsRegex.MatchString(domainName) {
		return NewInvalidDomainNameError("도메인명에 유효하지 않은 문자가 포함되어 있습니다")
	}

	// 하이픈으로 시작하거나 끝나면 안됨
	if strings.HasPrefix(domainName, "-") || strings.HasSuffix(domainName, "-") {
		return NewInvalidDomainNameError("도메인명은 하이픈으로 시작하거나 끝날 수 없습니다")
	}

	return nil
}

// ValidateID 는 ID를 검증합니다.
func ValidateID(idStr string) error {
	if idStr == "" {
		return NewInvalidIDError("ID가 비어있습니다")
	}

	if len(idStr) > MaxIDLength {
		return NewTooLongError("ID가 너무 깁니다")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return NewInvalidIDError("ID는 유효한 정수여야 합니다")
	}

	if id <= 0 {
		return NewInvalidIDError("ID는 양의 정수여야 합니다")
	}

	return nil
}

// ValidateCompositeKey 는 합성키 전체를 검증합니다.
func ValidateCompositeKey(compositeKey string) error {
	// 1. 기본 형식 검증
	if err := ValidateFormat(compositeKey); err != nil {
		return err
	}

	// 2. 구성 요소 분해
	parts := strings.Split(compositeKey, ":")
	toolName := parts[0]
	domainName := parts[1]
	idStr := parts[2]

	// 3. 각 구성 요소 검증
	if err := ValidateToolName(toolName); err != nil {
		return err
	}

	if err := ValidateDomainName(domainName); err != nil {
		return err
	}

	if err := ValidateID(idStr); err != nil {
		return err
	}

	return nil
}

// ValidateCompositeKeyStruct 는 CompositeKey 구조체를 검증합니다.
func ValidateCompositeKeyStruct(ck CompositeKey) error {
	if err := ValidateToolName(ck.ToolName); err != nil {
		return err
	}

	if err := ValidateDomainName(ck.DomainName); err != nil {
		return err
	}

	if ck.ID <= 0 {
		return NewInvalidIDError("ID는 양의 정수여야 합니다")
	}

	return nil
}
