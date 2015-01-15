// Copyright 2012, Kevin Ko <kevin@faveset.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// +build amd64,!appengine

#include "textflag.h"

TEXT ·checkHaveSsse3(SB),7,$0
	XORQ AX, AX
	INCL AX
	CPUID
	SHRQ $9, CX
	ANDQ $1, CX
	MOVB CX, ret+0(FP)
	RET

/*
Aligned reads:
16-byte aligned: 0xffffffff fffffff0

Determine amount to consume in order to become aligned.

  amount = addr & 0xf

If CX < amount, perform a partial read.
If CX > amount, perform a partial read.  Then, switch over to aligned reads
until it is necessary to perform partial reads again.

Downsides:
  - This is only meaningful for large pop count inputs (several 16-byte words
  in length).
*/

/*
  Define the dest array that contains the pop count for the corresponding
  index value:
    dest[16] = {
	    0,  // 0  = 0000
            1,  // 1  = 0001
	    1,  // 2  = 0010
	    2,  // 3  = 0011
	    1,  // 4  = 0100
	    2,  // 5  = 0101
	    2,  // 6  = 0110
	    3,  // 7  = 0111
	    1,  // 8  = 1000
            2,  // 9  = 1001
	    2,  // 10 = 1010
	    3,  // 11 = 1011
	    2,  // 12 = 1100
	    3,  // 13 = 1101
	    3,  // 14 = 1110
	    4,  // 15 = 1111
    }
    We can store dest as two 64-bit immediates, each holding 8 bytes
    destLow  = 0x0302020102010100  // dest[0...7]
    destHigh = 0x0403030203020201  // dest[8...15]

  Place the nibbles in the source array:
    src[16] = {n0, n1, ..., n15}

  Be sure to mask just the lower 4 bits, since PSHUFB will set an entry to 0
  if the most significant bit is set:
    mask = 0x0f0f0f0f0f0f0f0f

  PSHUFB uses the lower 4 bits of each src[i] to permute dest:
    dest[i] = dest[src[i]{0...3} & mask]

  After one round, dest will hold the population count for each nibble.
*/

DATA ·maskArray<>+0(SB)/8, $0x0f0f0f0f0f0f0f0f
DATA ·maskArray<>+8(SB)/8, $0x0f0f0f0f0f0f0f0f
GLOBL ·maskArray<>(SB), RODATA, $16

DATA ·popCountArray<>+0(SB)/8, $0x0302020102010100
DATA ·popCountArray<>+8(SB)/8, $0x0403030203020201
GLOBL ·popCountArray<>(SB), RODATA, $16

// func popCountSsse3(v []byte) int
TEXT ·popCountSsse3(SB),4,$0-32
	// 0(FP) v
	// 8(FP) len(v)
	// 16(FP) cap(v)
	// 24(FP) rval
	// 32(FP) ret

	// SI = v.  This points to the current position in v.
	MOVQ v_base+0(FP), SI
	// CX = remaining_len
	MOVQ v_len+8(FP), CX

	// Load the pre-calculated pop counts into X0 for the duration of this
	// procedure.
	MOVOU ·popCountArray<>(SB), X0

	// Load the mask for lower nibbles into X5.
	MOVOU ·maskArray<>(SB), X5

	// X7 is always 0 for convenience.
	PXOR X7, X7

	// X6[0] accumulates the result.
	PXOR X6, X6

	CMPQ CX, $16
	JL cleanup

loop:
	// Process 128-bits at a time.  This may be unaligned.
	MOVOU (SI), X1

	// Prepare the PSHUFB control mask.  X1 holds the lower nibbles.
	MOVO X1, X2
	// X2 holds the upper nibbles.
	PSRLW $4, X2

	// Mask the upper nibbles to keep PSHUFB from setting values to 0
	// errantly.
	PAND X5, X1
	PAND X5, X2

	// Prepare the PSHUFB dest for manipulation.  Use copies to preserve X0.
	MOVO X0, X3
	MOVO X0, X4

	// X3 holds the pop counts of nibbles in X1.
	// PSHUFB X1, X3
	BYTE $0x66
	LONG $0xD900380F
	// X4 holds the pop counts of nibbles X2.
	// PSHUFB X2, X4
	BYTE $0x66
	LONG $0xE200380F

	// X4 += X3
	PADDB X3, X4

	// NOTE: you can unroll and keep accumulating in X4 up to
	// 255/8 = 31 times before performing an absolute difference calc.

	// Sum the popcounts in X4 using sum of absolute differences with
	// the 0 vector (X7).
	// X4 = (highSum, lowSum)
	PSADBW X7, X4

	MOVHLPS X4, X3
	// Update the result; we only care about the lower word.
	PADDW X4, X6
	PADDW X3, X6

	SUBQ $16, CX
	JZ done
	ADDQ $16, SI

	CMPQ CX, $16
	JGE loop

cleanup:
	// Clean-up smaller than 128-bits by partially loading an xmm reg.
	// CX < 16 at this point.

	// We'll copy the < 128-bit remainder into X1 and
	// the final 64-bits will be shifted into AX.
	PXOR X1, X1
	XORQ AX, AX

	CMPB CX, $8
	JL cleanup8

	// Read 64-bits.
	MOVQ (SI), X1

	SUBB $8, CX
	JZ cleanup0
	ADDQ $8, SI

cleanup8:
	// CX < 8
	CMPB CX, $4
	JL cleanup4

	// Read 32-bits.
	MOVL (SI), AX

	SUBB $4, CX
	JZ cleanup0
	ADDQ $4, SI

cleanup4:
	// CX < 4
	CMPB CX, $2
	JL cleanup2

	// Shift the 16-bit word into AX.
	SHLQ $16, AX
	MOVW (SI), AX

	SUBB $2, CX
	JZ cleanup0
	ADDQ $2, SI

cleanup2:
	CMPB CX, $0
	JE cleanup0

	// CX < 2 => CX == 1
	// Shift the final byte into AX.
	SHLQ $8, AX
	MOVB (SI), AX

cleanup0:

	// AX holds shifted in partial bytes.
	MOVQ AX, X2
	// Now merge with X1.
	MOVLHPS X2, X1

	MOVO X1, X2
	// X2 holds the upper nibbles.
	PSRLW $4, X2

	// Mask upper nibbles.
	PAND X5, X1
	PAND X5, X2

	MOVO X0, X4

	// PSHUFB X1, X0.  At this point, it's safe to sacrifice X0.
	// X0 holds the pop counts of the lower nibbles from X1.
	BYTE $0x66
	LONG $0xC100380F

	// PSHUFB X2, X4.  X4 holds the pop counts of the upper nibbles from X2.
	BYTE $0x66
	LONG $0xE200380F

	PADDB X0, X4

	// Sum the pop counts.  X4 = (highSum, lowSum)
	PSADBW X7, X4

	MOVHLPS X4, X3
	// Update the result in X6.  We only care about the lower quad.
	PADDW X4, X6
	PADDW X3, X6

done:
	MOVQ X6, ret+24(FP)
	RET
