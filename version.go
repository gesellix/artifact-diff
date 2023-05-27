package artifact_diff

import "fmt"

var (
	version        = "dev"
	commit         = "n/a"
	buildTimestamp = "n/a"
)

func BuildVersion() string {
	return fmt.Sprintf("%s-%s (%s)", version, commit, buildTimestamp)
}
