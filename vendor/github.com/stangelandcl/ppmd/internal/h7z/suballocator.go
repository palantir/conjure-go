package h7z

type subAllocator struct {
	glueCount, subAllocatorSize, unitsStart, freeListPos uint32
	heapStart, loUnit, hiUnit, tempMemBlockPos           uint32
	HeapEnd, PText, FakeUnitsStart                       uint32
	units2Indx                                           [128]uint32
	indx2Units                                           [nIndexes]uint32
	Heap                                                 *heap
	freeList                                             [nIndexes]rarNode
}

func newSubAllocator() *subAllocator {
	return &subAllocator{}
}

func (s *subAllocator) Clean() {
	s.subAllocatorSize = 0
}

func (s *subAllocator) InsertNode(p, indx uint32) {
	x := newRarNode(s.Heap)
	x.Address = p
	x.SetNext(s.freeList[indx].Next())
	s.freeList[indx].SetNext(x.Address)
}

func (s *subAllocator) IncPText() {
	s.PText++
}

func (s *subAllocator) RemoveNode(indx uint32) uint32 {
	r := s.freeList[indx].Next()
	x := newRarNode(s.Heap)
	x.Address = r
	s.freeList[indx].SetNext(x.Next())
	return r
}

func (s *subAllocator) U2B(nu uint32) uint32 {
	return unitSize * nu
}

func (s *subAllocator) MbPtr(basePtr, items uint32) uint32 {
	return basePtr + s.U2B(items)
}

func (s *subAllocator) SplitBlock(pv, oldIndx, newIndx uint32) {
	uDiff := s.indx2Units[oldIndx] - s.indx2Units[newIndx]
	p := pv + s.U2B(s.indx2Units[newIndx])
	i := s.units2Indx[uDiff-1]
	if s.indx2Units[i] != uDiff {
		i--
		s.InsertNode(p, i)
		i = s.indx2Units[i]
		p += s.U2B(i)
		uDiff -= i
	}
	s.InsertNode(p, s.units2Indx[uDiff-1])
}

func (s *subAllocator) AllocateMemory() uint32 {
	return s.subAllocatorSize
}

func (s *subAllocator) StartSubAllocator(saSize uint32) {
	allocSize := saSize/fixedUnitSize*unitSize + unitSize
	realAllocSize := 1 + allocSize + 4*nIndexes
	s.tempMemBlockPos = realAllocSize
	realAllocSize += rarMemBlockSize
	s.Heap = newHeap(realAllocSize)
	s.heapStart = 1
	s.HeapEnd = s.heapStart + allocSize - unitSize
	s.subAllocatorSize = saSize

	s.freeListPos = s.heapStart + allocSize

	pos := s.freeListPos
	for i := 0; i < len(s.freeList); i++ {
		r := newRarNode(s.Heap)
		r.Address = pos
		s.freeList[i] = r
		pos += rarNodeSize
	}
}

func (s *subAllocator) GlueFreeBlocks() {
	s0 := newRarMemBlock(s.Heap)
	s0.Address = s.tempMemBlockPos
	p := newRarMemBlock(s.Heap)
	p1 := newRarMemBlock(s.Heap)

	if s.loUnit != s.hiUnit {
		s.Heap.PutByte(s.loUnit, 0)
	}
	s0.SetPrev(s0.Address)
	s0.SetNext(s0.Address)
	for i := uint32(0); i < nIndexes; i++ {
		for s.freeList[i].Next() != 0 {
			p.Address = s.RemoveNode(i)
			p.InsertAt(s0)
			p.SetStamp(0xFFFF)
			p.SetNu(s.indx2Units[i])
		}
	}

	for p.Address = s0.Next(); p.Address != s0.Address; p.Address = p.Next() {
		p1.Address = s.MbPtr(p.Address, p.Nu())
		for p1.Stamp() == 0xFFFF && p.Nu()+p1.Nu() < 0x10000 {
			p1.Remove()
			p.SetNu(p.Nu() + p1.Nu())
			p1.Address = s.MbPtr(p.Address, p.Nu())
		}
	}

	p.Address = s0.Next()
	for p.Address != s0.Address {
		p.Remove()
		var sz uint32
		for sz = p.Nu(); sz > 128; sz -= 128 {
			s.InsertNode(p.Address, nIndexes-1)
			p.Address = s.MbPtr(p.Address, 128)
		}
		i := s.units2Indx[sz-1]
		if s.indx2Units[i] != sz {
			i--
			k := sz - s.indx2Units[i]
			s.InsertNode(s.MbPtr(p.Address, sz-k), k-1)
		}
		s.InsertNode(p.Address, i)
		p.Address = s0.Next()
	}
}

func (s *subAllocator) AllocUnitsRare(indx uint32) uint32 {
	if s.glueCount == 0 {
		s.glueCount = 255
		s.GlueFreeBlocks()
		if s.freeList[indx].Next() != 0 {
			return s.RemoveNode(indx)
		}
	}

	i := indx
	for {
		i++
		if i == nIndexes {
			s.glueCount--
			i = s.U2B(s.indx2Units[indx])
			j := fixedUnitSize * s.indx2Units[indx]
			if s.FakeUnitsStart-s.PText > j {
				s.FakeUnitsStart -= j
				s.unitsStart -= i
				return s.unitsStart
			}
			return 0
		}

		if s.freeList[i].Next() != 0 {
			break
		}
	}

	r := s.RemoveNode(i)
	s.SplitBlock(r, i, indx)
	return r
}

func (s *subAllocator) AllocUnits(nu uint32) uint32 {
	indx := s.units2Indx[nu-1]
	if s.freeList[indx].Next() != 0 {
		return s.RemoveNode(indx)
	}
	r := s.loUnit
	s.loUnit += s.U2B(s.indx2Units[indx])
	if s.loUnit <= s.hiUnit {
		return r
	}
	s.loUnit -= s.U2B(s.indx2Units[indx])
	return s.AllocUnitsRare(indx)
}

func (s *subAllocator) AllocContext() uint32 {
	if s.hiUnit != s.loUnit {
		s.hiUnit -= unitSize
		return s.hiUnit
	}
	if s.freeList[0].Next() != 0 {
		return s.RemoveNode(0)
	}
	return s.AllocUnitsRare(0)
}

func (s *subAllocator) ExpandUnits(oldPtr, oldNu uint32) uint32 {
	i0 := s.units2Indx[oldNu-1]
	i1 := s.units2Indx[oldNu-1+1]
	if i0 == i1 {
		return oldPtr
	}

	ptr := s.AllocUnits(oldNu + 1)
	if ptr != 0 {
		s.Heap.Copy(ptr, oldPtr, s.U2B(oldNu))
		s.InsertNode(oldPtr, i0)
	}
	return ptr
}

func (s *subAllocator) ShrinkUnits(oldPtr, oldNu, newNu uint32) uint32 {
	i0 := s.units2Indx[oldNu-1]
	i1 := s.units2Indx[newNu-1]
	if i0 == i1 {
		return oldPtr
	}

	if s.freeList[i1].Next() != 0 {
		ptr := s.RemoveNode(i1)
		s.Heap.Copy(ptr, oldPtr, s.U2B(newNu))
		s.InsertNode(oldPtr, i0)
		return ptr
	}
	s.SplitBlock(oldPtr, i0, i1)
	return oldPtr
}

func (s *subAllocator) FreeUnits(ptr, oldNu uint32) {
	s.InsertNode(ptr, s.units2Indx[oldNu-1])
}

func (s *subAllocator) DecPText(dptext uint32) {
	s.PText -= dptext
}

func (s *subAllocator) InitSubAllocator() {
	s.Heap.Clear(s.freeListPos, s.SizeOfFreeList())
	s.PText = s.heapStart

	size2 := fixedUnitSize * (s.subAllocatorSize / 8 / fixedUnitSize * 7)
	realSize2 := size2 / fixedUnitSize * unitSize
	size1 := s.subAllocatorSize - size2
	realSize1 := size1/fixedUnitSize*unitSize + size1%fixedUnitSize
	s.hiUnit = s.heapStart + s.subAllocatorSize
	s.loUnit = s.heapStart + realSize1
	s.unitsStart = s.loUnit
	s.FakeUnitsStart = s.heapStart + size1
	s.hiUnit = s.loUnit + realSize2

	k := uint32(1)
	i := uint32(0)
	for i = 0; i < n1; i++ {
		s.indx2Units[i] = k & 0xff
		k++
	}

	for k++; i < n1+n2; i++ {
		s.indx2Units[i] = k & 0xff
		k += 2
	}

	for k++; i < n1+n2+n3; i++ {
		s.indx2Units[i] = k & 0xff
		k += 3
	}

	for k++; i < n1+n2+n3+n4; i++ {
		s.indx2Units[i] = k & 0xff
		k += 4
	}

	s.glueCount = 0
	i = 0
	for k = 0; k < 128; k++ {
		if s.indx2Units[i] < k+1 {
			i++
		}
		s.units2Indx[k] = i & 0xff
	}
}

func (s *subAllocator) SizeOfFreeList() uint32 {
	return uint32(len(s.freeList)) * rarNodeSize
}
