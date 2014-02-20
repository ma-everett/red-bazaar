/* red-bazaar/cache
 * service.go
 * mae 02014-02
 */

package cache

import (
	"sync"
	"time"
)

var (
	once = new(sync.Once)
	queue = make(chan *service,1)
)

type service struct {

	*Client
}

func start_service() {

	s := new(service)
	s.Client = NewClient()
	s.Client.Dial("loopback:7070") /* FIXME, this should be configurable */

	queue <- s
}

func Add(key string,content []byte,ttl time.Duration) error {

	once.Do(start_service)
	
	s := <- queue

	err := s.Client.Add(key,content,ttl)
	
	queue <- s

	return err
}

func Set(key string,content []byte,ttl time.Duration) error {

	once.Do(start_service)

	s := <- queue

	err := s.Client.Set(key,content,ttl)
	
	queue <- s

	return err
}

func Get(key string) ([]byte,bool,error) {

	once.Do(start_service)

	s := <- queue

	content,ok,err := s.Client.Get(key)

	queue <- s

	return content,ok,err
}

func Remove(key string) error {

	once.Do(start_service)

	s := <- queue

	err := s.Client.Remove(key)

	queue <- s

	return err
}
