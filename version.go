package monolith

import "time"

var (
	Version   = "v0.0.0"
	Commit    = "unknown"
	DateBuilt = time.Now().Format(time.RFC3339) // Default to current time if not set
)

type VersionInfo struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	DateBuilt string `json:"dateBuilt"`
}

func GetVersionInfo() VersionInfo {
	return VersionInfo{
		Version:   Version,
		Commit:    Commit,
		DateBuilt: DateBuilt,
	}
}
