package handlers

import (
	"embed"
	"html/template"

	"github.com/mynameismaxz/httpdebugger/dtos"
)

//go:embed templates/debug.html
var templatesFS embed.FS

// debugTmpl is parsed once at startup; a broken template fails fast here rather
// than per request.
var debugTmpl = template.Must(template.ParseFS(templatesFS, "templates/debug.html"))

// debugView is the template-facing model. It carries the raw debug info plus
// key-sorted slices of headers and query params (templates can't sort maps).
type debugView struct {
	*dtos.RequestDebugInfo
	HeadersSorted []kv
	QuerySorted   []kv
}

// newDebugView builds the view-model for HTML rendering.
func newDebugView(info *dtos.RequestDebugInfo) debugView {
	return debugView{
		RequestDebugInfo: info,
		HeadersSorted:    sortedPairs(info.Headers),
		QuerySorted:      sortedPairs(info.Query),
	}
}
