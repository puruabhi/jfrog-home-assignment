package config

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/go-playground/validator/v10"
)

// NewConfig initializes and validates a new Config instance.
func NewConfig() (*Config, error) {
	config := &Config{}
	config.build()

	if err := validator.New().Struct(config); err != nil {
		return nil, fmt.Errorf("caught err while building config: %w", err)
	}

	configJSON, _ := json.Marshal(config)
	fmt.Println("Parsed config:", string(configJSON))

	return config, nil
}

func (c *Config) build() {
	c.buildCmdLineArgs()
	c.buildReadConfig()
	c.buildWriteConfig()
}

func (c *Config) buildCmdLineArgs() {
	filepath := flag.String("csv-file", "", "CSV File path")
	outDir := flag.String("out-dir", "", "Directory to store downloaded files")

	flag.Parse()

	c.Cmd = cmdLineArgs{
		FilePath: *filepath,
		OutDir:   *outDir,
	}
}

func (c *Config) buildReadConfig() {
	c.Read = ReadConfig{
		FilePath: c.Cmd.FilePath,
	}
}

func (c *Config) buildWriteConfig() {
	c.Write = WriteConfig{
		WriteDir: c.Cmd.OutDir,
	}
}
