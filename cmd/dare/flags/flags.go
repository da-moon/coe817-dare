package flags

import (
	"flag"
	"os"
	"path/filepath"
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

// EncryptorPathFlag ...
func EncryptorPathFlag(f *flag.FlagSet) *string {
	result := os.Getenv("DARE_ENCRYPTOR_PLUGIN")
	if result == "" {
		path, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		result = filepath.Join(path, "bin/encryptor-plugin")
	}
	return f.String("encryptor-path", result,
		"encryptor plugin path.")
}

// DecryptorPathFlag ...
func DecryptorPathFlag(f *flag.FlagSet) *string {
	result := os.Getenv("DARE_DECRYPTOR_PLUGIN")
	if result == "" {
		path, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		result = filepath.Join(path, "bin/decryptor-plugin")
	}
	return f.String("decryptor-path", result,
		"decryptor plugin path.")
}

// APIAddrFlag ...
func APIAddrFlag(f *flag.FlagSet) *string {
	result := os.Getenv("DARE_API_ADDRESS")
	if result == "" {
		result = "127.0.0.1:8080"
	}
	return f.String("api-addr", result,
		"api address to bind the daemon to.")
}

// APIPasswordFlag ...
func APIPasswordFlag(f *flag.FlagSet) *string {
	result := os.Getenv("DARE_API_PASSWORD")
	return f.String("api-password", result,
		"daemon's authorization header value, used for authentication.")
}

// LogLevelFlag ...
func LogLevelFlag(f *flag.FlagSet) *string {
	result := os.Getenv("DARE_LOG_LEVEL")
	if result == "" {
		result = "INFO"
	}
	return f.String("log-level", result,
		"flag used to indicate log level")
}

// DevFlag ...
func DevFlag(f *flag.FlagSet) *bool {
	// its false by default
	var result bool
	return f.Bool("dev", result,
		"Enable development mode.")
}
