// Copyright (c) 2022 Hirotsuna Mizuno. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package tdma_test

import (
	"math"
	"testing"

	"github.com/tunabay/go-tdma"
)

func TestMatrix_TDMA_1(t *testing.T) {
	equal := func(x, y []float64) bool {
		const eps = 1e-3
		if len(x) != len(y) {
			return false
		}
		for i, xe := range x {
			if eps < math.Abs(y[i]-xe) {
				return false
			}
		}
		return true
	}
	tdata := [][]float64{
		// #1
		{
			2, 1,
			1, 2, 1,
			1, 2, 1,
			1, 2,
		},
		{4, 8, 12, 11},
		{1, 2, 3, 4},

		// #2
		{
			3, 1,
			1, 4, 2,
			2, 5,
		},
		{5, 15, 19},
		{1, 2, 3},

		// #3
		{
			1, 2,
			3, 4, 5,
			6, 7, 8,
			9, 1, 2,
			3, 4,
		},
		{1, 2, 3, 4, 5},
		{-0.7229, 0.8614, 0.1446, -0.3976, 1.5482},

		// #4
		{
			6, 0,
			1, 4, 1,
			1, 4, 1,
			1, 4, 1,
			1, 4, 1,
			1, 4, 1,
			0, 6,
		},
		{0, 1, 2, -6, 2, 1, 0},
		{0, 0, 1, -2, 1, 0, 0},

		// #5
		{
			2, 1,
			1, 3, 2,
			1, 3, 1,
			7, 2, 6,
			6, 2, 1,
			3, 4, 3,
			8, 1, 5,
			6, 2, 7,
			5, 4, 3,
			4, 5,
		},
		{1, 2, 6, 34, 10, 1, 4, 22, 25, 3},
		{1, -1, 2, 1, 3, -2, 0, 4, 2, -1},

		// #6 fail
		// TODO: find out how to solve this
		// https://www.scirp.org/pdf/AM_2014021111074341.pdf
		/*
			{
				1, 1,
				1, 1, 10,
				7, 1, 2,
				2, 11, 1,
				2, 3, 7,
				3, 1, 2,
				-1, 2, 2,
				2, 1, 1,
				5, 2, 4,
				1, 5,
			},
			{4, 14, 26, 25, 0, 2, 1, 3, 10, 8},
			{1, 3, 1, 2, 1, -1, 0, 0, 3, 1},
		*/
	}
	for i := 0; i < len(tdata); i += 3 {
		tno := i/3 + 1
		m, err := tdma.New(tdata[i])
		if err != nil {
			t.Fatalf("#%d: invalid test data: %v", tno, err)
		}
		x, err := m.TDMA(tdata[i+1])
		if err != nil {
			t.Errorf("#%d: TDMA failed: %v", tno, err)
			continue
		}
		if !equal(x, tdata[i+2]) {
			t.Errorf("#%d: got %+v, want %+v", tno, x, tdata[i+2])
			continue
		}
		t.Logf("#%d: passed: %+v", tno, x)
	}
}
