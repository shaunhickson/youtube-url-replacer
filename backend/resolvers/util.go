package resolvers

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// SSRF Protection: Private and Reserved IP ranges
var privateIPBlocks []*net.IPNet

func init() {
	for _, cidr := range []string{
		"127.0.0.0/8",    // IPv4 loopback
		"10.0.0.0/8",     // RFC1918
		"172.16.0.0/12",  // RFC1918
		"192.168.0.0/16", // RFC1918
		"169.254.0.0/16", // RFC3927 link-local
		"::1/128",        // IPv6 loopback
		"fe80::/10",      // IPv6 link-local
		"fc00::/7",       // IPv6 unique local addr
	} {
		_, block, _ := net.ParseCIDR(cidr)
		privateIPBlocks = append(privateIPBlocks, block)
	}
}

func isPrivateIP(ip net.IP) bool {
	if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
		return true
	}
	for _, block := range privateIPBlocks {
		if block.Contains(ip) {
			return true
		}
	}
	return false
}

// SafeHttpClient returns an http.Client with SSRF protection
func SafeHttpClient(timeout time.Duration) *http.Client {
	dialer := &net.Dialer{
		Timeout:   timeout,
		KeepAlive: 30 * time.Second,
	}

	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			host, port, err := net.SplitHostPort(addr)
			if err != nil {
				return nil, err
			}

			ips, err := net.DefaultResolver.LookupIPAddr(ctx, host)
			if err != nil {
				return nil, err
			}

			for _, ip := range ips {
				if isPrivateIP(ip.IP) {
					return nil, fmt.Errorf("SSRF protection: access to private IP %s is blocked", ip.IP.String())
				}
			}

			// Use the first safe IP
			return dialer.DialContext(ctx, network, net.JoinHostPort(ips[0].IP.String(), port))
		},
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	return &http.Client{
		Transport: transport,
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
