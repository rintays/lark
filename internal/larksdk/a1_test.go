package larksdk

import "testing"

func TestA1RangeShape(t *testing.T) {
	rows, cols := a1RangeShape("Sheet1!A1:B2")
	if rows != 2 || cols != 2 {
		t.Fatalf("expected 2x2, got %dx%d", rows, cols)
	}

	rows, cols = a1RangeShape("A1")
	if rows != 1 || cols != 1 {
		t.Fatalf("expected 1x1, got %dx%d", rows, cols)
	}

	rows, cols = a1RangeShape("B2:B4")
	if rows != 3 || cols != 1 {
		t.Fatalf("expected 3x1, got %dx%d", rows, cols)
	}
}
