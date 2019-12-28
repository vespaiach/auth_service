package storage

import (
	"time"

	"github.com/vespaiach/auth_service/pkg/share"
)

//Key model
type Key struct {
	ID        int64
	Name      string
	Desc      string
	UpdatedAt time.Time
}

//CreateKey model
type CreateKey struct {
	Name string
	Desc string
}

//UpdateKey model
type UpdateKey struct {
	ID   int64
	Name string
	Desc string
}

//QueryKey model
type QueryKey struct {
	Limit  int64
	Offset int64
	Name   string
	Desc   string
	From   time.Time
	To     time.Time
}

//SortKey model
type SortKey struct {
	Name      share.Direction
	Desc      share.Direction
	UpdatedAt share.Direction
}

//KeyStorer defines fundamental functions to interact with storage repository
type KeyStorer interface {
	Insert(k CreateKey) (int64, error)
	Update(k UpdateKey) error
	Delete(id int64) error
	Get(id int64) (*Key, error)
	GetByName(name string) (*Key, error)
	Query(queries QueryKey, sorts SortKey) ([]*Key, error)
}
