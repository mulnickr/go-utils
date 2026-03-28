package stream

import (
	"testing"
)

type Object struct {
	ID       int
	Value    string
	Factor   float32
	IsActive bool
}

var values = []*Object{
	{
		ID:       1,
		Value:    "Test",
		Factor:   1.25,
		IsActive: true,
	},
	{
		ID:       2,
		Value:    "Another",
		Factor:   1.75,
		IsActive: true,
	},
	{
		ID:       3,
		Value:    "More",
		Factor:   1.75,
		IsActive: true,
	},
	{
		ID:       4,
		Value:    "Last",
		Factor:   1.00,
		IsActive: false,
	},
	{
		ID:       5,
		Value:    "Last",
		Factor:   1.25,
		IsActive: true,
	},
}

type TestObject[T any] struct {
	want    T
	actual  T
	Compare func(want, actual T) bool
}

func TestForEach(t *testing.T) {
	sum := 0
	fn := func(i int, s *Object) {
		sum += s.ID
	}
	ForEach(values, fn)
	if sum != 15 {
		t.Errorf("Want: 15, Actual: %v", sum)
	}
}

func TestMap(t *testing.T) {
	type First struct {
		Value string
	}

	type Second struct {
		Other string
	}

	firsts := []First{
		{Value: "test"},
	}

	results := Map(firsts, func(f First) Second {
		return Second{
			Other: f.Value,
		}
	})

	if len(results) < 1 {
		t.Errorf("Want: 1, Actual: %v", len(results))
	}

	if results[0].Other != firsts[0].Value {
		t.Errorf("Want: 'test', Actual: '%v'", results[0].Other)
	}
}

func TestBooleans(t *testing.T) {
	compare := func(s *Object) bool {
		return s.Value == "Test"
	}

	mapper := func(s *Object) string {
		return s.Value
	}

	s := Some(values, compare)
	if s != true {
		t.Errorf("Want: Test, Actual: %v", values)
	}

	e := Every(values, compare)
	if e != false {
		t.Errorf("Want: Test, Actual: %v", values)
	}

	n := Never(values, compare)
	if n != false {
		t.Errorf("Want: !Test, Actual: %v", Map(values, mapper))
	}
}

func TestFind(t *testing.T) {
	fn := func(s *Object) bool {
		return s.Value == "Last"
	}

	only := func(s *Object) bool {
		return s.Value == "Test"
	}

	r, err := Find(values, fn)
	if err != nil {
		t.Errorf("Unexpected error: %v", err.Error())
	}
	if len(r) != 2 {
		t.Errorf("Want: 2, Actual: %v", len(r))
	}

	o, err := FindOnly(values, only)
	if err != nil {
		t.Errorf("Unexpected error: %v", err.Error())
	}
	if o.Value != "Test" {
		t.Errorf("Want: Test, Actual: %v", o.Value)
	}

	f, err := FindFirst(values, fn)
	if err != nil {
		t.Errorf("Unexpected error: %v", err.Error())
	}
	if f.Value != "Last" || f.ID != 4 {
		t.Errorf("Want: Last|ID4, Actual: %vID%v", f.Value, f.ID)
	}

	l, err := FindLast(values, fn)
	if err != nil {
		t.Errorf("Unexpected error: %v", err.Error())
	}
	if l.Value != "Last" || l.ID != 5 {
		t.Errorf("Want: Last|ID5, Actual: %vID%v", f.Value, f.ID)
	}
}
