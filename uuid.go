package data

import (
	cr "crypto/rand"
)

type UUID [16]byte

var halfbyte2hexchar = []byte{
	48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 97, 98, 99, 100, 101, 102,
}

func (u UUID) String() string {
	b := [36]byte{}

	for i, n := range []int{
		0, 2, 4, 6,
		9, 11,
		14, 16,
		19, 21,
		24, 26, 28, 30, 32, 34,
	} {
		b[n] = halfbyte2hexchar[(u[i]>>4)&0x0f]
		b[n+1] = halfbyte2hexchar[u[i]&0x0f]
	}

	b[8] = '-'
	b[13] = '-'
	b[18] = '-'
	b[23] = '-'

	return string(b[:])
}

func V4() (UUID, error) {
	u := UUID{}

	_, err := cr.Read(u[:])
	if err != nil {
		return u, err
	}

	u[8] = (u[8] | 0x80) & 0xBF
	u[6] = (u[6] | 0x40) & 0x4F

	return u, nil
}

func V4Quick() string {
	u, err := V4()
	if err != nil {
		return err.Error()
	}
	return u.String()
}
