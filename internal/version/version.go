package version

import "fmt"

// Build information, set via ldflags.
var (
	Version   = "dev"
	Commit    = "none"
	BuildDate = "unknown"
)

// Info holds version metadata.
type Info struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	BuildDate string `json:"build_date"`
}

// Get returns the current version info.
func Get() Info {
	return Info{
		Version:   Version,
		Commit:    Commit,
		BuildDate: BuildDate,
	}
}

// String returns a human-readable version string.
func (i Info) String() string {
	return fmt.Sprintf("vaultpull %s (commit: %s, built: %s)", i.Version, i.Commit, i.BuildDate)
}

// Short returns just the version number.
func (i Info) Short() string {
	return i.Version
}
