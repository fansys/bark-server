package misc

var (
	Version   string
	BuildDate string
	CommitID  string
)

func InitVersion(version string, buildDate string, commitId string) {
	Version = version
	BuildDate = buildDate
	CommitID = commitId
}
