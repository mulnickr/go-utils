package optional

import "testing"

type Reference struct {
	ID      string
	Data    map[string]any
	Version int
}

func TestOptional(t *testing.T) {
	ref := &Reference{
		ID: "test",
		Data: map[string]any{
			"key": "value",
		},
		Version: 0,
	}

	o := Ok(ref)
	if o.IsErr() || !o.IsPresent() {
		t.Errorf("Unexpected error")
	}

	if o.Err() != nil {
		t.Errorf("Want: nil, Actual: %v", o.Err().Error())
	}

	if o.Value() != ref {
		t.Errorf("Want: %v, Actual: %v", ref, o.Value())
	}
}
