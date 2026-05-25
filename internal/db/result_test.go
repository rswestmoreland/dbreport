package db

import (
	"testing"
	"time"
)

func TestConvertCell(t *testing.T) {
	if got := ConvertCell([]byte("hello")); got.Text != "hello" || got.IsNull {
		t.Fatalf("unexpected byte cell: %#v", got)
	}
	if got := ConvertCell(nil); got.Text != "" || !got.IsNull {
		t.Fatalf("unexpected nil cell: %#v", got)
	}
	when := time.Date(2026, 5, 25, 14, 30, 0, 0, time.UTC)
	if got := ConvertCell(when); got.Text != "2026-05-25 14:30:00" {
		t.Fatalf("unexpected time cell: %#v", got)
	}
}
