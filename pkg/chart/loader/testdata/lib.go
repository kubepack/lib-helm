package testdata

import (
	"embed"
)

//go:embed ** **/.helmignore **/charts/_ignore_me
var FS embed.FS
