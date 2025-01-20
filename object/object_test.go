package object

import "testing"

func TestStringHashKey(t *testing.T) {
	u1 := &String{Value: "whatup"}
	u2 := &String{Value: "whatup"}
	d1 := &String{Value: "whatdown"}
	d2 := &String{Value: "whatdown"}

	if u1.Hash() != u2.Hash() {
		t.Errorf("Identical strings do not deterministically hash")
	}
	if d1.Hash() != d2.Hash() {
		t.Errorf("Identical strings do not deterministically hash")
	}
	if u1.Hash() == d1.Hash() {
		t.Errorf("Hash collision on unique inputs")
	}
}
