package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
	"strings"

	diff "github.com/gesellix/artifact-diff"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println(fmt.Sprintf("Usage: %s <report directory> <path1> [path2]", os.Args[0]))
		os.Exit(1)
	}

	reportDir, err := prepareReportDirectory(os.Args[1])
	if err != nil {
		panic(err)
	}

	leftRoot := filepath.Clean(os.Args[2])
	leftResult, err := collectFileInfos(leftRoot)
	if err != nil {
		panic(err)
	}
	writeReport(reportDir, fmt.Sprintf("path1-%s", leftRoot), leftResult)

	if len(os.Args) > 3 {
		rightRoot := filepath.Clean(os.Args[3])
		rightResult, err := collectFileInfos(rightRoot)
		if err != nil {
			panic(err)
		}
		writeReport(reportDir, fmt.Sprintf("path2-%s", rightRoot), rightResult)

		sharedPaths := make([]string, 0)
		for k := range leftResult.FileInfos {
			if _, ok := rightResult.FileInfos[k]; ok {
				if leftResult.FileInfos[k].Checksum == rightResult.FileInfos[k].Checksum {
					sharedPaths = append(sharedPaths, leftResult.FileInfos[k].Path)
				}
			}
		}
		//log.Println("shared", len(sharedPaths), strings.Join(sharedPaths, ","))
		for _, k := range sharedPaths {
			delete(leftResult.FileInfos, k)
			delete(rightResult.FileInfos, k)
		}
	}
}

func prepareReportDirectory(reportDir string) (string, error) {
	dir, err := filepath.Abs(filepath.Clean(reportDir))
	if err != nil {
		return "", err
	}
	if _, e := os.Stat(dir); os.IsNotExist(e) {
		err := os.MkdirAll(dir, 0777)
		if err != nil {
			return "", err
		}
	}
	//log.Println("Reports will be written to", dir)
	return dir, nil
}

func collectFileInfos(path string) (*diff.Diff, error) {
	log.Println("Scanning", path)

	stat, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	var result *diff.Diff
	if stat.IsDir() {
		result, err = diff.WalkTree(
			path,
			os.DirFS(path),
		)
		if err != nil {
			return nil, err
		}
		normalizedFiles := getNormalizedPaths(path, "", result)
		result = normalizedFiles
	} else if strings.HasSuffix(path, ".zip") || strings.HasSuffix(path, ".jar") {
		result, err = diff.WalkArchive(
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

func getNormalizedPaths(prefix string, replacement string, result *diff.Diff) *diff.Diff {
	normalizedFiles := &diff.Diff{
		Count:     result.Count,
		FileInfos: diff.FileInfos{},
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

func writeReport(reportDir string, path string, infos *diff.Diff) {
	log.Println("Writing report to", reportDir)

	flat := infos.WithFlattenedAndSortedFileInfos()

	_, file := filepath.Split(path)
	writeJson(reportDir, file, flat)
	writeYaml(reportDir, file, flat)
}

func writeYaml(reportDir string, file string, content interface{}) {
	yamlResult, _ := yaml.Marshal(content)
	filename := fmt.Sprintf("%s.yaml", filepath.Join(reportDir, file))
	os.WriteFile(filename, yamlResult, 0644)
	log.Println("Report (yaml) written to", filename)
	//log.Println(string(yamlResult))
}

func writeJson(reportDir string, file string, content interface{}) {
	jsonResult, _ := json.Marshal(content)
	filename := fmt.Sprintf("%s.json", filepath.Join(reportDir, file))
	os.WriteFile(filename, jsonResult, 0644)
	log.Println("Report (json) written to", filename)
	//fmt.Println(string(jsonResult))
}
