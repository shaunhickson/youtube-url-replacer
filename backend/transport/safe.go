package transport

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"
)

var (
	// Private IP blocks (RFC 1918, RFC 4193, RFC 4291)
	privateBlocks = []string{
		"127.0.0.0/8",    // IPv4 loopback
		"10.0.0.0/8",     // RFC1918
		"172.16.0.0/12",  // RFC1918
		"192.168.0.0/16", // RFC1918
		"169.254.0.0/16", // RFC3927 link-local
		"::1/128",        // IPv6 loopback
		"fc00::/7",       // IPv6 unique local
		"fe80::/10",      // IPv6 link-local
	}

	cidrs []*net.IPNet

	// AllowLocalIPs should only be true during testing
	AllowLocalIPs = false
)

func init() {
	for _, b := range privateBlocks {
		_, cidr, err := net.ParseCIDR(b)
		if err != nil {
			panic(err) // Should never happen with constant strings
		}
		cidrs = append(cidrs, cidr)
	}
}

func isPrivateIP(ip net.IP) bool {
	if AllowLocalIPs {
		return false
	}
	if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
		return true
	}
	for _, block := range cidrs {
		if block.Contains(ip) {
			return true
		}
	}
	return false
}

// SafeDialer returns a dial function that blocks private IPs
func SafeDialer(dialer *net.Dialer) func(ctx context.Context, network, addr string) (net.Conn, error) {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		host, port, err := net.SplitHostPort(addr)
		if err != nil {
			return nil, err
		}

		// Resolve IP
		ips, err := dialer.Resolver.LookupIPAddr(ctx, host)
		if err != nil {
			return nil, err
		}

		if len(ips) == 0 {
			return nil, errors.New("no IP addresses found")
		}

		// Check first IP (or iterate)
		// For strict safety, we dial the first valid one we find, but we must validate it.
		var targetIP net.IP
		for _, ip := range ips {
			if isPrivateIP(ip.IP) {
				continue // Skip private IPs
			}
			targetIP = ip.IP
			break
		}

		if targetIP == nil {
			return nil, errors.New("blocked: resolves to private/local IP")
		}

		// Dial the specific IP
		// We reconstruct the address using the validated IP
		// Note: This prevents DNS rebinding because we validated THIS IP.
		// However, for TLS (HTTPS), we need the hostname for SNI.
		// net.Dialer handles this if we pass the original hostname?
		// No, dialer takes (network, address). If we pass IP:Port, SNI might break.
		// But SafeDialer is used for the TCP connection. TLS handshake happens ON TOP of this connection.
		// The http.Transport handles SNI using the Request.URL.Host.
		// So dialing IP:Port is safe for the TCP layer.
		
		return dialer.DialContext(ctx, network, net.JoinHostPort(targetIP.String(), port))
	}
}

// NewSafeTransport returns an http.Transport configured for security
func NewSafeTransport() *http.Transport {
	dialer := &net.Dialer{
		Timeout:   2 * time.Second, // Fast connect timeout
		KeepAlive: 30 * time.Second,
	}

	return &http.Transport{
		DialContext:           SafeDialer(dialer),
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   5 * time.Second, // Fast TLS timeout
		ExpectContinueTimeout: 1 * time.Second,
	}
}
