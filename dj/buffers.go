package dj

import (
	"github.com/palantir/pkg/bytesbuffers"
)

var (
	BytesBuffer1k  = bytesbuffers.NewSyncPool(1024)
	BytesBuffer32k = bytesbuffers.NewSyncPool(32 * 1024)
)
