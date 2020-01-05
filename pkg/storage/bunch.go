package storage

import (
	"time"

	"github.com/vespaiach/auth_service/pkg/share"
)

//Bunch model
type Bunch struct {
	ID        int64
	Name      string
	Desc      string
	Active    bool
	UpdatedAt time.Time
}

//CreateBunch model
type CreateBunch struct {
	Name string
	Desc string
}

//UpdateBunch model
type UpdateBunch struct {
	ID     int64
	Name   string
	Desc   string
	Active share.Boolean
}

//QueryBunch model
type QueryBunch struct {
	Limit  int64
	Offset int64
	Name   string
	Desc   string
	Active share.Boolean
	From   time.Time
	To     time.Time
}

//SortBunch model
type SortBunch struct {
	Name      share.Direction
	Desc      share.Direction
	Active    share.Direction
	UpdatedAt share.Direction
}

//BunchKey model
type BunchKey struct {
	ID        int64
	BunchID   int64
	KeyID     int64
	UpdatedAt time.Time
}

//QueryBunchKey model
type QueryBunchKey struct {
	Limit       int64
	Offset      int64
	BunchName   string
	KeyName     string
	BunchActive share.Boolean
}

//SortBunchKey model
type SortBunchKey struct {
	BunchName   share.Direction
	KeyName     share.Direction
	BunchActive share.Direction
}

//AggregateBunchKey model
type AggregateBunchKey struct {
	*BunchKey
	*Key
	*Bunch
}

//BunchStorer defines fundamental functions to interact with storage repository
type BunchStorer interface {
	Insert(b CreateBunch) (*Bunch, error)
	Update(b UpdateBunch) (*Bunch, error)
	Get(id int64) (*Bunch, error)
	GetByName(name string) (*Bunch, error)
	Query(queries QueryBunch, sorts SortBunch) ([]*Bunch, int64, error)
}

//BunchKeyStorer defines fundamental functions to interact with storage repository
type BunchKeyStorer interface {
	Insert(bk BunchKey) (*BunchKey, error)
	Delete(id int64) error
	Query(queries QueryBunchKey, sorts SortBunchKey) ([]*AggregateBunchKey, error)
}
