package configuration

import (
	"io/fs"
	"os"
	"time"
)

type RunescapeConfiguration struct {
	DataDir       string        `required:"true"`
	PollFrequency time.Duration `default:"5m"`
}

func (config RunescapeConfiguration) DataDirFS() fs.FS {
	return os.DirFS(config.DataDir)
}
