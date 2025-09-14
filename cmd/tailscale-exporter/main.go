package main

// Version information (set by build).
var (
	version   = "dev"
	commit    = "unknown"
	buildTime = "unknown"
)

func main() {
	// Set version information in the cmd package
	SetVersionInfo(version, commit, buildTime)

	// Execute the root command
	Execute()
}
