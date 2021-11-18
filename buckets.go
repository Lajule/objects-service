package main

import (
	"os"

	"github.com/spf13/afero"
	"go.uber.org/zap"
)

func createBucketIfNotExists(bucketPath string) error {
	bucketExists, err := afs.Exists(bucketPath)
	if err != nil {
		logger.Error("Can not check if bucket exists", zap.Error(err))
		return err
	}

	if !bucketExists {
		if err := afs.MkdirAll(bucketPath, 0755); err != nil {
			logger.Fatal("Can not create bucket",
				zap.String("bucketPath", bucketPath),
				zap.Error(err))
			return err
		}
	}

	return nil
}

func createOrOpenObject(objectPath string) (afero.File, error) {
	objectExists, err := afs.Exists(objectPath)
	if err != nil {
		logger.Error("Can not check if object exists", zap.Error(err))
		return nil, err
	}

	var f afero.File

	if objectExists {
		f, err = afs.OpenFile(objectPath, os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			logger.Error("Can not open object", zap.Error(err))
			return nil, err
		}
	} else {
		f, err = afs.Create(objectPath)
		if err != nil {
			logger.Error("Can not create object", zap.Error(err))
			return nil, err
		}
	}

	return f, nil
}

func getObjectIfExists(objectPath string) (afero.File, error) {
	objectExists, err := afs.Exists(objectPath)
	if err != nil {
		logger.Error("Can not check if object exists", zap.Error(err))
		return nil, err
	}

	if !objectExists {
		return nil, nil
	}

	f, err := afs.Open(objectPath)
	if err != nil {
		logger.Error("Can not open object", zap.Error(err))
		return nil, err
	}

	return f, nil
}

func removeObjectIfExists(objectPath string) (bool, error) {
	objectExists, err := afs.Exists(objectPath)
	if err != nil {
		logger.Error("Can not check if object exists", zap.Error(err))
		return false, err
	}

	if !objectExists {
		return false, nil
	}

	if err := afs.Remove(objectPath); err != nil {
		logger.Error("Can remove object", zap.Error(err))
		return false, err
	}

	return true, nil
}
