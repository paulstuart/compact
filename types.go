package compact

import (
	"encoding/binary"
	"fmt"
	"math"
	"strconv"
)

type DataType int

const (
	DT_Unknown DataType = iota
	DT_Byte
	DT_F1
	DT_F16
	DT_F32
	DT_F64
	DT_X16
	DT_X32
	DT_X64
)

type HoldsByte byte

// var pp interface{} = &P{}

// validate Record interface is met
var (
	_ Record = (*HoldsByte)(nil)
	_ Record = (*FP1)(nil)
	_ Record = (*F16)(nil)
	_ Record = (*F32)(nil)
	_ Record = (*F64)(nil)
)

// var _ Record = (*HoldsF16)(nil)  //.Record()
// var _ MyInterface = (*MyType)(nil)
type RecFun func() Record

func getF16() *F16 {
	var f F16
	return &f
}

func foobar() Record {
	return getF16()
}

var (
	dtMap = map[DataType]RecFun{
		DT_Byte: foobar,
	}
)

// var (
// 	dtMap = map[DataType]RecFun{
// 		DT_Byte: func() {
// 			return &F16{}
// 		},
// 	}
// )

type RecMux interface {
	Mux() Record
}

type HoldsF16 F16

func (h HoldsByte) Size() int {
	return 1
}

func (h HoldsByte) String() string {
	return strconv.Itoa(int(h))
}

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

// type HoldsI16 uint16

// func (h HoldsI16) Encode(b []byte) error {
// 	binary.LittleEndian.PutUint16(b, uint16(h))
// 	return nil
// }

// func (h *HoldsI16) Decode(b []byte) error {
// 	*h = HoldsI16(b[0])
// 	return nil
// }

// func (h HoldsI16) Input(s string) error {
// 	var v int64
// 	if _, err := fmt.Scanf("%d", &v); err != nil {
// 		return err
// 	}
// 	if v > math.MaxUint16 {
// 		return fmt.Errorf("boo hoo - %w", ErrExceeds)
// 	}
// 	h = HoldsI16(v)
// 	return nil
// }

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

// =============================================================================================== //

type F32 float32

const F32_max = math.MaxFloat32
const _F32_size = 4

func (f F32) Size() int {
	return _F32_size
}

func (f *F32) Decode(b []byte) error {
	if len(b) < _F32_size {
		return ErrTooSmall
	}
	*f = F32(math.Float32frombits(binary.LittleEndian.Uint32(b)))
	return nil
}

func (f F32) Encode(b []byte) error {
	binary.LittleEndian.PutUint32(b, math.Float32bits(float32(f)))
	return nil
}

func (f *F32) FromFloat64(v float64) error {
	*f = F32(v)
	return nil
}

func (f F32) String() string {
	return strconv.FormatFloat(float64(f), 'f', 5, 32)
}

func (f *F32) Input(s string) error {
	v, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return fmt.Errorf("parse fail for F32 %q -- %w", s, err)
	}
	// return (*F32)(f).FromFloat64(v)
	return f.FromFloat64(v)
}

// ==========================================================================================//

type F64 float64

const F64_max = math.MaxFloat64
const _F64_size = 8

func (f F64) Size() int {
	return _F64_size
}

func (f *F64) Decode(b []byte) error {
	if len(b) < _F64_size {
		return ErrTooSmall
	}
	*f = F64(math.Float64frombits(binary.LittleEndian.Uint64(b)))
	return nil
}

func (f F64) Encode(b []byte) error {
	binary.LittleEndian.PutUint64(b, math.Float64bits(float64(f)))
	return nil
}

func (f *F64) FromFloat64(v float64) error {
	*f = F64(v)
	return nil
}

func (f F64) String() string {
	return strconv.FormatFloat(float64(f), 'f', 9, 32)
}

func (f *F64) Input(s string) error {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return fmt.Errorf("parse fail for F64 %q -- %w", s, err)
	}
	return (*F64)(f).FromFloat64(v)
}

// ===============================================================================/

type NX32 int32

const NX32_max = math.MaxInt32
const _NX32_size = 4

func (f NX32) Size() int {
	return _NX32_size
}

func (f *NX32) Decode(b []byte) error {
	if len(b) < _NX32_size {
		return ErrTooSmall
	}
	*f = NX32(binary.LittleEndian.Uint32(b))
	return nil
}

func (f NX32) Encode(b []byte) error {
	binary.LittleEndian.PutUint32(b, uint32(f))
	return nil
}

func (f NX32) String() string {
	return strconv.FormatInt(int64(f), 32)
}

func (f *NX32) Input(s string) error {
	v, err := strconv.Atoi(s)
	if err != nil {
		return fmt.Errorf("parse fail for NX32 %q -- %w", s, err)
	}
	*f = NX32(v)

	return nil
}

// ================================================================================== //

type NX64 int64

const NX64_max = math.MaxInt64
const _NX64_size = 8

func (f NX64) Size() int {
	return _NX64_size
}

func (f *NX64) Decode(b []byte) error {
	if len(b) < _NX64_size {
		return ErrTooSmall
	}
	*f = NX64(binary.LittleEndian.Uint64(b))
	return nil
}

func (f NX64) Encode(b []byte) error {
	binary.LittleEndian.PutUint64(b, uint64(f))
	return nil
}

func (f NX64) String() string {
	return strconv.FormatInt(int64(f), 64)
}

func (f *NX64) Input(s string) error {
	v, err := strconv.Atoi(s)
	if err != nil {
		return fmt.Errorf("parse fail for NX64 %q -- %w", s, err)
	}
	*f = NX64(v)

	return nil
}
