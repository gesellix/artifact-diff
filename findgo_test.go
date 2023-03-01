package artifact_diff_test

import (
	"archive/zip"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	diff "github.com/gesellix/artifact-diff"
)

func TestFilesOnDisk(t *testing.T) {
	t.Parallel()
	leftRoot := "testdata/findgo"
	leftFS := os.DirFS(leftRoot)
	want := 4
	got, _ := diff.WalkTree(leftRoot, leftFS)
	if want != got.Count {
		t.Errorf("want %d, got %d", want, got.Count)
	}
}

func TestFilesInMemory(t *testing.T) {
	t.Parallel()
	leftRoot := "."
	leftFS := fstest.MapFS{
		"file.go":                {},
		"subfolder/subfolder.go": {},
		"subfolder2/another.go":  {},
		"subfolder2/file.go":     {},
	}
	want := 4
	got, _ := diff.WalkTree(leftRoot, leftFS)
	if want != got.Count {
		t.Errorf("want %d, got %d", want, got.Count)
	}
}

func TestFilesInZIP(t *testing.T) {
	t.Parallel()
	leftRoot := "testdata/findgo.zip"
	leftFS, err := zip.OpenReader(leftRoot)
	if err != nil {
		t.Fatal(err)
	}
	want := 4
	got, _ := diff.WalkTree(leftRoot, leftFS)
	if want != got.Count {
		t.Errorf("want %d, got %d", want, got.Count)
	}
}

func TestFilesOnDiskAndInZIP(t *testing.T) {
	t.Parallel()
	leftRoot := "testdata/findgo-with-zip"
	leftFS := os.DirFS(leftRoot)
	want := 16
	got, _ := diff.WalkTree(leftRoot, leftFS)
	if want != got.Count {
		t.Errorf("want %d, got %d", want, got.Count)
	}
}

func TestFileInfoOnDiskAndInZIP(t *testing.T) {
	t.Parallel()
	leftRoot := "testdata/findgo-with-zip"
	leftFS := os.DirFS(leftRoot)
	want := int64(13)
	got, _ := diff.WalkTree(leftRoot, leftFS)
	//fmt.Println(fmt.Sprintf("infos: %s", got.FileInfos))
	info := got.FileInfos[diff.Md5hash(filepath.Clean("testdata/findgo-with-zip/subfolder2/file.go"))]
	if want != info.Filesize {
		t.Errorf("want %d, got %d", want, got.Count)
	}
	found := false
	for _, v := range got.FileInfos {
		normalized := filepath.ToSlash(v.Path)
		if strings.HasPrefix(normalized, "findgo.zip") &&
			strings.Contains(normalized, "/findgo/subfolder2/another.go") {
			found = true
		}
	}
	if !found {
		t.Errorf("didn't find entry matching %s", "findgo.zip!/findgo/subfolder2/another.go")
	}
}

func TestWalkArchive(t *testing.T) {
	t.Parallel()
	leftRoot := filepath.Clean("testdata/findgo-with-zip/findgo.zip")
	got, _ := diff.WalkArchive(".", leftRoot)
	found := false
	for _, v := range got.FileInfos {
		normalized := filepath.ToSlash(v.Path)
		if strings.HasPrefix(normalized, "testdata/findgo-with-zip/findgo.zip") &&
			strings.Contains(normalized, "/findgo/subfolder2/another.go") {
			found = true
		}
	}
	if !found {
		t.Errorf("didn't find entry matching %s", "findgo.zip!/findgo/subfolder2/another.go")
	}
}

func BenchmarkFilesOnDisk(b *testing.B) {
	leftRoot := "testdata/findgo"
	leftFS := os.DirFS(leftRoot)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		diff.WalkTree(leftRoot, leftFS)
	}
}

func BenchmarkFilesInMemory(b *testing.B) {
	leftRoot := "."
	leftFS := fstest.MapFS{
		"file.go":                {},
		"subfolder/subfolder.go": {},
		"subfolder2/another.go":  {},
		"subfolder2/file.go":     {},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		diff.WalkTree(leftRoot, leftFS)
	}
}
