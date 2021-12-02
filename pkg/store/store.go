package store

import (
	"os"
	"path/filepath"

	"github.com/spf13/afero"
	"go.uber.org/zap"
)

type Store struct {
	afs     *afero.Afero
	log     *zap.Logger
	rootDir string
}

func NewStore(logger *zap.Logger, memory bool, rootDir string) *Store {
	var fs afero.Fs

	logger.Info("Creating store",
		zap.Bool("memory", memory),
		zap.String("rootDir", rootDir))

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

	return &Store{
		afs:     &afero.Afero{Fs: fs},
		log:     logger.Named("store"),
		rootDir: rootDir,
	}
}

func (s *Store) CreateBucketIfNotExists(bucket string) error {
	bucketPath := filepath.Join(s.rootDir, bucket)

	bucketExists, err := s.afs.Exists(bucketPath)
	if err != nil {
		s.log.Error("Can not check if bucket exists", zap.Error(err))
		return err
	}

	if !bucketExists {
		if err := s.afs.MkdirAll(bucketPath, 0755); err != nil {
			s.log.Fatal("Can not create bucket",
				zap.String("bucketPath", bucketPath),
				zap.Error(err))
			return err
		}
	}

	return nil
}

func (s *Store) CreateOrOpenObject(bucket, objectID string) (afero.File, error) {
	objectPath := filepath.Join(s.rootDir, bucket, objectID)

	objectExists, err := s.afs.Exists(objectPath)
	if err != nil {
		s.log.Error("Can not check if object exists", zap.Error(err))
		return nil, err
	}

	var f afero.File

	if objectExists {
		f, err = s.afs.OpenFile(objectPath, os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			s.log.Error("Can not open object", zap.Error(err))
			return nil, err
		}
	} else {
		f, err = s.afs.Create(objectPath)
		if err != nil {
			s.log.Error("Can not create object", zap.Error(err))
			return nil, err
		}
	}

	return f, nil
}

func (s *Store) GetObjectIfExists(bucket, objectID string) (afero.File, error) {
	objectPath := filepath.Join(s.rootDir, bucket, objectID)

	objectExists, err := s.afs.Exists(objectPath)
	if err != nil {
		s.log.Error("Can not check if object exists", zap.Error(err))
		return nil, err
	}

	if !objectExists {
		return nil, nil
	}

	f, err := s.afs.Open(objectPath)
	if err != nil {
		s.log.Error("Can not open object", zap.Error(err))
		return nil, err
	}

	return f, nil
}

func (s *Store) RemoveObjectIfExists(bucket, objectID string) (bool, error) {
	objectPath := filepath.Join(s.rootDir, bucket, objectID)

	objectExists, err := s.afs.Exists(objectPath)
	if err != nil {
		s.log.Error("Can not check if object exists", zap.Error(err))
		return false, err
	}

	if !objectExists {
		return false, nil
	}

	if err := s.afs.Remove(objectPath); err != nil {
		s.log.Error("Can remove object", zap.Error(err))
		return false, err
	}

	return true, nil
}
