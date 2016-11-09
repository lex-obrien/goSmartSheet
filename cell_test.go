package goSmartSheet

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshalJSON(t *testing.T) {
	assert := assert.New(t)
	var c1, c2, c3 CellValue
	var err error
	var b []byte

	c1.IntVal = 5
	if b, err = c1.MarshalJSON(); err != nil {
		t.Error(err)
	} else {
		assert.Equal(b, []byte{53}, "Should be equal to 53 or \"5\"")
	}

	c2.StringVal = "HEY"
	if b, err = c2.MarshalJSON(); err != nil {
		t.Error(err)
	} else {
		assert.Equal(b, []byte{34, 72, 69, 89, 34}, "Should be equal to {34 72 69 89 34} or \"HEY\"")
	}

	c3.FloatVal = 1.34
	if b, err = c3.MarshalJSON(); err != nil {
		t.Error(err)
	} else {
		assert.Equal(b, []byte{49, 46, 51, 52}, "Should be equal to {49, 46, 51, 52} or \"1.34\"")
	}
}
