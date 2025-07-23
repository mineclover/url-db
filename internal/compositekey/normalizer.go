package compositekey

import (
	"regexp"
	"strings"
)

import "url-db/internal/constants"

// constants 패키지에서 가져온 상수 사용
const (
	MaxToolNameLength   = constants.MaxToolNameLength
	MaxDomainNameLength = constants.MaxDomainNameLength
	MaxIDLength         = constants.MaxIDLength
)

// 특수문자를 하이픈으로 변환하는 정규표현식
var (
	// 영문자, 숫자, 하이픈, 언더스코어가 아닌 문자를 매칭
	invalidCharsRegex = regexp.MustCompile(`[^a-zA-Z0-9\-_]`)
	// 연속된 하이픈이나 언더스코어를 매칭
	multipleDelimiterRegex = regexp.MustCompile(`[-_]+`)
)

// NormalizeToolName 은 도구명을 정규화합니다.
func NormalizeToolName(toolName string) (string, error) {
	normalized := normalizeString(toolName)

	if len(normalized) == 0 {
		return "", NewInvalidToolNameError("도구명이 비어있습니다")
	}

	if len(normalized) > MaxToolNameLength {
		return "", NewTooLongError("도구명이 너무 깁니다")
	}

	return normalized, nil
}

// NormalizeDomainName 은 도메인명을 정규화합니다.
func NormalizeDomainName(domainName string) (string, error) {
	normalized := normalizeString(domainName)

	if len(normalized) == 0 {
		return "", NewInvalidDomainNameError("도메인명이 비어있습니다")
	}

	if len(normalized) > MaxDomainNameLength {
		return "", NewTooLongError("도메인명이 너무 깁니다")
	}

	return normalized, nil
}

// normalizeString 은 문자열을 정규화합니다.
func normalizeString(input string) string {
	if input == "" {
		return ""
	}

	// 1. 앞뒤 공백 제거
	normalized := strings.TrimSpace(input)

	// 2. 소문자로 변환
	normalized = strings.ToLower(normalized)

	// 3. 특수문자를 하이픈으로 변환
	normalized = invalidCharsRegex.ReplaceAllString(normalized, "-")

	// 4. 연속된 구분자를 단일 하이픈으로 변환
	normalized = multipleDelimiterRegex.ReplaceAllString(normalized, "-")

	// 5. 앞뒤 하이픈 제거
	normalized = strings.Trim(normalized, "-")

	return normalized
}

// CreateNormalized 는 정규화된 구성 요소로 합성키를 생성합니다.
func CreateNormalized(toolName, domainName string, id int) (CompositeKey, error) {
	normalizedToolName, err := NormalizeToolName(toolName)
	if err != nil {
		return CompositeKey{}, err
	}

	normalizedDomainName, err := NormalizeDomainName(domainName)
	if err != nil {
		return CompositeKey{}, err
	}

	if id <= 0 {
		return CompositeKey{}, NewInvalidIDError("ID는 양의 정수여야 합니다")
	}

	return CompositeKey{
		ToolName:   normalizedToolName,
		DomainName: normalizedDomainName,
		ID:         id,
	}, nil
}
