package h7z

import (
	"io"
)

type decoder struct {
	rng, code uint
	total     int
	r         io.Reader
	pos       int
	counter   int
	err       error
}

func newDecoder(r io.Reader) (decoder, error) {
	d := decoder{}
	d.r = r
	d.rng = 0xFFFFFFFF
	for i := 0; i < 5; i++ {
		b, err := d.byte()
		if err != nil {
			return d, err
		}
		d.code = (d.code << 8) | uint(b)
	}
	d.total = 5
	return d, nil
}

func (d *decoder) byte() (byte, error) {
	if d.err != nil {
		return 0, d.err
	}
	buf := [1]byte{}
	//fmt.Println("reading", d.pos)
	n, err := d.r.Read(buf[:])
	if err != nil {
		d.err = err
		if n == 0 {
			return 0, d.err
		}
	}
	d.pos++
	return buf[0], nil
}

func (d *decoder) normalize() error {
	//fmt.Printf("norm1: %v\t%v\t%v\t%v\n", d.counter, d.rng, d.code, d.total)
	if d.rng < kTopValue {
		b, err := d.byte()
		if err != nil {
			return err
		}
		d.code = (d.code << 8) | uint(b)
		d.rng <<= 8
		d.total++
		if d.rng < kTopValue {
			b, err := d.byte()
			if err != nil {
				return nil
			}
			d.code = (d.code << 8) | uint(b)
			d.rng <<= 8
			d.total++
		}
	}
	return nil
	//fmt.Printf("norm2: %v\t%v\t%v\t%v\n", d.counter, d.rng, d.code, d.total)
}

func (d *decoder) Threshold(total uint) uint {
	//fmt.Printf("threshold: %v\t%v\t%v\t%v\n", d.counter, d.rng, total, d.code)

	d.counter++
	d.rng /= total
	return d.code / d.rng
}

func (d *decoder) Decode(start, size uint) error {
	//fmt.Printf("decode: %v\t%v\t%v\t%v\t%v\n", d.counter, d.rng, start, d.code, size)
	d.code -= start * d.rng
	d.rng *= size
	return d.normalize()
}

func (d *decoder) DecodeBit(size0 uint, numTotalBits int) (uint, error) {
	//fmt.Printf("decodebit: %v\t%v\t%v\t%v\t%v\n", d.counter, d.rng, numTotalBits, d.code, size0)
	/*
		if d.counter == 396 {
			debug.PrintStack()
			os.Exit(1)
		}
	*/
	newbound := (d.rng >> uint(numTotalBits)) * size0
	var symbol uint
	if d.code < newbound {
		d.rng = newbound
	} else {
		symbol = 1
		d.code -= newbound
		d.rng -= newbound
	}
	err := d.normalize()
	return symbol, err
}
