package storage

import (
	"time"

	"github.com/vespaiach/auth_service/pkg/share"
)

//User model
type User struct {
	ID        int64
	FullName  string
	Username  string
	Email     string
	Hash      string
	Salt      string
	Active    share.Boolean
	UpdatedAt time.Time
}

//CreateUser model
type CreateUser struct {
	FullName string
	Username string
	Email    string
	Hash     string
	Salt     string
}

//UpdateUser model
type UpdateUser struct {
	ID       int64
	FullName string
	Username string
	Email    string
	Hash     string
	Salt     string
	Active   share.Boolean
}

//QueryUser model
type QueryUser struct {
	Limit    int64
	Offset   int64
	FullName string
	Username string
	Email    string
	Active   share.Boolean
	From     time.Time
	To       time.Time
}

//SortUser model
type SortUser struct {
	FullName  share.Direction
	Username  share.Direction
	Email     share.Direction
	Active    share.Direction
	UpdatedAt share.Direction
}

//UserBunch model
type UserBunch struct {
	ID        int64
	UserID    int64
	BunchID   int64
	UpdatedAt time.Time
}

//AggregateUserBunch model
type AggregateUserBunch struct {
	*User
	*Bunch
	*UserBunch
}

//CreateUserBunch model
type CreateUserBunch struct {
	UserID  int64
	BunchID int64
}

//QueryUserBunch model
type QueryUserBunch struct {
	Limit       int64
	Offset      int64
	Username    string
	BunchName   string
	UserActive  share.Boolean
	BunchActive share.Boolean
}

//SortUserBunch model
type SortUserBunch struct {
	Username  share.Direction
	BunchName share.Direction
}

//UserStorer defines fundamental functions to interact with storage repository
type UserStorer interface {
	Insert(u CreateUser) (*User, error)
	Update(u UpdateUser) (*User, error)
	Get(id int64) (*User, error)
	GetByName(username string) (*User, error)
	GetByEmail(email string) (*User, error)
	Query(queries QueryUser, sorts SortUser) ([]*User, int64, error)
}

//UserBunchStorer defines fundamental functions to interact with storage repository
type UserBunchStorer interface {
	Insert(ub CreateUserBunch) (*UserBunch, error)
	Delete(id int64) error
	Query(queries QueryUserBunch, sorts SortUserBunch) ([]*AggregateUserBunch, int64, error)
}
