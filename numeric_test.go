package compact

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOdds(t *testing.T) {

}

func TestOddsRange(t *testing.T) {
	samples := []float32{
		0,
		0.0001,
		0.3333,
		0.5000,
		0.7777,
		0.9999,
		1.0000,
	}
	var threshold float32 = 0.0001
	for i, sample := range samples {
		odd := SetOdds(sample)
		dec := odd.Float32()
		diff := sample - dec
		assert.Less(t, diff, threshold)
		t.Logf("[%d] %.4f -> (%5d) -> %.4f", i, sample, odd, dec)
	}
}

func TestOddsCodex(t *testing.T) {
	samples := []float32{
		0,
		0.0001,
		0.1235,
		0.3333,
		0.5000,
		0.7777,
		0.9999,
		1.0000,
	}
	// var threshold float32 = 0.0001
	b := make([]byte, 2)
	for i, sample := range samples {
		odd := SetOdds(sample)
		err := odd.Encode(b)
		assert.NoError(t, err)
		var restored Odds
		err = (&restored).Decode(b)
		assert.NoError(t, err)
		assert.Equal(t, odd, restored)
		t.Logf("[%d] %.4f -> (%5d) -> %.4f", i, sample, odd, restored.Float32())
	}
}
