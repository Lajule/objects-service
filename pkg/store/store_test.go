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
		t.Errorf("CreateBucketIfNotExists(\"bucket\") = %#v ", err)
	}
}
