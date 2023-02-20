package graph

type bitID uint

// bitset is used for representing a set of integers with O(1) operations and O(n/64) merges.
type bitset struct {
	bits []uint64
	size int
}

func newBitset(n int) bitset {
	return bitset{
		size: n,
		// Number of uint64 is the number of bits required divided by 64 rounded up.
		bits: make([]uint64, (n+63)/64),
	}
}

func (bs *bitset) turnBitOn(i bitID) {
	bs.bits[i/64] |= 1 << (i % 64)
}

func (bs bitset) getBit(i bitID) bool {
	return (bs.bits[i/64] & (1 << (i % 64))) != 0
}

func (bs bitset) merge(o bitset) bitset {
	ret := newBitset(bs.size)
	for i := range ret.bits {
		ret.bits[i] = bs.bits[i] | o.bits[i]
	}
	return ret
}
