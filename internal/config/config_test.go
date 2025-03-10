package config

import (
	"flag"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const csvFileArg = "csv-file"
const outDirArg = "out-dir"

func resetFlags() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	flag.CommandLine.SetOutput(io.Discard)
}

func TestNewConfig(t *testing.T) {
	resetFlags()
	// Set up command line arguments

	testCSVFile := "/path/to/dummy/dir/test.csv"
	testOutDir := "/path/to/dummy/dir/output"

	os.Args = []string{
		"dummy", // first arg is path to file
		fmt.Sprintf("--%s=%s", csvFileArg, testCSVFile),
		fmt.Sprintf("--%s=%s", outDirArg, testOutDir),
	}

	config, err := NewConfig()
	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, testCSVFile, config.Cmd.FilePath)
	assert.Equal(t, testOutDir, config.Cmd.OutDir)
	assert.Equal(t, testCSVFile, config.Read.FilePath)
	assert.Equal(t, testOutDir, config.Write.WriteDir)
}

func TestNewConfigError(t *testing.T) {
	resetFlags()
	// Set up command line arguments

	testCSVFile := "/path/to/dummy/dir/test.csv"

	os.Args = []string{
		"dummy", // first arg is path to file
		fmt.Sprintf("--%s=%s", csvFileArg, testCSVFile),
	}

	config, err := NewConfig()
	assert.Error(t, err)
	assert.Nil(t, config)
}
