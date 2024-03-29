package compact

import (
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"fmt"
	"os"
)

// GobDump saves the object in a gzipped GOB encoded file
func GobDump(filename string, obj interface{}) error {
	if obj == nil {
		return fmt.Errorf("nil object of type %T", obj)
	}
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	w := gzip.NewWriter(f)
	enc := gob.NewEncoder(w)
	enc.Encode(obj)
	if err := w.Close(); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}

// GobLoad populates the object from a gzipped GOB encoded file
func GobLoad(filename string, obj interface{}) error {
	if obj == nil {
		return fmt.Errorf("nil object of type %T", obj)
	}
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	r, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer r.Close()
	dec := gob.NewDecoder(r)
	return dec.Decode(obj)
}

func gobEncode(obj interface{}) ([]byte, error) {
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	enc := gob.NewEncoder(w)
	if err := enc.Encode(obj); err != nil {
		return nil, fmt.Errorf("encoder fail: %w", err)
	}
	return buf.Bytes(), nil
}

func gobDecode(obj interface{}, b []byte) error {
	src := bytes.NewReader(b)
	r, err := gzip.NewReader(src)
	if err != nil {
		return err
	}
	dec := gob.NewDecoder(r)
	return dec.Decode(obj)
}
