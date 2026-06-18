package h7z

type stateRef struct {
	Successor, symbol, freq uint32
}

func newStateRef(s *state) stateRef {
	r := stateRef{}
	r.SetFreq(s.Freq())
	r.SetSymbol(s.Symbol())
	r.Successor = s.Successor()
	return r
}

func (s *stateRef) IncrementFreq(dfreq uint32) {
	s.freq = (s.freq + dfreq) & 0xff
}

func (s *stateRef) DecrementFreq(dfreq uint32) {
	s.freq = (s.freq - dfreq) & 0xff
}

func (s *stateRef) Freq() uint32 {
	return s.freq
}

func (s *stateRef) SetFreq(f uint32) {
	s.freq = f & 0xff
}

func (s *stateRef) Symbol() uint32 {
	return s.symbol
}

func (s *stateRef) SetSymbol(sym uint32) {
	s.symbol = sym & 0xff
}
