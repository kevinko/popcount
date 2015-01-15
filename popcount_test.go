// Copyright 2012, Kevin Ko <kevin@faveset.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package popcount

import (
	"fmt"
	"math/rand"
	"testing"
)

func benchPopCountRand(b *testing.B, size int) {
	b.StopTimer()

	v := makeBuff(size)
	if PopCount(v) != PopCountData(v) {
		b.Error("!=")
	}

	b.StartTimer()

	for ii := 0; ii < b.N; ii++ {
		PopCount(v)
	}
}

func benchPopCountDataRand(b *testing.B, size int) {
	b.StopTimer()

	v := makeBuff(size)
	if PopCount(v) != PopCountData(v) {
		b.Error("!=")
	}

	b.StartTimer()

	for ii := 0; ii < b.N; ii++ {
		PopCountData(v)
	}
}

func Benchmark_PopCount16(b *testing.B) {
	benchPopCountRand(b, 16)
}

func Benchmark_PopCount20(b *testing.B) {
	benchPopCountRand(b, 20)
}

func Benchmark_PopCount250(b *testing.B) {
	benchPopCountRand(b, 250)
}

func Benchmark_PopCountData16(b *testing.B) {
	benchPopCountDataRand(b, 16)
}

func Benchmark_PopCountData20(b *testing.B) {
	benchPopCountDataRand(b, 20)
}

func Benchmark_PopCountData250(b *testing.B) {
	benchPopCountDataRand(b, 250)
}

func Benchmark_PopCountSimple(b *testing.B) {
	b.StopTimer()

	v := []byte{0xde, 0xad, 0xbe, 0xef}
	vInt := uint32(0xefbeadde)
	if PopCount32(vInt) != PopCount(v) {
		b.Error("!=")
	}

	b.StartTimer()

	for ii := 0; ii < b.N; ii++ {
		PopCount(v)
	}
}

func Benchmark_PopCountDataSimple(b *testing.B) {
	b.StopTimer()

	v := []byte{0xde, 0xad, 0xbe, 0xef}
	if PopCountData(v) != PopCount(v) {
		b.Error("!=")
	}

	b.StartTimer()

	for ii := 0; ii < b.N; ii++ {
		PopCountData(v)
	}
}

func Benchmark_PopCount32(b *testing.B) {
	b.StopTimer()

	v := []byte{0xde, 0xad, 0xbe, 0xef}
	vInt := uint32(0xefbeadde)
	if PopCount32(vInt) != PopCount(v) {
		b.Error("!=")
	}

	b.StartTimer()

	for ii := 0; ii < b.N; ii++ {
		PopCount32(vInt)
	}
}

func Benchmark_PopCount64(b *testing.B) {
	b.StopTimer()

	v := []byte{0xde, 0xad, 0xbe, 0xef}
	vInt := uint64(0xefbeadde)
	if PopCount64(vInt) != PopCount(v) {
		b.Error("!=")
	}

	b.StartTimer()

	for ii := 0; ii < b.N; ii++ {
		PopCount64(vInt)
	}
}

func makeBuffFixed(size int, v byte) []byte {
	buff := make([]byte, 0)
	for ii := 0; ii < size; ii++ {
		buff = append(buff, v)
	}
	return buff
}

func makeBuff(size int) []byte {
	r := rand.New(rand.NewSource(0))

	buff := make([]byte, 0)
	for ii := 0; ii < size; ii++ {
		buff = append(buff, byte(r.Int()))
	}
	return buff
}

func Test_PopCount(t *testing.T) {
	for ii := 0; ii < 256; ii++ {
		x := []byte{byte(ii)}
		if r := PopCount(x); r != PopCount32(uint32(ii)) {
			t.Error(fmt.Sprintf("0x%x %d, expected %d", ii, r, PopCount32(uint32(ii))))
		}

		x = []byte{byte(ii), 0xff, 0xff, 0xff,
			0xff, byte(ii), 0xff, 0xff,
			0xff, 0xff, byte(ii), 0xff,
			0xff, 0xff, 0xff, byte(ii),
		}
		if r := PopCount(x); r != 4*PopCount32(uint32(ii))+96 {
			t.Error(fmt.Sprintf("0x%x %d", ii, r))
		}

		// Exercise all partial paths.
		x = []byte{byte(ii), 0xff, 0xff, 0xff,
			0xff, byte(ii), 0xff, 0xff,
			0xff, 0xff, byte(ii), 0xff,
			0xff, 0xff, 0xff, byte(ii),
			byte(ii), byte(ii), byte(ii), byte(ii),
			byte(ii), byte(ii), byte(ii), byte(ii),
			byte(ii), byte(ii), byte(ii), byte(ii),
			byte(ii), byte(ii), byte(ii),
		}
		if r := PopCount(x); r != 19*PopCount32(uint32(ii))+96 {
			t.Error(fmt.Sprintf("0x%x %d", ii, r))
		}
	}
	x := []byte{0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff}
	if r := PopCount(x); r != 128 {
		t.Error(fmt.Sprintf("0x%x %d", x, r))
	}
	x = []byte{0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff}
	if r := PopCount(x); r != 256 {
		t.Error(fmt.Sprintf("0x%x %d", x, r))
	}

	// Test variable length clean-up.
	for ii := 0; ii < 2048; ii++ {
		x = makeBuffFixed(ii, 0xff)
		if r := PopCount(x); r != PopCountData(x) {
			t.Error(fmt.Sprintf("%d 0x%x, got %d, expected %d", ii, x, r, PopCountData(x)))
		}

		x = makeBuff(ii)
		if r := PopCount(x); r != PopCountData(x) {
			t.Error(fmt.Sprintf("%d 0x%x", ii, x))
		}
	}
}

func Test_PopCountData(t *testing.T) {
	for ii := 0; ii < 256; ii++ {
		x := []byte{byte(ii)}
		if r := PopCountData(x); r != PopCount32(uint32(ii)) {
			t.Error(fmt.Sprintf("0x%x %d", ii, r))
		}

		x = []byte{byte(ii), 0xff, 0xff, 0xff,
			0xff, byte(ii), 0xff, 0xff,
			0xff, 0xff, byte(ii), 0xff,
			0xff, 0xff, 0xff, byte(ii),
		}
		if r := PopCountData(x); r != 4*PopCount32(uint32(ii))+96 {
			t.Error(fmt.Sprintf("0x%x %d", ii, r))
		}

		// Exercise all partial paths.
		x = []byte{byte(ii), 0xff, 0xff, 0xff,
			0xff, byte(ii), 0xff, 0xff,
			0xff, 0xff, byte(ii), 0xff,
			0xff, 0xff, 0xff, byte(ii),
			byte(ii), byte(ii), byte(ii), byte(ii),
			byte(ii), byte(ii), byte(ii), byte(ii),
			byte(ii), byte(ii), byte(ii), byte(ii),
			byte(ii), byte(ii), byte(ii),
		}
		if r := PopCountData(x); r != 19*PopCount32(uint32(ii))+96 {
			t.Error(fmt.Sprintf("0x%x %d", ii, r))
		}
	}
	x := []byte{0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff}
	if r := PopCountData(x); r != 128 {
		t.Error(fmt.Sprintf("0x%x %d", x, r))
	}
	x = []byte{0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff}
	if r := PopCountData(x); r != 256 {
		t.Error(fmt.Sprintf("0x%x %d", x, r))
	}
}

func Test_PopCount32(t *testing.T) {
	if PopCount32(0x0) != 0 {
		t.Error("0")
	}
	if PopCount32(0x1) != 1 {
		t.Error("1")
	}
	if PopCount32(0x2) != 1 {
		t.Error("2")
	}
	if PopCount32(0x3) != 2 {
		t.Error("3")
	}
	if PopCount32(0xf0f0f0f0) != 16 {
		t.Error("0xf0f0f0f0")
	}
	if p := PopCount32(0x0001f); p != 5 {
		t.Error("0x1f", p)
	}
	if PopCount32(0xffff0000) != 16 {
		t.Error("0xffff0000")
	}
	if PopCount32(0xffffffff) != 32 {
		t.Error("0xffffffff")
	}
}

func Test_PopCount64(t *testing.T) {
	if PopCount64(0x0) != 0 {
		t.Error("0")
	}
	if PopCount64(0x1) != 1 {
		t.Error("1")
	}
	if PopCount64(0x2) != 1 {
		t.Error("2")
	}
	if PopCount64(0x3) != 2 {
		t.Error("3")
	}
	if PopCount64(0xf0f0f0f0) != 16 {
		t.Error("0xf0f0f0f0")
	}
	if p := PopCount64(0x0001f); p != 5 {
		t.Error("0x1f", p)
	}
	if PopCount64(0xffff0000) != 16 {
		t.Error("0xffff0000")
	}
	if PopCount64(0xffffffff) != 32 {
		t.Error("0xffffffff")
	}
}
