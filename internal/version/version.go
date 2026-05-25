package version

const (
	AppName     = "dbreport"
	Description = "Generate clean HTML reports from MariaDB queries"
	Author      = "Richard S. Westmoreland"
	Email       = "dev@rswestmore.land"
	License     = "MIT"
	Copyright   = "Copyright (c) 2026 Richard S. Westmoreland"
)

var (
	Version = "0.1.0-dev"
	Commit  = "unknown"
	Date    = "unknown"
)

func Short() string {
	return AppName + " " + Version
}
