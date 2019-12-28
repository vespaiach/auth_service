package helper

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vespaiach/auth_service/pkg/share"
)

func TestCreateKey(t *testing.T) {
	t.Parallel()

	t.Run("create_key", func(t *testing.T) {
		t.Parallel()

		name := fmt.Sprintf("key_name_%d", inc.new())
		key, err := test.keyService.CreateKey(name, "just test one")
		require.Nil(t, err)
		require.Equal(t, name, key.Name)
		require.Equal(t, "just test one", key.Desc)
	})

	t.Run("miss_name", func(t *testing.T) {
		t.Parallel()

		key, err := test.keyService.CreateKey("", "just test one")
		require.NotNil(t, err)
		require.Nil(t, key)
	})

	t.Run("duplicated_name", func(t *testing.T) {
		t.Parallel()

		name := fmt.Sprintf("key_name_%d", inc.new())
		key, err := test.keyService.CreateKey(name, "just test one")
		require.Nil(t, err)
		require.Equal(t, name, key.Name)
		require.Equal(t, "just test one", key.Desc)

		duplicated, err := test.keyService.CreateKey(name, "just test one")
		require.NotNil(t, err)
		require.Nil(t, duplicated)

	})
}

func TestUpdateKey(t *testing.T) {
	t.Parallel()

	t.Run("update_name", func(t *testing.T) {
		t.Parallel()

		name := fmt.Sprintf("key_name_%d", inc.new())
		key, err := test.keyService.CreateKey(name, "just test one")
		require.Nil(t, err)
		require.Equal(t, name, key.Name)
		require.Equal(t, "just test one", key.Desc)

		updatedName := fmt.Sprintf("key_name_%d", inc.new())
		updated, err := test.keyService.UpdateKey(key.ID, updatedName, "")
		require.Nil(t, err)
		require.Equal(t, updatedName, updated.Name)
		require.Equal(t, "just test one", updated.Desc)
	})

	t.Run("update_desc", func(t *testing.T) {
		t.Parallel()

		name := fmt.Sprintf("key_name_%d", inc.new())
		key, err := test.keyService.CreateKey(name, "just test one")
		require.Nil(t, err)
		require.Equal(t, name, key.Name)
		require.Equal(t, "just test one", key.Desc)

		updated, err := test.keyService.UpdateKey(key.ID, "", "just test one updated")
		require.Nil(t, err)
		require.Equal(t, name, updated.Name)
		require.Equal(t, "just test one updated", updated.Desc)
	})

	t.Run("duplicated_name", func(t *testing.T) {
		t.Parallel()

		name := fmt.Sprintf("key_name_%d", inc.new())
		key, err := test.keyService.CreateKey(name, "just test one")
		require.Nil(t, err)
		require.Equal(t, name, key.Name)
		require.Equal(t, "just test one", key.Desc)

		updated, err := test.keyService.UpdateKey(key.ID, name, "just test one updated")
		require.NotNil(t, err)
		require.Nil(t, updated)
	})
}

func TestDeleteKey(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("key_name_%d", inc.new())
	key, err := test.keyService.CreateKey(name, "just test one")
	require.Nil(t, err)
	require.Equal(t, name, key.Name)
	require.Equal(t, "just test one", key.Desc)

	err = test.keyService.DeleteKey(name)
	require.Nil(t, err)

	notExistingName := fmt.Sprintf("key_name_%d", inc.new())
	err = test.keyService.DeleteKey(notExistingName)
	require.Nil(t, err)
}

func TestGetKey(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("key_name_%d", inc.new())
	key, err := test.keyService.CreateKey(name, "just test one")
	require.Nil(t, err)
	require.Equal(t, name, key.Name)
	require.Equal(t, "just test one", key.Desc)

	t.Run("existing_key", func(t *testing.T) {
		t.Parallel()

		key, err := test.keyService.GetKey(name)
		require.Nil(t, err)
		require.Equal(t, name, key.Name)
		require.Equal(t, "just test one", key.Desc)
	})

	t.Run("not_existing_key", func(t *testing.T) {
		t.Parallel()

		notExistingName := fmt.Sprintf("key_name_%d", inc.new())
		key, err := test.keyService.GetKey(notExistingName)
		require.Nil(t, err)
		require.Nil(t, key)
	})
}

func TestListKey(t *testing.T) {
	t.Parallel()

	for i := 0; i < 20; i++ {
		name := fmt.Sprintf("list_key_name_%d", inc.new())
		key, err := test.keyService.CreateKey(name, fmt.Sprintf("just test one - %d", i))
		require.Nil(t, err)
		require.Equal(t, name, key.Name)
		require.Equal(t, fmt.Sprintf("just test one - %d", i), key.Desc)
	}

	t.Run("sort_by_name", func(t *testing.T) {
		t.Parallel()

		keys, count, err := test.keyService.ListKey(KeyQueryInput{
			Limit:  10,
			Offset: 0,
			Name:   "list_key",
		}, KeySortInput{
			Name: share.Descendant,
		})
		require.Nil(t, err)
		require.Equal(t, 20, count)
		require.Len(t, keys, 10)
	})

	t.Run("sort_by_name", func(t *testing.T) {
		t.Parallel()

		keys, count, err := test.keyService.ListKey(KeyQueryInput{
			Limit:  10,
			Offset: 0,
			Name:   "list_key",
		}, KeySortInput{
			Desc: share.Ascendant,
		})
		require.Nil(t, err)
		require.Equal(t, 20, count)
		require.Len(t, keys, 10)
	})
}
