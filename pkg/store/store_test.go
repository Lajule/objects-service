package store_test

import (
	"testing"

	"go.uber.org/zap"

	"github.com/Lajule/objects-service/store"
)

func TestCreateBucketIfNotExists(t *testing.T) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	st := store.NewStore(logger, true, "test")

	if err := st.CreateBucketIfNotExists("bucket"); err != nil {
		t.Errorf("CreateBucketIfNotExists(\"bucket\") = %#v ", err)
	}
}
