package main

import (
	"crypto/elliptic"
	"crypto/sha256"
	"hash"
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

// Transcript : --------- Step:2  Fiat-Shamir Transcript (deterministic)
type Transcript struct {
	h hash.Hash
}

// NewTranscript constructs a fresh,empty transcript
func NewTranscript() *Transcript {
	return &Transcript{h: sha256.New()}
}

// Absorb ingests a labeled sequence of byte slices.
// Each label and each part is length-prefixed to avoid ambiguity.
// No randomness(nothing up my sleeves); fully determined by inputs and their order.
func (t *Transcript) Absorb(label string, parts ...[]byte) {
	// label length (1 byte, for demo purpose, intentionally kept short)
	if len(label) > 255 {
		panic("label too long")
	}
	t.h.Write([]byte{byte(len(label))})
	t.h.Write([]byte(label))

	// number of parts (1 byte, demo simplicity)
	if len(parts) > 255 {
		panic("too many parts; keep small for demo")
	}
	t.h.Write([]byte{byte(len(parts))})

	// each part: 4-byte big-endian length + bytes
	for _, p := range parts {
		lp := []byte{
			byte(len(p) >> 24),
			byte(len(p) >> 16),
			byte(len(p) >> 8),
			byte(len(p)),
		}
		t.h.Write(lp)
		t.h.Write(p)
	}
}

/*
Step 2 goal:
	- Add a transparent Fiat–Shamir transcript:
  	- Absorb(label, data...) with explicit length-prefixing
  	- ChallengeScalar(mod n): deterministically derive e ∈ {1..n-1}
	- No randomness, no hidden seeds.

We’ll *use* the EC helpers from Step 1 to create two points (G and 2G),
absorb (domain, msg, P=G, R=2G), and show that flipping ANY byte changes e.
*/

// ----------------- task 2: binding sensitivity
func main() {
	s := NewP256Suite()

	one := modN(big.NewInt(1), s.N) // 1.G = G
	two := modN(big.NewInt(2), s.N) // 2.G

	gx, gy := baseMul(s.C, one) // G
	hx, hy := baseMul(s.C, two) // 2G

	P := encodePoint(s.C, gx, gy) // P
	R := encodePoint(s.C, hx, hy) // R

	domain := []byte("qfold-ec/v1")
	message := []byte("register: user=haseeb, nonce=42")

	// in Schnorr/qFold-style proofs, the verifier computes e = H(transcript byes)
	// and then checks z.G = R + e.P

	// Transcript #1
	tr1 := NewTranscript()
	tr1.Absorb("domain", domain)
	tr1.Absorb("message", message)
	tr1.Absorb("P", P)
	tr1.Absorb("R", R)

}
