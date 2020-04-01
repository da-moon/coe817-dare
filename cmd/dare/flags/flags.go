package flags

import (
	"flag"
	"os"
)

// MasterKeyFlag ...
func MasterKeyFlag(f *flag.FlagSet) *string {
	result := os.Getenv("DARE_MASTER_KEY")
	if result == "" {
		result = "b6c4bba7a385aef779965cb0b7d66316ab091704042606797871"
	}
	return f.String("master-key", result,
		"Master Key used in encryption-decryption process.")
}

// LogLevelFlag ...
func LogLevelFlag(f *flag.FlagSet) *string {
	result := os.Getenv("DARE_LOG_LEVEL")
	if result == "" {
		result = "INFO"
	}
	return f.String("log", result,
		"flag used to indicate log level")
}

// DevFlag ...
func DevFlag(f *flag.FlagSet) *bool {
	// its false by default
	var result bool
	return f.Bool("dev", result,
		"Enable development mode.")
}
