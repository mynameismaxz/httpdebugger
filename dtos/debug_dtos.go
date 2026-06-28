package dtos

import (
	"encoding/json"
)

// RequestDebugInfo holds the captured details of an incoming HTTP request.
// It is the payload the debugger echoes back to the client (as JSON, or
// rendered into HTML).
type RequestDebugInfo struct {
	Method        string              `json:"method"`
	Path          string              `json:"path"`
	Query         map[string][]string `json:"query"`
	Headers       map[string][]string `json:"headers"`
	Host          string              `json:"host"`
	Proto         string              `json:"proto"`
	RemoteAddr    string              `json:"remote_addr"`
	ClientIP      string              `json:"client_ip"`
	TLS           *TLSInfo            `json:"tls,omitempty"`
	ContentType   string              `json:"content_type"`
	BodySize      int                 `json:"body_size"`
	BodyTruncated bool                `json:"body_truncated"`
	Body          string              `json:"body"`
	BodyIsBinary  bool                `json:"body_is_binary"`
	Timestamp     string              `json:"timestamp"`
}

// TLSInfo describes the TLS connection details when the request arrived over
// HTTPS. It is omitted entirely for plaintext requests.
type TLSInfo struct {
	Version     string `json:"version"`
	CipherSuite string `json:"cipher_suite"`
	ServerName  string `json:"server_name"`
}

// ToJSON marshals the debug info with indentation for readable debug output.
func (r *RequestDebugInfo) ToJSON() []byte {
	resp, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return []byte("{}")
	}
	return resp
}
