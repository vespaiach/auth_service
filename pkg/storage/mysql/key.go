package mysql

import (
	"fmt"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/vespaiach/auth_service/pkg/share"
	"github.com/vespaiach/auth_service/pkg/storage"
)

// KeyMysqlStorer implements key's storages in mysql db
type KeyMysqlStorer struct {
	db *sqlx.DB
}

// KeyMysqlStorer creates a new instance of KeyMysqlStorer
func NewKeyMysqlStorer(db *sqlx.DB) *KeyMysqlStorer {
	return &KeyMysqlStorer{
		db,
	}
}

func (st *KeyMysqlStorer) Insert(k storage.CreateKey) (int64, error) {
	sql := "INSERT INTO `keys` (`name`, `desc`, updated_at) VALUES (?, ?, ?);"

	stmt, err := st.db.Prepare(sql)
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(k.Name, k.Desc, time.Now())
	if err != nil {
		return 0, err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastID, nil
}

func (st *KeyMysqlStorer) Update(k storage.UpdateKey) error {
	var (
		sql      string = "UPDATE `keys` SET %s WHERE id = :id;"
		fields   string
		prefix   string
		updating = make(map[string]interface{})
	)

	if len(k.Name) > 0 {
		fields += prefix + " `name` = :name "
		prefix = ","
		updating["name"] = k.Name
	}
	if len(k.Desc) > 0 {
		fields += prefix + " `desc` = :desc "
		updating["desc"] = k.Desc
	}

	if len(updating) > 0 {
		fields += prefix + " updated_at = :updated_at "
		updating["updated_at"] = time.Now()
		updating["id"] = k.ID

		_, err := st.db.NamedExec(fmt.Sprintf(sql, fields), updating)
		if err != nil {
			return err
		}
	}

	return nil
}

func (st *KeyMysqlStorer) Delete(id int64) error {
	sqlbunchkey := "DELETE FROM `bunch_keys` WHERE key_id=:id"
	sqlkey := "DELETE FROM `keys` WHERE id=:id"
	deleting := map[string]interface{}{"id": id}

	tx, err := st.db.Beginx()
	if err != nil {
		return err
	}

	_, err = tx.NamedExec(sqlbunchkey, deleting)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.NamedExec(sqlkey, deleting)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func (st *KeyMysqlStorer) Get(id int64) (*storage.Key, error) {
	sql := "SELECT id, `name`, `desc`, updated_at FROM `keys` WHERE id = ? LIMIT 1;"

	rows, err := st.db.Queryx(sql, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	key := new(storage.Key)
	if err := rows.Scan(&key.ID, &key.Name, &key.Desc, &key.UpdatedAt); err != nil {
		return nil, err
	}

	return key, nil
}

func (st *KeyMysqlStorer) GetByName(name string) (*storage.Key, error) {
	sql := "SELECT id, `name`, `desc`, updated_at FROM `keys` WHERE `name` = ? LIMIT 1;"

	rows, err := st.db.Queryx(sql, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	key := new(storage.Key)
	if err := rows.Scan(&key.ID, &key.Name, &key.Desc, &key.UpdatedAt); err != nil {
		return nil, err
	}

	return key, nil
}

func (st *KeyMysqlStorer) Query(queries storage.QueryKey, sorts storage.SortKey) ([]*storage.Key, int64, error) {
	var (
		sql      string = "SELECT id, `name`, `desc`, updated_at FROM `keys` %s ORDER BY %s LIMIT :offset, :limit;"
		sqlcount string = "SELECT count(id) FROM `keys` %s;"

		orderPrefix   string
		order         string
		wherePrefix   = "WHERE "
		where         string
		wg            sync.WaitGroup
		queryErr      error
		countTotalErr error
		results       []*storage.Key
		total         int64
	)

	filter := map[string]interface{}{"limit": queries.Limit, "offset": queries.Offset}
	if queries.Limit == 0 {
		filter["queries"] = share.DefaultLimit
	}

	if len(queries.Name) > 0 {
		filter["name"] = "%" + queries.Name + "%"
		where += wherePrefix + "`name` LIKE :name"
		wherePrefix = " AND "
	}

	if len(queries.Desc) > 0 {
		filter["desc"] = "%" + queries.Desc + "%"
		where += wherePrefix + "`desc` LIKE :desc"
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

	if sorts.Name != share.BiDirection {
		order += orderPrefix + fmt.Sprintf("`name` %s", getOrderDirection(sorts.Name))
		orderPrefix = " , "
	}

	if sorts.Desc != share.BiDirection {
		order += orderPrefix + fmt.Sprintf("`desc` %s", getOrderDirection(sorts.Desc))
		orderPrefix = " , "
	}

	if sorts.UpdatedAt != share.BiDirection {
		order += orderPrefix + fmt.Sprintf("`updated_at` %s", getOrderDirection(sorts.UpdatedAt))
		orderPrefix = " , "
	}

	if len(order) == 0 {
		order = "id DESC"
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

		results = make([]*storage.Key, 0, queries.Limit)
		for rows.Next() {
			key := new(storage.Key)
			err := rows.Scan(&key.ID, &key.Name, &key.Desc, &key.UpdatedAt)
			if err != nil {
				queryErr = err
				return
			}
			results = append(results, key)
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
