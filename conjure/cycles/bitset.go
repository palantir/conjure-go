// Copyright (c) 2018 Palantir Technologies. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cycles

type bitID uint

// bitset is used for representing a set of integers in the range [0, size-1]
// Consumes O(size/64) memory, has O(1) set and get operations and O(size/64) union set operations.
type bitset struct {
	bits []uint64
	size int
}

func newBitset(size int) bitset {
	return bitset{
		size: size,
		// Number of uint64 is the number of bits required divided by 64 rounded up.
		bits: make([]uint64, (size+63)/64),
	}
}

func (bs *bitset) turnBitOn(i bitID) {
	// i/64; just need to move 6 bits to the right
	iRem := i >> 6
	// i%64: turn off all bits but the last 6
	iMod := i & 63
	// i == iRem*64 + iMod
	// Turn on i-th bit in the bitset by turning on the iMod-th bit in the iRem-th uint64
	bs.bits[iRem] |= 1 << iMod
}

func (bs bitset) getBit(i bitID) bool {
	// i/64; just need to move 6 bits to the right
	iRem := i >> 6
	// i%64: turn off all bits but the last 6
	iMod := i & 63
	// i == iRem*64 + iMod
	// Check if i-th bit in the bitset is on by checking the iMod-th bit in the iRem-th uint64
	return (bs.bits[iRem] & (1 << iMod)) != 0
}

func (bs bitset) merge(o bitset) bitset {
	if bs.size != o.size {
		panic("tried to merge two bitsets of different sizes")
	}

	ret := newBitset(bs.size)
	for i := range ret.bits {
		ret.bits[i] = bs.bits[i] | o.bits[i]
	}
	return ret
}
