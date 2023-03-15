package cycles

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertBitset(t *testing.T, bs bitset, set []bitID) {
	var bsSlice []bitID
	for i := 0; i < bs.size; i++ {
		if bs.has(bitID(i)) {
			bsSlice = append(bsSlice, bitID(i))
		}
	}
	assert.ElementsMatch(t, set, bsSlice)
}

func TestBitset(t *testing.T) {
	bs := newBitset(10)
	assertBitset(t, bs, []bitID{})

	bs.add(3)
	assertBitset(t, bs, []bitID{3})

	bs.add(5)
	assertBitset(t, bs, []bitID{3, 5})

	bs.add(9)
	assertBitset(t, bs, []bitID{3, 5, 9})

	bs.remove(5)
	assertBitset(t, bs, []bitID{3, 9})

	bs.add(6)
	assertBitset(t, bs, []bitID{3, 6, 9})

	bs.remove(3)
	assertBitset(t, bs, []bitID{6, 9})

	bs.remove(9)
	assertBitset(t, bs, []bitID{6})

	bs.add(1)
	assertBitset(t, bs, []bitID{1, 6})

	bs.remove(1)
	assertBitset(t, bs, []bitID{6})

	bs.remove(6)
	assertBitset(t, bs, []bitID{})
}
