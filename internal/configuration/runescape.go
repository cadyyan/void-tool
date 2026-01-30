package configuration

import (
	"io/fs"
	"os"
)

type RunescapeConfiguration struct {
	DataDir string `required:"true"`
}

func (config RunescapeConfiguration) DataDirFS() fs.FS {
	return os.DirFS(config.DataDir)
}
