package mysql

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vespaiach/auth_service/pkg/share"
	"github.com/vespaiach/auth_service/pkg/storage"
)

func TestUserMysqlStorage_Insert(t *testing.T) {
	t.Parallel()

	t.Run("success_add_a_user", func(t *testing.T) {
		t.Parallel()

		name := test.mig.createUniqueString("username")
		email := test.mig.createUniqueString("email")

		id, err := test.ust.Insert(storage.CreateUser{
			FullName: "full name",
			Username: name,
			Email:    email,
			Hash:     "hash",
			Salt:     "salt",
		})
		require.Nil(t, err)
		require.NotZero(t, id)
	})

	t.Run("success_add_a_duplicated_username", func(t *testing.T) {
		t.Parallel()

		name := test.mig.createUniqueString("username")
		email := test.mig.createUniqueString("email")
		id := test.mig.createSeedingUser(func(fields map[string]interface{}) {
			fields["username"] = name
		})

		id, err := test.ust.Insert(storage.CreateUser{
			FullName: "full name",
			Username: name,
			Email:    email,
			Hash:     "hash",
			Salt:     "salt",
		})
		require.NotNil(t, err)
		require.Zero(t, id)
	})

	t.Run("success_add_a_duplicated_email", func(t *testing.T) {
		t.Parallel()

		name := test.mig.createUniqueString("username")
		email := test.mig.createUniqueString("email")
		id := test.mig.createSeedingUser(func(fields map[string]interface{}) {
			fields["email"] = email
		})

		id, err := test.ust.Insert(storage.CreateUser{
			FullName: "full name",
			Username: name,
			Email:    email,
			Hash:     "hash",
			Salt:     "salt",
		})
		require.NotNil(t, err)
		require.Zero(t, id)
	})
}

func TestUserMysqlStorage_Update(t *testing.T) {
	t.Parallel()

	t.Run("success_update_a_user", func(t *testing.T) {
		t.Parallel()

		id := test.mig.createSeedingUser(nil)
		newname := test.mig.createUniqueString("username")
		newemail := test.mig.createUniqueString("email")

		err := test.ust.Update(storage.UpdateUser{
			ID:       id,
			FullName: "full name",
			Username: newname,
			Email:    newemail,
			Hash:     "hash_updated",
			Salt:     "2",
			Active:   share.Boolean{IsSet: true, Bool: false},
		})
		require.Nil(t, err)

		name, email, hash, active := test.mig.getUserByID(id)
		require.Equal(t, name, newname)
		require.Equal(t, email, newemail)
		require.Equal(t, hash, "hash_updated")
		require.False(t, active)
	})
}

func TestUserMysqlStorage_Get(t *testing.T) {
	t.Parallel()

	t.Run("success_get_a_user_by_id", func(t *testing.T) {
		t.Parallel()

		id := test.mig.createSeedingUser(nil)

		user, err := test.ust.Get(id)
		require.Nil(t, err)
		require.NotNil(t, user)
	})
}

func TestUserMysqlStorage_GetByName(t *testing.T) {
	t.Parallel()

	t.Run("success_get_a_user_by_name", func(t *testing.T) {
		t.Parallel()

		name := test.mig.createUniqueString("name")
		id := test.mig.createSeedingUser(func(fields map[string]interface{}) {
			fields["username"] = name
		})

		user, err := test.ust.GetByName(name)
		require.Nil(t, err)
		require.NotNil(t, user)
		require.Equal(t, id, user.ID)
	})
}

func TestUserMysqlStorage_GetByEmail(t *testing.T) {
	t.Parallel()

	t.Run("success_get_a_user_by_email", func(t *testing.T) {
		t.Parallel()

		email := test.mig.createUniqueString("email")
		id := test.mig.createSeedingUser(func(fields map[string]interface{}) {
			fields["email"] = email
		})

		user, err := test.ust.GetByEmail(email)
		require.Nil(t, err)
		require.NotNil(t, user)
		require.Equal(t, id, user.ID)
	})
}

func TestUserMysqlStorage_Query(t *testing.T) {
	t.Parallel()

	t.Run("success_query_users", func(t *testing.T) {
		t.Parallel()

		name1 := test.mig.createUniqueString("user1ame")
		name2 := test.mig.createUniqueString("user1ame")
		name3 := test.mig.createUniqueString("user1ame")
		name4 := test.mig.createUniqueString("user1ame")
		name5 := test.mig.createUniqueString("user1ame")

		test.mig.createSeedingUser(func(field map[string]interface{}) { field["username"] = name1 })
		test.mig.createSeedingUser(func(field map[string]interface{}) { field["username"] = name2 })
		test.mig.createSeedingUser(func(field map[string]interface{}) { field["username"] = name3 })
		test.mig.createSeedingUser(func(field map[string]interface{}) { field["username"] = name4 })
		test.mig.createSeedingUser(func(field map[string]interface{}) {
			field["username"] = name5
			field["active"] = false
		})

		users, total, err := test.ust.Query(storage.QueryUser{
			Limit:    2,
			Offset:   2,
			Username: "user1ame",
			Active:   share.Boolean{IsSet: true, Bool: true},
		}, storage.SortUser{
			FullName:  share.Ascendant,
			UpdatedAt: share.Descendant,
			Email:     share.Descendant,
		})
		require.Nil(t, err)
		require.NotNil(t, users)
		require.Equal(t, int64(4), total)
		require.Len(t, users, 2)
	})
}

func TestUserBunchMysqlStorage_Insert(t *testing.T) {
	t.Parallel()

	t.Run("success_add_a_user_bunch", func(t *testing.T) {
		t.Parallel()

		userID := test.mig.createSeedingUser(nil)
		buncheID := test.mig.createSeedingBunch(nil)

		id, err := test.ubst.Insert(storage.CreateUserBunch{
			UserID:  userID,
			BunchID: buncheID,
		})
		require.Nil(t, err)
		require.NotZero(t, id)
	})

	t.Run("fail_add_a_user_bunch", func(t *testing.T) {
		t.Parallel()

		id, err := test.ubst.Insert(storage.CreateUserBunch{
			UserID:  -1,
			BunchID: -10,
		})
		require.NotNil(t, err)
		require.Zero(t, id)
	})
}

func TestUserBunchMysqlStorage_Delete(t *testing.T) {
	t.Parallel()

	t.Run("success_delete_a_user_bunch", func(t *testing.T) {
		t.Parallel()

		userID := test.mig.createSeedingUser(nil)
		buncheID := test.mig.createSeedingBunch(nil)

		id, err := test.ubst.Insert(storage.CreateUserBunch{
			UserID:  userID,
			BunchID: buncheID,
		})
		require.Nil(t, err)
		require.NotZero(t, id)

		err = test.ubst.Delete(id)
		require.Nil(t, err)
	})
}

func TestUserBunchMysqlStorage_Query(t *testing.T) {
	t.Parallel()

	t.Run("success_query_user_bunches", func(t *testing.T) {
		t.Parallel()

		name1 := test.mig.createUniqueString("pre1")
		name2 := test.mig.createUniqueString("pre1")
		name3 := test.mig.createUniqueString("pre1")
		name4 := test.mig.createUniqueString("pre1")

		userID := test.mig.createSeedingUser(func(fields map[string]interface{}) { fields["username"] = name1 })

		bunchID1 := test.mig.createSeedingBunch(func(fields map[string]interface{}) { fields["name"] = name2 })
		bunchID2 := test.mig.createSeedingBunch(func(fields map[string]interface{}) { fields["name"] = name3 })
		bunchID3 := test.mig.createSeedingBunch(func(fields map[string]interface{}) {
			fields["name"] = name4
			fields["active"] = false
		})

		_, err := test.ubst.Insert(storage.CreateUserBunch{
			UserID:  userID,
			BunchID: bunchID1,
		})
		require.Nil(t, err)

		_, err = test.ubst.Insert(storage.CreateUserBunch{
			UserID:  userID,
			BunchID: bunchID2,
		})
		require.Nil(t, err)

		_, err = test.ubst.Insert(storage.CreateUserBunch{
			UserID:  userID,
			BunchID: bunchID3,
		})
		require.Nil(t, err)

		rows, total, err := test.ubst.Query(storage.QueryUserBunch{
			Limit:       2,
			Offset:      0,
			Username:    name1,
			BunchActive: share.Boolean{IsSet: true, Bool: true},
		}, storage.SortUserBunch{
			Username:  share.Ascendant,
			BunchName: share.Descendant,
		})
		require.Nil(t, err)
		require.Equal(t, int64(2), total)
		require.NotNil(t, rows)
		require.Len(t, rows, 2)
	})
}
