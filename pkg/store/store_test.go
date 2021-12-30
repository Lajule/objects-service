package store_test

import (
	"testing"

	"github.com/matryer/is"
	"go.uber.org/zap"

	"github.com/Lajule/objects-service/pkg/store"
)

func TestStore(t *testing.T) {
	is := is.New(t)

	logger, _ := zap.NewProduction()
	st := store.NewStore("test", true, logger)

	t.Run("CreateBucketIfNotExists", func(t *testing.T) {
		err := st.CreateBucketIfNotExists("bucket")
		is.NoErr(err)
	})

	t.Run("CreateOrOpenObject", func(t *testing.T) {
		_, err := st.CreateOrOpenObject("bucket", "object")
		is.NoErr(err)
	})

	t.Run("GetObjectIfExists", func(t *testing.T) {
		_, err := st.CreateOrOpenObject("bucket", "object")
		is.NoErr(err)

		f, err := st.GetObjectIfExists("bucket", "object")
		is.NoErr(err)
		is.True(f != nil)

		f2, err := st.GetObjectIfExists("bucket", "object2")
		is.NoErr(err)
		is.True(f2 == nil)
	})

	t.Run("RemoveObjectIfExists", func(t *testing.T) {
		_, err := st.CreateOrOpenObject("bucket", "object")
		is.NoErr(err)

		removed, err := st.RemoveObjectIfExists("bucket", "object")
		is.NoErr(err)
		is.True(removed)

		removed, err = st.RemoveObjectIfExists("bucket", "object")
		is.NoErr(err)
		is.True(!removed)
	})
}
