package gen

import "testing"

func TestCleanString(t *testing.T) {
	if s := cleanString("100|  97"); s != "100 97" {
		t.Errorf("Failed to convert '100|  99', got '%s'\n", s)
	}
	if s := cleanString("100|  97 123"); s != "100 97 123" {
		t.Errorf("Failed to convert '100|  99', got '%s'\n", s)
	}
}
