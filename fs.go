package main

import (
	"github.com/spf13/afero"
	"go.uber.org/zap"
)

func initFs() *afero.Afero {
	var fs afero.Fs
	if memory {
		fs = afero.NewMemMapFs()
	} else {
		fs = afero.NewOsFs()
	}

	if err := fs.MkdirAll(rootDir, 0755); err != nil {
		logger.Fatal("Can not create root directory",
			zap.String("rootDir", rootDir),
			zap.Error(err))
	}

	return &afero.Afero{Fs: fs}
}
