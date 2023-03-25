package artifact_diff

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func CollectFileInfos(path string) (*ArtifactInfo, error) {
	log.Println("Scanning", path)

	stat, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	var result *ArtifactInfo
	if stat.IsDir() {
		result, err = WalkTree(
			path,
			os.DirFS(path),
		)
		if err != nil {
			return nil, err
		}
		normalizedFiles := getNormalizedPaths(path, "", result)
		result = normalizedFiles
	} else if strings.HasSuffix(path, ".zip") || strings.HasSuffix(path, ".jar") {
		result, err = WalkArchive(
			".",
			path,
		)
		if err != nil {
			return nil, err
		}
		normalizedFiles := getNormalizedPaths(path, "zip", result)
		result = normalizedFiles
	}
	result.Path = path
	return result, nil
}

func getNormalizedPaths(prefix string, replacement string, result *ArtifactInfo) *ArtifactInfo {
	normalizedFiles := &ArtifactInfo{
		Count:     result.Count,
		FileInfos: FileInfos{},
	}
	for _, v := range result.FileInfos {
		// make the f.Path look like 'path/to/content.txt'
		// so that we can compare relative paths regardless of the prefix
		//v.Path = fmt.Sprintf("%s%s", replacement, filepath.ToSlash(strings.TrimPrefix(v.Path, prefix)))
		v.Path = fmt.Sprintf("%s%s", replacement, filepath.ToSlash(strings.TrimPrefix(filepath.ToSlash(v.Path), filepath.ToSlash(prefix))))
		normalizedFiles.AddFileInfo(v.Path, v)
	}
	return normalizedFiles
}

func WalkTree(root string, rootFS fs.FS) (*ArtifactInfo, error) {
	diff := ArtifactInfo{
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

func WalkArchive(root string, zip string) (*ArtifactInfo, error) {
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
	normalizedFiles := &ArtifactInfo{
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
