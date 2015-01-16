// Copyright 2012, Kevin Ko <kevin@faveset.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// +build 386 arm appengine

package popcount

func PopCount(data []byte) int {
	return popCountGeneric(data)
}
