//go:build noui

package frontend

import (
	"embed"
)

// Stub for testing/CI when dist/ doesn't exist.
// Build with `noui` tag to skip embedding.
var embeddedFiles embed.FS
