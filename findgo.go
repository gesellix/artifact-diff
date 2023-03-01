package artifact_diff

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type FileInfo struct {
	Path     string `json:"path" yaml:"path"`
	Filesize int64  `json:"filesize" yaml:"filesize"`
	Checksum string `json:"checksum" yaml:"checksum"`
}

func (fi *FileInfo) String() string {
	return fmt.Sprintf("path=%s, filesize=%v, checksum=%s", fi.Path, fi.Filesize, fi.Checksum)
}

type FileInfos map[string]*FileInfo

func (i FileInfos) String() string {
	values := make([]string, 0, len(i))
	for _, info := range i {
		values = append(values, info.String())
	}
	return fmt.Sprintf(strings.Join(values, ""))
}

type Diff struct {
	Path      string    `json:"path" yaml:"path"`
	Count     int       `json:"count" yaml:"count"`
	FileInfos FileInfos `json:"fileInfos" yaml:"fileInfos"`
}

type FlatDiff struct {
	Path      string      `json:"path" yaml:"path"`
	Count     int         `json:"count" yaml:"count"`
	FileInfos []*FileInfo `json:"fileInfos" yaml:"fileInfos"`
}

func (d *Diff) WithFlattenedAndSortedFileInfos() *FlatDiff {
	flat := &FlatDiff{
		Path:  d.Path,
		Count: d.Count,
	}

	infos := make([]*FileInfo, 0, len(d.FileInfos))
	for _, f := range d.FileInfos {
		infos = append(infos, f)
	}
	sort.Slice(infos, func(i, j int) bool {
		return infos[i].Path < infos[j].Path
	})
	flat.FileInfos = infos
	return flat
}

func (d *Diff) AddFileInfo(path string, info *FileInfo) {
	pathMd5 := Md5hash(path)
	if _, ok := d.FileInfos[pathMd5]; ok {
		log.Println(fmt.Sprintf("duplicate path for %s?", path))
	}
	d.FileInfos[pathMd5] = info
}

func (d *Diff) String() string {
	return fmt.Sprintf("Count=%v, FileInfos=%v", d.Count, d.FileInfos)
}

func WalkTree(root string, rootFS fs.FS) (*Diff, error) {
	diff := Diff{
		FileInfos: FileInfos{},
	}
	err := fs.WalkDir(rootFS, ".", func(p string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		diff.Count++
		//log.Println(fmt.Sprintf("file[%v]: %s", diff.Count, p))

		path := filepath.Clean(fmt.Sprintf("%s/%s", root, p))

		info := getFileInfo(path)
		diff.AddFileInfo(path, info)

		if filepath.Ext(p) == ".zip" || filepath.Ext(p) == ".jar" {
			// TODO zip in zip won't work?
			//archiveFs, err := zip.OpenReader(p)
			//if err != nil {
			//	log.Println(fmt.Sprintf("error reading archive %s at tree %s: %v", p, root, err))
			//	return err
			//}
			// //var AppFs = afero.NewMemMapFs()
			// //var AppFs = afero.NewOsFs()

			archiveDiff, err := WalkArchive(root, p)
			if err != nil {
				log.Println(fmt.Sprintf("error walking archive: %v", err))
				return err
			}
			diff.Count = diff.Count + archiveDiff.Count

			for _, f := range archiveDiff.FileInfos {
				diff.AddFileInfo(f.Path, f)
			}
		}
		return nil
	})
	return &diff, err
}

func WalkArchive(root string, zip string) (*Diff, error) {
	temp, err := os.MkdirTemp(
		"",
		fmt.Sprintf("artifact-diff_unzipped_%s",
			strings.ReplaceAll(
				strings.ReplaceAll(filepath.ToSlash(zip), "/", "_"),
				":", "_")))
	if err != nil {
		log.Println(fmt.Sprintf("error creating temporary directory for unzipping %s: %v", zip, err))
		return nil, nil
	}
	log.Println(fmt.Sprintf("tmp dir: %s", temp))

	absSrc := zip
	if !filepath.IsAbs(zip) {
		absSrc, _ = filepath.Abs(fmt.Sprintf("%s/%s", root, zip))
	}
	err = Unzip(absSrc, temp)
	if err != nil {
		log.Println(fmt.Sprintf("unzip err %v", err))
		return nil, nil
	}
	defer func() {
		err = os.RemoveAll(temp)
		if err != nil {
			log.Println(fmt.Sprintf("error removing temporary files at %s: %v", temp, err))
		}
	}()

	tmpDirFS := os.DirFS(temp)
	files, err := WalkTree(temp, tmpDirFS)
	if err != nil {
		return nil, err
	}
	normalizedFiles := &Diff{
		Count:     files.Count,
		FileInfos: FileInfos{},
	}
	for _, v := range files.FileInfos {
		// make the f.Path look like 'file.zip!/path/to/content.txt
		// so that we remove the system specific temp directory prefix
		v.Path = fmt.Sprintf("%s!%s", filepath.ToSlash(zip), filepath.ToSlash(strings.TrimPrefix(v.Path, temp)))
		normalizedFiles.AddFileInfo(v.Path, v)
	}
	files = nil
	return normalizedFiles, nil
}

func getFileInfo(path string) *FileInfo {
	info := &FileInfo{
		Path: path,
	}

	stat, err := os.Stat(path)
	if err != nil {
		//return err
	} else {
		info.Filesize = stat.Size()
	}

	hash, err := checksum(path)
	if err != nil {
		//return err
	} else {
		info.Checksum = hash
	}
	//log.Println(fmt.Sprintf("info: %s", info))
	return info
}
