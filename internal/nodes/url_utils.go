package nodes

import (
	"net/url"
	"strings"
)

func GenerateTitleFromURL(rawURL string) string {
	if rawURL == "" {
		return ""
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return extractTitleFromRawURL(rawURL)
	}

	// Extract domain name
	host := parsedURL.Host
	if host == "" {
		return extractTitleFromRawURL(rawURL)
	}

	// Remove www. prefix if present
	if strings.HasPrefix(host, "www.") {
		host = host[4:]
	}

	// Get path for additional context
	path := parsedURL.Path
	if path == "" || path == "/" {
		return host
	}

	// Handle trailing slash - if path ends with /, treat it as no meaningful path
	if strings.HasSuffix(path, "/") && len(path) > 1 {
		return host
	}

	// Extract meaningful part from path
	pathParts := strings.Split(strings.Trim(path, "/"), "/")
	if len(pathParts) > 0 {
		lastPart := pathParts[len(pathParts)-1]
		if lastPart != "" {
			// Remove file extension if present
			if dotIndex := strings.LastIndex(lastPart, "."); dotIndex > 0 {
				lastPart = lastPart[:dotIndex]
			}

			// Replace dashes and underscores with spaces
			lastPart = strings.ReplaceAll(lastPart, "-", " ")
			lastPart = strings.ReplaceAll(lastPart, "_", " ")

			// Capitalize first letter of each word
			words := strings.Fields(lastPart)
			for i, word := range words {
				if len(word) > 0 {
					words[i] = strings.ToUpper(word[:1]) + word[1:]
				}
			}

			return strings.Join(words, " ") + " - " + host
		}
	}

	return host
}

func extractTitleFromRawURL(rawURL string) string {
	// If URL parsing fails, try to extract meaningful parts manually
	if strings.Contains(rawURL, "://") {
		parts := strings.Split(rawURL, "://")
		if len(parts) > 1 {
			remaining := parts[1]
			if slashIndex := strings.Index(remaining, "/"); slashIndex >= 0 {
				host := remaining[:slashIndex]
				path := remaining[slashIndex+1:]

				// Remove www. prefix
				if strings.HasPrefix(host, "www.") {
					host = host[4:]
				}

				if path != "" {
					pathParts := strings.Split(path, "/")
					if len(pathParts) > 0 {
						lastPart := pathParts[len(pathParts)-1]
						if lastPart != "" {
							return lastPart + " - " + host
						}
					}
				}

				return host
			} else {
				// No slash found, entire remaining part is the host
				host := remaining
				// Remove www. prefix
				if strings.HasPrefix(host, "www.") {
					host = host[4:]
				}
				return host
			}
		}
	}

	return rawURL
}

func ValidateURL(rawURL string) error {
	if rawURL == "" {
		return ErrNodeURLInvalid
	}

	if len(rawURL) > 2048 {
		return ErrNodeURLInvalid
	}

	// Basic URL validation - just check if it looks like a URL
	// We're being lenient here as per the requirements
	if !strings.Contains(rawURL, "://") && !strings.HasPrefix(rawURL, "//") {
		return ErrNodeURLInvalid
	}

	return nil
}
