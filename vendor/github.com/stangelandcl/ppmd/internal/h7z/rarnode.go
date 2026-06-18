package h7z

type rarNode struct {
	heap    *heap
	Address uint32
}

func newRarNode(buf *heap) rarNode {
	return rarNode{heap: buf}
}

func (r *rarNode) Next() uint32 {
	return r.heap.UInt32(r.Address)
}

func (r *rarNode) SetNext(next uint32) {
	r.heap.PutUInt32(r.Address, next)
}
