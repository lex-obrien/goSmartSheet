package goSmartSheet

import (
	"reflect"
	"testing"
)

func TestMarshalJSON(t *testing.T) {
	var c1, c2, c3 CellValue
	c1.IntVal = 5
	c2.StringVal = "HEY"
	c3.FloatVal = 1.34
	var err error
	var b []byte
	if b, err = c1.MarshalJSON(); err != nil {
		t.Error(err)
	} else {
		expected := []byte{53} // [53] = "5"
		if !reflect.DeepEqual(b, expected) {
			t.Errorf("Expected %v and got %v", expected, b)
		}
	}
	if b, err = c2.MarshalJSON(); err != nil {
		t.Error(err)
	} else {
		expected := []byte{34, 72, 69, 89, 34} // [34 72 69 89 34] = \"HEY\"
		if !reflect.DeepEqual(b, expected) {
			t.Errorf("Expected %v and got %v", expected, b)
		}
	}
	if b, err = c3.MarshalJSON(); err != nil {
		t.Error(err)
	} else {
		expected := []byte{49, 46, 51, 52} // [49 46 51 52] = "1.34"
		if !reflect.DeepEqual(b, expected) {
			t.Errorf("Expected %v and got %v", expected, b)
		}
	}
}
