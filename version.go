package monolith

import "time"

var (
	Version   = "dev"
	Commit    = "dev"
	BuildTime = time.Now().Format(time.RFC3339) // Default to current time if not set
)

type VersionInfo struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	BuildTime string `json:"buildTime"`
}

func GetVersionInfo() VersionInfo {
	return VersionInfo{
		Version:   Version,
		Commit:    Commit,
		BuildTime: BuildTime,
	}
}

func IsDevEnv() bool {
	return Version == "dev"
}
