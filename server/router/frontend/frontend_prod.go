//go:build !noui

package frontend

import (
	"embed"
)

//go:embed dist/*
var embeddedFiles embed.FS
