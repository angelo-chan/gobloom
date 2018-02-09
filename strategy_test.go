package gobloom

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestMURMUR128MITZ64 the expected data are exported from guava test results.
func TestMURMUR128MITZ64(t *testing.T) {
	expected1 := []uint{3942, 6555, 9168, 2181, 4794, 7407, 420}
	s1 := "this_is_a_test_string"
	strategy := &MURMUR128MITZ64{}
	indexes := strategy.Indexes([]byte(s1), 9600, 7)
	assert.Equal(t, expected1, indexes)
}
