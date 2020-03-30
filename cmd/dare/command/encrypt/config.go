package encrypt

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	flags "github.com/da-moon/coe817-dare/cmd/dare/flags"
	dare "github.com/da-moon/coe817-dare/dare"
	mapstructure "github.com/mitchellh/mapstructure"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// DefaultBindPort ...

type dirEnts []os.FileInfo

// Config ...
type Config struct {
	MasterKey       string `mapstructure:"master_key"`
	LogLevel        string `mapstructure:"log_level"`
	Protocol        int    `mapstructure:"protocol"`
	DevelopmentMode bool   `mapstructure:"development_mode"`
}

// DefaultConfig ...
func DefaultConfig() *Config {
	return &Config{
		LogLevel:        "INFO",
		DevelopmentMode: false,
		Protocol:        dare.CoreVersionMax,
	}
}

func (c *Command) readConfig() *Config {
	var cmdConfig Config
	var configFiles []string
	const entrypoint = "encrypt"
	cmdFlags := flag.NewFlagSet(entrypoint, flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()) }
	cmdFlags.Var((*flags.AppendSliceValue)(&configFiles), "config-file",
		"json file to read config from")
	cmdFlags.Var((*flags.AppendSliceValue)(&configFiles), "config-dir",
		"directory of json files to read")
	masterKey := flags.MasterKeyFlag(cmdFlags)
	logLevl := flags.LogLevelFlag(cmdFlags)
	dev := flags.DevFlag(cmdFlags)
	if err := cmdFlags.Parse(c.args); err != nil {
		return nil
	}
	cmdConfig.DevelopmentMode = *dev
	cmdConfig.LogLevel = *logLevl
	cmdConfig.MasterKey = *masterKey
	config := DefaultConfig()
	if len(configFiles) > 0 {
		fileConfig, err := ReadConfigPaths(configFiles)
		if err != nil {
			c.Ui.Error(err.Error())
			return nil
		}
		config = MergeConfig(config, fileConfig)
	}
	config = MergeConfig(config, &cmdConfig)
	if !config.DevelopmentMode {
		config.LogLevel = "INFO"
	}
	return config
}

// EncryptBytes returns the encryption key configured.
func (c *Config) EncryptBytes() ([]byte, error) {
	return base64.StdEncoding.DecodeString(c.MasterKey)
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

	if b.MasterKey != "" {
		result.MasterKey = b.MasterKey
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
