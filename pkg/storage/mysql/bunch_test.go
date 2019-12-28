package mysql

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vespaiach/auth_service/pkg/share"
	"github.com/vespaiach/auth_service/pkg/storage"
)

func TestBunchMysqlStorer_Insert(t *testing.T) {
	t.Parallel()

	t.Run("success_add_a_bunch", func(t *testing.T) {
		t.Parallel()

		bunch := test.mig.createUniqueString("bunch")
		desc := test.mig.createUniqueString("desc")

		id, err := test.bst.Insert(storage.CreateBunch{Name: bunch, Desc: desc})
		require.Nil(t, err)
		require.NotZero(t, id)
	})

	t.Run("fail_add_a_duplicated_bunch", func(t *testing.T) {
		t.Parallel()

		bunch := test.mig.createUniqueString("bunch")
		desc := test.mig.createUniqueString("desc")

		test.mig.createSeedingBunch(func(fields map[string]interface{}) {
			fields["name"] = bunch
			fields["desc"] = desc
		})

		id, err := test.bst.Insert(storage.CreateBunch{Name: bunch, Desc: desc})
		require.NotNil(t, err)
		require.Zero(t, id)
	})
}

func TestBunchMysqlStorer_Update(t *testing.T) {
	t.Parallel()

	t.Run("success_update_a_bunch", func(t *testing.T) {
		t.Parallel()

		bunch := test.mig.createUniqueString("bunch")
		desc := test.mig.createUniqueString("desc")
		id := test.mig.createSeedingBunch(func(fields map[string]interface{}) {
			fields["name"] = bunch
			fields["desc"] = desc
		})

		err := test.bst.Update(storage.UpdateBunch{ID: id, Name: bunch + "updated", Desc: desc})
		require.Nil(t, err)
	})

	t.Run("fail_update_a_duplicated_bunch", func(t *testing.T) {
		t.Parallel()

		bunch := test.mig.createUniqueString("bunch")
		desc := test.mig.createUniqueString("desc")
		test.mig.createSeedingBunch(func(fields map[string]interface{}) {
			fields["name"] = bunch
			fields["desc"] = desc
		})
		id := test.mig.createSeedingBunch(nil)

		err := test.bst.Update(storage.UpdateBunch{ID: id, Name: bunch, Desc: desc})
		require.NotNil(t, err)
	})
}

func TestBunchMysqlStorer_Get(t *testing.T) {
	t.Parallel()

	t.Run("success_get_a_bunch_by_id", func(t *testing.T) {
		t.Parallel()

		name := test.mig.createUniqueString("bunch")

		id := test.mig.createSeedingBunch(func(fields map[string]interface{}) {
			fields["name"] = name
		})

		bunch, err := test.bst.Get(id)
		require.Nil(t, err)
		require.NotNil(t, bunch)
		require.Equal(t, id, bunch.ID)
		require.Equal(t, name, bunch.Name)
	})
}

func TestBunchMysqlStorer_GetByName(t *testing.T) {
	t.Parallel()

	t.Run("success_get_a_bunch_by_name", func(t *testing.T) {
		t.Parallel()

		name := test.mig.createUniqueString("bunch")

		id := test.mig.createSeedingBunch(func(fields map[string]interface{}) {
			fields["name"] = name
		})

		bunch, err := test.bst.GetByName(name)
		require.Nil(t, err)
		require.NotNil(t, bunch)
		require.Equal(t, id, bunch.ID)
	})
}

func TestBunchMysqlStorer_Query(t *testing.T) {
	t.Parallel()

	t.Run("success_query_bunches", func(t *testing.T) {
		t.Parallel()

		prefix := test.mig.createUniqueString("prefix")
		name1 := test.mig.createUniqueString(prefix)
		name2 := test.mig.createUniqueString(prefix)
		name3 := test.mig.createUniqueString(prefix)
		name4 := test.mig.createUniqueString(prefix)
		name5 := test.mig.createUniqueString(prefix)
		name6 := test.mig.createUniqueString(prefix)
		name7 := test.mig.createUniqueString(prefix)

		test.mig.createSeedingBunch(func(fields map[string]interface{}) { fields["name"] = name1 })
		test.mig.createSeedingBunch(func(fields map[string]interface{}) { fields["name"] = name2 })
		test.mig.createSeedingBunch(func(fields map[string]interface{}) { fields["name"] = name3 })
		test.mig.createSeedingBunch(func(fields map[string]interface{}) { fields["name"] = name4 })
		test.mig.createSeedingBunch(func(fields map[string]interface{}) { fields["name"] = name5 })
		test.mig.createSeedingBunch(func(fields map[string]interface{}) { fields["name"] = name6 })
		test.mig.createSeedingBunch(func(fields map[string]interface{}) {
			fields["name"] = name7
			fields["active"] = false
		})

		bunches, total, err := test.bst.Query(storage.QueryBunch{
			Limit:  2,
			Offset: 2,
			Name:   prefix,
			Active: share.Boolean{IsSet: true, Bool: true},
		}, storage.SortBunch{})
		require.Nil(t, err)
		require.NotNil(t, bunches)
		require.Equal(t, int64(6), total)
		require.Len(t, bunches, 2)
	})
}

func TestBunchKeyMysqlStorer_Insert(t *testing.T) {
	t.Parallel()

	t.Run("success_insert_a_bunch_key", func(t *testing.T) {
		t.Parallel()

		bunchID := test.mig.createSeedingBunch(nil)
		keyID := test.mig.createSeedingServiceKey(nil)

		id, err := test.bkst.Insert(storage.BunchKey{BunchID: bunchID, KeyID: keyID})
		require.Nil(t, err)
		require.NotZero(t, id)
	})

	t.Run("fail_insert_a_bunch_key", func(t *testing.T) {
		t.Parallel()

		id, err := test.bkst.Insert(storage.BunchKey{BunchID: -1, KeyID: -2})
		require.NotNil(t, err)
		require.Zero(t, id)
	})
}

func TestBunchKeyMysqlStorer_Delete(t *testing.T) {
	t.Parallel()

	t.Run("success_delete_a_bunch_key", func(t *testing.T) {
		t.Parallel()

		bunchID := test.mig.createSeedingBunch(nil)
		keyID := test.mig.createSeedingServiceKey(nil)

		id, err := test.bkst.Insert(storage.BunchKey{BunchID: bunchID, KeyID: keyID})
		require.Nil(t, err)
		require.NotZero(t, id)

		err = test.bkst.Delete(id)
		require.Nil(t, err)
	})
}

func TestBunchKeyMysqlStorer_Query(t *testing.T) {
	t.Parallel()

	t.Run("success_query_bunch_keys", func(t *testing.T) {
		t.Parallel()
		prefix := test.mig.createUniqueString("pre")
		name1 := test.mig.createUniqueString(prefix)
		name2 := test.mig.createUniqueString(prefix)
		name3 := test.mig.createUniqueString(prefix)
		name4 := test.mig.createUniqueString(prefix)
		name5 := test.mig.createUniqueString(prefix)
		name6 := test.mig.createUniqueString(prefix)

		bunchID1 := test.mig.createSeedingBunch(func(fields map[string]interface{}) { fields["name"] = name1 })

		keyID1 := test.mig.createSeedingServiceKey(func(fields map[string]interface{}) { fields["name"] = name3 })
		keyID2 := test.mig.createSeedingServiceKey(func(fields map[string]interface{}) { fields["name"] = name4 })
		keyID3 := test.mig.createSeedingServiceKey(func(fields map[string]interface{}) { fields["name"] = name5 })
		keyID4 := test.mig.createSeedingServiceKey(func(fields map[string]interface{}) { fields["name"] = name6 })
		keyID5 := test.mig.createSeedingServiceKey(func(fields map[string]interface{}) { fields["name"] = name2 })

		_, err := test.bkst.Insert(storage.BunchKey{BunchID: bunchID1, KeyID: keyID1})
		require.Nil(t, err)

		_, err = test.bkst.Insert(storage.BunchKey{BunchID: bunchID1, KeyID: keyID2})
		require.Nil(t, err)

		_, err = test.bkst.Insert(storage.BunchKey{BunchID: bunchID1, KeyID: keyID3})
		require.Nil(t, err)

		_, err = test.bkst.Insert(storage.BunchKey{BunchID: bunchID1, KeyID: keyID4})
		require.Nil(t, err)

		_, err = test.bkst.Insert(storage.BunchKey{BunchID: bunchID1, KeyID: keyID5})
		require.Nil(t, err)

		rows, total, err := test.bkst.Query(storage.QueryBunchKey{
			Limit:     2,
			Offset:    1,
			BunchName: name1,
		}, storage.SortBunchKey{})
		require.Nil(t, err)
		require.NotNil(t, rows)
		require.Equal(t, int64(5), total)
		require.Len(t, rows, 2)
	})
}
