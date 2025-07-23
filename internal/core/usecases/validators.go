package attributes

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

// AttributeValueValidator validates attribute values based on their type
type AttributeValueValidator interface {
	Validate(value string, orderIndex *int) error
}

// TagValidator validates tag attribute values
type TagValidator struct{}

func (v *TagValidator) Validate(value string, orderIndex *int) error {
	if len(value) == 0 {
		return ErrValueRequired
	}
	if len(value) > 255 {
		return ErrValueTooLong
	}
	return nil
}

// OrderedTagValidator validates ordered tag attribute values
type OrderedTagValidator struct{}

func (v *OrderedTagValidator) Validate(value string, orderIndex *int) error {
	if len(value) == 0 {
		return ErrValueRequired
	}
	if len(value) > 255 {
		return ErrValueTooLong
	}
	if orderIndex == nil {
		return ErrOrderIndexRequired
	}
	return nil
}

// NumberValidator validates number attribute values
type NumberValidator struct{}

func (v *NumberValidator) Validate(value string, orderIndex *int) error {
	if len(value) == 0 {
		return ErrValueRequired
	}

	// Try to parse as float64 to validate it's a number
	_, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return ErrInvalidNumber
	}

	return nil
}

// StringValidator validates string attribute values
type StringValidator struct{}

func (v *StringValidator) Validate(value string, orderIndex *int) error {
	if len(value) == 0 {
		return ErrValueRequired
	}
	if len(value) > 2048 {
		return ErrValueTooLong
	}
	return nil
}

// MarkdownValidator validates markdown attribute values
type MarkdownValidator struct{}

var (
	// Basic markdown patterns - we'll do simple validation
	markdownHeaderPattern = regexp.MustCompile(`^#{1,6}\s+.+$`)
	markdownLinkPattern   = regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
	markdownCodePattern   = regexp.MustCompile("```[\\s\\S]*?```|`[^`]+`")
	markdownBoldPattern   = regexp.MustCompile(`\*\*[^*]+\*\*|__[^_]+__`)
	markdownItalicPattern = regexp.MustCompile(`\*[^*]+\*|_[^_]+_`)
)

func (v *MarkdownValidator) Validate(value string, orderIndex *int) error {
	if len(value) == 0 {
		return ErrValueRequired
	}
	if len(value) > 10000 { // Larger limit for markdown
		return ErrValueTooLong
	}

	// Basic markdown syntax validation
	lines := strings.Split(value, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Check for malformed markdown links
		if strings.Contains(line, "](") {
			matches := markdownLinkPattern.FindAllString(line, -1)
			for _, match := range matches {
				if !isValidMarkdownLink(match) {
					return ErrInvalidMarkdown
				}
			}
		}
	}

	return nil
}

func isValidMarkdownLink(link string) bool {
	// Extract URL from markdown link pattern [text](url)
	re := regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
	matches := re.FindStringSubmatch(link)
	if len(matches) < 3 {
		return false
	}

	urlStr := matches[2]

	// Allow relative URLs and anchors
	if strings.HasPrefix(urlStr, "#") || strings.HasPrefix(urlStr, "/") || strings.HasPrefix(urlStr, "./") || strings.HasPrefix(urlStr, "../") {
		return true
	}

	// Validate absolute URLs
	_, err := url.Parse(urlStr)
	return err == nil
}

// ImageValidator validates image URL attribute values
type ImageValidator struct{}

var (
	// Common image extensions
	imageExtensions = []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp", ".svg", ".ico"}
)

func (v *ImageValidator) Validate(value string, orderIndex *int) error {
	if len(value) == 0 {
		return ErrValueRequired
	}
	if len(value) > 2048 {
		return ErrValueTooLong
	}

	// Parse URL
	parsedURL, err := url.Parse(value)
	if err != nil {
		return ErrInvalidURL
	}

	// Must be absolute URL for images
	if !parsedURL.IsAbs() {
		return ErrInvalidURL
	}

	// Check scheme
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return ErrInvalidURL
	}

	// Check if it looks like an image (has image extension or is from known image hosting)
	if !isValidImageURL(parsedURL) {
		return ErrInvalidURL
	}

	return nil
}

func isValidImageURL(parsedURL *url.URL) bool {
	path := strings.ToLower(parsedURL.Path)

	// Check file extension
	for _, ext := range imageExtensions {
		if strings.HasSuffix(path, ext) {
			return true
		}
	}

	// Check for known image hosting services
	host := strings.ToLower(parsedURL.Host)
	imageHosts := []string{
		"imgur.com", "i.imgur.com",
		"images.unsplash.com", "unsplash.com",
		"pixabay.com", "images.pexels.com",
		"cdn.pixabay.com", "images.pixabay.com",
		"gravatar.com", "s.gravatar.com",
		"githubusercontent.com",
		"googleusercontent.com",
		"cloudinary.com",
		"amazonaws.com", // S3 buckets
		"cloudfront.net",
	}

	for _, imageHost := range imageHosts {
		if strings.Contains(host, imageHost) {
			return true
		}
	}

	return false
}

// ValidatorFactory creates validators for different attribute types
type ValidatorFactory struct{}

func NewValidatorFactory() *ValidatorFactory {
	return &ValidatorFactory{}
}

func (f *ValidatorFactory) GetValidator(attrType AttributeType) (AttributeValueValidator, error) {
	switch attrType {
	case AttributeTypeTag:
		return &TagValidator{}, nil
	case AttributeTypeOrderedTag:
		return &OrderedTagValidator{}, nil
	case AttributeTypeNumber:
		return &NumberValidator{}, nil
	case AttributeTypeString:
		return &StringValidator{}, nil
	case AttributeTypeMarkdown:
		return &MarkdownValidator{}, nil
	case AttributeTypeImage:
		return &ImageValidator{}, nil
	default:
		return nil, fmt.Errorf("unsupported attribute type: %s", attrType)
	}
}

// ValidateAttributeValue validates a value against the given attribute type
func ValidateAttributeValue(attrType AttributeType, value string, orderIndex *int) error {
	factory := NewValidatorFactory()
	validator, err := factory.GetValidator(attrType)
	if err != nil {
		return err
	}

	return validator.Validate(value, orderIndex)
}
