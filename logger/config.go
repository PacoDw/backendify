package logger

import (
	"os"

	"github.com/spf13/cast"
)

// EnvironmentMode represent the type of the environment that the logger is running.
type EnvironmentMode string

const (
	// Production is used in production that removes unnecessary logs.
	Production EnvironmentMode = "PRODUCTION"

	// Development by default this is enabled that displays debug logs.
	Development EnvironmentMode = "DEVELOPENT"
)

// Config represents the Logger configuration.
type Config struct {
	// OutputFile is the file to write logs to.  Backup log files will be retained
	// in the same directory.  It uses <processname>-lumberjack.log in
	// os.TempDir() if empty.
	OutputFile string

	// MaxSize is the maximum size in megabytes of the log file before it gets
	// rotated. It defaults to 100 megabytes.
	MaxSize int

	// MaxBackups is the maximum number of old log files to retain.  The default
	// is to retain all old log files (though MaxAge may still cause them to get
	// deleted.)
	MaxBackups int

	// MaxAge is the maximum number of days to retain old log files based on the
	// timestamp encoded in their filename.  Note that a day is defined as 24
	// hours and may not exactly correspond to calendar days due to daylight
	// savings, leap seconds, etc. The default is not to remove old log files
	// based on age.
	MaxAge int

	// Environment determines if the logger is on mode Development or Production, by default
	// it is set as Development.
	Environment EnvironmentMode
}

// DefaultEnvLoggerConfig gets the set env variables to create a Config.
func DefaultEnvLoggerConfig() *Config {
	path := os.Getenv("OUTPUT_FILE")

	if path == "" {
		path = "./logfile.log"
	}

	envMode := Development
	if os.Getenv("ENVIRONMENT") == "PRODUCTION" {
		envMode = Production
	}

	return &Config{
		OutputFile:  path,
		MaxSize:     cast.ToInt(os.Getenv("MAX_SIZE")),
		MaxBackups:  cast.ToInt(os.Getenv("MAX_BACKUPS")),
		MaxAge:      cast.ToInt(os.Getenv("MAX_AGE")),
		Environment: envMode,
	}
}
