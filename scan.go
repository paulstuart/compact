package compact

import "fmt"

type AColumn struct {
	Name       string
	Keys       map[string]int
	FMax, FMin float64
	IMax, IMin int64
}

type Columnar interface {
	Evaluate(string) error
	Column() string
}

type ColumnKey struct {
	name string
	keys map[string]int
}

const MaxKeyCount = 65535

var ErrKeyOverload = fmt.Errorf("key cardinality exceeeds max size")

func (ck *ColumnKey) Evaluate(s string) error {
	ck.keys[s]++
	if len(ck.keys) > MaxKeyCount {
		return ErrKeyOverload
	}
	return nil
}

func (ck *ColumnKey) Column() string {
	return ck.name
}

func AnalyzeFile(filename string) {
	// return ""
	// get scanner

	// info, err := FileInfo()
	// for  _, col := range columns {

	// }

}
