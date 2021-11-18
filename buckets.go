package main

import (
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
