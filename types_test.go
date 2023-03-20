package compact

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestF16(t *testing.T) {
	const expects = 123.34
	s := fmt.Sprint(expects)
	var h, k F16
	t.Logf("expects %.2f with string %q", expects, s)
	err := (&h).Input(s)
	if err != nil {
		t.Fatal(err)
	}
	if F16(h).Float32() != expects {
		t.Fatalf("want: %d -- have %d", h, k)
	}
	t.Logf("F16 val: %v", h)
	buf := make([]byte, 4)
	_ = h.Encode(buf)
	_ = (&k).Decode(buf)
	if h != k {
		t.Fatalf("want: %d -- have %d", h, k)
	}
	t.Logf("want: %d -- have %d", h, k)
}

func TestI8(t *testing.T) {
	const expects = 123
	s := fmt.Sprint(expects)
	var h, k HoldsByte
	// t.Logf("expects %d with string %q", expects, s)
	err := h.Input(s)
	if err != nil {
		t.Fatal(err)
	}
	buf := make([]byte, 4)
	_ = h.Encode(buf)
	_ = (&k).Decode(buf)
	if h != k {
		t.Fatalf("want: %d -- have %d", h, k)
	}
}

func TestI32(t *testing.T) {
	const expects = 1234
	s := fmt.Sprint(expects)
	var h, k HoldsI32
	t.Logf("expects %d with string %q", expects, s)
	err := (&h).Input(s)
	if err != nil {
		t.Fatal(err)
	}
	if h != expects {
		t.Fatalf("we want: %d -- have %d", expects, k)
	}
	buf := make([]byte, 4)
	h.Encode(buf)
	t.Logf("encoded: %v", buf)
	(&k).Decode(buf)
	if k != expects {
		t.Fatalf("and want: %d -- have %d", expects, k)
	}
	if h != k {
		t.Fatalf("still want: %d -- have %d", h, k)
	}
}

func TestI64(t *testing.T) {
	const expects = 1234
	s := fmt.Sprint(expects)
	var h, k HoldsI64
	t.Logf("expects %d with string %q", expects, s)
	err := (&h).Input(s)
	if err != nil {
		t.Fatal(err)
	}
	if h != expects {
		t.Fatalf("we want: %d -- have %d", expects, k)
	}
	buf := make([]byte, 8)
	h.Encode(buf)
	t.Logf("encoded: %v", buf)
	(&k).Decode(buf)
	if k != expects {
		t.Fatalf("and want: %d -- have %d", expects, k)
	}
	if h != k {
		t.Fatalf("still want: %d -- have %d", h, k)
	}
}

func TestText(t *testing.T) {
	const expects = "hi"
	s := fmt.Sprint(expects)
	const max = 2
	h := NewText(max)
	// var h, k HoldsText
	// h.size = 2
	// t.Logf("expects %d with string %q", expects, s)
	err := (&h).Input(s)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("STR: %s", h)
	t.Logf("TXT: %+v", h)
	// if false {
	// 	t.Logf("K: %+v", k)
	// }
	/*
		buf := make([]byte, 8)
		h.Encode(buf)
		t.Logf("encoded: %v", buf)
		(&k).Decode(buf)
		if k != expects {
			t.Fatalf("and want: %d -- have %d", expects, k)
		}
		if h != k {
			t.Fatalf("still want: %d -- have %d", h, k)
		}
	*/
}

func FuzzF64(f *testing.F) {
	testcases := []float64{1.23, 101.0, 23.42, 1999.0}
	for _, tc := range testcases {
		f.Add(tc) // Use f.Add to provide a seed corpus
	}
	f.Fuzz(func(t *testing.T, in float64) {
		t.Logf("fuzzy wuzzy: %f", in)
		var orig, dupe F64
		orig = F64(in)
		b := make([]byte, orig.Size())
		err := orig.Encode(b)
		assert.NoError(t, err)
		err = (&dupe).Decode(b)
		assert.NoError(t, err)
		assert.Equal(t, orig, dupe)
		t.Logf("Before: %q, after: %q", orig, dupe)
		// if orig != dupe {
		// 	t.Errorf("Before: %q, after: %q", orig, dupe)
		// }
		// if utf8.ValidString(orig) && !utf8.ValidString(rev) {
		// 	t.Errorf("Reverse produced invalid UTF-8 string %q", rev)
		// }
	})
}
