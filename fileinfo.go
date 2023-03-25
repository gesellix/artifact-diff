package artifact_diff

import (
	"fmt"
	"log"
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

type ArtifactInfo struct {
	Path      string    `json:"path" yaml:"path"`
	Count     int       `json:"count" yaml:"count"`
	FileInfos FileInfos `json:"fileInfos" yaml:"fileInfos"`
}

type FlatArtifactInfo struct {
	Path      string      `json:"path" yaml:"path"`
	Count     int         `json:"count" yaml:"count"`
	FileInfos []*FileInfo `json:"fileInfos" yaml:"fileInfos"`
}

func (d *ArtifactInfo) WithFlattenedAndSortedFileInfos() *FlatArtifactInfo {
	flat := &FlatArtifactInfo{
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

func (d *ArtifactInfo) AddFileInfo(path string, info *FileInfo) {
	pathMd5 := Md5hash(path)
	if _, ok := d.FileInfos[pathMd5]; ok {
		log.Println(fmt.Sprintf("duplicate path for %s?", path))
	}
	d.FileInfos[pathMd5] = info
}

func (d *ArtifactInfo) String() string {
	return fmt.Sprintf("Count=%v, FileInfos=%v", d.Count, d.FileInfos)
}
