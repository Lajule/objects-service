package store_test

import (
	"testing"

	"go.uber.org/zap"

	"github.com/Lajule/objects-service/pkg/store"
)

func TestCreateBucketIfNotExists(t *testing.T) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	st := store.NewStore("test", true, logger)

	if err := st.CreateBucketIfNotExists("bucket"); err != nil {
		t.Errorf("CreateBucketIfNotExists(\"bucket\") = %#v", err)
	}
}

func TestCreateOrOpenObject(t *testing.T) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	st := store.NewStore("test", true, logger)

	_, err := st.CreateOrOpenObject("bucket", "object")
	if err != nil {
		t.Errorf("CreateOrOpenObject(\"bucket\", \"object\") = %#v", err)
	}
}
