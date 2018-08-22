// Package version defines "hello-world" version.
package version

var (
	// GitCommit is the git commit on build.
	GitCommit = ""
	// ReleaseVersion is the release version.
	ReleaseVersion = ""
	// BuildTime is the build timestamp.
	BuildTime = ""
)

// Version represents the version information.
type Version struct {
	GitCommit      string `json:"gitCommit"`
	ReleaseVersion string `json:"releaseVersion"`
	BuildTime      string `json:"buildTime"`
	HostName       string `json:"hostname"`
}
