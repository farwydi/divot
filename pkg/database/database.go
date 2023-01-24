package database

import "fmt"

var ErrKeyNotFound = fmt.Errorf("key not found")

type Database interface {
	Load(key []byte) ([]byte, error)
	Write(key []byte, value []byte) error
	Scan(prefix []byte, openFn func(key, value []byte) error) error
}
