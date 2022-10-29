package compact

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestScan(t *testing.T) {
	var a, b, c, d float64
	const format = "%f,%f,%f,%f"
	const text = "12.3,45,99,.1"
	_, err := fmt.Sscanf(text, format, &a, &b, &c, &d)
	require.NoError(t, err)
	t.Logf("%f,%f,%f,%f", a, b, c, d)
}

// func TestLineReader(t *testing.T) {
// 	const testcsv = "testdata/test1.csv"
// 	f, err := os.Open(testcsv)
// 	assert.NoError(t, err)
// 	require.NoError(t, err)
// 	var pers FirePersisted
// 	fn := func(s string) error {
// 		var fp FirePersist
// 		err := (&fp).Import(s)
// 		require.NoError(t, err, s)
// 		pers = append(pers, fp)
// 		return err
// 	}
// 	err = LineReader(f, fn)
// 	require.NoError(t, err)
// }

/*
func TestTranspose(t *testing.T) {
	const testcsv = "testdata/test1.csv"
	const testout = "testdata/test1.dat.gz"
	err := TransposeFirePersist(testcsv, testout)
	require.NoError(t, err)
}
*/
