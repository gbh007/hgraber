package static

import "embed"

//go:embed *.js *.html *.css *.ico
var StaticDir embed.FS
