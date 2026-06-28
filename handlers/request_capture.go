package handlers

import (
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/mynameismaxz/httpdebugger/dtos"
)

// maxBodyBytes caps how much of the request body the debugger reads. Anything
// beyond this is dropped and the body is flagged as truncated.
const maxBodyBytes = 1 << 20 // 1 MiB

// captureRequest reads the incoming request and builds a RequestDebugInfo.
// It consumes the body (up to maxBodyBytes); since this is a terminal echo
// handler nothing downstream needs the body, so it is not restored.
func captureRequest(w http.ResponseWriter, r *http.Request) *dtos.RequestDebugInfo {
	body, size, truncated := readBody(w, r)

	isBinary := false
	bodyStr := ""
	if !utf8.Valid(body) {
		isBinary = true
	} else {
		bodyStr = string(body)
	}

	info := &dtos.RequestDebugInfo{
		Method:        r.Method,
		Path:          r.URL.Path,
		Query:         map[string][]string(r.URL.Query()),
		Headers:       map[string][]string(r.Header),
		Host:          r.Host,
		Proto:         r.Proto,
		RemoteAddr:    r.RemoteAddr,
		ClientIP:      clientIP(r),
		TLS:           tlsInfo(r),
		ContentType:   r.Header.Get("Content-Type"),
		BodySize:      size,
		BodyTruncated: truncated,
		Body:          bodyStr,
		BodyIsBinary:  isBinary,
		Timestamp:     time.Now().UTC().Format(time.RFC3339),
	}
	return info
}

// readBody reads the request body up to maxBodyBytes. It returns the bytes
// read, their length, and whether the body was truncated at the cap.
func readBody(w http.ResponseWriter, r *http.Request) (data []byte, size int, truncated bool) {
	if r.Body == nil {
		return nil, 0, false
	}
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
	data, err := io.ReadAll(r.Body)
	if err != nil {
		// MaxBytesReader returns *http.MaxBytesError once the cap is exceeded;
		// data still holds the bytes read up to that point.
		var maxErr *http.MaxBytesError
		if asMaxBytesError(err, &maxErr) {
			truncated = true
		}
		// For any other read error we keep whatever bytes were returned.
	}
	return data, len(data), truncated
}

// asMaxBytesError reports whether err is an *http.MaxBytesError, storing it in
// target when so. Wrapped in a helper to keep readBody readable.
func asMaxBytesError(err error, target **http.MaxBytesError) bool {
	if me, ok := err.(*http.MaxBytesError); ok {
		*target = me
		return true
	}
	return false
}

// clientIP resolves the originating client IP, preferring proxy headers since
// the debugger typically runs behind Cloud Run's front end.
func clientIP(r *http.Request) string {
	// X-Forwarded-For: the leftmost entry is the original client.
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		if ip := strings.TrimSpace(strings.Split(xff, ",")[0]); ip != "" {
			return ip
		}
	}
	if xr := strings.TrimSpace(r.Header.Get("X-Real-IP")); xr != "" {
		return xr
	}
	if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		return host
	}
	return r.RemoteAddr
}

// tlsInfo extracts TLS connection details, or nil for plaintext requests.
func tlsInfo(r *http.Request) *dtos.TLSInfo {
	if r.TLS == nil {
		return nil
	}
	return &dtos.TLSInfo{
		Version:     tlsVersionName(r.TLS.Version),
		CipherSuite: tls.CipherSuiteName(r.TLS.CipherSuite),
		ServerName:  r.TLS.ServerName,
	}
}

func tlsVersionName(v uint16) string {
	switch v {
	case tls.VersionTLS13:
		return "TLS 1.3"
	case tls.VersionTLS12:
		return "TLS 1.2"
	case tls.VersionTLS11:
		return "TLS 1.1"
	case tls.VersionTLS10:
		return "TLS 1.0"
	default:
		return "unknown"
	}
}

// wantsJSON decides whether to respond with JSON instead of HTML.
// Precedence: explicit ?format= wins, then the Accept header, else HTML.
func wantsJSON(r *http.Request) bool {
	if f := r.URL.Query().Get("format"); f != "" {
		return strings.EqualFold(f, "json")
	}
	return strings.Contains(r.Header.Get("Accept"), "application/json")
}

// kv is a sorted key/values pair used by the HTML template, which cannot sort
// maps on its own.
type kv struct {
	Key    string
	Values []string
}

// sortedPairs converts a multi-valued map into a slice sorted by key.
func sortedPairs(m map[string][]string) []kv {
	pairs := make([]kv, 0, len(m))
	for k, v := range m {
		pairs = append(pairs, kv{Key: k, Values: v})
	}
	sort.Slice(pairs, func(i, j int) bool { return pairs[i].Key < pairs[j].Key })
	return pairs
}
