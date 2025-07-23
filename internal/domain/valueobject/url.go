package valueobject

import (
	"errors"
	"net/url"
	"strings"
)

// URL represents a URL value object
type URL struct {
	value string
}

// NewURL creates a new URL value object with validation
func NewURL(urlString string) (*URL, error) {
	if urlString == "" {
		return nil, errors.New("URL cannot be empty")
	}

	if len(urlString) > 2048 {
		return nil, errors.New("URL cannot exceed 2048 characters")
	}

	// Parse and validate URL
	parsedURL, err := url.Parse(urlString)
	if err != nil {
		return nil, errors.New("invalid URL format")
	}

	if parsedURL.Scheme == "" {
		return nil, errors.New("URL must have a scheme (http:// or https://)")
	}

	if parsedURL.Host == "" {
		return nil, errors.New("URL must have a host")
	}

	// Normalize URL
	normalizedURL := normalizeURL(parsedURL)

	return &URL{
		value: normalizedURL,
	}, nil
}

// Value returns the URL string
func (u *URL) Value() string {
	return u.value
}

// String returns the URL string (implements Stringer interface)
func (u *URL) String() string {
	return u.value
}

// IsValid checks if the URL is in a valid state
func (u *URL) IsValid() bool {
	return u.value != "" && len(u.value) <= 2048
}

// normalizeURL normalizes the URL by removing trailing slashes and normalizing scheme
func normalizeURL(u *url.URL) string {
	// Ensure scheme is lowercase
	u.Scheme = strings.ToLower(u.Scheme)

	// Ensure host is lowercase
	u.Host = strings.ToLower(u.Host)

	// Remove trailing slash from path (except for root path)
	if u.Path != "/" {
		u.Path = strings.TrimSuffix(u.Path, "/")
	}

	return u.String()
}
