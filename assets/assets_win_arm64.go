//go:build windows && amr64
// +build windows,amr64

package assets

import (
	_ "embed"
)

//go:embed node_windows_amr64.zst
var EmbeddedNode []byte
