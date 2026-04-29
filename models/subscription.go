// Package models defines the core data structures used throughout sub2api.
package models

import "time"

// ProxyType represents the protocol type of a proxy node.
type ProxyType string

const (
	ProxyTypeVMess    ProxyType = "vmess"
	ProxyTypeVLess    ProxyType = "vless"
	ProxyTypeTrojan   ProxyType = "trojan"
	ProxyTypeShadowsocks ProxyType = "ss"
	ProxyTypeHysteria ProxyType = "hysteria"
	ProxyTypeHysteria2 ProxyType = "hysteria2"
	ProxyTypeTUIC     ProxyType = "tuic"
)

// ProxyNode represents a single proxy server configuration.
type ProxyNode struct {
	// Name is the human-readable label for this node.
	Name string `json:"name"`

	// Type is the proxy protocol type.
	Type ProxyType `json:"type"`

	// Server is the hostname or IP address of the proxy server.
	Server string `json:"server"`

	// Port is the port number the proxy server listens on.
	Port int `json:"port"`

	// UUID is used by VMess/VLess protocols for authentication.
	UUID string `json:"uuid,omitempty"`

	// Password is used by Trojan/Shadowsocks protocols.
	Password string `json:"password,omitempty"`

	// Cipher is the encryption method (used by Shadowsocks).
	Cipher string `json:"cipher,omitempty"`

	// Network is the transport layer type (tcp, ws, grpc, etc.).
	Network string `json:"network,omitempty"`

	// TLS indicates whether TLS is enabled.
	TLS bool `json:"tls,omitempty"`

	// SNI is the Server Name Indication value for TLS connections.
	SNI string `json:"sni,omitempty"`

	// ALPN is a list of application-layer protocol negotiation values.
	ALPN []string `json:"alpn,omitempty"`

	// Path is the WebSocket or HTTP/2 path.
	Path string `json:"path,omitempty"`

	// Host is the HTTP host header value.
	Host string `json:"host,omitempty"`

	// SkipCertVerify disables TLS certificate verification when true.
	SkipCertVerify bool `json:"skip-cert-verify,omitempty"`

	// Extra holds protocol-specific fields not covered by common fields.
	Extra map[string]interface{} `json:"extra,omitempty"`
}

// Subscription represents a parsed proxy subscription containing metadata
// and a list of proxy nodes.
type Subscription struct {
	// URL is the original subscription URL this was fetched from.
	URL string `json:"url"`

	// Nodes holds all proxy nodes parsed from the subscription.
	Nodes []ProxyNode `json:"nodes"`

	// FetchedAt is the timestamp when the subscription was last fetched.
	FetchedAt time.Time `json:"fetched_at"`

	// UserAgent is the HTTP User-Agent used when fetching the subscription.
	UserAgent string `json:"user_agent,omitempty"`

	// RawContent stores the raw subscription body before parsing.
	RawContent []byte `json:"-"`

	// UploadBytes reports the upload traffic quota from the subscription header.
	UploadBytes int64 `json:"upload_bytes,omitempty"`

	// DownloadBytes reports the download traffic quota from the subscription header.
	DownloadBytes int64 `json:"download_bytes,omitempty"`

	// TotalBytes reports the total traffic quota from the subscription header.
	TotalBytes int64 `json:"total_bytes,omitempty"`

	// ExpireAt is the subscription expiry time parsed from the header.
	ExpireAt *time.Time `json:"expire_at,omitempty"`
}

// NodeCount returns the total number of proxy nodes in the subscription.
func (s *Subscription) NodeCount() int {
	return len(s.Nodes)
}

// IsExpired reports whether the subscription has passed its expiry time.
func (s *Subscription) IsExpired() bool {
	if s.ExpireAt == nil {
		return false
	}
	return time.Now().After(*s.ExpireAt)
}
