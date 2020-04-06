package leveldb

import (
	"fmt"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

// Manager struct
type Manager struct {
	db *leveldb.DB
}

// Item struct
type Item struct {
	Key []byte
	Value []byte
}

// IManager - Manager interface
type IManager interface {
	Close()
	GetAllItems() (items []Item, err error)
	AddNewItem(item Item) (status int, err error)
	FindItemByKey(key []byte) (item *Item, err error) 
}

// Error - custom leveldb error struct
type Error struct {
	Key []byte
	err error
}

// NewErr - create new custom error instance
func NewErr(key []byte, err error) error {
	return &Error{
		Key: key,
		err: err,
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("Error at Key: %s | Message: %s", string(e.Key), e.err.Error())
}

// Close database connection,
// (*Manager).db cannot be used after this function called.
func (mgr *Manager) Close() {
	mgr.db.Close()
}

// FindItemByKey - Return Item object with key and value.
func (mgr *Manager) FindItemByKey(key []byte) (item *Item, err error) {
	var value []byte
	value, err = mgr.db.Get(key, nil)
	if err != nil {
		err = NewErr(key, err)
		return nil, err
	}
	item = &Item {
		Key: key,
		Value: value,
	}
	return item, nil
}

// AddNewItem add new item into database
func (mgr *Manager) AddNewItem(item Item) (status int, err error) {
	err =  mgr.db.Put(item.Key, item.Value, nil)
	if err != nil {
		return 0, NewErr(item.Key, err)
	}
	return 1, nil
}

// GetAllItems - Get all items in database
func (mgr *Manager) GetAllItems() (listItems []Item, err error) {
	iter := mgr.db.NewIterator(nil, nil)
	for iter.Next() {
		listItems = append(listItems, Item{
			Key: iter.Key(),
			Value: iter.Value(),
		})
	}
	iter.Release()
	err = iter.Error() 
	return listItems, err 
}

// New - create new leveldb database manager instance
func New(path string, options *opt.Options) (manager *Manager, err error) {
	db, err := leveldb.OpenFile(path, options)
	if err != nil {
		return nil, err
	}
	mgr := &Manager {
		db: db,
	}
	return mgr, nil 
}