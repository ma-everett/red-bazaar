
/* red-bazaar/cache
 * protocol.go
 * mae 02014-02
 */

package cache

import (
	"errors"
	"time"
)

var (
	NotImplemented = errors.New("Not Implemented")
	OperationInvalid = errors.New("Operation Invalid")
	DataInvalid = errors.New("Data Invalid")

	NoConnection = errors.New("No Connection")
	ClientDropped = errors.New("Client Dropped")

	CacheMiss = errors.New("CacheMiss")
)

const (
	TAddToCache      = byte(1)
	RAddToCache      = byte(2)
	TSetToCache      = byte(3)
	RSetToCache      = byte(4)
	TGetFromCache    = byte(5)
	RGetFromCache    = byte(6)
	TRemoveFromCache = byte(7)
	RRemoveFromCache = byte(8)
)


type CacheServer interface {

	Add(key string,value []byte,ttl time.Duration) error
	Set(key string,value []byte,ttl time.Duration) error
	Get(key string) ([]byte,bool)
	Remove(key string) error
}



