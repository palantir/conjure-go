package h7z

const (
	rarNodeSize     = 4
	freqDataSize    = 6
	rarMemBlockSize = 12
	stateSize       = 6
	see2ContextSize = 4
	fixedUnitSize   = 12
	unitSize        = 12
	ppmContextSize  = 12

	n1       = 4
	n2       = 4
	n3       = 4
	n4       = (128 + 3 - 1*n1 - 2*n2 - 3*n3) / 4
	nIndexes = n1 + n2 + n3 + n4

	maxO       = 64
	intBits    = 7
	periodBits = 7
	totBits    = intBits + periodBits
	interval   = 1 << intBits
	binScale   = 1 << totBits
	maxFreq    = 124

	kTopValue = 1 << 24
)
