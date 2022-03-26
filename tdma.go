// Copyright (c) 2022 Hirotsuna Mizuno. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package tdma

import (
	"errors"
	"fmt"
)

// ErrInvalidTridiagonalMatrix is the error thrown when the matrix data is not
// valid.
var ErrInvalidTridiagonalMatrix = errors.New("invalid tridiagonal matrix")

// ErrTDMA is the error thrown when TDMA operation is failed.
var ErrTDMA = errors.New("TDMA failure")

// Matrix represents a tridiagonal matrix.
// https://en.wikipedia.org/wiki/Tridiagonal_matrix
//
//     (n=4)
//     +-                   -+
//     | e[0] e[1]         0 |
//     | e[2] e[3] e[4]      |
//     |      e[5] e[6] e[7] |
//     |    0      e[8] e[9] |
//     +-                   -+
//
//     (4<=n)
//     +-                                         -+
//     | e[ 0] e[ 1]                             0 |
//     | e[ 2] e[ 3] e[ 4]                         |
//     |                 ...                       |
//     |                   e[3n-7] e[3n-6] e[3n-5] |
//     |     0                     e[3n-4] e[3n-3] |
//     +-                                         -+
type Matrix struct {
	n int       // n x n tridiagonal matrix
	m []float64 // (3n - 2) elements
}

// New creates a new tridiagonal matrix from the slice of elements d.
// The size of d must be exactly (3n - 2) for the n x n matrix.
func New(d []float64) (*Matrix, error) {
	if len(d)%3 != 1 {
		return nil, fmt.Errorf("%w: length %d must 3n-2",
			ErrInvalidTridiagonalMatrix, len(d))
	}
	m := &Matrix{
		n: (len(d)-1)/3 + 1,
		m: d,
	}

	return m, nil
}

// TDMA solves the system of equations M * x = r using TDMA (tridiagonal matrix
// algorithm, aka Thomas algorithm) and returns the x solved. r must have
// exactly n elements for the n x n tridiagonal matrix.
// https://en.wikipedia.org/wiki/Tridiagonal_matrix_algorithm
func (m *Matrix) TDMA(r []float64) ([]float64, error) {
	if len(r) != m.n {
		return nil, fmt.Errorf("%w: r must have exactly %d elements",
			ErrTDMA, m.n)
	}
	c := make([]float64, m.n-1)
	if m.m[0] == 0 {
		return nil, fmt.Errorf("%w: m[0] is zero", ErrTDMA)
	}
	c[0] = m.m[1] / m.m[0]
	for i := 1; i < m.n-1; i++ {
		i3 := i * 3
		if m.m[i3] == m.m[i3-1]*c[i-1] {
			// TODO: find out how to solve this case
			return nil, fmt.Errorf("%w: m[%d] == m[%d]c[%d] == %f",
				ErrTDMA, i3, i3-1, i-1, m.m[i3])
		}
		c[i] = m.m[i3+1] / (m.m[i3] - m.m[i3-1]*c[i-1])
	}
	if i, i3 := m.n-1, (m.n-1)*3; m.m[i3] == m.m[i3-1]*c[i-1] {
		// TODO: find out how to solve this case
		return nil, fmt.Errorf("%w: m[%d] == m[%d]c[%d] == %f",
			ErrTDMA, i3, i3-1, i-1, m.m[i3])
	}
	d := make([]float64, m.n)
	d[0] = r[0] / m.m[0]
	for i := 1; i < m.n; i++ {
		i3 := i * 3
		d[i] = (r[i] - m.m[i3-1]*d[i-1]) / (m.m[i3] - m.m[i3-1]*c[i-1])
	}
	x := make([]float64, m.n)
	x[m.n-1] = d[m.n-1]
	for i := m.n - 2; 0 <= i; i-- {
		x[i] = d[i] - c[i]*x[i+1]
	}

	return x, nil
}

//
func detF(elem []float64, n int) float64 {
	switch n {
	case 0:
		return elem[0]
	case 1:
		return elem[3]*elem[0] - elem[2]*elem[1]
	}
	idx := n * 3
	return elem[idx]*detF(elem, n-1) - elem[idx-1]*elem[idx-2]*detF(elem, n-2)
}

// Determinant calculates the determinant of the tridiagonal matrix.
func (m *Matrix) Determinant() float64 { return detF(m.m, m.n-1) }
