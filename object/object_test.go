package object

import "testing"

func TestStringHashKey(t *testing.T) {
	hello1 := &String{Value: "Hello World"}
	hello2 := &String{Value: "Hello World"}
	diff1 := &String{Value: "brown fox"}
	diff2 := &String{Value: "brown fox"}

	if hello1.HashKey() != hello2.HashKey() {
		t.Errorf("Expected %s and %s to have the same hash key. (%+v), (%+v)", hello1.Value, hello2.Value, hello1, hello2)
	}
	if diff1.HashKey() != diff2.HashKey() {
		t.Errorf("Expected %s and %s to have the same hash key. (%+v), (%+v)", diff1.Value, diff2.Value, diff1, diff2)
	}
	if hello1.HashKey() == diff2.HashKey() {
		t.Errorf("Expected %s and %s to have different hash key. (%+v), (%+v)", hello1.Value, diff2.Value, hello1, diff2)
	}
}

func TestIntegerHashKey(t *testing.T) {
	hello1 := &Integer{Value: 1}
	hello2 := &Integer{Value: 1}
	diff1 := &Integer{Value: 11111111111}
	diff2 := &Integer{Value: 11111111111}

	if hello1.HashKey() != hello2.HashKey() {
		t.Errorf("Expected %d and %d to have the same hash key. (%+v), (%+v)", hello1.Value, hello2.Value, hello1, hello2)
	}
	if diff1.HashKey() != diff2.HashKey() {
		t.Errorf("Expected %d and %d to have the same hash key. (%+v), (%+v)", diff1.Value, diff2.Value, diff1, diff2)
	}
	if hello1.HashKey() == diff2.HashKey() {
		t.Errorf("Expected %d and %d to have different hash key. (%+v), (%+v)", hello1.Value, diff2.Value, hello1, diff2)
	}
}

func TestBooleanHashKey(t *testing.T) {
	hello1 := &Boolean{Value: true}
	hello2 := &Boolean{Value: true}
	diff1 := &Boolean{Value: false}
	diff2 := &Boolean{Value: false}

	if hello1.HashKey() != hello2.HashKey() {
		t.Errorf("Expected %t and %t to have the same hash key. (%+v), (%+v)", hello1.Value, hello2.Value, hello1, hello2)
	}
	if diff1.HashKey() != diff2.HashKey() {
		t.Errorf("Expected %t and %t to have the same hash key. (%+v), (%+v)", diff1.Value, diff2.Value, diff1, diff2)
	}
	if hello1.HashKey() == diff2.HashKey() {
		t.Errorf("Expected %t and %t to have different hash key. (%+v), (%+v)", hello1.Value, diff2.Value, hello1, diff2)
	}
}
