package mysql

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
)

type uniqInt struct {
	order int
	mux   sync.Mutex
}

var inc *uniqInt = &uniqInt{order: 1}

func (unq *uniqInt) New() int {
	unq.mux.Lock()
	defer unq.mux.Unlock()
	unq.order++
	return unq.order
}

// createUniqueString is to create unique string for testing
func (m *Migrator) createUniqueString(prefix string) string {
	return fmt.Sprintf("%s%s", prefix, strconv.Itoa(inc.New()))
}

// Script migration script
type Script struct {
	Name string
	Text string
}

// Migrator struct
type Migrator struct {
	db   *sqlx.DB
	init []*Script
	drop []*Script
	seed []*Script
}

// NewMigrator return struct instance
func NewMigrator(db *sqlx.DB) *Migrator {

	var initScripts = []*Script{
		&Script{Name: "init_database", Text: initDatabase},
	}

	var dropScripts = []*Script{
		&Script{Name: "drop_database", Text: dropDatabase},
	}

	var seedScripts = []*Script{
		&Script{Name: "seed_database", Text: seedingData},
	}

	return &Migrator{
		db,
		initScripts,
		dropScripts,
		seedScripts,
	}
}

// Init database
func (m *Migrator) Init() {
	tx := m.db.MustBegin()

	for _, s := range m.init {
		tx.MustExec(santizeSQL(s.Text))
	}

	tx.Commit()
}

// Drop database
func (m *Migrator) Drop() {
	tx := m.db.MustBegin()

	for i := len(m.drop) - 1; i >= 0; i-- {
		tx.MustExec(santizeSQL(m.drop[i].Text))
	}

	tx.Commit()
}

// Seed database
func (m *Migrator) Seed() {
	tx := m.db.MustBegin()

	for i := len(m.seed) - 1; i >= 0; i-- {
		tx.MustExec(santizeSQL(m.seed[i].Text))
	}

	tx.Commit()
}

func santizeSQL(sql string) string {
	return strings.Replace(sql, `"`, "`", -1)
}

func (m *Migrator) createSeedingServiceKey(beforeCreate func(map[string]interface{})) int64 {
	fields := map[string]interface{}{
		"name":       fmt.Sprintf("name_%s", strconv.Itoa(inc.New())),
		"desc":       fmt.Sprintf("desc_%s", strconv.Itoa(inc.New())),
		"updated_at": time.Now(),
	}

	if beforeCreate != nil {
		beforeCreate(fields)
	}

	result, _ := m.db.NamedExec("INSERT INTO `keys` (`name`, `desc`, updated_at) VALUES (:name, :desc, :updated_at);", fields)
	id, _ := result.LastInsertId()

	return id
}

func (m *Migrator) createSeedingBunch(beforeCreate func(map[string]interface{})) int64 {
	fields := map[string]interface{}{
		"name":   fmt.Sprintf("name_%s", strconv.Itoa(inc.New())),
		"desc":   fmt.Sprintf("desc_%s", strconv.Itoa(inc.New())),
		"active": true,
	}

	if beforeCreate != nil {
		beforeCreate(fields)
	}

	result, _ := m.db.NamedExec("INSERT INTO `bunches` (`name`, `desc`, active) VALUES (:name, :desc, :active);",
		fields)
	id, _ := result.LastInsertId()

	return id
}

func (m *Migrator) createSeedingUser(beforeCreate func(map[string]interface{})) int64 {
	fields := map[string]interface{}{
		"full_name":  fmt.Sprintf("full_name_%s", strconv.Itoa(inc.New())),
		"username":   fmt.Sprintf("username_%s", strconv.Itoa(inc.New())),
		"email":      fmt.Sprintf("email_%s", strconv.Itoa(inc.New())),
		"hash":       fmt.Sprintf("hash_%s", strconv.Itoa(inc.New())),
		"salt":       fmt.Sprintf("salt_%s", strconv.Itoa(inc.New())),
		"active":     true,
		"updated_at": time.Now(),
	}

	if beforeCreate != nil {
		beforeCreate(fields)
	}

	result, _ := m.db.NamedExec("INSERT INTO users(full_name, `username`, `email`, `hash`, `salt`, `active`, updated_at) "+
		"VALUES(:full_name, :username, :email, :hash, :salt, :active, :updated_at);", fields)
	id, _ := result.LastInsertId()

	return id
}

func (m *Migrator) getServiceKeyByID(id int64) (key string, desc string) {
	rows, err := m.db.Queryx("SELECT `key`, `desc` FROM `keys` WHERE id = ?", id)
	defer rows.Close()

	if err == nil && rows.Next() {
		rows.Scan(&key, &desc)
	}

	return
}

func (m *Migrator) getBunchByID(id int64) (name string, desc string, active bool) {
	rows, err := m.db.Queryx("SELECT `name`, `desc`, `active` FROM bunches WHERE id = ?", id)
	defer rows.Close()

	if err == nil && rows.Next() {
		rows.Scan(&name, &desc, &active)
	}

	return
}

func (m *Migrator) getUserByID(id int64) (username string, email string, hash string, active bool) {
	rows, err := m.db.Queryx("Select `username`, `email`, `hash`, `active` FROM `users` WHERE id = ?", id)
	defer rows.Close()

	if err == nil && rows.Next() {
		rows.Scan(&username, &email, &hash, &active)
	}

	return
}

func (m *Migrator) getKeyIDByBunchID(id int64) []int64 {
	rows, _ := m.db.Queryx("SELECT key_id FROM `bunch_keys` WHERE bunch_id = ?", id)
	defer rows.Close()

	results := make([]int64, 0)

	for rows.Next() {
		var id int64
		rows.Scan(&id)
		results = append(results, id)
	}

	return results
}
