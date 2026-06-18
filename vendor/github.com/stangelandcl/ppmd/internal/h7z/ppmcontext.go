package h7z

type ppmContext struct {
	ps       [256]uint32
	Memory   *heap
	FreqData freqData
	OneState state
	address  uint32
	tmp      *ppmContext
}

var expEscape = []byte{25, 14, 9, 7, 5, 5, 4, 4, 4, 3, 3, 3, 2, 2, 2, 2}

func newPpmContext(mem *heap) *ppmContext {
	ctx := &ppmContext{
		Memory: mem,
		tmp: &ppmContext{
			Memory:   mem,
			FreqData: newFreqData(mem),
			OneState: newState(mem),
		},
		FreqData: newFreqData(mem),
		OneState: newState(mem),
	}
	ctx.reset()
	ctx.tmp.reset()
	return ctx
}

func (c *ppmContext) reset() {
	c.SetAddress(0)
}

func (p *ppmContext) NumStats() uint32 {
	return p.Memory.UInt16(p.address)
}

func (p *ppmContext) SetNumStats(num uint32) {
	p.Memory.PutUInt16(p.address, num)
}

func (p *ppmContext) SetOneState(r stateRef) {
	p.OneState.SetRef(r)
}

func (p *ppmContext) Suffix() uint32 {
	return p.Memory.UInt32(p.address + 8)
}

func (p *ppmContext) SetSuffix(suffix uint32) {
	p.Memory.PutUInt32(p.address+8, suffix)
}

func (p *ppmContext) SetAddress(addr uint32) {
	p.address = addr
	p.OneState.Address = addr + 2
	p.FreqData.Address = addr + 2
}

func (p *ppmContext) GetMean(summ, shift, round uint32) uint32 {
	return urshift(summ+(1<<(shift-round)), shift)
}

func (p *ppmContext) CreateChild(m *ModelPpm, pstats state, firstState stateRef) uint32 {
	//p.tmp.reset(m.SubAlloc.Heap)
	p.tmp.SetAddress(m.SubAlloc.AllocContext())
	p.tmp.SetNumStats(1)
	p.tmp.SetOneState(firstState)
	p.tmp.SetSuffix(p.address)
	pstats.SetSuccessor(p.tmp.address)
	return p.tmp.address
}

func (c *ppmContext) Rescale(model *ModelPpm) {
	oldNs := c.NumStats()
	i := c.NumStats() - 1

	// STATE* p1, * p;
	p1 := newState(model.SubAlloc.Heap)
	p := newState(model.SubAlloc.Heap)
	temp := newState(model.SubAlloc.Heap)

	for p.Address = model.FoundState.Address; p.Address != c.FreqData.Stats(); p.DecrementAddress() {
		temp.Address = p.Address - stateSize
		ppmdSwap(&p, &temp)
	}
	temp.Address = c.FreqData.Stats()
	temp.IncrementFreq(4)
	c.FreqData.IncrementSummFreq(4)
	escFreq := c.FreqData.SummFreq() - p.Freq()
	adder := uint32(0)
	if model.OrderFall != 0 {
		adder = 1
	}
	p.SetFreq(urshift(p.Freq()+adder, 1))
	c.FreqData.SetSummFreq(p.Freq())

	for {
		p.IncrementAddress()
		escFreq -= p.Freq()
		p.SetFreq(urshift(p.Freq()+adder, 1))
		c.FreqData.IncrementSummFreq(p.Freq())
		temp.Address = p.Address - stateSize
		if p.Freq() > temp.Freq() {
			p1.Address = p.Address
			tmp := newStateRef(&p1)
			temp2 := newState(model.SubAlloc.Heap)
			temp3 := newState(model.SubAlloc.Heap)
			for {
				temp2.Address = p1.Address - stateSize
				p1.SetValues(temp2)
				p1.DecrementAddress()
				temp3.Address = p1.Address - stateSize
				if p1.Address == c.FreqData.Stats() || tmp.Freq() <= temp3.Freq() {
					break
				}
			}
			p1.SetRef(tmp)
		}
		i--
		if i == 0 {
			break
		}
	}

	if p.Freq() == 0 {
		for {
			i++
			p.DecrementAddress()
			if p.Freq() != 0 {
				break
			}
		}
		escFreq += i
		c.SetNumStats(c.NumStats() - i)
		if c.NumStats() == 1 {
			temp.Address = c.FreqData.Stats()
			var tmp = newStateRef(&temp)

			// STATE tmp=*U.Stats;
			for {
				// tmp.Freq-=(tmp.Freq >> 1)
				tmp.DecrementFreq(urshift(tmp.Freq(), 1))
				escFreq = urshift(escFreq, 1)
				if escFreq <= 1 {
					break
				}
			}
			model.SubAlloc.FreeUnits(c.FreqData.Stats(), urshift(oldNs+1, 1))
			c.OneState.SetRef(tmp)
			model.FoundState.Address = c.OneState.Address
			return
		}
	}
	escFreq -= urshift(escFreq, 1)
	c.FreqData.IncrementSummFreq(escFreq)
	n0 := urshift(oldNs+1, 1)
	n1 := urshift(c.NumStats()+1, 1)
	if n0 != n1 {
		c.FreqData.SetStats(model.SubAlloc.ShrinkUnits(c.FreqData.Stats(), n0, n1))
	}
	model.FoundState.Address = c.FreqData.Stats()
}

func (p *ppmContext) GetArrayIndex(model *ModelPpm, rs state) uint32 {
	//p.tmp.reset(model.SubAlloc.Heap)
	p.tmp.SetAddress(p.Suffix())
	ret := uint32(0)
	ret += model.PrevSuccess()
	ret += model.Ns2BsIndex[p.tmp.NumStats()-1]
	ret += (model.hiBitsFlag & 0xff) + 2*model.Hb2Flag[rs.Symbol()]
	ret += (urshift(model.runLength, 26)) & 0x20
	return ret
}

func (c *ppmContext) Update1(model *ModelPpm, p uint32) {
	model.FoundState.Address = p
	model.FoundState.IncrementFreq(4)
	c.FreqData.IncrementSummFreq(4)
	p0 := newState(model.SubAlloc.Heap)
	p1 := newState(model.SubAlloc.Heap)
	p0.Address = p
	p1.Address = p - stateSize
	if p0.Freq() > p1.Freq() {
		ppmdSwap(&p0, &p1)
		model.FoundState.Address = p1.Address
		if p1.Freq() > maxFreq {
			c.Rescale(model)
		}
	}
}

func (c *ppmContext) update1_0(model *ModelPpm, p uint32) {
	model.FoundState.Address = p
	x := uint32(0)
	if 2*model.FoundState.Freq() > c.FreqData.SummFreq() {
		x = 1
	}
	model.SetPrevSuccess(x & 0xff)
	model.IncRunLength(model.PrevSuccess())
	c.FreqData.IncrementSummFreq(4)
	model.FoundState.IncrementFreq(4)
	if model.FoundState.Freq() > maxFreq {
		c.Rescale(model)
	}
}

func (c *ppmContext) Update2(model *ModelPpm, p uint32) {
	temp := newState(model.SubAlloc.Heap)
	temp.Address = p
	model.FoundState.Address = p
	model.FoundState.IncrementFreq(4)
	c.FreqData.IncrementSummFreq(4)
	if temp.Freq() > maxFreq {
		c.Rescale(model)
	}
	model.IncEscCount(1)
	model.runLength = model.InitRl
}

func (c *ppmContext) MakeEscFreq(model *ModelPpm, numMasked uint32) (*see2context, uint32) {
	//if model.decoder.counter == 701 {
	//debug.PrintStack()
	//}
	numStats := c.NumStats()
	nonMasked := numStats - numMasked
	if numStats != 256 {
		//c.tmp.reset(model.SubAlloc.Heap)
		c.tmp.SetAddress(c.Suffix())
		idx1 := model.Ns2Index[nonMasked-1]
		idx2 := uint32(0)
		if nonMasked < c.tmp.NumStats()-numStats {
			idx2++
		}
		if c.FreqData.SummFreq() < 11*numStats {
			idx2 += 2
		}
		if numMasked > nonMasked {
			idx2 += 4
		}
		idx2 += model.hiBitsFlag & 0xff
		psee2C := &model.See2Cont[idx1][idx2]
		escFreq := psee2C.Mean()
		return psee2C, escFreq
	} else {
		return &model.DummySee2Cont, 1
	}
}
