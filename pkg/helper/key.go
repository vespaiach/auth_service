package helper

import (
	"errors"
	"time"

	"github.com/vespaiach/auth_service/pkg/share"
	"github.com/vespaiach/auth_service/pkg/storage"
)

//KeyModel model
type KeyModel struct {
	ID        int64
	Name      string
	Desc      string
	UpdatedAt time.Time
}

//KeyQueryInput query a list of key
type KeyQueryInput struct {
	Limit  int64
	Offset int64
	Name   string
	Desc   string
	From   time.Time
	To     time.Time
}

//KeySortInput sort a list of key
type KeySortInput struct {
	Name      share.Direction
	Desc      share.Direction
	UpdatedAt share.Direction
}

//KeyHelper define all helper functions to interact with key
type KeyHelper interface {
	CreateKey(name string, desc string) (*KeyModel, error)
	UpdateKey(id int64, name string, desc string) (*KeyModel, error)
	DeleteKey(name string) error
	GetKey(name string) (*KeyModel, error)
	ListKey(query KeyQueryInput, sort KeySortInput) ([]*KeyModel, error)
}

//KeyService will implement KeyHelper interface
type KeyService struct {
	keyStorer storage.KeyStorer
}

//NewKeyHelper create an instance of KeyService
func NewKeyHelper(storer storage.KeyStorer) *KeyService {
	return &KeyService{storer}
}

//CreateKey creates a new key
func (s *KeyService) CreateKey(name string, desc string) (*KeyModel, error) {
	id, err := s.keyStorer.Insert(storage.CreateKey{Name: name, Desc: desc})
	if err != nil {
		return nil, err
	}
	if id == 0 {
		return nil, errors.New("cannot create a new key")
	}

	key, err := s.keyStorer.Get(id)
	if err != nil {
		return nil, err
	}
	return &KeyModel{key.ID, key.Name, key.Desc, key.UpdatedAt}, nil
}

//UpdateKey updates a key
func (s *KeyService) UpdateKey(id int64, name string, desc string) (*KeyModel, error) {
	err := s.keyStorer.Update(storage.UpdateKey{ID: id, Name: name, Desc: desc})
	if err != nil {
		return nil, err
	}

	key, err := s.keyStorer.Get(id)
	if err != nil {
		return nil, err
	}
	return &KeyModel{key.ID, key.Name, key.Desc, key.UpdatedAt}, nil
}

//DeleteKey deletes a key
func (s *KeyService) DeleteKey(name string) error {
	k, err := s.keyStorer.GetByName(name)
	if err != nil {
		return err
	}
	if k == nil {
		return nil
	}

	return s.keyStorer.Delete(k.ID)
}

//GetKey gets a key
func (s *KeyService) GetKey(name string) (*KeyModel, error) {
	k, err := s.keyStorer.GetByName(name)
	if err != nil {
		return nil, err
	}

	return &KeyModel{
		ID:        k.ID,
		Name:      k.Name,
		Desc:      k.Desc,
		UpdatedAt: k.UpdatedAt,
	}, nil
}

//ListKey list all keys
func (s *KeyService) ListKey(query KeyQueryInput, sort KeySortInput) ([]*KeyModel, int64, error) {
	keys, total, err := s.keyStorer.Query(storage.QueryKey{
		Limit:  query.Limit,
		Offset: query.Offset,
		Name:   query.Name,
		Desc:   query.Desc,
		From:   query.From,
		To:     query.To,
	}, storage.SortKey{
		Name:      sort.Name,
		Desc:      sort.Desc,
		UpdatedAt: sort.UpdatedAt,
	})
	if err != nil {
		return nil, 0, err
	}
	if keys == nil {
		return nil, 0, nil
	}

	result := make([]*KeyModel, 0, len(keys))
	for _, k := range keys {
		result = append(result, &KeyModel{
			ID:        k.ID,
			Name:      k.Name,
			Desc:      k.Desc,
			UpdatedAt: k.UpdatedAt,
		})
	}

	return result, total, nil
}
