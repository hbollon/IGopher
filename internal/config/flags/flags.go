package flags

import "flag"

// Flags declarations
var Flags = struct {
	// LogLevelFlag set loglevel threshold
	// If undefined or wrong set it to INFO level
	LogLevelFlag *string

	// ForceDlFlag force re-download of all dependencies
	ForceDlFlag *bool

	// DebugFlag set selenium debug mode and display its logging to stderr
	DebugFlag *bool

	// DevToolsFlag launch Electron gui with devtools openned
	DevToolsFlag *bool

	// IgnoreDependenciesFlag disable dependencies manager on startup
	IgnoreDependenciesFlag *bool

	// BackgroundFlag IGopher as background task with actual configuration and ignore TUI
	BackgroundFlag *bool

	// HeadlessFlag execute Selenium webdriver in headless mode
	HeadlessFlag *bool

	// PortFlag specifie custom communication port for Selenium and web drivers
	PortFlag *int
}{
	LogLevelFlag:           flag.String("loglevel", "info", "Log level threshold"),
	ForceDlFlag:            flag.Bool("force-download", false, "Force redownload of all dependencies even if exists"),
	DebugFlag:              flag.Bool("debug", false, "Display debug and selenium output"),
	DevToolsFlag:           flag.Bool("dev-tools", false, "Launch Electron gui with dev tools openned"),
	IgnoreDependenciesFlag: flag.Bool("ignore-dependencies", false, "Skip dependencies management"),
	HeadlessFlag:           flag.Bool("headless", false, "Run WebDriver with frame buffer"),
	PortFlag:               flag.Int("port", 8080, "Specify custom communication port"),
}
