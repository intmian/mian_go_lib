package cipher

import "testing"

func TestSha2562String(t *testing.T) {
	s := Sha2562String("1234567")
	if s != "8bb0cf6eb9b17d0f7d22b456f121257dc1254e1f01665370476383ea776df414" {
		t.Error("Sha2562String error")
	}
}
