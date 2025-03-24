package assets

import (
	_ "embed"
)

//go:embed sub-store.bundle.js.zst
var EmbeddedSubStore []byte
