package utils

import (
	"fmt"
	"net/url"
	"strings"
)

func ValidateURL(rawURL string) error {
	if rawURL == "" {
		return fmt.Errorf("URL cannot be empty")
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}

	if u.Scheme == "" {
		return fmt.Errorf("URL must include scheme (http:// or https://)")
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("only http and https schemes are supported")
	}

	if u.Host == "" {
		return fmt.Errorf("URL must include host")
	}

	if strings.HasPrefix(u.Host, "file://") || u.Scheme == "file" {
		return fmt.Errorf("file:// URLs are not allowed for security reasons")
	}

	return nil
}