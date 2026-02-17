package resolvers

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/sph/youtube-url-replacer/backend/transport"
)

// SafeHttpClient returns an http.Client with SSRF protection
func SafeHttpClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Transport: transport.NewSafeTransport(),
		Timeout:   timeout,
	}
}

// ExtractMetadata parses HTML and attempts to find a title
func ExtractMetadata(r io.Reader) (*Result, error) {
	// Read a limited amount of data to avoid memory issues (e.g., 512KB)
	limitReader := io.LimitReader(r, 512*1024)
	body, err := io.ReadAll(limitReader)
	if err != nil {
		return nil, err
	}

	html := string(body)
	res := &Result{}

	// Try OpenGraph title first
	ogTitleRegex := regexp.MustCompile(`(?i)<meta\s+property=["']og:title["']\s+content=["']([^"']+)["']`)
	if matches := ogTitleRegex.FindStringSubmatch(html); len(matches) > 1 {
		res.Title = matches[1]
	}

	// Fallback to <title> tag
	if res.Title == "" {
		titleRegex := regexp.MustCompile(`(?i)<title>(.*?)</title>`)
		if matches := titleRegex.FindStringSubmatch(html); len(matches) > 1 {
			res.Title = matches[1]
		}
	}

	// Try OpenGraph description
	ogDescRegex := regexp.MustCompile(`(?i)<meta\s+property=["']og:description["']\s+content=["']([^"']+)["']`)
	if matches := ogDescRegex.FindStringSubmatch(html); len(matches) > 1 {
		res.Description = matches[1]
	}

	// Clean up title (unescape common entities, remove extra whitespace)
	res.Title = strings.TrimSpace(res.Title)
	res.Title = strings.ReplaceAll(res.Title, "&amp;", "&")
	res.Title = strings.ReplaceAll(res.Title, "&quot;", "\"")
	res.Title = strings.ReplaceAll(res.Title, "&#39;", "'")
	res.Title = strings.ReplaceAll(res.Title, "&lt;", "<")
	res.Title = strings.ReplaceAll(res.Title, "&gt;", ">")

	if res.Title == "" {
		return nil, fmt.Errorf("no title found")
	}

	return res, nil
}
