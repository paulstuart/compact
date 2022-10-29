package compact

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"math"
)

const fp1Mul = 10
const fp1Max = 6553.5

// FP1 allows up to 2 decimal places on a positive number < 655
// Storage is 2 bytes
type FP1 uint16

func (f FP1) Size() int {
	return 2
}

func F32ToFP1(v float32) FP1 {
	if v > fp1Max {
		log.Fatalf("value %.1f exceeds max of %.1f", v, fp1Max)
	}
	return FP1(math.Round(float64(v * fp1Mul)))
}

func (f FP1) Float32() float32 {
	return float32(f) / fp1Mul
}

func (f FP1) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%.1f", f.Float32())
	return buf.Bytes(), nil
}

// we only care about 1 sig dig
// so in a uint16 w/ max of 65535,
// that allows for values up to 655.35
// which should never be exceeded
// TODO: VERIFY THIS LIMIT!!!!!
const f16Mul = 100.0

// F16 allows up to 2 decimal places on a number < 655
// (as we pack) it times 100 as an int
// NOTE: the number is expected to be *positive*
type F16 uint16

const F16_max = math.MaxUint16 / f16Mul

func (f F16) Size() int {
	return 2
}

var (
	ErrTooSmall = errors.New("buffer is too small")
	ErrExceeds  = errors.New("value exceeds max size")
)

func (f *F16) Decode(b []byte) error {
	if len(b) < 2 {
		return ErrTooSmall
	}
	*f = F16(binary.LittleEndian.Uint16(b))
	return nil
}

func F32ToF16(v float32) F16 {
	return F16(math.Round(float64(v * f16Mul)))
}

func (f *F16) FromFloat64(v float64) error {
	if v > F16_max {
		return ErrExceeds
	}
	*f = F16(v * f16Mul)
	return nil
}

func (f F16) Float32() float32 {
	return float32(f) / f16Mul
}

// Odds represents a number between 0.0 and 1.0
// with accuracy of 4 decimal places
// It is able to do so in 2 bytes, for compactness
type Odds uint16

const oddity = 65535

func (oo *Odds) Set(v float32) {
	*oo = Odds(v * oddity)
}

func (oo *Odds) Get() float32 {
	return float32(*oo) / oddity
}

func SetOdds(v float32) Odds {
	return Odds(v * oddity)
}

func GetOdds(oo Odds) float32 {
	return float32(oo) / oddity
}

func (oo *Odds) String() string {
	return fmt.Sprintf("%.4f", oo.Get())
}

func (oo *Odds) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%.4f", oo.Get())
	return buf.Bytes(), nil
}

func (oo *Odds) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%.4f", oo.Get())
	return buf.Bytes(), nil
}
