package ppmd

import (
	"fmt"
	"io"

	"github.com/stangelandcl/ppmd/internal/h7z"
)

type Reader struct {
	uncompressedSize, uncompressed int
	r                              io.Reader
	m                              *h7z.ModelPpm
}

/*
create PPMD variant H with 7-zip compatible modifications decompressor
(like PPMD7 in 7-zip source code)
*/
func NewH7zReader(r io.Reader, order, memorySize, uncompressedSize int) (Reader, error) {
	d := Reader{}
	if order < 2 || order > 64 {
		return d, fmt.Errorf("order out of range: %v. must be in [2, 64]", order)
	}
	if memorySize < (1<<11) || memorySize > (0xFFFFFFFF-12*3) {
		return d, fmt.Errorf("memory size out of range: %v", memorySize)
	}
	d.uncompressedSize = uncompressedSize
	d.r = r
	m, err := h7z.NewModelPpm(uint32(order), uint32(memorySize), r)
	d.m = m
	return d, err
}

func (d *Reader) Read(buf []byte) (int, error) {
	if d.uncompressed >= d.uncompressedSize {
		return 0, io.EOF
	}
	n := len(buf)
	remain := d.uncompressedSize - d.uncompressed
	if remain < n {
		n = remain
	}
	i := 0
	for i < n {
		c, err := d.m.DecodeChar()
		if err != nil {
			return i, err
		}
		buf[i] = byte(c)
		i++
		d.uncompressed++
	}
	return i, nil
}
