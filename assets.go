// c:\Users\LabAdmin\Documents\darenFE\daren2\assets.go
package daren // Or your main package name at the root

import (
	"embed"
)

//go:embed all:web/static
var EmbeddedAssets embed.FS
