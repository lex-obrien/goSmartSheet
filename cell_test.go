package goSmartSheet

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCellValue_Settings(t *testing.T) {
	assert := assert.New(t)
	var cv CellValue

	cv.SetInt(53)
	assert.NotNil(cv.IntVal)
	assert.Equal(*(cv.IntVal), 53)
	assert.Nil(cv.FloatVal)
	assert.Nil(cv.StringVal)

	cv.SetString("HEY")
	assert.NotNil(cv.StringVal)
	assert.Equal("HEY", *(cv.StringVal))
	assert.Nil(cv.IntVal)
	assert.Nil(cv.FloatVal)

	cv.SetFloat(1.34)
	assert.NotNil(cv.FloatVal)
	assert.Equal(1.34, *(cv.FloatVal))
	assert.Nil(cv.IntVal)
	assert.Nil(cv.StringVal)

	cv.SetInt(231)
	assert.NotNil(cv.IntVal)
	assert.Equal(231, *(cv.IntVal))
	assert.Nil(cv.FloatVal)
	assert.Nil(cv.StringVal)

	cv.SetString("tqwtf2t")
	assert.NotNil(cv.StringVal)
	assert.Equal("tqwtf2t", *(cv.StringVal))
	assert.Nil(cv.IntVal)
	assert.Nil(cv.FloatVal)

	cv.SetFloat(6.26)
	assert.NotNil(cv.FloatVal)
	assert.Equal(6.26, *(cv.FloatVal))
	assert.Nil(cv.IntVal)
	assert.Nil(cv.StringVal)
}

func TestCellValue_String(t *testing.T) {
	assert := assert.New(t)
	var cv CellValue
	var v string
	var e error

	cv.SetInt(53)
	v, e = cv.String()
	assert.NoError(e)
	assert.Equal("53", v)

	cv.SetInt(141453)
	v, e = cv.String()
	assert.NoError(e)
	assert.Equal("141453", v)

	cv.SetInt(2)
	v, e = cv.String()
	assert.NoError(e)
	assert.Equal("2", v)

	cv.SetString("BOB")
	v, e = cv.String()
	assert.NoError(e)
	assert.Equal("BOB", v)

	cv.SetString("Be3twg")
	v, e = cv.String()
	assert.NoError(e)
	assert.Equal("Be3twg", v)

	cv.SetFloat(1.346)
	v, e = cv.String()
	assert.NoError(e)
	assert.Equal("1.346", v)

	cv.SetFloat(1.34600)
	v, e = cv.String()
	assert.NoError(e)
	assert.Equal("1.346", v)

	cv.SetFloat(1.30000)
	v, e = cv.String()
	assert.NoError(e)
	assert.Equal("1.3", v)

	cv.SetFloat(0001.30000)
	v, e = cv.String()
	assert.NoError(e)
	assert.Equal("1.3", v)

	cv.SetFloat(0001.30003)
	v, e = cv.String()
	assert.NoError(e)
	assert.Equal("1.30003", v)
}

func TestMarshalJSON(t *testing.T) {
	assert := assert.New(t)
	var c1, c2, c3 CellValue
	var err error
	var b []byte

	c1.SetInt(5)
	if b, err = c1.MarshalJSON(); err != nil {
		t.Error(err)
	} else {
		assert.Equal(b, []byte{53}, "Should be equal to 53 or \"3\"")
	}

	c2.SetString("HEY")
	if b, err = c2.MarshalJSON(); err != nil {
		t.Error(err)
	} else {
		assert.Equal(b, []byte{34, 72, 69, 89, 34}, "Should be equal to {34 72 69 89 34} or \"HEY\"")
	}

	c3.SetFloat(1.34)
	if b, err = c3.MarshalJSON(); err != nil {
		t.Error(err)
	} else {
		assert.Equal(b, []byte{49, 46, 51, 52}, "Should be equal to {49, 46, 51, 52} or \"1.34\"")
	}
}

func TestUnMarshalJSON(t *testing.T) {
	assert := assert.New(t)
	var cv CellValue
	var err error
	var b []byte

	b = []byte(`"HEY"`)
	if err = cv.UnmarshalJSON(b); err != nil {
		t.Error(err)
	} else {
		v, _ := cv.String()
		assert.Equal("HEY", v)
		assert.Nil(cv.FloatVal)
		assert.Nil(cv.IntVal)
	}

	b = []byte(`"1.2"`) //as string
	if err = cv.UnmarshalJSON(b); err != nil {
		t.Error(err)
	} else {
		assert.NotNil(cv.StringVal)
		v, _ := cv.String()
		assert.Equal("1.2", v)
		assert.Nil(cv.FloatVal)
		assert.Nil(cv.IntVal)
	}

	b = []byte(`1.2`)
	if err = cv.UnmarshalJSON(b); err != nil {
		t.Error(err)
	} else {
		assert.NotNil(cv.FloatVal)
		v, _ := cv.Float()
		assert.Equal(1.2, v)
		assert.Nil(cv.StringVal)
		assert.Nil(cv.IntVal)
	}

	b = []byte(`25`)
	if err = cv.UnmarshalJSON(b); err != nil {
		t.Error(err)
	} else {
		assert.NotNil(cv.IntVal)
		v, _ := cv.Int()
		assert.Equal(25, v)
		assert.Nil(cv.FloatVal)
		assert.Nil(cv.StringVal)
	}
}

func BenchmarkUnmarshallInt(b *testing.B) {
	var cv CellValue
	d := []byte(`25`)
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		cv.UnmarshalJSON(d)
	}
}

func BenchmarkUnmarshallString(b *testing.B) {
	var cv CellValue
	d := []byte(`"vtGdj"`)
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		cv.UnmarshalJSON(d)
	}
}

func BenchmarkUnmarshallStringLong(b *testing.B) {
	var cv CellValue
	d := []byte(`"vtj1413#SDG2352tw45dj"`)
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		cv.UnmarshalJSON(d)
	}
}

func BenchmarkUnmarshallFloat(b *testing.B) {
	var cv CellValue
	d := []byte(`1.25`)
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		cv.UnmarshalJSON(d)
	}
}
