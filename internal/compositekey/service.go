package compositekey

import "strings"

// Service 는 합성키 서비스를 나타냅니다.
type Service struct {
	defaultToolName string
}

// NewService 는 새로운 합성키 서비스를 생성합니다.
func NewService(defaultToolName string) *Service {
	return &Service{
		defaultToolName: defaultToolName,
	}
}

// Create 는 주어진 구성 요소로 합성키를 생성합니다.
func (s *Service) Create(domainName string, id int) (string, error) {
	compositeKey, err := CreateNormalized(s.defaultToolName, domainName, id)
	if err != nil {
		return "", err
	}
	
	return compositeKey.String(), nil
}

// CreateWithTool 은 도구명을 포함하여 합성키를 생성합니다.
func (s *Service) CreateWithTool(toolName, domainName string, id int) (string, error) {
	compositeKey, err := CreateNormalized(toolName, domainName, id)
	if err != nil {
		return "", err
	}
	
	return compositeKey.String(), nil
}

// Parse 는 합성키 문자열을 파싱하여 CompositeKey 구조체로 변환합니다.
func (s *Service) Parse(compositeKey string) (CompositeKey, error) {
	// 먼저 검증
	if err := ValidateCompositeKey(compositeKey); err != nil {
		return CompositeKey{}, err
	}
	
	// 파싱
	return Parse(compositeKey)
}

// Validate 는 합성키의 유효성을 검증합니다.
func (s *Service) Validate(compositeKey string) bool {
	return ValidateCompositeKey(compositeKey) == nil
}

// ParseComponents 는 합성키를 구성 요소로 분해하여 반환합니다.
func (s *Service) ParseComponents(compositeKey string) (toolName, domainName string, id int, err error) {
	ck, err := s.Parse(compositeKey)
	if err != nil {
		return "", "", 0, err
	}
	
	return ck.ToolName, ck.DomainName, ck.ID, nil
}

// GetToolName 은 합성키에서 도구명을 추출합니다.
func (s *Service) GetToolName(compositeKey string) (string, error) {
	parts := strings.Split(compositeKey, ":")
	if len(parts) != 3 {
		return "", NewInvalidFormatError("합성키 형식이 올바르지 않습니다")
	}
	
	return parts[0], nil
}

// GetDomainName 은 합성키에서 도메인명을 추출합니다.
func (s *Service) GetDomainName(compositeKey string) (string, error) {
	parts := strings.Split(compositeKey, ":")
	if len(parts) != 3 {
		return "", NewInvalidFormatError("합성키 형식이 올바르지 않습니다")
	}
	
	return parts[1], nil
}

// GetID 는 합성키에서 ID를 추출합니다.
func (s *Service) GetID(compositeKey string) (int, error) {
	ck, err := s.Parse(compositeKey)
	if err != nil {
		return 0, err
	}
	
	return ck.ID, nil
}

// IsValidFormat 은 합성키의 기본 형식만 검증합니다.
func (s *Service) IsValidFormat(compositeKey string) bool {
	return ValidateFormat(compositeKey) == nil
}

// NormalizeComponents 는 구성 요소를 정규화합니다.
func (s *Service) NormalizeComponents(toolName, domainName string) (string, string, error) {
	normalizedToolName, err := NormalizeToolName(toolName)
	if err != nil {
		return "", "", err
	}
	
	normalizedDomainName, err := NormalizeDomainName(domainName)
	if err != nil {
		return "", "", err
	}
	
	return normalizedToolName, normalizedDomainName, nil
}

// ValidateComponents 는 구성 요소들을 개별적으로 검증합니다.
func (s *Service) ValidateComponents(toolName, domainName string, id int) error {
	if err := ValidateToolName(toolName); err != nil {
		return err
	}
	
	if err := ValidateDomainName(domainName); err != nil {
		return err
	}
	
	if id <= 0 {
		return NewInvalidIDError("ID는 양의 정수여야 합니다")
	}
	
	return nil
}