package main

import (
	"bytes"
	"crypto/elliptic"
	"fmt"
	"math/big"
)

type ECSuite struct {
	C elliptic.Curve
	N *big.Int
}

func NewP256Suite() *ECSuite {
	curve := elliptic.P256()
	N := curve.Params().N
	return &ECSuite{
		C: curve,
		N: N,
	}
}

// modN reduces s mod N (in non integer representation)
func modN(s, N *big.Int) *big.Int {
	r := new(big.Int).Mod(s, N)
	// if sign is negative, add N
	if r.Sign() < 0 {
		r.Add(r, N)
	}
	return r
}

// baseMul returns z.G (scalar multiplication by the base point)
func baseMul(c elliptic.Curve, z *big.Int) (x, y *big.Int) {
	return c.ScalarBaseMult(z.Bytes())
}

// pointAdd returns A + B returns the affine coordinates
func pointAdd(c elliptic.Curve, ax, ay, bx, by *big.Int) (x, y *big.Int) {
	return c.Add(ax, ay, bx, by)
}

// eqPoints checks the affine equality of 2 points
func eqPoints(ax, ay, bx, by *big.Int) bool {
	if ax == nil || ay == nil || bx == nil || by == nil {
		return false
	}
	return ax.Cmp(bx) == 0 && ay.Cmp(by) == 0
}

// ------------------ EC point encoding ----------------------

// encodePoint returns SEC1 uncompressed encoding (0x04 || x || y).
// 0x04: prefix byte as defined by Standards for Efficient Cryptography
//
//	it means what you see is what you get (the 0x04 part), here are both the points (followed by points)
//
// 0x04 || x || y: here are both the x and y coordinates in their entirety
func encodePoint(c elliptic.Curve, x, y *big.Int) []byte {
	return elliptic.Marshal(c, x, y)
}

// decodePoint parses a SEC1 uncompressed point into affine (x,y).
func decodePoint(c elliptic.Curve, b []byte) (x, y *big.Int, ok bool) {
	x, y = elliptic.Unmarshal(c, b)
	return x, y, x != nil && y != nil
}

// ---------------------------------------------- 2.G == G + G
func main() {
	s := NewP256Suite()

	one := modN(big.NewInt(1), s.N) // 1.G = G
	// scalar 2 mod n
	two := modN(big.NewInt(2), s.N)

	// Left side: 2.G
	lx, ly := baseMul(s.C, two)

	// right side: G+G (1.G)
	one := modN(big.NewInt(1), s.N)
	gx, gy := baseMul(s.C, one)             // G
	rx, ry := pointAdd(s.C, gx, gy, gx, gy) // G+G

	// check equals
	equals := eqPoints(lx, ly, rx, ry)
	if !equals {
		fmt.Println("not equals: x,y")
	}

	fmt.Println("lx:", lx, "\tly:", ly)
	fmt.Println("rx:", rx, "\try:", ry)

	LeftEncoded := encodePoint(s.C, lx, ly)
	RightEncoded := encodePoint(s.C, rx, ry)

	if !bytes.Equal(LeftEncoded, RightEncoded) {
		fmt.Println("not equals: leftEncoded, rightEncoded")
	}
}
