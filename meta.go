package compact

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"unsafe"

	"golang.org/x/exp/constraints"
)

type Number interface {
	constraints.Integer | constraints.Float
}

type Record interface {
	Decode([]byte) error
	Encode([]byte) error
	Size() int
}

func FMe[T Number](b []byte, num T) error {
	return nil
}

// type Kind int
type Size int
type What struct {
	kind Kind
	size Size
}

type Handler func([]byte) error

type Runner interface {
	Handle([]byte)
	Columns() []string
	MarshalJSON() ([]byte, error)
}

/*
	build list of decoders with offset
	 - decode bytes
	 - assign to array of float64
	 -

*/
func buildHandler(ss ...What) Handler {
	// var fns []Record
	// var off int
	// for i, s := range ss {
	// 	log.Printf("[%d]: %s\n", i, s)
	// 	fn := func(b []byte) error {
	// 		// Decode
	// 		fns := append(fns, Record{})
	// 	}
	// }
	/*
		What := range through types
			decode value from buf[off:]
			off += size
	*/
	return func(b []byte) error {
		return nil
	}
}

func dump(obj interface{}) {
	json.NewEncoder(os.Stdout).Encode(obj)
}

type Kind int

const (
	DTUnknown Kind = iota
	DTI64     Kind = 1 << iota
	DTI32
	DTI16
	DTF64
	DTF32
	DTF1
	DTF16
	DTByte
	DTOmitEmpty
)

// func init() {
// 	fmt.Printf("UNK: %v 64 (%T): %v\n", DTUnknown, DTI64, DTI64)
// }

type Decipher func([]byte) (interface{}, error)
type Encipher func([]byte, interface{}) error

//Holder manages a binary packed table
// layout on disk is ordered to do word alignment if possible
type Holder struct {
	width   int // binary record width
	columns []string
	sizes   []int
	layout  []int // column order that goes on disk
	offset  []int
	// kinds    []Kind
	decoders  []Decipher
	encoders  []Encipher
	recorders []Record
	omitEmpty []bool
}

// RecordDeck is an instance of a Record
// to be used per individual goroutine
// type RecordDeck struct {
// 	recs []Record
// 	size int // record size
// 	data []byte
// }

type RecordHolder struct {
	hold   *Holder
	values []Record
}

func (h *Holder) NewRecord() *RecordHolder {
	rh := &RecordHolder{
		hold:   h,
		values: make([]Record, len(h.columns)),
	}
	copy(rh.values, h.recorders)
	return rh
}

// Decode converts a byte slice to a record
func (rh *RecordHolder) Decode(b []byte) error {
	/*
		if len(b) < rh.hold.width {
			return fmt.Errorf("size wants %d has %d - %w", rh.hold.width, len(b), ErrTooSmall)
		}
		off := 0
		for i, x := range rh.hold.layout {
			// fn := rh.hold.decoders[x]
			fn := rh.hold.recorders[x]
			val, err := fn(rhy.b[off:])
			if err != nil {
				return fmt.Errorf("[%d/%d] decode fail: %w", i, x, err)
			}
			rh.values[i] = val
			off += rh.hold.sizes[x]
		}
	*/
	return nil
}

func Serialized[T ~int | ~float64](b []byte, x interface{}) error {
	val := int64(b[0])
	fmt.Println(val)
	return nil //val, nil
}

// SerialByte writes the value of x into a byte in b
func SerialByte(b []byte, x interface{}) error {
	/*
		val, ok := x.(int)
		if !ok {
			return fmt.Errorf("%T is not an int", x)
		}
		// val := int64(b[0])
	*/
	return nil
}

// FetchByte reads a byte from the buffer and returns an int
func FetchByte(b []byte) (interface{}, error) {
	val := int64(b[0])
	return val, nil
}

func FetchI16(b []byte) (interface{}, error) {
	f16 := F16(binary.LittleEndian.Uint16(b))
	return float64(f16.Float32()), nil
	// Lat := math.Float64frombits(binary.LittleEndian.Uint64(buf))
}

func FetchI32(b []byte) (interface{}, error) {
	return int32(binary.LittleEndian.Uint32(b)), nil
}

func FetchI64(b []byte) (interface{}, error) {
	return int64(binary.LittleEndian.Uint64(b)), nil
}

func FetchF16(b []byte) (interface{}, error) {
	f16 := F16(binary.LittleEndian.Uint16(b))
	return float64(f16.Float32()), nil
	// Lat := math.Float64frombits(binary.LittleEndian.Uint64(buf))
}

func FetchF1(b []byte) (interface{}, error) {
	f16 := FP1(binary.LittleEndian.Uint16(b))
	return float64(f16.Float32()), nil
}

func FetchF32(b []byte) (interface{}, error) {
	Lat := math.Float32frombits(binary.LittleEndian.Uint32(b))
	return Lat, nil
}

func FetchF64(b []byte) (interface{}, error) {
	Lat := math.Float64frombits(binary.LittleEndian.Uint64(b))
	return Lat, nil
}

func FetchString(size int) Decipher {
	return func(b []byte) (interface{}, error) {
		buf := b[:size]
		return byteString(buf), nil
	}
}

// if copying from half full buffer
func byteString(b []byte) string {
	for i := 0; i < len(b); i++ {
		if b[i] == 0 {
			return string(b[:i])
		}
	}
	return string(b)
}

type Column struct {
	Name      string
	Width     int
	Precision int
}

type HeadFlag int

const (
	thisVersion = 1

	HeadsUnknown HeadFlag = iota
	Heads64      HeadFlag = 1 << iota
	HeadsLon              // Data is lon, lat
)

type HolderHeader struct {
	Magic   [4]byte
	Version uint32
	Start   uint32
	Flags   HeadFlag
	// ColumnGlob []byte
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

*/
func (hh HolderHeader) Encode(b []byte) error {
	n := copy(b, hh.Magic[:])
	log.Printf("copy %d", n)
	binary.LittleEndian.PutUint32(b[4:], uint32(hh.Version))
	binary.LittleEndian.PutUint32(b[8:], uint32(hh.Start))
	binary.LittleEndian.PutUint32(b[12:], uint32(hh.Flags))
	return nil
}

func (hh *HolderHeader) Decode(b []byte) error {
	for i := 0; i < 4; i++ {
		hh.Magic[i] = b[i]
	}
	hh.Version = binary.LittleEndian.Uint32(b[4:])
	hh.Start = binary.LittleEndian.Uint32(b[8:])
	hh.Flags = HeadFlag(binary.LittleEndian.Uint32(b[12:]))
	return nil
}

const HolderHeaderSize = int(unsafe.Sizeof(HolderHeader{}))

type LiveHeader struct {
	Head    HolderHeader
	Columns []Column
}

// MagicHead is the 4 byte header prefix to identify the file type
func MagicHead() [4]byte {
	return [4]byte{'p', '@', 'k', 'R'}
}

func ReadHeader(r io.Reader) (*LiveHeader, error) {
	b := make([]byte, HolderHeaderSize)
	n, err := r.Read(b)
	if err != nil {
		return nil, fmt.Errorf("read fail: %w", err)
	}
	if n < len(b) {
		panic("add goto to fix this")
	}
	return &LiveHeader{}, fmt.Errorf("foo not you")
}

func MakeHeader(columns ...Column) LiveHeader {
	b, err := gobEncode(columns)
	if err != nil {
		log.Fatalf("gob encode fail: %v", err)
	}
	size := HolderHeaderSize + len(b)
	size = RoundUp(size)
	Head := HolderHeader{
		Magic:   MagicHead(),
		Version: thisVersion,
		Start:   uint32(size),
	}
	return LiveHeader{
		Head:    Head,
		Columns: columns,
	}
}

func (hh HolderHeader) Equal(x HolderHeader) error {
	if hh.Magic != x.Magic {
		return fmt.Errorf("bad magic: %v vs %v", hh.Magic, x.Magic)
	}
	if hh.Version != x.Version {
		return fmt.Errorf("bad version: %v vs %v", hh.Version, x.Version)
	}
	if hh.Start != x.Start {
		return fmt.Errorf("bad start: %v vs %v", hh.Start, x.Start)
	}
	return nil
}

// NewHolder takes columns with optional formatting/range info
// columnName:maxSize:precision
// if precision == "s" the field is a fixed with string of maxSize width
func NewHolder(columns ...string) (*Holder, error) {
	hold := &Holder{
		columns: make([]string, len(columns)),
		sizes:   make([]int, len(columns)),
		// kinds:   make([]Kind, len(columns)),
		layout: make([]int, len(columns)),
		offset: make([]int, len(columns)),
		// values:  make([]interface{}, len(columns)),
	}
	for i, col := range columns {
		var fn Decipher
		// var fx Encipher
		hold.layout[i] = i
		// dt := DTF64
		size := 8
		fore, aft, found := strings.Cut(col, ":")
		if !found {
			continue
		}
		columns[i] = fore
		fore, aft, found = strings.Cut(aft, ":")
		if !found {
			if val, err := strconv.Atoi(fore); err == nil {
				switch {
				case val < 256:
					size = 1
					fn = FetchByte
					// fx = SerialByte
				case val < 65536:
					size = 2
					fn = FetchI16
					// fx = Serial16
					// dt = DTI16
				case val < math.MaxInt32:
					size = 4
					fn = FetchI32
					// dt = DTI32
				default:
					size = 8
					fn = FetchI32
				}
			}
			hold.decoders[i] = fn
			hold.sizes[i] = size
			// hold.kinds[i] = dt
			// continue
		} else {
			width, err := strconv.Atoi(fore)
			if err != nil {
				return nil, fmt.Errorf("parse fore fail: %w", err)
			}
			if aft == "s" {
				fn = FetchString(width)
				hold.decoders[i] = fn
				hold.sizes[i] = size
				continue
			}
			precision, err := strconv.Atoi(aft)
			if err != nil {
				return nil, fmt.Errorf("parse aft of %q fail: %w", aft, err)
			}
			switch {
			case width == 1 && precision <= 4:
				fn = FetchF1
				size = 2
			case width < 655 && precision <= 2:
				fn = FetchF16
				size = 2
			}

		}
		// aft should have precision
	}
	size := 0
	for i, _ := range hold.sizes {
		off := hold.layout[i]
		size += hold.sizes[off]
		hold.offset[i] += size
	}
	return hold, nil
}

func JDump(obj interface{}) {
	log.Printf("DUMP (%T): %+v", obj, obj)
	byte, err := json.Marshal(obj)
	if err != nil {
		log.Printf("marshal fail: %v", err)
	} else {
		fmt.Println(string(byte))
	}
}

func (h Holder) Show() {
	// b, err := h.MarshalJSON()
	// if err != nil {
	// 	log.Printf("show fail: %v", err)
	// } else {
	// 	fmt.Println(string(b))
	// }
}
func (rh RecordHolder) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	// log.Println("START")
	(&buf).WriteString(`{`)
	for i, col := range rh.hold.columns {
		if false {
			log.Printf("COL: %d", i)
		}
		if i > 0 {
			fmt.Fprintf(&buf, ",")
		}
		fmt.Fprintf(&buf, "%q", col)
		fmt.Fprintf(&buf, ":")
		// kind := h.kinds[i]
		switch value := rh.values[i].(type) {
		// case string:
		// 	fmt.Fprintf(&buf, "%q", value)
		// case int, float64:
		// 	fmt.Fprint(&buf, value)
		// fmt.Println("INT VALUE:", value)
		default:
			// log.Printf("HUH? (%T): %+v\n", value, value)
			fmt.Fprint(&buf, value)
		}
	}
	(&buf).WriteString(`}`)
	return buf.Bytes(), nil
	// func hey(what ...interface{}) {
	// 	for _, hm := range what {
	// 		fmt.Printf("%v", hm)
	// 	}
	// 	fmt.Println()
}

func (h *Holder) UnMarshal(b []byte) error {
	return nil
}

// Importer makes a function that persists a slice of strings to disk
func (h Holder) Importer(w io.Writer) func(ss ...string) error {
	// for _, k := range h.kinds {
	// 	fmt.Println("K:", k)
	// }
	// for i, r := range h.recorders {
	// 	i
	// }
	return func(ss ...string) error {
		for i, s := range ss {
			fmt.Printf("i: %d %s\n", i, s)
		}
		return nil
	}
}
func (h Holder) MakeEncoder() func(b []byte) error {
	// for _, k := range h.kinds {
	// 	fmt.Println("K:", k)
	// }
	return func(b []byte) error {
		return nil
	}
}

// RoundUp expects to be 32bit or less
func RoundUp(n int) int {
	n--
	n |= (n >> 1)
	n |= (n >> 2)
	n |= (n >> 4)
	n |= (n >> 8)
	n |= (n >> 16)
	return n + 1
}
