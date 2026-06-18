package h7z

import (
	"fmt"
	"io"
	"math"
)

type ModelPpm struct {
	charMask               [256]uint32
	Ns2Index               [256]uint32
	Ns2BsIndex             [256]uint32
	Hb2Flag                [256]uint32
	binSumm                [128][64]uint32
	ps                     [maxO]uint32
	See2Cont               [25][16]see2context
	escCount, _prevSuccess uint32
	SubAlloc               *subAllocator
	FoundState             state
	DummySee2Cont          see2context
	// temp contexts are used because they are created often and benchmarking
	// proved resetting was about 2x fast as creating each time
	minContext, maxContext, pc, pc2, upBranch, successor *ppmContext
	initEsc, maxOrder, runLength                         uint32
	InitRl, OrderFall                                    uint32

	decoder    decoder
	hiBitsFlag uint32
}

var InitBinEsc = []uint32{0x3CDD, 0x1F3F, 0x59BF, 0x48F3, 0x64A1, 0x5ABC, 0x6632, 0x6051}

func NewModelPpm(order, memsize uint32, r io.Reader) (*ModelPpm, error) {
	m := &ModelPpm{SubAlloc: newSubAllocator()}
	m.SubAlloc.StartSubAllocator(memsize)
	m.minContext = newPpmContext(m.SubAlloc.Heap)
	m.maxContext = newPpmContext(m.SubAlloc.Heap)
	m.pc = newPpmContext(m.SubAlloc.Heap)
	m.pc2 = newPpmContext(m.SubAlloc.Heap)
	m.upBranch = newPpmContext(m.SubAlloc.Heap)
	m.successor = newPpmContext(m.SubAlloc.Heap)
	m.FoundState = newState(m.SubAlloc.Heap)
	m.DummySee2Cont = see2context{}
	m.StartModelRare(order)
	if m.minContext.address == 0 {
		return nil, fmt.Errorf("mincontext is zero")
	}
	var err error
	m.decoder, err = newDecoder(r)
	return m, err
}

func (m *ModelPpm) PrevSuccess() uint32 {
	return m._prevSuccess & 0xff
}

func (m *ModelPpm) SetPrevSuccess(x uint32) {
	/*
		fmt.Printf("prev: %v\n", x)
		if m.decoder.counter == 150 {
			debug.PrintStack()
		}
	*/
	m._prevSuccess = x & 0xff
}

func (m *ModelPpm) RestartModelRare() {
	for i := 0; i < len(m.charMask); i++ {
		m.charMask[i] = 0
	}
	m.SubAlloc.InitSubAllocator()
	var initR1 uint32
	if m.maxOrder < 12 {
		initR1 = m.maxOrder
	} else {
		initR1 = 12
	}
	m.InitRl = -initR1 - 1
	addr := m.SubAlloc.AllocContext()
	m.minContext.SetAddress(addr)
	m.maxContext.SetAddress(addr)
	m.minContext.SetSuffix(0)
	m.OrderFall = m.maxOrder
	m.minContext.SetNumStats(256)
	m.minContext.FreqData.SetSummFreq(m.minContext.NumStats() + 1)

	addr = m.SubAlloc.AllocUnits(256 / 2)
	m.FoundState.Address = addr
	m.minContext.FreqData.SetStats(addr)

	state := newState(m.SubAlloc.Heap)
	addr = m.minContext.FreqData.Stats()
	m.runLength = m.InitRl
	m.SetPrevSuccess(0)
	for i := uint32(0); i < 256; i++ {
		state.Address = addr + i*stateSize
		state.SetSymbol(i)
		state.SetFreq(1)
		state.SetSuccessor(0)
	}

	for i := uint32(0); i < 128; i++ {
		for k := 0; k < 8; k++ {
			for n := 0; n < 64; n += 8 {
				m.binSumm[i][k+n] = binScale - InitBinEsc[k]/(i+2)
			}
		}
	}
	for i := uint32(0); i < 25; i++ {
		for k := 0; k < 16; k++ {
			m.See2Cont[i][k] = newSee2Context(5*i + 10)
		}
	}
}

func (m *ModelPpm) StartModelRare(maxOrder uint32) {
	var i, k, n, step uint32
	m.escCount = 1
	m.maxOrder = maxOrder
	m.RestartModelRare()

	// Bug Fixed
	m.Ns2BsIndex[0] = 0
	m.Ns2BsIndex[1] = 2
	for j := 0; j < 9; j++ {
		m.Ns2BsIndex[2+j] = 4
	}
	for j := 0; j < 256-11; j++ {
		m.Ns2BsIndex[11+j] = 6
	}
	for i = 0; i < 3; i++ {
		m.Ns2Index[i] = i
	}
	k = 1
	step = 1
	for n = i; i < 256; i++ {
		m.Ns2Index[i] = n
		k--
		if k == 0 {
			step++
			k = step
			n++
		}
	}
	for j := 0; j < 0x40; j++ {
		m.Hb2Flag[j] = 0
	}
	for j := 0; j < 0x100-0x40; j++ {
		m.Hb2Flag[0x40+j] = 0x08
	}
	m.DummySee2Cont.SetShift(periodBits)
}

func (m *ModelPpm) IncEscCount(dEscCount uint32) {
	m.escCount = (m.escCount + dEscCount) & 0xff
}

func (m *ModelPpm) IncRunLength(rl uint32) {
	m.runLength += rl
}

func (m *ModelPpm) CreateSuccessors(skip bool, p1 state) uint32 {
	upState := stateRef{}
	tempState := newState(m.SubAlloc.Heap)
	//m.pc2.reset()
	m.pc2.SetAddress(m.minContext.address)
	//m.upBranch.reset()
	m.upBranch.SetAddress(m.FoundState.Successor())

	p := newState(m.SubAlloc.Heap)
	pps := 0

	noLoop := false

	if !skip {
		m.ps[pps] = m.FoundState.Address
		pps++
		if m.pc2.Suffix() == 0 {
			noLoop = true
		}
	}
	if !noLoop {
		loopEntry := false
		if p1.Address != 0 {
			p.Address = p1.Address
			m.pc2.SetAddress(m.pc2.Suffix())
			loopEntry = true
		}
		for {
			if !loopEntry {
				m.pc2.SetAddress(m.pc2.Suffix())
				if m.pc2.NumStats() != 1 {
					p.Address = m.pc2.FreqData.Stats()
					if p.Symbol() != m.FoundState.Symbol() {
						for {
							p.IncrementAddress()
							if p.Symbol() == m.FoundState.Symbol() {
								break
							}
						}
					}
				} else {
					p.Address = m.pc2.OneState.Address
				}
			}
			loopEntry = false
			if p.Successor() != m.upBranch.address {
				m.pc2.SetAddress(p.Successor())
				break
			}
			m.ps[pps] = p.Address
			pps++

			if m.pc2.Suffix() == 0 {
				break
			}
		}
	}
	if pps == 0 {
		return m.pc2.address
	}
	upState.SetSymbol(m.SubAlloc.Heap.Byte(m.upBranch.address))
	upState.Successor = m.upBranch.address + 1 //TODO check if +1 necessary
	if m.pc2.NumStats() != 1 {
		if m.pc2.address <= m.SubAlloc.PText {
			return 0
		}
		p.Address = m.pc2.FreqData.Stats()
		if p.Symbol() != upState.Symbol() {
			for {
				p.IncrementAddress()
				if p.Symbol() == upState.Symbol() {
					break
				}
			}
		}
		cf := p.Freq() - 1
		s0 := m.pc2.FreqData.SummFreq() - m.pc2.NumStats() - cf

		tmp := uint32(0)
		if 2*cf <= s0 {
			if 5*cf > s0 {
				tmp = 1
			}
		} else {
			tmp = (2*cf + 3*s0 - 1) / (2 * s0)
		}
		upState.SetFreq(tmp + 1)
	} else {
		upState.SetFreq(m.pc2.OneState.Freq())
	}
	for {
		pps--
		tempState.Address = m.ps[pps]
		m.pc2.SetAddress(m.pc2.CreateChild(m, tempState, upState))
		if m.pc2.address == 0 {
			return 0
		}
		if pps == 0 {
			break
		}
	}
	return m.pc2.address
}

func (m *ModelPpm) UpdateModelRestart() {
	m.RestartModelRare()
	m.escCount = 0
}

func (m *ModelPpm) UpdateModel() {
	fs := newStateRef(&m.FoundState)
	p := newState(m.SubAlloc.Heap)
	tempState := newState(m.SubAlloc.Heap)

	//m.pc.reset()
	//m.successor.reset()

	var ns1, ns, cf, sf, s0 uint32
	m.pc.SetAddress(m.minContext.Suffix())
	if fs.Freq() < maxFreq/4 && m.pc.address != 0 {
		if m.pc.NumStats() != 1 {
			p.Address = m.pc.FreqData.Stats()
			if p.Symbol() != fs.Symbol() {
				for {
					p.IncrementAddress()
					if p.Symbol() == fs.Symbol() {
						break
					}
				}
				tempState.Address = p.Address - stateSize
				if p.Freq() >= tempState.Freq() {
					ppmdSwap(&p, &tempState)
					p.DecrementAddress()
				}
			}
			if p.Freq() < maxFreq-9 {
				p.IncrementFreq(2)
				m.pc.FreqData.IncrementSummFreq(2)
			}
		} else {
			p.Address = m.pc.OneState.Address
			if p.Freq() < 32 {
				p.IncrementFreq(1)
			}
		}
	}
	if m.OrderFall == 0 {
		m.FoundState.SetSuccessor(m.CreateSuccessors(true, p))
		m.minContext.SetAddress(m.FoundState.Successor())
		m.maxContext.SetAddress(m.FoundState.Successor())
		if m.minContext.address == 0 {
			m.UpdateModelRestart()
			return
		}
		return
	}
	m.SubAlloc.Heap.PutByte(m.SubAlloc.PText, byte(fs.Symbol()))
	m.SubAlloc.IncPText()
	m.successor.SetAddress(m.SubAlloc.PText)
	if m.SubAlloc.PText >= m.SubAlloc.FakeUnitsStart {
		m.UpdateModelRestart()
		return
	}

	if fs.Successor != 0 {
		if fs.Successor <= m.SubAlloc.PText {
			fs.Successor = m.CreateSuccessors(false, p)
			if fs.Successor == 0 {
				m.UpdateModelRestart()
				return
			}
		}
		m.OrderFall--
		if m.OrderFall == 0 {
			m.successor.SetAddress(fs.Successor)
			if m.maxContext.address != m.minContext.address {
				m.SubAlloc.DecPText(1)
			}
		}
	} else {
		m.FoundState.SetSuccessor(m.successor.address)
		fs.Successor = m.minContext.address
	}

	ns = m.minContext.NumStats()
	s0 = m.minContext.FreqData.SummFreq() - ns - (fs.Freq() - 1)
	m.pc.SetAddress(m.maxContext.address)
	for ; m.pc.address != m.minContext.address; m.pc.SetAddress(m.pc.Suffix()) {
		ns1 = m.pc.NumStats()
		if ns1 != 1 {
			if (ns1 & 1) == 0 {
				m.pc.FreqData.SetStats(m.SubAlloc.ExpandUnits(m.pc.FreqData.Stats(), urshift(ns1, 1)))
				if m.pc.FreqData.Stats() == 0 {
					m.UpdateModelRestart()
					return
				}
			}
			sum := uint32(0)
			if 2*ns1 < ns {
				sum = 1
			}
			sum2 := uint32(0)
			if 4*ns1 <= ns {
				sum2 = 1
			}
			sum3 := uint32(0)
			if m.pc.FreqData.SummFreq() <= 8*ns1 {
				sum3 = 1
			}
			sum4 := 2 * (sum2 & sum3)
			sum += sum4
			m.pc.FreqData.IncrementSummFreq(sum)
		} else {
			p.Address = m.SubAlloc.AllocUnits(1)
			if p.Address == 0 {
				m.UpdateModelRestart()
				return
			}
			p.SetValues(m.pc.OneState)
			m.pc.FreqData.SetStats(p.Address)
			if p.Freq() < maxFreq/4-1 {
				p.IncrementFreq(p.Freq())
			} else {
				p.SetFreq(maxFreq - 4)
			}
			freq := p.Freq() + m.initEsc
			if ns > 3 {
				freq++
			}
			m.pc.FreqData.SetSummFreq(freq)
		}
		cf = 2 * fs.Freq() * (m.pc.FreqData.SummFreq() + 6)
		sf = s0 + m.pc.FreqData.SummFreq()
		if cf < 6*sf {
			cf1 := uint32(0)
			if cf > sf {
				cf1 = 1
			}
			cf2 := uint32(0)
			if cf >= 4*sf {
				cf2 = 1
			}
			cf = 1 + cf1 + cf2
			m.pc.FreqData.IncrementSummFreq(3)
		} else {
			cf1 := uint32(0)
			if cf >= 9*sf {
				cf1 = 1
			}
			cf2 := uint32(0)
			if cf >= 12*sf {
				cf2 = 1
			}
			cf3 := uint32(0)
			if cf >= 15*sf {
				cf3 = 1
			}
			cf = 4 + cf1 + cf2 + cf3
			m.pc.FreqData.IncrementSummFreq(cf)
		}
		p.Address = m.pc.FreqData.Stats() + ns1*stateSize
		p.SetSuccessor(m.successor.address)
		p.SetSymbol(fs.Symbol())
		p.SetFreq(cf)
		ns1++
		m.pc.SetNumStats(ns1)
	}

	address := fs.Successor
	m.maxContext.SetAddress(address)
	m.minContext.SetAddress(address)
}

func (m *ModelPpm) NextContext() {
	addr := m.FoundState.Successor()
	if m.OrderFall == 0 && addr > m.SubAlloc.PText {
		m.minContext.SetAddress(addr)
		m.maxContext.SetAddress(addr)
	} else {
		m.UpdateModel()
	}
}

func (m *ModelPpm) DecodeChar() (byte, error) {
	if m.minContext.NumStats() != 1 {
		s := newState(m.SubAlloc.Heap)
		s.Address = m.minContext.FreqData.Stats()
		var i uint32
		sumfreq := m.minContext.FreqData.SummFreq()
		count := m.decoder.Threshold(uint(sumfreq))
		hiCnt := s.Freq()
		if uint32(count) < hiCnt {
			m.decoder.Decode(0, uint(s.Freq()))
			symbol := byte(s.Symbol())
			m.minContext.update1_0(m, s.Address)
			m.NextContext()
			return symbol, nil
		}
		m.SetPrevSuccess(0)
		i = m.minContext.NumStats() - 1
		for {
			s.IncrementAddress()
			hiCnt += s.Freq()
			if hiCnt > uint32(count) {
				m.decoder.Decode(uint(hiCnt-s.Freq()), uint(s.Freq()))
				symbol := byte(s.Symbol())
				m.minContext.Update1(m, s.Address)
				m.NextContext()
				return symbol, nil
			}
			i--
			if i <= 0 {
				break
			}
		}
		if uint32(count) >= m.minContext.FreqData.SummFreq() {
			return 0, fmt.Errorf("invalid SummFreq < count")
		}
		m.hiBitsFlag = m.Hb2Flag[m.FoundState.Symbol()] & 0xff
		m.decoder.Decode(uint(hiCnt), uint(m.minContext.FreqData.SummFreq()-hiCnt))
		for i = 0; i < 256; i++ {
			m.charMask[i] = math.MaxUint32
		}
		m.charMask[s.Symbol()] = 0
		i = m.minContext.NumStats() - 1
		for {
			s.DecrementAddress()
			m.charMask[s.Symbol()] = 0
			i--
			if i <= 0 {
				break
			}
		}
	} else {
		rs := newState(m.SubAlloc.Heap)
		rs.Address = m.minContext.OneState.Address
		m.hiBitsFlag = m.Hb2Flag[m.FoundState.Symbol()]
		off1 := rs.Freq() - 1
		off2 := m.minContext.GetArrayIndex(m, rs)
		bs := m.binSumm[off1][off2]
		symd, err := m.decoder.DecodeBit(uint(bs), 14)
		if err != nil {
			return 0, err
		}
		if symd == 0 {
			m.binSumm[off1][off2] = (bs + interval - m.minContext.GetMean(bs, periodBits, 2)) & 0xFFFF
			m.FoundState.Address = rs.Address
			symbol := byte(rs.Symbol())
			freq := uint32(0)
			if rs.Freq() < 128 {
				freq = 1
			}
			rs.IncrementFreq(freq)
			m.SetPrevSuccess(1)
			m.IncRunLength(1)
			m.NextContext()
			return symbol, nil
		}
		bs = (bs - m.minContext.GetMean(bs, periodBits, 2)) & 0xFFFF
		m.binSumm[off1][off2] = bs
		m.initEsc = uint32(expEscape[urshift(bs, 10)])
		for i := 0; i < 256; i++ {
			m.charMask[i] = math.MaxUint32
		}
		m.charMask[rs.Symbol()] = 0
		m.SetPrevSuccess(0)
	}
	for {
		s := newState(m.SubAlloc.Heap)
		numMasked := m.minContext.NumStats()
		for {
			m.OrderFall++
			m.minContext.SetAddress(m.minContext.Suffix())
			if m.minContext.address <= m.SubAlloc.PText || m.minContext.address > m.SubAlloc.HeapEnd {
				return 0, fmt.Errorf("invalid mincontext.address")
			}
			if m.minContext.NumStats() != numMasked {
				break
			}
		}
		hiCnt := uint32(0)
		s.Address = m.minContext.FreqData.Stats()
		i := uint32(0)
		num := m.minContext.NumStats() - numMasked
		for {
			k := m.charMask[s.Symbol()]
			hiCnt += s.Freq() & k
			m.minContext.ps[i] = s.Address
			s.IncrementAddress()
			i -= k
			if i == num {
				break
			}
		}

		see, freqSum := m.minContext.MakeEscFreq(m, numMasked)
		freqSum += hiCnt
		count := uint32(m.decoder.Threshold(uint(freqSum)))

		if count < hiCnt {
			ps := newState(m.SubAlloc.Heap)
			hiCnt = 0
			i = 0
			ps.Address = m.minContext.ps[i]
			hiCnt += ps.Freq()
			for hiCnt <= count {
				i++
				ps.Address = m.minContext.ps[i]
				hiCnt += ps.Freq()
			}
			s.Address = ps.Address
			m.decoder.Decode(uint(hiCnt-s.Freq()), uint(s.Freq()))
			see.Update()
			symbol := byte(s.Symbol())
			m.minContext.Update2(m, s.Address)
			m.UpdateModel()
			return symbol, nil
		}
		if count >= freqSum {
			return 0, fmt.Errorf("invalid ppmd state. count >= freqsum")
		}
		m.decoder.Decode(uint(hiCnt), uint(freqSum-hiCnt))
		see.SetSumm(see.Summ() + freqSum)
		for {
			i--
			s.Address = m.minContext.ps[i]
			m.charMask[s.Symbol()] = 0
			if i == 0 {
				break
			}
		}
	}

}
