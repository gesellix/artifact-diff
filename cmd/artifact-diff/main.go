package main

import (
	"encoding/json"
	"fmt"
	"github.com/urfave/cli/v2"
	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"

	diff "github.com/gesellix/artifact-diff"
)

func main() {
	supportedFormats := []string{"json", "yaml"}
	cliFlags := []cli.Flag{
		&cli.StringFlag{
			Name:    "reportdir",
			Aliases: []string{"t"},
			Usage:   "Path to the report directory `DIR`",
			EnvVars: []string{"TARGETDIR", "TARGET_DIR", "REPORTDIR", "REPORT_DIR"},
			Value:   "./reports",
			Action: func(ctx *cli.Context, dir string) error {
				dir, err := prepareReportDirectory(dir)
				if err != nil {
					return fmt.Errorf("failed to prepare the report directory %s: %w", dir, err)
				}
				log.Println("the report will be written to", dir)
				return nil
			},
		},
		&cli.StringSliceFlag{
			Name:     "reportformat",
			Aliases:  []string{"f"},
			Usage:    fmt.Sprintf("Report format. Supported formats are: %v", supportedFormats),
			EnvVars:  []string{"REPORT_FORMAT", "REPORTFORMAT"},
			Required: false,
			Value:    cli.NewStringSlice("yaml"),
			Action: func(ctx *cli.Context, formats []string) error {
				for _, format := range formats {
					if !slices.Contains(supportedFormats, format) {
						return fmt.Errorf("unsupported report format %s. Supported formats: %v", format, supportedFormats)
					}
					log.Println("the report will be written as:", format)
				}
				return nil
			},
		},
	}
	scanFlags := []cli.Flag{
		&cli.StringSliceFlag{
			Name:     "sourcepath",
			Aliases:  []string{"s"},
			Usage:    "Path(s) to to be scanned. `PATH` may be a path or a zip-compatible file",
			EnvVars:  []string{"SOURCEPATH", "SOURCE_PATH"},
			Required: true,
			Action: func(ctx *cli.Context, paths []string) error {
				for _, path := range paths {
					_, err := os.Stat(filepath.Clean(path))
					if err != nil {
						if os.IsNotExist(err) {
							return fmt.Errorf("path not found: %s", path)
						}
						return fmt.Errorf("invalid path %s: %w", path, err)
					}
					log.Println("path to be scanned:", path)
				}
				return nil
			},
		},
	}
	actionScan := func(ctx *cli.Context) error {
		reportDir := ctx.String("reportdir")
		reportFormats := ctx.StringSlice("reportformat")
		sourcepaths := ctx.StringSlice("sourcepath")

		for _, path := range sourcepaths {
			leftRoot := filepath.Clean(path)
			leftResult, err := diff.CollectFileInfos(leftRoot)
			if err != nil {
				return err
			}
			err = writeReport(reportDir, reportFormats, fmt.Sprintf("path1-%s", leftRoot), leftResult)
			if err != nil {
				return err
			}
		}
		return nil
	}
	cliCommands := []*cli.Command{
		{
			Name:    "scan",
			Aliases: []string{"s"},
			Usage:   "scans one ore more paths and collects metadata",
			Flags:   scanFlags,
			Action:  actionScan,
		},
		{
			Name:    "diff",
			Aliases: []string{"d"},
			Usage:   "extracts differences of the reports in `REPORTDIR`",
			Action: func(ctx *cli.Context) error {
				return fmt.Errorf("not yet implemented")
				//if len(os.Args) > 3 {
				//	sharedPaths := make([]string, 0)
				//	for k := range leftResult.FileInfos {
				//		if _, ok := rightResult.FileInfos[k]; ok {
				//			if leftResult.FileInfos[k].Checksum == rightResult.FileInfos[k].Checksum {
				//				sharedPaths = append(sharedPaths, leftResult.FileInfos[k].Path)
				//			}
				//		}
				//	}
				//	//log.Println("shared", len(sharedPaths), strings.Join(sharedPaths, ","))
				//	for _, k := range sharedPaths {
				//		delete(leftResult.FileInfos, k)
				//		delete(rightResult.FileInfos, k)
				//	}
				//}
			},
		},
	}
	app := &cli.App{
		Name:  "artifact-diff",
		Usage: "Compare directories and zip/jar artifacts",
		Action: func(ctx *cli.Context) error {
			return fmt.Errorf("please choose one of the commands, see the help (-h) for details")
		},
		Version:  diff.BuildVersion(),
		Commands: cliCommands,
		Flags:    cliFlags,
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
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
	//log.Println("reports will be written to", dir)
	return dir, nil
}

func writeReport(reportDir string, reportFormats []string, path string, infos *diff.ArtifactInfo) error {
	log.Println("writing report to", reportDir)

	flat := infos.WithFlattenedAndSortedFileInfos()

	_, file := filepath.Split(path)

	if slices.Contains(reportFormats, "yaml") {
		err := writeYaml(reportDir, file, flat)
		if err != nil {
			return err
		}
	}

	if slices.Contains(reportFormats, "json") {
		err := writeJson(reportDir, file, flat)
		if err != nil {
			return err
		}
	}

	return nil
}

func writeYaml(reportDir string, file string, content interface{}) error {
	yamlResult, _ := yaml.Marshal(content)
	filename := fmt.Sprintf("%s.yaml", filepath.Join(reportDir, file))
	err := os.WriteFile(filename, yamlResult, 0644)
	if err != nil {
		return err
	}
	log.Println("report (yaml) written to", filename)
	//log.Println(string(yamlResult))
	return nil
}

func writeJson(reportDir string, file string, content interface{}) error {
	jsonResult, _ := json.Marshal(content)
	filename := fmt.Sprintf("%s.json", filepath.Join(reportDir, file))
	err := os.WriteFile(filename, jsonResult, 0644)
	if err != nil {
		return err
	}
	log.Println("report (json) written to", filename)
	//fmt.Println(string(jsonResult))
	return nil
}
