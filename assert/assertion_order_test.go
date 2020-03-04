package assert

import (
	"reflect"
	"testing"
)

func TestCompare(t *testing.T) {
	for _, currCase := range []struct {
		less    interface{}
		greater interface{}
		cType   string
	}{
		{less: "a", greater: "b", cType: "string"},
		{less: int(1), greater: int(2), cType: "int"},
		{less: int8(1), greater: int8(2), cType: "int8"},
		{less: int16(1), greater: int16(2), cType: "int16"},
		{less: int32(1), greater: int32(2), cType: "int32"},
		{less: int64(1), greater: int64(2), cType: "int64"},
		{less: uint8(1), greater: uint8(2), cType: "uint8"},
		{less: uint16(1), greater: uint16(2), cType: "uint16"},
		{less: uint32(1), greater: uint32(2), cType: "uint32"},
		{less: uint64(1), greater: uint64(2), cType: "uint64"},
		{less: float32(1), greater: float32(2), cType: "float32"},
		{less: float64(1), greater: float64(2), cType: "float64"},
	} {
		resLess, isComparable := compare(currCase.less, currCase.greater, reflect.ValueOf(currCase.less).Kind())
		if !isComparable {
			t.Error("object should be comparable for type " + currCase.cType)
		}

		if resLess != compareLess {
			t.Errorf("object less should be less than greater for type " + currCase.cType)
		}

		resGreater, isComparable := compare(currCase.greater, currCase.less, reflect.ValueOf(currCase.less).Kind())
		if !isComparable {
			t.Error("object are comparable for type " + currCase.cType)
		}

		if resGreater != compareGreater {
			t.Errorf("object greater should be greater than less for type " + currCase.cType)
		}

		resEqual, isComparable := compare(currCase.less, currCase.less, reflect.ValueOf(currCase.less).Kind())
		if !isComparable {
			t.Error("object are comparable for type " + currCase.cType)
		}

		if resEqual != 0 {
			t.Errorf("objects should be equal for type " + currCase.cType)
		}
	}
}

func TestGreater(t *testing.T) {
	mockT := new(testing.T)

	if !Greater(mockT, 2, 1) {
		t.Error("Greater should return true")
	}

	if Greater(mockT, 1, 1) {
		t.Error("Greater should return false")
	}

	if Greater(mockT, 1, 2) {
		t.Error("Greater should return false")
	}
}

func TestGreaterOrEqual(t *testing.T) {
	mockT := new(testing.T)

	if !GreaterOrEqual(mockT, 2, 1) {
		t.Error("Greater should return true")
	}

	if !GreaterOrEqual(mockT, 1, 1) {
		t.Error("Greater should return true")
	}

	if GreaterOrEqual(mockT, 1, 2) {
		t.Error("Greater should return false")
	}
}

func TestLess(t *testing.T) {
	mockT := new(testing.T)

	if !Less(mockT, 1, 2) {
		t.Error("Less should return true")
	}

	if Less(mockT, 1, 1) {
		t.Error("Less should return false")
	}

	if Less(mockT, 2, 1) {
		t.Error("Less should return false")
	}
}

func TestLessOrEqual(t *testing.T) {
	mockT := new(testing.T)

	if !LessOrEqual(mockT, 1, 2) {
		t.Error("Greater should return true")
	}

	if !LessOrEqual(mockT, 1, 1) {
		t.Error("Greater should return true")
	}

	if LessOrEqual(mockT, 2, 1) {
		t.Error("Greater should return false")
	}
}

func Test_compareTwoValuesDifferentValuesTypes(t *testing.T) {
	mockT := new(testing.T)

	for _, currCase := range []struct {
		v1            interface{}
		v2            interface{}
		compareResult bool
	}{
		{v1: 123, v2: "abc"},
		{v1: "abc", v2: 123456},
		{v1: float64(12), v2: "123"},
		{v1: "float(12)", v2: float64(1)},
	} {
		compareResult := compareTwoValues(mockT, currCase.v1, currCase.v2, []CompareType{compareLess, compareEqual, compareGreater}, "testFailMessage")
		if compareResult {
			t.Errorf("Values %s and %s should be different kinds", currCase.v1, currCase.v2)
		}
	}
}

func Test_compareTwoValuesNotComparableValues(t *testing.T) {
	mockT := new(testing.T)

	type CompareStruct struct {
	}

	for _, currCase := range []struct {
		v1 interface{}
		v2 interface{}
	}{
		{v1: CompareStruct{}, v2: CompareStruct{}},
		{v1: map[string]int{}, v2: map[string]int{}},
		{v1: make([]int, 5, 5), v2: make([]int, 5, 5)},
	} {
		compareResult := compareTwoValues(mockT, currCase.v1, currCase.v2, []CompareType{compareLess, compareEqual, compareGreater}, "testFailMessage")
		if compareResult {
			t.Errorf("Values %s and %s should be not comparable", currCase.v1, currCase.v2)
		}
	}
}

func Test_compareTwoValuesCrrectCompareResult(t *testing.T) {
	mockT := new(testing.T)

	for _, currCase := range []struct {
		v1           interface{}
		v2           interface{}
		compareTypes []CompareType
	}{
		{v1: 1, v2: 2, compareTypes: []CompareType{compareLess}},
		{v1: 1, v2: 2, compareTypes: []CompareType{compareLess, compareEqual}},
		{v1: 2, v2: 2, compareTypes: []CompareType{compareGreater, compareEqual}},
		{v1: 2, v2: 2, compareTypes: []CompareType{compareEqual}},
		{v1: 2, v2: 1, compareTypes: []CompareType{compareEqual, compareGreater}},
		{v1: 2, v2: 1, compareTypes: []CompareType{compareGreater}},
	} {
		compareResult := compareTwoValues(mockT, currCase.v1, currCase.v2, currCase.compareTypes, "testFailMessage")
		if !compareResult {
			t.Errorf("Values %d and %d is not compared correctly", currCase.v1, currCase.v2)
		}
	}
}

func Test_containsValue(t *testing.T) {
	for _, currCase := range []struct {
		values []CompareType
		value  CompareType
		result bool
	}{
		{values: []CompareType{compareGreater}, value: compareGreater, result: true},
		{values: []CompareType{compareGreater, compareLess}, value: compareGreater, result: true},
		{values: []CompareType{compareGreater, compareLess}, value: compareLess, result: true},
		{values: []CompareType{compareGreater, compareLess}, value: compareEqual, result: false},
	} {
		compareResult := containsValue(currCase.values, currCase.value)
		if compareResult != currCase.result {
			t.Error("Value not in list")
		}
	}
}
