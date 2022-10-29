package compact

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMeta(t *testing.T) {
	hand := Holder{
		columns: []string{"first", "second", "3rd", "force"},
		// values:  []interface{}{1, 23.45, 3},
		// values: []interface{}{1, 23.45, true, "hey"},
	}
	JDump(hand)
	// hand.Show()
	// dump(hand)
}

func TestHeader(t *testing.T) {
	columns := []Column{
		{
			Name:  "joebob",
			Width: 2,
		},
	}
	head := MakeHeader(columns...)
	// var buf bytes.Buffer
	b := make([]byte, HolderHeaderSize)
	err := head.Head.Encode(b)
	require.NoError(t, err)
	t.Logf("HEAD: %+v", head)
	var hind HolderHeader
	err = (&hind).Decode(b)
	require.NoError(t, err)
	t.Logf("HIND: %+v", hind)
}
