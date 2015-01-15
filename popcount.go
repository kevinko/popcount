// Copyright 2012, Kevin Ko <kevin@faveset.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package popcount

func PopCount32(v uint32) int {
	/*
	  This is taken from bithacks.

	   0xf = 1111
	   0x5 = 0101
	   0x3 = 0011
	   0x1 = 0001

	  We wish to map each 2-bit pair as follows to its pop count
	  stored in 2 counter bits:

	   nibble count
	    00 -> 00
	    01 -> 01
	    10 -> 01
	    11 -> 10

	  Let v=[ab] be the bit string consisting of bits a and b.
	  Then,

	   [ab] - (a & 1) = v - ((v >> 1) & 0x2)

	  satisfies the mapping.

	  Generalized to 32-bits:
	    v - (v >> 1) & 0x55555555

	  Then, add the counts stored in each 2-bit pair and store in
	  the 4-bit nibble:

	       v = abcd efgh

	    => (ab + cd)  (ef + gh)

	    v & 0x3333 + ((v >> 2) & 0x3333) satisfies this for 8-bits.

	  There is no chance of overflow, since we are summing a 2 bit value
	  into a 4-bit storage unit.  At this point, we've pop counted 4 bits.

	  Next, sum the nibbles:
	       v = abcd efgh ijkl mnop
	   =>  (abcd + efgh) (ijkl + mnop)

	   v = (v + (v >> 4)) & 0x0F0F) for the above

	   For 32-bits:
	     v = (v + (v >> 4)) & 0x0F0F0F0F

	  Note that the upper nibbles of each byte will have garbage while
	  we add lower nibbles.  However, the lower nibble will not overflow,
	  since we will have pop counted just 8 bits, which
	  can fit into the 4-bits of storage in the lower nibble.  The final
	  mask is necessary to clear the upper nibbles.

	  Finally, we want to sum the lower nibbles of each byte.

	  Observe:
	    v = v0 v1 v2 v3 v4 v5 v6 v7 (each v_i is a nibble)

	  We want count = v0 + v2 + v4 + v6

	    v * 0x01010101
	    = v * (0x01 + 0x0100 + 0x010000 + 0x01000000)
	    = (v7 * 0x10000000 +
	       v6 * 0x01000000 +
	       v5 * 0x00100000 +
	       v4 * 0x00010000 +
	       v3 * 0x00001000 +
	       v2 * 0x00000100 +
	       v1 * 0x00000010 +
	       v0 * 0x00000001) *
	      (0x01 + 0x0100 + 0x010000 + 0x01000000)

	  At this point, v7, v5, v3, v1 are 0
	    = (v0 + v2 + v4 + v6) * 0x01000000 + ...


	  We can then extract the result:

	  c = (((v + (v >> 4)) & 0x0F0F0F0F) * 0x01010101) >> 24

	  NOTE: operator precedence in golang is different from that of C.
	  In particular, bitwise AND has greater precedence than binary +!

	*/
	v = v - ((v >> 1) & 0x55555555)
	v = (v & 0x33333333) + ((v >> 2) & 0x33333333)
	c := (((v + (v >> 4)) & 0x0F0F0F0F) * 0x01010101) >> 24
	return int(c)
}

func PopCount64(v uint64) int {
	v = v - ((v >> 1) & 0x5555555555555555)
	v = (v & 0x3333333333333333) + ((v >> 2) & 0x3333333333333333)
	c := (((v + (v >> 4)) & 0x0F0F0F0F0F0F0F0F) * 0x0101010101010101) >> 56
	return int(c)
}

func PopCountData(data []byte) int {
	count := 0
	ii := 0
	currLen := len(data)

	for currLen >= 8 {
		v := uint64(data[ii]) |
			(uint64(data[ii+1]) << 8) |
			(uint64(data[ii+2]) << 16) |
			(uint64(data[ii+3]) << 24) |
			(uint64(data[ii+4]) << 32) |
			(uint64(data[ii+5]) << 40) |
			(uint64(data[ii+6]) << 48) |
			(uint64(data[ii+7]) << 56)
		count += PopCount64(v)
		ii += 8
		currLen -= 8
	}

	for currLen >= 4 {
		v := uint32(data[ii]) |
			(uint32(data[ii+1]) << 8) |
			(uint32(data[ii+2]) << 16) |
			(uint32(data[ii+3]) << 24)
		count += PopCount32(v)
		ii += 4
		currLen -= 4
	}

	for kk := 0; kk < currLen; kk++ {
		count += PopCount32(uint32(data[ii+kk]))
	}

	return count
}

func popCountGeneric(data []byte) int {
	sum := 0
	for _, v := range data {
		sum += PopCount32(uint32(v))
	}
	return sum
}
