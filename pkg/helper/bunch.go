package helper

import (
	"time"

	"github.com/vespaiach/auth_service/pkg/share"
)

//BunchModel model
type BunchModel struct {
	ID        int64
	Name      string
	Desc      string
	Active    bool
	UpdatedAt time.Time
}

//BunchQueryInput query a list of key
type BunchQueryInput struct {
	Limit  int64
	Offset int64
	Name   string
	Desc   string
	From   time.Time
	To     time.Time
}

//BunchSortInput sort a list of key
type BunchSortInput struct {
	Name      share.Direction
	Desc      share.Direction
	UpdatedAt share.Direction
}

//CreateBunchInput input
type CreateBunchInput struct {
	Name string
	Desc string
}

//UpdateBunchInput input
type UpdateBunchInput struct {
	ID     int64
	Name   string
	Active share.Boolean
	Desc   string
}

//BunchHelper define all helper functions to interact with bunch
type BunchHelper interface {
	CreateBunch(input CreateBunchInput) (*BunchModel, error)
	UpdateBunch(input UpdateBunchInput) (*BunchModel, error)
	DeleteBunch(name string) error
	GetBunch(name string) (*BunchModel, error)
	ListBunch(query BunchQueryInput, sort BunchSortInput) ([]*BunchModel, error)
}
