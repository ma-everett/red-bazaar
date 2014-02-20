
/* red-bazaar/cache 
 * simple HTTP cache
 * mae 02014-02
 */

package main

import (
	"flag"
	"log"
	"time"
	"net"
	"github.com/pmylund/go-cache"

	cachep ".." /* include cache protocol */
)

type Context struct {

	*cache.Cache
}

/* FIXME: 
 * todo, the cache needs to be written to file on shutdown
 *       then read from file on startup
 *       which will need syscall support, remember to to setuuid as well 
 */

func main() { 
	/* a basic TCP wrapper around an inmemory cache */
	
	addr := flag.String("a","loopback:7070","TCP address")
	
	flag.Parse()

	/* create context with cache */
	
	ctx := new(Context)
	c := cache.New(5 * time.Minute,30 * time.Second) 
	ctx.Cache = c

	/* create TCP listen */

	log.Printf("serving on %s\n",*addr)

	naddr,err := net.ResolveTCPAddr("tcp",*addr)
	if err != nil {

		log.Fatalf("unable to resolve TCP addr %s - %v\n",*addr,err)
	}

	ln, err := net.ListenTCP(naddr.Network(), naddr)
	if err != nil {
		
		log.Fatalf("unable to listen to TCP addr %s - %v\n",*addr,err)
	}

	defer ln.Close()

	queue := make(chan *Context,1)
	queue <- ctx

	for {
		conn, err := ln.AcceptTCP()
		if err != nil {
			
			log.Printf("connection error - %v\n",err)
			continue
		}
		go handle(conn,queue)
	}
}

type Serve struct {

	q chan *Context
}

func (s *Serve) Add(key string,data []byte,ttl time.Duration) error {
	
	log.Printf("--> add %s - expires in %v\n",key,ttl)

	ctx := <- s.q
	
	err := ctx.Cache.Add(key,data,ttl)
	
	s.q <- ctx

	return err
}

func (s *Serve) Set(key string,data []byte,ttl time.Duration) error {

	go func() {

		log.Printf("--> set %s - expires in %v\n",key,ttl)

		ctx := <- s.q

		ctx.Cache.Set(key,data,ttl)

		s.q <- ctx
	}()

	return nil
}

func (s *Serve) Get(key string) ([]byte,bool) {

	log.Printf("--> get %s\n",key)

	ctx := <- s.q

	x,found := ctx.Cache.Get(key)

	s.q <- ctx

	if !found {
		return nil,false
	}

	if d,ok := x.([]byte); ok {

		return d,true
	}

	return nil,false
}

func (s *Serve) Remove(key string) error {

	go func() {

		log.Printf("--> remove %s\n",key)
		
		ctx := <- s.q
		
		ctx.Cache.Delete(key)
		
		s.q <- ctx
	}()

	return nil
}

func NewServe(queue chan *Context) *Serve {

	s := new(Serve)
	s.q = queue
	return s
}

func handle(conn *net.TCPConn,queue chan *Context) {

	s,err := cachep.NewServer(NewServe(queue),conn)
	if err != nil {

		log.Printf("new server error - %v\n",err)
		return
	}

	if err := s.Handle(); err != nil {
		log.Printf("handle client - %v\n",err) /* FIXME, add remote IP */
	}
}




