package h7z

type see2context struct {
	summ, shift, count uint32
}

func newSee2Context(val uint32) see2context {
	//fmt.Printf("newsee2: %v\n", val)
	s := see2context{}
	s.shift = (periodBits - 4) & 0xff
	s.summ = (val << s.shift) & 0xffff
	s.count = 4
	return s
}

func (s *see2context) Mean() uint32 {
	r := urshift(s.summ, s.shift)
	//fmt.Printf("mean: %v\t%v\t%v\n", r, s.summ, s.shift)
	s.summ -= r
	if r == 0 {
		r++
	}
	return r
}

func (s *see2context) Count() uint32 {
	return s.count
}

func (s *see2context) SetCount(c uint32) {
	s.count = c & 0xff
}

func (s *see2context) Shift() uint32 {
	return s.shift
}

func (s *see2context) SetShift(shift uint32) {
	//fmt.Printf("shift: %v\n", shift)
	s.shift = shift & 0xff
}

func (s *see2context) Summ() uint32 {
	return s.summ
}

func (s *see2context) SetSumm(summ uint32) {
	s.summ = summ & 0xffff
}

func (s *see2context) IncSumm(dsum uint32) {
	s.summ += dsum
}

func (s *see2context) Update() {
	//fmt.Printf("update: %v\n", s.shift)
	if s.shift < periodBits {
		s.count--
		if s.count == 0 {
			s.summ += s.summ
			s.count = 3 << s.shift
			s.shift++
		}
	}
	s.summ &= 0xffff
	s.count &= 0xff
	s.shift &= 0xff
}
