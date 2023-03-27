package daemon

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	flags "github.com/da-moon/dare-cli/cmd/dare/flags"
	daemon "github.com/da-moon/dare-cli/daemon"
	mapstructure "github.com/mitchellh/mapstructure"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const devNonce = "85d2bdff2b3c8b83814dd20da55ec1d5449f815d21b11d95"
const devKey = "48868145ed69ef5b7d37346470ade518c38f47e91a8102f1b83419d69bb6b835"
const jwtSecretLen = 32

type dirEnts []os.FileInfo

// Config ...
type Config struct {
	EncryptorPath   string `mapstructure:"encryptor_path"`
	DecryptorPath   string `mapstructure:"decryptor_path"`
	APIAddr         string `mapstructure:"api_addr"`
	APIPassword     string `mapstructure:"api_password"`
	LogLevel        string `mapstructure:"log_level"`
	Protocol        int    `mapstructure:"protocol"`
	DevelopmentMode bool   `mapstructure:"development_mode"`
}

// DefaultConfig ...
func DefaultConfig() *Config {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	bytes := make([]byte, jwtSecretLen)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)

	}

	return &Config{
		LogLevel:        "INFO",
		DevelopmentMode: false,
		Protocol:        daemon.CoreVersionMax,
		APIAddr:         "127.0.0.1:8080",
		APIPassword:     base64.StdEncoding.EncodeToString(bytes),
		EncryptorPath:   filepath.Join(path, "encryptor"),
		DecryptorPath:   filepath.Join(path, "decryptor"),
	}
}

func (c *Command) readConfig() *Config {
	var cmdConfig Config
	var configFiles []string
	const entrypoint = "daemon"
	cmdFlags := flag.NewFlagSet(entrypoint, flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()) }
	cmdFlags.Var((*flags.AppendSliceValue)(&configFiles), "config-file",
		"json file to read config from")
	cmdFlags.Var((*flags.AppendSliceValue)(&configFiles), "config-dir",
		"directory of json files to read")
	logLevel := flags.LogLevelFlag(cmdFlags)
	encryptorPath := flags.EncryptorPathFlag(cmdFlags)
	decryptorPath := flags.DecryptorPathFlag(cmdFlags)
	apiAddr := flags.APIAddrFlag(cmdFlags)
	apiPassword := flags.APIPasswordFlag(cmdFlags)
	dev := flags.DevFlag(cmdFlags)
	if err := cmdFlags.Parse(c.args); err != nil {
		return nil
	}
	cmdConfig.DevelopmentMode = *dev
	cmdConfig.LogLevel = *logLevel
	cmdConfig.EncryptorPath = *encryptorPath
	cmdConfig.DecryptorPath = *decryptorPath
	cmdConfig.APIAddr = *apiAddr
	cmdConfig.APIPassword = *apiPassword
	if len(cmdConfig.APIPassword) == 0 {
		c.Ui.Warn("[WARN] Daemon API password was not given. Generating a random one")
		bytes := make([]byte, jwtSecretLen)
		if _, err := rand.Read(bytes); err != nil {
			c.Ui.Error(fmt.Sprintf("[ERROR]: %s", err.Error()))
			return nil
		}
		cmdConfig.APIPassword = base64.StdEncoding.EncodeToString(bytes)
	}
	config := DefaultConfig()
	if len(configFiles) > 0 {
		fileConfig, err := ReadConfigPaths(configFiles)
		if err != nil {
			c.Ui.Error(fmt.Sprintf("[ERROR]: %s", err.Error()))
			return nil
		}
		config = MergeConfig(config, fileConfig)
	}
	config = MergeConfig(config, &cmdConfig)
	return config
}

// DecodeConfig ...
func DecodeConfig(r io.Reader) (*Config, error) {
	var raw interface{}
	dec := json.NewDecoder(r)
	if err := dec.Decode(&raw); err != nil {
		return nil, err
	}

	var md mapstructure.Metadata
	var result Config
	msdec, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Metadata:    &md,
		Result:      &result,
		ErrorUnused: true,
	})
	if err != nil {
		return nil, err
	}

	if err := msdec.Decode(raw); err != nil {
		return nil, err
	}

	return &result, nil
}

func containsKey(keys []string, key string) bool {
	for _, k := range keys {
		if k == key {
			return true
		}
	}
	return false
}

// MergeConfig ...
func MergeConfig(a, b *Config) *Config {
	result := *a

	if b.EncryptorPath != "" {
		result.EncryptorPath = b.EncryptorPath
	}
	if b.DecryptorPath != "" {
		result.DecryptorPath = b.DecryptorPath
	}
	if b.APIAddr != "" {
		result.APIAddr = b.APIAddr
	}
	if b.LogLevel != "" {
		result.LogLevel = b.LogLevel
	}
	if b.Protocol > 0 {
		result.Protocol = b.Protocol
	}
	result.DevelopmentMode = b.DevelopmentMode

	return &result
}

// ReadConfigPaths reads the paths in the given order to load configurations.
func ReadConfigPaths(paths []string) (*Config, error) {
	result := new(Config)
	for _, path := range paths {
		f, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("Error reading '%s': %s", path, err)
		}

		fi, err := f.Stat()
		if err != nil {
			f.Close()
			return nil, fmt.Errorf("Error reading '%s': %s", path, err)
		}

		if !fi.IsDir() {
			config, err := DecodeConfig(f)
			f.Close()

			if err != nil {
				return nil, fmt.Errorf("Error decoding '%s': %s", path, err)
			}

			result = MergeConfig(result, config)
			continue
		}

		contents, err := f.Readdir(-1)
		f.Close()
		if err != nil {
			return nil, fmt.Errorf("Error reading '%s': %s", path, err)
		}

		sort.Sort(dirEnts(contents))

		for _, fi := range contents {
			if fi.IsDir() {
				continue
			}

			if !strings.HasSuffix(fi.Name(), ".json") {
				continue
			}

			subpath := filepath.Join(path, fi.Name())
			f, err := os.Open(subpath)
			if err != nil {
				return nil, fmt.Errorf("Error reading '%s': %s", subpath, err)
			}

			config, err := DecodeConfig(f)
			f.Close()

			if err != nil {
				return nil, fmt.Errorf("Error decoding '%s': %s", subpath, err)
			}

			result = MergeConfig(result, config)
		}
	}

	return result, nil
}

// Len ...
func (d dirEnts) Len() int {
	return len(d)
}

// Less ...
func (d dirEnts) Less(i, j int) bool {
	return d[i].Name() < d[j].Name()
}

// Swap ...
func (d dirEnts) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}
