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
		dec := GetOdds(odd)
		diff := sample - dec
		assert.Less(t, diff, threshold)
		t.Logf("[%d] %.4f -> (%5d) -> %.4f", i, sample, odd, dec)
	}
}
