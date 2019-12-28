package mysql

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vespaiach/auth_service/pkg/share"
	"github.com/vespaiach/auth_service/pkg/storage"
)

func TestKeyMysqlStorer_Insert(t *testing.T) {
	t.Parallel()

	t.Run("success_add_a_key", func(t *testing.T) {
		t.Parallel()

		key := test.mig.createUniqueString("key")
		desc := test.mig.createUniqueString("desc")

		id, err := test.kst.Insert(storage.CreateKey{Name: key, Desc: desc})
		require.Nil(t, err)
		require.NotZero(t, id)
	})

	t.Run("fail_add_a_duplicated_key", func(t *testing.T) {
		t.Parallel()

		key := test.mig.createUniqueString("key")
		desc1 := test.mig.createUniqueString("desc")
		desc2 := test.mig.createUniqueString("desc")

		test.mig.createSeedingServiceKey(func(fields map[string]interface{}) {
			fields["name"] = key
			fields["desc"] = desc1
		})

		dupID, errDup := test.kst.Insert(storage.CreateKey{Name: key, Desc: desc2})
		require.NotNil(t, errDup)
		require.Zero(t, dupID)
	})
}

func TestKeyMysqlStorer_Update(t *testing.T) {
	t.Parallel()

	t.Run("success_update_a_key", func(t *testing.T) {
		t.Parallel()

		id := test.mig.createSeedingServiceKey(nil)
		key := test.mig.createUniqueString("key")
		desc := test.mig.createUniqueString("desc")

		err := test.kst.Update(storage.UpdateKey{
			ID:   id,
			Name: key,
			Desc: desc,
		})
		require.Nil(t, err)
	})
}

func TestKeyMysqlStorer_Delete(t *testing.T) {
	t.Parallel()

	t.Run("success_delete_a_key", func(t *testing.T) {
		t.Parallel()

		id := test.mig.createSeedingServiceKey(nil)

		err := test.kst.Delete(id)
		require.Nil(t, err)
	})
}

func TestKeyMysqlStorer_Get(t *testing.T) {
	t.Parallel()

	t.Run("success_get_a_key_by_id", func(t *testing.T) {
		t.Parallel()

		key := test.mig.createUniqueString("key")
		desc := test.mig.createUniqueString("desc")

		id := test.mig.createSeedingServiceKey(func(fields map[string]interface{}) {
			fields["name"] = key
			fields["desc"] = desc
		})

		found, err := test.kst.Get(id)
		require.Nil(t, err)
		require.NotNil(t, found)
		require.Equal(t, found.ID, id)
		require.Equal(t, found.Name, key)
		require.Equal(t, found.Desc, desc)
	})
}

func TestKeyMysqlStorer_GetByName(t *testing.T) {
	t.Parallel()

	t.Run("success_get_a_key_by_name", func(t *testing.T) {
		t.Parallel()

		key := test.mig.createUniqueString("key")
		desc := test.mig.createUniqueString("desc")

		id := test.mig.createSeedingServiceKey(func(fields map[string]interface{}) {
			fields["name"] = key
			fields["desc"] = desc
		})

		found, err := test.kst.GetByName(key)
		require.Nil(t, err)
		require.NotNil(t, found)
		require.Equal(t, found.ID, id)
		require.Equal(t, found.Name, key)
		require.Equal(t, found.Desc, desc)
	})
}

func TestKeyMysqlStorer_Query(t *testing.T) {
	t.Parallel()

	t.Run("success_query_keys", func(t *testing.T) {
		t.Parallel()
		prefix := test.mig.createUniqueString("pre")
		name1 := test.mig.createUniqueString(prefix)
		name2 := test.mig.createUniqueString(prefix)
		name3 := test.mig.createUniqueString(prefix)
		name4 := test.mig.createUniqueString(prefix)
		name5 := test.mig.createUniqueString(prefix)
		name6 := test.mig.createUniqueString(prefix)

		test.mig.createSeedingServiceKey(func(fields map[string]interface{}) { fields["name"] = name1 })
		test.mig.createSeedingServiceKey(func(fields map[string]interface{}) { fields["name"] = name2 })
		test.mig.createSeedingServiceKey(func(fields map[string]interface{}) { fields["name"] = name3 })
		test.mig.createSeedingServiceKey(func(fields map[string]interface{}) { fields["name"] = name4 })
		test.mig.createSeedingServiceKey(func(fields map[string]interface{}) { fields["name"] = name5 })
		test.mig.createSeedingServiceKey(func(fields map[string]interface{}) { fields["name"] = name6 })

		rows, total, err := test.kst.Query(storage.QueryKey{
			Limit:  2,
			Offset: 2,
			Name:   prefix,
		}, storage.SortKey{
			Name: share.Ascendant,
		})
		require.Nil(t, err)
		require.NotNil(t, rows)
		require.Equal(t, int64(6), total)
		require.Len(t, rows, 2)
	})
}
