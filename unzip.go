package artifact_diff

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// From https://stackoverflow.com/a/24792688/372019
func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	os.MkdirAll(dest, 0755)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		// Check for ZipSlip (Directory traversal)
		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", path)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode()|0777)
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode()|0777)
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			//_, err = io.CopyN(os.Stdout, rc, int64(f.UncompressedSize64))
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		//fmt.Println(fmt.Sprintf("Found in jar: %s, mode=%v, size=%v", f.Name, f.Mode(), f.UncompressedSize64))
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}
