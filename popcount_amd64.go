// Copyright 2012, Kevin Ko <kevin@faveset.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// +build amd64,!appengine

package popcount

var haveSsse3 = checkHaveSsse3()

// Tests for Supplemental SSE3 support.
func checkHaveSsse3() bool

func PopCount(v []byte) int {
	if haveSsse3 {
		return popCountSsse3(v)
	}
	return popCountGeneric(v)
}

func popCountSsse3(v []byte) int
