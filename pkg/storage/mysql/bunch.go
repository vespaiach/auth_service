package mysql

import (
	"fmt"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/vespaiach/auth_service/pkg/share"
	"github.com/vespaiach/auth_service/pkg/storage"
)

// BunchMysqlStorer implements db's storage for bunch
type BunchMysqlStorer struct {
	db *sqlx.DB
}

// BunchKeyMysqlStorer implements db's storage for bunch-key
type BunchKeyMysqlStorer struct {
	db *sqlx.DB
}

// NewBunchMysqlStorer create new instance of BunchMysqlStorer
func NewBunchMysqlStorer(db *sqlx.DB) *BunchMysqlStorer {
	return &BunchMysqlStorer{
		db,
	}
}

// BunchKeyMysqlStorer create new instance of BunchKeyMysqlStorer
func NewBunchKeyMysqlStorer(db *sqlx.DB) *BunchKeyMysqlStorer {
	return &BunchKeyMysqlStorer{
		db,
	}
}

func (st *BunchMysqlStorer) Insert(u storage.CreateBunch) (int64, error) {
	sql := "INSERT INTO bunches (`name`, `desc`, `active`, updated_at) VALUES (?, ?, ?, ?);"

	stmt, err := st.db.Prepare(sql)
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(u.Name, u.Desc, true, time.Now())
	if err != nil {
		return 0, err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastID, nil
}

func (st *BunchMysqlStorer) Update(u storage.UpdateBunch) error {
	var (
		sql      string = "UPDATE bunches SET %s WHERE id = :id;"
		fields   string
		prefix   string
		updating = make(map[string]interface{})
	)

	if len(u.Name) > 0 {
		fields += prefix + " `name` = :name "
		prefix = ","
		updating["name"] = u.Name
	}

	if len(u.Desc) > 0 {
		fields += prefix + " `desc` = :desc "
		updating["desc"] = u.Desc
	}

	if u.Active.IsSet {
		fields += prefix + " `active` = :active "
		updating["active"] = u.Active.Bool
	}

	if len(updating) > 0 {
		fields += prefix + " updated_at = :updated_at "
		updating["updated_at"] = time.Now()
		updating["id"] = u.ID

		_, err := st.db.NamedExec(fmt.Sprintf(sql, fields), updating)
		if err != nil {
			return err
		}
	}

	return nil
}

func (st *BunchMysqlStorer) Get(id int64) (*storage.Bunch, error) {
	sql := "SELECT id, `name`, `desc`, active, updated_at FROM `bunches` WHERE id = ? LIMIT 1;"
	rows, err := st.db.Queryx(sql, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	b := &storage.Bunch{Active: share.Boolean{IsSet: true}}
	if err := rows.Scan(&b.ID, &b.Name, &b.Desc, &b.Active.Bool, &b.UpdatedAt); err != nil {
		return nil, err
	}

	return b, nil
}

func (st *BunchMysqlStorer) GetByName(name string) (*storage.Bunch, error) {
	sql := "SELECT id, `name`, `desc`, active, updated_at FROM `bunches` WHERE name = ? LIMIT 1;"
	rows, err := st.db.Queryx(sql, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	b := &storage.Bunch{Active: share.Boolean{IsSet: true}}
	if err := rows.Scan(&b.ID, &b.Name, &b.Desc, &b.Active.Bool, &b.UpdatedAt); err != nil {
		return nil, err
	}

	return b, nil
}

func (st *BunchMysqlStorer) Query(queries storage.QueryBunch, sorts storage.SortBunch) ([]*storage.Bunch, int64, error) {
	var (
		sql           = "SELECT id, `name`, `desc`, active, updated_at FROM `bunches` %s ORDER BY %s LIMIT :offset, :limit;"
		sqlcount      = "SELECT count(id) FROM `bunches` %s;"
		orderPrefix   string
		order         string
		wherePrefix   = "WHERE "
		where         string
		wg            sync.WaitGroup
		queryErr      error
		countTotalErr error
		results       []*storage.Bunch
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

		results = make([]*storage.Bunch, 0, queries.Limit)
		for rows.Next() {
			b := &storage.Bunch{Active: share.Boolean{IsSet: true}}
			err := rows.Scan(&b.ID, &b.Name, &b.Desc, &b.Active.Bool, &b.UpdatedAt)
			if err != nil {
				queryErr = err
				return
			}
			results = append(results, b)
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

func (st *BunchKeyMysqlStorer) Insert(bk storage.BunchKey) (int64, error) {
	sql := "INSERT INTO `bunch_keys` (bunch_id, key_id, updated_at) VALUES (?, ?, ?);"

	stmt, err := st.db.Prepare(sql)
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(bk.BunchID, bk.KeyID, time.Now())
	if err != nil {
		return 0, err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastID, nil
}

func (st *BunchKeyMysqlStorer) Delete(id int64) error {
	sql := "DELETE FROM `bunch_keys` WHERE id=?"

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

func (st *BunchKeyMysqlStorer) Query(queries storage.QueryBunchKey, sorts storage.SortBunchKey) ([]*storage.AggregateBunchKey, int64, error) {
	var (
		sql = "SELECT `keys`.id, `keys`.`name`, `keys`.`desc`, `keys`.updated_at, " +
			"bunches.`id`, bunches.`name`, bunches.`desc`, bunches.`active`, bunches.updated_at, " +
			"bunch_keys.`id`, bunch_keys.bunch_id, bunch_keys.key_id, bunch_keys.updated_at " +
			"FROM bunch_keys " +
			"INNER JOIN `keys` ON `keys`.id = bunch_keys.key_id " +
			"INNER JOIN `bunches` ON `bunches`.id = bunch_keys.bunch_id " +
			"%s ORDER BY %s LIMIT :offset, :limit;"
		sqlcount = "SELECT count(bunch_keys.`id`) " +
			"FROM bunch_keys " +
			"INNER JOIN `keys` ON `keys`.`id` = bunch_keys.key_id " +
			"INNER JOIN `bunches` ON `bunches`.id = bunch_keys.bunch_id %s"
		orderPrefix   string
		order         string
		wherePrefix   = "WHERE "
		where         string
		wg            sync.WaitGroup
		queryErr      error
		countTotalErr error
		results       []*storage.AggregateBunchKey
		total         int64
	)

	filter := map[string]interface{}{"limit": queries.Limit, "offset": queries.Offset}
	if queries.Limit == 0 {
		filter["queries"] = share.DefaultLimit
	}

	if len(queries.BunchName) > 0 {
		filter["bunch_name"] = queries.BunchName
		where += wherePrefix + "bunches.`name` = :bunch_name"
		wherePrefix = " AND "
	}

	if len(queries.KeyName) > 0 {
		filter["key_name"] = queries.BunchName
		where += wherePrefix + "`keys`.`name` = :key_name"
		wherePrefix = " AND "
	}

	if queries.BunchActive.IsSet {
		filter["active"] = queries.BunchActive.Bool
		where += wherePrefix + "bunches.`active` = :active"
		wherePrefix = " AND "
	}

	if sorts.BunchName != share.BiDirection {
		order += orderPrefix + fmt.Sprintf("bunches.`name` %s", getOrderDirection(sorts.BunchName))
		orderPrefix = " , "
	}

	if sorts.KeyName != share.BiDirection {
		order += orderPrefix + fmt.Sprintf("`keys`.`name` %s", getOrderDirection(sorts.KeyName))
		orderPrefix = " , "
	}

	if sorts.BunchActive != share.BiDirection {
		order += orderPrefix + fmt.Sprintf("bunches.`active` %s", getOrderDirection(sorts.BunchActive))
		orderPrefix = " , "
	}

	if len(order) == 0 {
		order = "bunch_keys.`id` DESC"
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

		results = make([]*storage.AggregateBunchKey, 0, queries.Limit)

		for rows.Next() {
			k := new(storage.Key)
			b := &storage.Bunch{Active: share.Boolean{IsSet: true}}
			bk := new(storage.BunchKey)

			err := rows.Scan(&k.ID, &k.Name, &k.Desc, &k.UpdatedAt,
				&b.ID, &b.Name, &b.Desc, &b.Active.Bool, &b.UpdatedAt,
				&bk.ID, &bk.BunchID, &bk.KeyID, &bk.UpdatedAt)
			if err != nil {
				queryErr = err
				return
			}
			results = append(results, &storage.AggregateBunchKey{bk, k, b})
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
