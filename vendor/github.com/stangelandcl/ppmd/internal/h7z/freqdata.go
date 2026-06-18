package h7z

type freqData struct {
	heap    *heap
	Address uint32
}

func newFreqData(buf *heap) freqData {
	return freqData{heap: buf}
}

func (f *freqData) Stats() uint32 {
	return f.heap.UInt32(f.Address + 2)
}

func (f *freqData) SetStats(state uint32) {
	f.heap.PutUInt32(f.Address+2, state)
}

func (f *freqData) SummFreq() uint32 {
	return f.heap.UInt16(f.Address)
}

func (f *freqData) SetSummFreq(freq uint32) {
	f.heap.PutUInt16(f.Address, freq)
}

func (f *freqData) IncrementSummFreq(fr uint32) {
	f.SetSummFreq(f.SummFreq() + fr)
}
