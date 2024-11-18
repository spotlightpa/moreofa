package static

import "embed"

//go:embed *.html *.ico
var FS embed.FS
