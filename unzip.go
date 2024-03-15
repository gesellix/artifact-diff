package artifact_diff

import (
	"golift.io/xtractr"
)

func Unzip(src, dest string) error {
	x := &xtractr.XFile{
		FilePath:  src,
		OutputDir: dest,
		FileMode:  0777,
		DirMode:   0755,
	}
	_, files, err := xtractr.ExtractZIP(x)
	if err != nil || files == nil {
		return err
	}
	//log.Println("Bytes written:", size, "Files Extracted:\n -", strings.Join(files, "\n -"))
	return nil
}

func Gunzip(src, dest string) error {
	x := &xtractr.XFile{
		FilePath:  src,
		OutputDir: dest,
		FileMode:  0777,
		DirMode:   0755,
	}
	_, files, err := xtractr.ExtractGzip(x)
	if err != nil || files == nil {
		return err
	}
	//log.Println("Bytes written:", size, "Files Extracted:\n -", strings.Join(files, "\n -"))
	return nil
}

func Untar(src, dest string) error {
	x := &xtractr.XFile{
		FilePath:  src,
		OutputDir: dest,
		FileMode:  0777,
		DirMode:   0755,
	}
	_, files, err := xtractr.ExtractTar(x)
	if err != nil || files == nil {
		return err
	}
	//log.Println("Bytes written:", size, "Files Extracted:\n -", strings.Join(files, "\n -"))
	return nil
}
