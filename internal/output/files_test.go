package output

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteFileCreatesParentDirectory(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "report.html")
	if err := WriteFile(path, []byte("hello")); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}
	if string(data) != "hello" {
		t.Fatalf("unexpected content: %q", string(data))
	}
}
