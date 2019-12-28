package mysql

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

type testApp struct {
	mig  *Migrator
	kst  *KeyMysqlStorer
	bst  *BunchMysqlStorer
	bkst *BunchKeyMysqlStorer
	ust  *UserMysqlStorage
	ubst *UserBunchMysqlStorage
}

var test *testApp

// TestMain setup testing env for mysql repository
func TestMain(m *testing.M) {
	db, err := initDb()
	if err != nil {
		log.Fatal(err)
	}

	test = &testApp{
		mig:  NewMigrator(db),
		kst:  NewKeyMysqlStorer(db),
		bst:  NewBunchMysqlStorer(db),
		bkst: NewBunchKeyMysqlStorer(db),
		ust:  NewUserMysqlStorage(db),
		ubst: NewUserBunchMysqlStorage(db),
	}

	test.mig.Drop()
	test.mig.Init()

	code := m.Run()

	db.Close()
	os.Exit(code)
}

func initDb() (*sqlx.DB, error) {
	db, err := sqlx.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", "root", "password",
		"127.0.0.1", "3306", "auth", "charset=utf8&parseTime=True&loc=Local&multiStatements=True&maxAllowedPacket=0"))
	if err != nil {
		return nil, err
	}

	return db, nil
}
