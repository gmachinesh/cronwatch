package api

import (
	"runtime"
	"time"
)

// BuildInfo holds version and build metadata for the running binary.
type BuildInfo struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	BuiltAt   string `json:"built_at"`
	GoVersion string `json:"go_version"`
	OS        string `json:"os"`
	Arch      string `json:"arch"`
}

// These variables are set at build time via -ldflags.
var (
	Version = "dev"
	Commit  = "none"
	BuiltAt = ""
)

func buildInfo() BuildInfo {
	builtAt := BuiltAt
	if builtAt == "" {
		builtAt = time.Now().UTC().Format(time.RFC3339)
	}
	return BuildInfo{
		Version:   Version,
		Commit:    Commit,
		BuiltAt:   builtAt,
		GoVersion: runtime.Version(),
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
	}
}

func (s *Server) handleVersion(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, buildInfo())
}
