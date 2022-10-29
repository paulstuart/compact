package compact

import (
	"encoding/binary"
	"fmt"
	"math"
	"strconv"
)

type HoldsByte byte

type HoldsF16 F16

func (h HoldsByte) Encode(b []byte) error {
	b[0] = byte(h)
	return nil
}

func (h *HoldsByte) Decode(b []byte) error {
	*h = HoldsByte(b[0])
	return nil
}

func (h *HoldsByte) Input(s string) error {
	v, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	if v > 255 {
		return fmt.Errorf("boo - %w", ErrExceeds)
	}
	*h = HoldsByte(v)
	return nil
}

// HoldsF16

func (h HoldsF16) Encode(b []byte) error {
	binary.LittleEndian.PutUint16(b, uint16(h))
	return nil
}

func (h *HoldsF16) Decode(b []byte) error {
	*h = HoldsF16(binary.LittleEndian.Uint16(b))
	return nil
}

func (h *HoldsF16) Input(s string) error {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return fmt.Errorf("parse fail for F16 %q -- %w", s, err)
	}
	return (*F16)(h).FromFloat64(v)
}

type HoldsI16 uint16

func (h HoldsI16) Encode(b []byte) error {
	binary.LittleEndian.PutUint16(b, uint16(h))
	return nil
}

func (h *HoldsI16) Decode(b []byte) error {
	*h = HoldsI16(b[0])
	return nil
}

func (h HoldsI16) Input(s string) error {
	var v int64
	if _, err := fmt.Scanf("%d", &v); err != nil {
		return err
	}
	if v > math.MaxUint16 {
		return fmt.Errorf("boo hoo - %w", ErrExceeds)
	}
	h = HoldsI16(v)
	return nil
}

type HoldsI32 uint32

func (h HoldsI32) Encode(b []byte) error {
	binary.LittleEndian.PutUint32(b, uint32(h))
	// fmt.Println("PUT:", h)
	return nil
}

func (h *HoldsI32) Decode(b []byte) error {
	*h = HoldsI32(binary.LittleEndian.Uint32(b))
	// fmt.Printf("GET (%v): %d", b, h)
	return nil
}

func (h *HoldsI32) Input(s string) error {
	v, err := strconv.Atoi(s)
	if err != nil {
		return fmt.Errorf("parse fail for %q -- %w", s, err)
	}
	// if _, err := fmt.Scanf("%d", &v); err != nil {
	// 	return fmt.Errorf("parse failed for %q -- %w", s, err)
	// }
	if v > math.MaxUint16 {
		return fmt.Errorf("boo hoo - %w", ErrExceeds)
	}
	// fmt.Printf("FROM %q to %d\n", s, v)
	*h = HoldsI32(v)
	return nil
}

type HoldsI64 int64

func (h HoldsI64) Encode(b []byte) error {
	binary.LittleEndian.PutUint32(b, uint32(h))
	// fmt.Println("PUT:", h)
	return nil
}

func (h *HoldsI64) Decode(b []byte) error {
	*h = HoldsI64(binary.LittleEndian.Uint64(b))
	// fmt.Printf("GET (%v): %d", b, h)
	return nil
}

func (h *HoldsI64) Input(s string) error {
	v, err := strconv.Atoi(s)
	if err != nil {
		return fmt.Errorf("parse fail for %q -- %w", s, err)
	}
	// if _, err := fmt.Scanf("%d", &v); err != nil {
	// 	return fmt.Errorf("parse failed for %q -- %w", s, err)
	// }
	if v > math.MaxUint16 {
		return fmt.Errorf("boo hoo - %w", ErrExceeds)
	}
	// fmt.Printf("FROM %q to %d\n", s, v)
	*h = HoldsI64(v)
	return nil
}

type HoldsText struct {
	size int
	text []byte
}

func NewText(size int) HoldsText {
	return HoldsText{
		size,
		make([]byte, size),
	}
}

// String assumes that the text is padded with spaces
func (h HoldsText) String() string {
	return "XX::" + string(h.text)
}

func (h HoldsText) GoString() string {
	return fmt.Sprintf("%+x", h)
}

func (h HoldsText) Encode(b []byte) error {
	copy(b, h.text)
	return nil
}

func (h *HoldsText) Decode(b []byte) error {
	copy(h.text, b)
	// pad empty columns with spaces
	for i := h.size - 1; i >= 0; i-- {
		if h.text[i] != 0 {
			break
		}
		h.text[i] = ' '
	}
	return nil
}

func (h *HoldsText) Input(s string) error {
	if len(s) > h.size {
		return fmt.Errorf("d'oh too big: %w", ErrExceeds)
	}
	n := copy(h.text, []byte(s))
	for i := n; i < h.size; i++ {
		h.text[i] = ' '
	}
	// for i := len(s); i < h.size; i++ {
	// 	h.text[i] = ' '
	// }
	return nil
}

/*
func (fp *FireJulySrc) Encode(buf []byte) error {
	if len(buf) < fp.Size() {
		log.Printf("buffer size: %d -- we need: %d", len(buf), fp.Size())
		return io.EOF
	}
	const off = 4
	idx := 0
	binary.LittleEndian.PutUint32(buf[idx:], math.Float32bits(float32(fp.Lat)))
	idx += 4
	binary.LittleEndian.PutUint32(buf[idx:], math.Float32bits(float32(fp.Lon)))
	idx += 4
	binary.LittleEndian.PutUint32(buf[idx:], fp.Score)
	idx += 4
	binary.LittleEndian.PutUint32(buf[idx:], math.Float32bits(float32(fp.Baseline_avg)))
	idx += 4
	binary.LittleEndian.PutUint32(buf[idx:], math.Float32bits(float32(fp.BP)))
	idx += 4
	binary.LittleEndian.PutUint32(buf[idx:], math.Float32bits(float32(fp.FLEP4)))
	idx += 4
	binary.LittleEndian.PutUint32(buf[idx:], math.Float32bits(float32(fp.FLEP8)))
	idx += 4
	binary.LittleEndian.PutUint32(buf[idx:], uint32(fp.Baseline_avg_counts))
	idx += 4
	binary.LittleEndian.PutUint32(buf[idx:], uint32(fp.Baseline_avg_magnitude))
	idx += 4
	binary.LittleEndian.PutUint32(buf[idx:], math.Float32bits(float32(fp.Future_avg_2050_Counts)))
	return nil
}

func (ss *FireJulySrc) Decode(buf []byte) error {
	if len(buf) < ss.Size() {
		log.Printf("buffer size: %d -- we need: %d", len(buf), ss.Size())
		return io.EOF
	}
	idx := 0
	const off = 4
	ss.Lat = GeoType(math.Float32frombits(	return (*F16)(h).FromFloat64(v)
(buf[idx:])))
	idx += off
	ss.Lon = GeoType(math.Float32frombits(binary.LittleEndian.Uint32(buf[idx:])))
	idx += off

	ss.Score = binary.LittleEndian.Uint32(buf[idx:])
	idx += off
	ss.Baseline_avg = decimal_5(math.Float32frombits(binary.LittleEndian.Uint32(buf[idx:])))
	idx += 4
	ss.BP = decimal_4(math.Float32frombits(binary.LittleEndian.Uint32(buf[idx:])))
	idx += 4
	ss.FLEP4 = decimal_4(math.Float32frombits(binary.LittleEndian.Uint32(buf[idx:])))
	idx += 4
	ss.FLEP8 = decimal_4(math.Float32frombits(binary.LittleEndian.Uint32(buf[idx:])))
	idx += 4
	ss.Baseline_avg_counts = int32(binary.LittleEndian.Uint32(buf[idx:]))
	idx += 4
	ss.Baseline_avg_magnitude = int32(binary.LittleEndian.Uint32(buf[idx:]))
	idx += 4
	ss.Future_avg_2050_Counts = decimal_5(math.Float32frombits(binary.LittleEndian.Uint32(buf[idx:])))
	return nil
}

*/
