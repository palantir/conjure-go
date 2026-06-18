package h7z

import "encoding/binary"

type heap struct {
	buffer  []byte
	counter int
}

func newHeap(size uint32) *heap {
	return &heap{
		buffer: make([]byte, size),
	}
}

func (h *heap) PutByte(address uint32, x byte) {
	/*fmt.Printf("%v\t%v\t%v\n", h.counter, address, x)
	if h.counter == 360152 {
		debug.PrintStack()
	}*/

	h.buffer[address] = x
	h.counter++
}

func (h *heap) Byte(address uint32) uint32 {
	return uint32(h.buffer[address])
}

func (h *heap) PutUInt32(address, x uint32) {
	binary.LittleEndian.PutUint32(h.buffer[address:], x)
	/* for debugging use this but it disables inlining */
	/*
		u := uint32(int32(x))
		h.PutByte(address, byte(u))
		h.PutByte(address+1, byte(u>>8))
		h.PutByte(address+2, byte(u>>16))
		h.PutByte(address+3, byte(u>>24))
	*/
}

func (h *heap) PutUInt16(address uint32, x uint32) {
	u := uint16(x)
	h.PutByte(address, byte(u))
	h.PutByte(address+1, byte(u>>8))
}

func (h *heap) UInt16(address uint32) uint32 {
	u := uint16(h.buffer[0+address]) | uint16(h.buffer[1+address])<<8
	return uint32(u)
}

func (h *heap) UInt32(address uint32) uint32 {
	u := uint32(h.buffer[0+address]) |
		uint32(h.buffer[1+address])<<8 |
		uint32(h.buffer[2+address])<<16 |
		uint32(h.buffer[3+address])<<24
	return u
}

func (h *heap) Copy(dst, src, size uint32) {
	copy(h.buffer[dst:], h.buffer[src:src+size])
	// use for debugging
	/*
		for i := 0; i < size; i++ {
			h.PutByte(dst+i, h.buffer[src+i])
		}
	*/
}

func (h *heap) Clear(address, size uint32) {
	for i := uint32(0); i < size; i++ {
		h.PutByte(address+i, 0)
	}
}
