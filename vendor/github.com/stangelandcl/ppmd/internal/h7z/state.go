package h7z

type state struct {
	heap    *heap
	Address uint32
}

func newState(buf *heap) state {
	return state{heap: buf}
}

func (s *state) Symbol() uint32 {
	return s.heap.Byte(s.Address)
}

func (s *state) SetSymbol(v uint32) {
	s.heap.PutByte(s.Address, byte(v))
}

func (s *state) Freq() uint32 {
	return s.heap.Byte(s.Address + 1)
}

func (s *state) SetFreq(freq uint32) {
	s.heap.PutByte(s.Address+1, byte(freq))
}

func (s *state) IncrementFreq(dfreq uint32) {
	s.heap.PutByte(s.Address+1, byte(s.heap.Byte(s.Address+1)+dfreq))
	//s.SetFreq(s.Freq() + dfreq)
}

func (s *state) Successor() uint32 {
	return s.heap.UInt32(s.Address + 2)
}

func (s *state) SetSuccessor(successor uint32) {
	s.heap.PutUInt32(s.Address+2, successor)
}

func (s *state) SetRef(r stateRef) {
	s.SetSymbol(r.symbol)
	s.SetFreq(r.Freq())
	s.SetSuccessor(r.Successor)
}

func (s *state) SetValues(ptr state) {
	s.heap.Copy(s.Address, ptr.Address, stateSize)
}

func (s *state) DecrementAddress() {
	s.Address -= stateSize
}

func (s *state) IncrementAddress() {
	s.Address += stateSize
}

func ppmdSwap(p1 *state, p2 *state) {
	for i := uint32(0); i < stateSize; i++ {
		x := p1.heap.Byte(p1.Address + i)
		p1.heap.PutByte(p1.Address+i, byte(p2.heap.Byte(p2.Address+i)))
		p2.heap.PutByte(p2.Address+i, byte(x))
	}
}
