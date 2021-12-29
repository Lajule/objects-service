package store_test

import (
	"testing"

	"github.com/matryer/is"
	"go.uber.org/zap"

	"github.com/Lajule/objects-service/pkg/store"
)

var (
	logger *zap.Logger

	st *store.Store
)

func init() {
	logger, _ = zap.NewProduction()
	st = store.NewStore("test", true, logger)
}

func TestCreateBucketIfNotExists(t *testing.T) {
	is := is.New(t)

	err := st.CreateBucketIfNotExists("bucket")
	is.NoErr(err)
}

func TestCreateOrOpenObject(t *testing.T) {
	is := is.New(t)

	_, err := st.CreateOrOpenObject("bucket", "object")
	is.NoErr(err)
}

func TestGetObjectIfExists(t *testing.T) {
	is := is.New(t)

	_, err := st.CreateOrOpenObject("bucket", "object")
	is.NoErr(err)

	f, err := st.GetObjectIfExists("bucket", "object")
	is.NoErr(err)
	is.True(f != nil)

	f2, err := st.GetObjectIfExists("bucket", "object2")
	is.NoErr(err)
	is.True(f2 == nil)
}

func TestRemoveObjectIfExists(t *testing.T) {
	is := is.New(t)

	_, err := st.CreateOrOpenObject("bucket", "object")
	is.NoErr(err)

	removed, err := st.RemoveObjectIfExists("bucket", "object")
	is.NoErr(err)
	is.True(removed)

	removed, err = st.RemoveObjectIfExists("bucket", "object")
	is.NoErr(err)
	is.True(!removed)
}
