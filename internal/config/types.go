package config

type Config struct {
	Read  ReadConfig  `json:"read" validate:"required"`
	Write WriteConfig `json:"write" validate:"required"`
	Cmd   cmdLineArgs `json:"cmd" validate:"required"`
}

type ReadConfig struct {
	FilePath string `json:"filePath" validate:"required"`
}

type WriteConfig struct {
	WriteDir string `json:"writeDir" validate:"required"`
}

type cmdLineArgs struct {
	FilePath string `json:"filePath" validate:"required"`
	OutDir   string `json:"outDir" validate:"required"`
}
