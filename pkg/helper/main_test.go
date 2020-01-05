package helper

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/vespaiach/auth_service/pkg/share"
	"github.com/vespaiach/auth_service/pkg/storage"
)

var ctx context.Context

type uniqInt struct {
	order int64
	mux   sync.Mutex
}

var inc *uniqInt = &uniqInt{order: 1}

func (unq *uniqInt) new() int64 {
	unq.mux.Lock()
	defer unq.mux.Unlock()
	unq.order++
	return unq.order
}

type testHelpers struct {
	keyService *KeyService
}

var test *testHelpers

func TestMain(m *testing.M) {
	test = &testHelpers{
		keyService: NewKeyHelper(&KeyStorerStub{Keys: make([]*storage.Key, 0)}),
	}
}

//-----------------------------------------------------------------------------

type KeyStorerStub struct {
	Keys   []*storage.Key
	sortby func(key1, key2 *storage.Key) bool
}

func (sort *KeyStorerStub) Len() int {
	return len(sort.Keys)
}

func (sort *KeyStorerStub) Swap(i, j int) {
	sort.Keys[i], sort.Keys[j] = sort.Keys[j], sort.Keys[i]
}

func (sort *KeyStorerStub) Less(i, j int) bool {
	return sort.sortby(sort.Keys[i], sort.Keys[j])
}

func (stub *KeyStorerStub) Insert(k storage.CreateKey) (int64, error) {
	id := inc.new()
	stub.Keys = append(stub.Keys, &storage.Key{
		ID:        id,
		Name:      k.Name,
		Desc:      k.Desc,
		UpdatedAt: time.Now(),
	})
	return id, nil
}

func (stub *KeyStorerStub) Update(k storage.UpdateKey) error {
	for _, key := range stub.Keys {
		if key.ID == k.ID {

			if len(k.Name) > 0 {
				key.Name = k.Name
			}

			if len(k.Desc) > 0 {
				key.Desc = k.Desc
			}

			return nil
		}
	}

	return fmt.Errorf("cannot find key with id: %d", k.ID)
}

func (stub *KeyStorerStub) Delete(id int64) error {
	i := -1

	for j, k := range stub.Keys {
		if k.ID == id {
			i = j
			break
		}
	}

	if i > 0 {
		stub.Keys = append(stub.Keys[:i], stub.Keys[i+1:]...)
		return nil
	}

	return fmt.Errorf("cannot find key with id: %d", id)
}

func (stub *KeyStorerStub) Get(id int64) (*storage.Key, error) {
	for _, k := range stub.Keys {
		if k.ID == id {
			return k, nil
		}
	}

	return nil, nil
}

func (stub *KeyStorerStub) GetByName(name string) (*storage.Key, error) {
	for _, k := range stub.Keys {
		if k.Name == name {
			return k, nil
		}
	}

	return nil, nil
}

func (stub *KeyStorerStub) Query(queries storage.QueryKey, sorts storage.SortKey) ([]*storage.Key, int64, error) {
	results := make([]*storage.Key, 0)
	for _, k := range stub.Keys {
		if (len(queries.Name) == 0 || strings.Contains(k.Name, queries.Name)) &&
			(len(queries.Desc) == 0 || strings.Contains(k.Desc, queries.Desc)) &&
			(queries.From.IsZero() || k.UpdatedAt.After(queries.From)) &&
			(queries.To.IsZero() || k.UpdatedAt.Before(queries.To)) {
			results = append(results, k)
		}
	}

	//sort by name
	stub.sortby = func(key1, key2 *storage.Key) bool {
		switch sorts.Name {
		case share.Ascendant:
			return key1.Name < key2.Name
		case share.Descendant:
			return key1.Name > key2.Name
		default:
			return false
		}
	}
	sort.Sort(stub)

	//sort by desc
	stub.sortby = func(key1, key2 *storage.Key) bool {
		switch sorts.Desc {
		case share.Ascendant:
			return key1.Desc < key2.Desc
		case share.Descendant:
			return key1.Desc > key2.Desc
		default:
			return false
		}
	}
	sort.Sort(stub)

	//sort by updatedat
	stub.sortby = func(key1, key2 *storage.Key) bool {
		switch sorts.UpdatedAt {
		case share.Ascendant:
			return key1.UpdatedAt.Before(key2.UpdatedAt)
		case share.Descendant:
			return key1.UpdatedAt.After(key2.UpdatedAt)
		default:
			return false
		}
	}
	sort.Sort(stub)

	return results[queries.Offset*queries.Limit : queries.Offset*queries.Limit+queries.Limit], int64(len(results)), nil
}

//-----------------------------------------------------------------------------

type bunchStorerStub struct {
	bunches []*storage.Bunch
	sortby  func(bunch1, bunch2 *storage.Bunch) bool
}

func (stub *bunchStorerStub) Len() int {
	return len(stub.bunches)
}

func (stub *bunchStorerStub) Swap(i, j int) {
	stub.bunches[i], stub.bunches[j] = stub.bunches[j], stub.bunches[i]
}

func (stub *bunchStorerStub) Less(i, j int) bool {
	return stub.sortby(stub.bunches[i], stub.bunches[j])
}

func (stub *bunchStorerStub) Insert(k storage.CreateBunch) (int64, error) {
	id := inc.new()
	stub.bunches = append(stub.bunches, &storage.Bunch{
		ID:        id,
		Name:      k.Name,
		Desc:      k.Desc,
		Active:    true,
		UpdatedAt: time.Now(),
	})
	return id, nil
}
