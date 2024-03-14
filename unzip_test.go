package artifact_diff_test

import (
	diff "github.com/gesellix/artifact-diff"
	"os"
	"testing"
)

func _TestWriteJarFile(t *testing.T) {
	buff := []byte{80, 75, 3, 4, 10, 0, 0, 0, 0, 0, 91, 109, 103, 78, 132, 225, 60, 127, 13, 0, 0, 0, 13, 0, 0, 0, 5, 0, 28, 0, 97, 46, 116, 120, 116, 85, 84, 9, 0, 3, 206, 17, 129, 92, 219, 17, 129, 92, 117, 120, 11, 0, 1, 4, 232, 3, 0, 0, 4, 232, 3, 0, 0, 72, 101, 108, 108, 111, 32, 71, 111, 112, 104, 101, 114, 10, 80, 75, 1, 2, 30, 3, 10, 0, 0, 0, 0, 0, 91, 109, 103, 78, 132, 225, 60, 127, 13, 0, 0, 0, 13, 0, 0, 0, 5, 0, 24, 0, 0, 0, 0, 0, 1, 0, 0, 0, 164, 129, 0, 0, 0, 0, 97, 46, 116, 120, 116, 85, 84, 5, 0, 3, 206, 17, 129, 92, 117, 120, 11, 0, 1, 4, 232, 3, 0, 0, 4, 232, 3, 0, 0, 80, 75, 5, 6, 0, 0, 0, 0, 1, 0, 1, 0, 75, 0, 0, 0, 76, 0, 0, 0, 0, 0}
	os.WriteFile("testdata/jarfiles/minimal.jar", buff, 0x666)
}

func TestUnzipJarFile(t *testing.T) {
	t.Parallel()
	err := diff.Unzip("testdata/jarfiles/minimal.jar", "testdata/_tmp")
	defer os.RemoveAll("testdata/_tmp")
	if err != nil {
		t.Errorf("unzip failed: %v", err)
	}
	stat, err := os.Stat("testdata/_tmp/a.txt")
	if err != nil {
		t.Errorf("stat on %s failed: %v", "testdata/_tmp/a.txt", err)
	}

	if stat.Size() != 13 {
		t.Errorf("expected %s to have a size of %v, but got %v", "testdata/_tmp/a.txt", 1, stat.Size())
	}
}

func TestUnzipTarFile(t *testing.T) {
	t.Parallel()
	err := diff.Untar("testdata/tarfiles/minimal.tar.gz", "testdata/_tmp")
	defer os.RemoveAll("testdata/_tmp")
	if err != nil {
		t.Errorf("untar failed: %v", err)
	}
	stat, err := os.Stat("testdata/_tmp/a.txt")
	if err != nil {
		t.Errorf("stat on %s failed: %v", "testdata/_tmp/a.txt", err)
	}

	if stat.Size() != 13 {
		t.Errorf("expected %s to have a size of %v, but got %v", "testdata/_tmp/a.txt", 1, stat.Size())
	}
}
