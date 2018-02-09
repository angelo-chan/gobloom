package gobloom

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBloomFilter_M_K(t *testing.T) {
	f := New(1000, 0.01)
	assert.Equal(t, uint(9600), f.M())
	assert.Equal(t, uint(7), f.K())

	f = New(10000, 0.03)
	assert.Equal(t, uint(73024), f.M())
	assert.Equal(t, uint(5), f.K())

	f = New(1024, 0.01)
	assert.Equal(t, uint(9856), f.M())
	assert.Equal(t, uint(7), f.K())
}

func TestBloomFilter(t *testing.T) {
	s1 := "this_is_a_test_string"
	s2 := "this_is_another_test_string"
	s3 := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ"
	s4 := "dyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ"

	f := New(100, 0.01)
	f.Put([]byte(s1)).Put([]byte(s3))
	assert.True(t, f.MightContain([]byte(s1)))
	assert.False(t, f.MightContain([]byte(s2)))
	assert.True(t, f.MightContain([]byte(s3)))
	assert.False(t, f.MightContain([]byte(s4)))
}

func TestBloomFilter_JSON(t *testing.T) {
	s1 := "this_is_a_test_string"
	s2 := "this_is_another_test_string"
	s3 := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ"
	s4 := "dyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ"

	f := New(100, 0.01)
	f.Put([]byte(s1)).Put([]byte(s3))
	exported, _ := f.MarshalJSON()
	bloomFilterJSON := &BloomFilterJSON{}
	json.Unmarshal(exported, bloomFilterJSON)
	newF := From(bloomFilterJSON.B, bloomFilterJSON.K)
	assert.True(t, newF.MightContain([]byte(s1)))
	assert.False(t, newF.MightContain([]byte(s2)))
	assert.True(t, newF.MightContain([]byte(s3)))
	assert.False(t, newF.MightContain([]byte(s4)))

}
