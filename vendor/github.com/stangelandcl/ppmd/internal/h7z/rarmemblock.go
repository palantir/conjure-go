package h7z

type rarMemBlock struct {
	heap    *heap
	Address uint32
}

func newRarMemBlock(heap *heap) rarMemBlock {
	return rarMemBlock{heap: heap}
}

func (r *rarMemBlock) Stamp() uint32 {
	return r.heap.UInt16(r.Address)
}

func (r *rarMemBlock) SetStamp(v uint32) {
	r.heap.PutUInt16(r.Address, v)
}

func (r *rarMemBlock) InsertAt(p rarMemBlock) {
	tmp := newRarMemBlock(r.heap)
	r.SetPrev(p.Address)
	tmp.Address = r.Prev()
	r.SetNext(tmp.Next())
	tmp.SetNext(r.Address)
	tmp.Address = r.Next()
	tmp.SetPrev(r.Address)
}

func (r *rarMemBlock) Remove() {
	tmp := newRarMemBlock(r.heap)
	tmp.Address = r.Prev()
	tmp.SetNext(r.Next())
	tmp.Address = r.Next()
	tmp.SetPrev(r.Prev())
}

func (r *rarMemBlock) Next() uint32 {
	return r.heap.UInt32(r.Address + 4)
}

func (r *rarMemBlock) SetNext(next uint32) {
	r.heap.PutUInt32(r.Address+4, next)
}

func (r *rarMemBlock) Prev() uint32 {
	return r.heap.UInt32(r.Address + 8)
}

func (r *rarMemBlock) SetPrev(prev uint32) {
	r.heap.PutUInt32(r.Address+8, prev)
}

func (r *rarMemBlock) Nu() uint32 {
	return r.heap.UInt16(r.Address + 2)
}

func (r *rarMemBlock) SetNu(prev uint32) {
	r.heap.PutUInt16(r.Address+2, prev)
}
