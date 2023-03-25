package artifact_diff

import "fmt"

var (
	version        = "dev"
	commitHash     = "n/a"
	buildTimestamp = "n/a"
)

func BuildVersion() string {
	return fmt.Sprintf("%s-%s (%s)", version, commitHash, buildTimestamp)
}
