package mysql

import (
	"fmt"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/vespaiach/auth_service/pkg/share"
	"github.com/vespaiach/auth_service/pkg/storage"
)

// UserMysqlStorage implements db's storage for user
type UserMysqlStorage struct {
	db *sqlx.DB
}

// UserBunchMysqlStorage implements db's storage for user
type UserBunchMysqlStorage struct {
	db *sqlx.DB
}

// NewUserMysqlStorage create new instance of UserMysqlStorage
func NewUserMysqlStorage(db *sqlx.DB) *UserMysqlStorage {
	return &UserMysqlStorage{
		db,
	}
}

// UserBunchMysqlStorage create new instance of UserBunchMysqlStorage
func NewUserBunchMysqlStorage(db *sqlx.DB) *UserBunchMysqlStorage {
	return &UserBunchMysqlStorage{
		db,
	}
}

func (st *UserMysqlStorage) Insert(u storage.CreateUser) (int64, error) {
	sql := "INSERT INTO users(full_name, `username`, `email`, `hash`, `salt`, updated_at) " +
		"VALUES(?, ?, ?, ?, ?, ?);"

	stmt, err := st.db.Prepare(sql)
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(u.FullName, u.Username, u.Email, u.Hash, u.Salt, time.Now())
	if err != nil {
		return 0, err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastID, nil
}

func (st *UserMysqlStorage) Update(u storage.UpdateUser) error {
	var (
		sql       = "UPDATE `users` SET %s	WHERE id = :id;"
		condition string
		prefix    string
	)

	updating := make(map[string]interface{})

	if len(u.FullName) > 0 {
		updating["full_name"] = u.FullName
		condition += prefix + "`full_name` = :full_name"
		prefix = ", "
	}

	if len(u.Username) > 0 {
		updating["username"] = u.Username
		condition += prefix + "`username` = :username"
		prefix = ", "
	}

	if len(u.Email) > 0 {
		updating["email"] = u.Email
		condition += prefix + "`email` = :email"
		prefix = ", "
	}

	if len(u.Hash) > 0 {
		updating["hash"] = u.Hash
		condition += prefix + "`hash` = :hash"
		prefix = ", "
	}

	if len(u.Salt) > 0 {
		updating["salt"] = u.Salt
		condition += prefix + "`salt` = :salt"
		prefix = ", "
	}

	if u.Active.IsSet {
		updating["active"] = u.Active.Bool
		condition += prefix + "`active` = :active"
		prefix = ", "
	}

	if len(updating) > 0 {
		updating["id"] = u.ID
		updating["updated_at"] = time.Now()
		condition += prefix + "`updated_at` = :updated_at"

		_, err := st.db.NamedExec(fmt.Sprintf(sql, condition), updating)
		if err != nil {
			return err
		}
	}

	return nil
}

func (st *UserMysqlStorage) Get(id int64) (*storage.User, error) {
	sql := "SELECT id, full_name, `username`, `email`, `hash`, `salt`, `active`, updated_at FROM `users` " +
		"WHERE `id` = ? LIMIT 1;"

	rows, err := st.db.Queryx(sql, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	u := &storage.User{Active: share.Boolean{IsSet: true}}
	err = rows.Scan(&u.ID, &u.FullName, &u.Username, &u.Email, &u.Hash, &u.Salt, &u.Active.Bool, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (st *UserMysqlStorage) GetByName(username string) (*storage.User, error) {
	sql := "SELECT id, full_name, `username`, `email`, `hash`, `salt`, `active`, updated_at FROM `users` " +
		"WHERE `username` = ? LIMIT 1;"

	rows, err := st.db.Queryx(sql, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	u := &storage.User{Active: share.Boolean{IsSet: true}}
	err = rows.Scan(&u.ID, &u.FullName, &u.Username, &u.Email, &u.Hash, &u.Salt, &u.Active.Bool, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (st *UserMysqlStorage) GetByEmail(email string) (*storage.User, error) {
	sql := "SELECT id, full_name, `username`, `email`, `hash`, `salt`, `active`, updated_at FROM `users` " +
		"WHERE `email` = ? LIMIT 1;"

	rows, err := st.db.Queryx(sql, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	u := &storage.User{Active: share.Boolean{IsSet: true}}
	err = rows.Scan(&u.ID, &u.FullName, &u.Username, &u.Email, &u.Hash, &u.Salt, &u.Active.Bool, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (st *UserMysqlStorage) Query(queries storage.QueryUser, sorts storage.SortUser) ([]*storage.User, int64, error) {
	var (
		sql           = "SELECT id, full_name, `username`, `email`, `hash`, `salt`, active, updated_at FROM `users` %s ORDER BY %s LIMIT :offset, :limit;"
		sqlcount      = "SELECT count(id) FROM `users` %s;"
		orderPrefix   string
		order         string
		wherePrefix   = "WHERE "
		where         string
		wg            sync.WaitGroup
		queryErr      error
		countTotalErr error
		results       []*storage.User
		total         int64
	)

	filter := map[string]interface{}{"limit": queries.Limit, "offset": queries.Offset}
	if queries.Limit == 0 {
		filter["queries"] = share.DefaultLimit
	}

	if len(queries.FullName) > 0 {
		filter["full_name"] = "%" + queries.FullName + "%"
		where += wherePrefix + "`full_name` LIKE :full_name"
		wherePrefix = " AND "
	}

	if len(queries.Username) > 0 {
		filter["username"] = "%" + queries.Username + "%"
		where += wherePrefix + "`username` LIKE :username"
		wherePrefix = " AND "
	}

	if len(queries.Email) > 0 {
		filter["desc"] = "%" + queries.Email + "%"
		where += wherePrefix + "`email` LIKE :email"
		wherePrefix = " AND "
	}

	if queries.Active.IsSet {
		filter["active"] = queries.Active.Bool
		where += wherePrefix + "`active` = :active"
		wherePrefix = " AND "
	}

	if !queries.From.IsZero() {
		filter["from"] = queries.From
		where += wherePrefix + "updated_at > :from"
		wherePrefix = " AND "
	}

	if !queries.To.IsZero() {
		filter["to"] = queries.To
		where += wherePrefix + "updated_at <= :to"
		wherePrefix = " AND "
	}

	if sorts.Username != share.BiDirection {
		order += orderPrefix + fmt.Sprintf("`username` %s", getOrderDirection(sorts.Username))
		orderPrefix = " , "
	}

	if sorts.FullName != share.BiDirection {
		order += orderPrefix + fmt.Sprintf("`full_name` %s", getOrderDirection(sorts.FullName))
		orderPrefix = " , "
	}

	if sorts.Email != share.BiDirection {
		order += orderPrefix + fmt.Sprintf("`email` %s", getOrderDirection(sorts.Email))
		orderPrefix = " , "
	}

	if sorts.UpdatedAt != share.BiDirection {
		order += orderPrefix + fmt.Sprintf("`updated_at` %s", getOrderDirection(sorts.UpdatedAt))
		orderPrefix = " , "
	}

	if sorts.Active != share.BiDirection {
		order += orderPrefix + fmt.Sprintf("`active` %s", getOrderDirection(sorts.Active))
		orderPrefix = " , "
	}

	if len(order) == 0 {
		order = "`id` DESC"
	}

	sql = fmt.Sprintf(sql, where, order)
	sqlcount = fmt.Sprintf(sqlcount, where)

	wg.Add(1)
	go func() {
		defer wg.Done()
		rows, err := st.db.NamedQuery(sql, filter)
		if err != nil {
			queryErr = err
			return
		}
		defer rows.Close()

		results = make([]*storage.User, 0, queries.Limit)
		for rows.Next() {
			u := &storage.User{Active: share.Boolean{IsSet: true}}
			err := rows.Scan(&u.ID, &u.FullName, &u.Username, &u.Email, &u.Hash, &u.Salt, &u.Active.Bool, &u.UpdatedAt)
			if err != nil {
				queryErr = err
				return
			}
			results = append(results, u)
		}

		if rows.Err() != nil {
			queryErr = rows.Err()
			return
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		rows, err := st.db.NamedQuery(sqlcount, filter)
		if err != nil {
			countTotalErr = err
			return
		}
		defer rows.Close()

		if rows.Next() {
			err := rows.Scan(&total)
			if err != nil {
				countTotalErr = err
				return
			}
		}
	}()

	wg.Wait()

	if queryErr != nil {
		return nil, 0, queryErr
	}
	if countTotalErr != nil {
		return nil, 0, countTotalErr
	}

	return results, total, nil
}

func (st *UserBunchMysqlStorage) Insert(u storage.CreateUserBunch) (int64, error) {
	sql := "INSERT INTO user_bunches (user_id, bunch_id, updated_at) VALUES(?, ?, ?);"

	stmt, err := st.db.Prepare(sql)
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(u.UserID, u.BunchID, time.Now())
	if err != nil {
		return 0, err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastID, nil
}

func (st *UserBunchMysqlStorage) Delete(id int64) error {
	sql := "DELETE FROM `user_bunches` WHERE id=?"

	stmt, err := st.db.Prepare(sql)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}

	return nil
}

func (st *UserBunchMysqlStorage) Query(queries storage.QueryUserBunch, sorts storage.SortUserBunch) ([]*storage.AggregateUserBunch, int64, error) {
	var (
		sql = "SELECT `users`.id, `users`.full_name, `users`.`username`, `users`.`email`, `users`.`hash`, " +
			"`users`.`salt`, `users`.`active`, `users`.updated_at, bunches.`id`, bunches.`name`, bunches.`desc`, " +
			"bunches.`active`, bunches.updated_at, user_bunches.`id`, user_bunches.user_id, user_bunches.bunch_id, " +
			"user_bunches.updated_at FROM `users` " +
			"INNER JOIN user_bunches ON `users`.id = user_bunches.user_id " +
			"INNER JOIN bunches ON user_bunches.bunch_id = bunches.`id` " +
			"%s ORDER BY %s LIMIT :offset, :limit;"
		sqlcount = "SELECT count(user_bunches.`id`) FROM `users` " +
			"INNER JOIN user_bunches ON `users`.id = user_bunches.user_id " +
			"INNER JOIN bunches ON user_bunches.bunch_id = bunches.`id` %s;"
		orderPrefix   string
		order         string
		wherePrefix   = "WHERE "
		where         string
		wg            sync.WaitGroup
		queryErr      error
		countTotalErr error
		results       []*storage.AggregateUserBunch
		total         int64
	)

	filter := map[string]interface{}{"limit": queries.Limit, "offset": queries.Offset}
	if queries.Limit == 0 {
		filter["queries"] = share.DefaultLimit
	}

	if len(queries.Username) > 0 {
		filter["username"] = "%" + queries.Username + "%"
		where += wherePrefix + "`users`.`username` LIKE :username"
		wherePrefix = " AND "
	}

	if len(queries.BunchName) > 0 {
		filter["name"] = "%" + queries.BunchName + "%"
		where += wherePrefix + "bunches.`name` LIKE :name"
		wherePrefix = " AND "
	}

	if queries.UserActive.IsSet {
		filter["user_active"] = queries.UserActive.Bool
		where += wherePrefix + "`users`.`active` = :user_active"
		wherePrefix = " AND "
	}

	if queries.BunchActive.IsSet {
		filter["user_active"] = queries.BunchActive.Bool
		where += wherePrefix + "bunches.`active` = :user_active"
		wherePrefix = " AND "
	}

	if sorts.Username != share.BiDirection {
		order += orderPrefix + fmt.Sprintf("`users`.`username` %s", getOrderDirection(sorts.Username))
		orderPrefix = " , "
	}

	if sorts.BunchName != share.BiDirection {
		order += orderPrefix + fmt.Sprintf("bunches.`name` %s", getOrderDirection(sorts.BunchName))
		orderPrefix = " , "
	}

	if len(order) == 0 {
		order = "user_bunches.`id` DESC"
	}

	sql = fmt.Sprintf(sql, where, order)
	sqlcount = fmt.Sprintf(sqlcount, where)

	wg.Add(1)
	go func() {
		defer wg.Done()
		rows, err := st.db.NamedQuery(sql, filter)
		if err != nil {
			queryErr = err
			return
		}
		defer rows.Close()

		results = make([]*storage.AggregateUserBunch, 0, queries.Limit)
		for rows.Next() {
			u := &storage.User{Active: share.Boolean{IsSet: true}}
			b := &storage.Bunch{Active: share.Boolean{IsSet: true}}
			ub := &storage.UserBunch{}

			err := rows.Scan(&u.ID, &u.FullName, &u.Username, &u.Email, &u.Hash, &u.Salt, &u.Active.Bool, &u.UpdatedAt,
				&b.ID, &b.Name, &b.Desc, &b.Active.Bool, &u.UpdatedAt,
				&ub.ID, &ub.UserID, &ub.BunchID, &ub.UpdatedAt)
			if err != nil {
				queryErr = err
				return
			}
			results = append(results, &storage.AggregateUserBunch{
				User:      u,
				Bunch:     b,
				UserBunch: ub,
			})
		}

		if rows.Err() != nil {
			queryErr = rows.Err()
			return
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		rows, err := st.db.NamedQuery(sqlcount, filter)
		if err != nil {
			countTotalErr = err
			return
		}
		defer rows.Close()

		if rows.Next() {
			err := rows.Scan(&total)
			if err != nil {
				countTotalErr = err
				return
			}
		}
	}()

	wg.Wait()

	if queryErr != nil {
		return nil, 0, queryErr
	}
	if countTotalErr != nil {
		return nil, 0, countTotalErr
	}

	return results, total, nil
}
