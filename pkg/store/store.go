package store

import (
	"os"
	"path/filepath"

	"github.com/spf13/afero"
	"go.uber.org/zap"
)

// Store contains file system abstraction
type Store struct {
	basePath string
	afs      *afero.Afero
	logger   *zap.Logger
}

// NewStore creates a new store
func NewStore(basePath string, memMapFs bool, logger *zap.Logger) *Store {
	var fs afero.Fs

	logger.Info("creating store",
		zap.String("basePath", basePath),
		zap.Bool("memMapFs", memMapFs))

	if memMapFs {
		fs = afero.NewMemMapFs()
	} else {
		fs = afero.NewOsFs()
	}

	if err := fs.MkdirAll(basePath, 0755); err != nil {
		logger.Fatal("can not create directory",
			zap.String("basePath", basePath),
			zap.Error(err))
	}

	return &Store{
		basePath: basePath,
		afs:      &afero.Afero{Fs: fs},
		logger:   logger.Named("store"),
	}
}

// CreateBucketIfNotExists creates a bucket
func (s *Store) CreateBucketIfNotExists(bucket string) error {
	bucketPath := filepath.Join(s.basePath, bucket)

	bucketExists, err := s.afs.Exists(bucketPath)
	if err != nil {
		s.logger.Error("failed to check if bucket exists", zap.Error(err))
		return err
	}

	if !bucketExists {
		if err := s.afs.MkdirAll(bucketPath, 0755); err != nil {
			s.logger.Fatal("can not create bucket",
				zap.String("bucketPath", bucketPath),
				zap.Error(err))
			return err
		}
	}

	return nil
}

// CreateOrOpenObject creates an object
func (s *Store) CreateOrOpenObject(bucket, objectID string) (afero.File, error) {
	objectPath := filepath.Join(s.basePath, bucket, objectID)

	objectExists, err := s.afs.Exists(objectPath)
	if err != nil {
		s.logger.Error("failed to check if object exists", zap.Error(err))
		return nil, err
	}

	var f afero.File

	if objectExists {
		f, err = s.afs.OpenFile(objectPath, os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			s.logger.Error("can not open object", zap.Error(err))
			return nil, err
		}
	} else {
		f, err = s.afs.Create(objectPath)
		if err != nil {
			s.logger.Error("can not create object", zap.Error(err))
			return nil, err
		}
	}

	return f, nil
}

// GetObjectIfExists get an object
func (s *Store) GetObjectIfExists(bucket, objectID string) (afero.File, error) {
	objectPath := filepath.Join(s.basePath, bucket, objectID)

	objectExists, err := s.afs.Exists(objectPath)
	if err != nil {
		s.logger.Error("failed to check if object exists", zap.Error(err))
		return nil, err
	}

	if !objectExists {
		return nil, nil
	}

	f, err := s.afs.Open(objectPath)
	if err != nil {
		s.logger.Error("can not open object", zap.Error(err))
		return nil, err
	}

	return f, nil
}

// RemoveObjectIfExists deletes an object
func (s *Store) RemoveObjectIfExists(bucket, objectID string) (bool, error) {
	objectPath := filepath.Join(s.basePath, bucket, objectID)

	objectExists, err := s.afs.Exists(objectPath)
	if err != nil {
		s.logger.Error("failed to check if object exists", zap.Error(err))
		return false, err
	}

	if !objectExists {
		return false, nil
	}

	if err := s.afs.Remove(objectPath); err != nil {
		s.logger.Error("can not remove object", zap.Error(err))
		return false, err
	}

	return true, nil
}
