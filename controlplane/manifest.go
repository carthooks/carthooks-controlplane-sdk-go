package controlplane

import _ "embed"

//go:embed arcubase-operations.manifest.json
var embeddedArcubaseOperationsManifest []byte

func EmbeddedArcubaseOperationsManifest() []byte {
	out := make([]byte, len(embeddedArcubaseOperationsManifest))
	copy(out, embeddedArcubaseOperationsManifest)
	return out
}
