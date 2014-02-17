
/* red-bazaar/http-proxy
 * host.go
 * mae 02014-02
 */
package main

import (
	"sort"
	"strings"
	"io"
	"fmt"
	"log"
	"time"
	"math/rand"
	"crypto/md5"
	"errors"
	"net/http"	
	"github.com/pmylund/go-cache"
)

var (
	HostAlreadyPresent = errors.New("Host Already Present")
	BackEndAlreadyPresent = errors.New("BackEnd Already Present")
	NoBackEndsOnline = errors.New("No BackEnds online")
)

type Host struct {

	CName string
	BackEnds map[string]*BackEnd
	
	XHeaders []XHeader
	
	/* add configuration here */

	/* add cache ?! */
	*cache.Cache
}

/* Weighted Random */
func (host *Host) WeightedBackEnd() (*BackEnd,error) {

	if len(host.BackEnds) == 0 {
		return nil,NoBackEndsOnline
	}
	
	totalWeight := 1 /* must be > 0 */
	running := make([]*BackEnd,0)

	for _,be := range host.BackEnds {

		if be.isOnline() {
			running = append(running,be)
			totalWeight += be.Weight
		}
	}

	if len(running) == 0 {
		return nil,NoBackEndsOnline
	}

	if len(running) == 1 { /* if we only have the 1 choice, then skip */
		return running[0],nil
	}

	sort.Sort(ByWeight(running))
	
	n := rand.Intn(totalWeight)
	
	for i := len(running) - 1; i >= 0; i-- {
		
		be := running[i]
		n -= be.Weight	

		if n < be.Weight {

			return be,nil
		}
	}	

	return nil,NoBackEndsOnline
}

func (host *Host) SelectBackEnd() (*BackEnd,error) {

	if len(host.BackEnds) == 0 {
		return nil,NoBackEndsOnline
	}

	running := make([]*BackEnd,0)

	for _,be := range host.BackEnds {

		if be.isOnline() {
			running = append(running,be)
		}
	}

	if len(running) == 0 {
		return nil,NoBackEndsOnline
	}
	
	if len(running) == 1 {
		return running[0],nil
	}

	sort.Sort(ByWeight(running))

	return running[0],nil
}

func (host *Host) ServeHTTP(w http.ResponseWriter,req *http.Request) {

	w.Header().Set("X-Host",host.CName)

	/* read the host, locate a backend
         * check for session in cache based on
         * ip-address and client agent 
         */
	userAgent := req.Header.Get("User-Agent")
	
	/* hash the user-agent string into a more compact form */
	remotehost := ""
	
	/* strip the port from the host string */
	parts := strings.Split(req.Host,":")
	if len(parts) != 2 {
		
		remotehost = req.Host
	} else {
		
		remotehost = parts[0]
	}
	
	h := md5.New()
	io.WriteString(h, userAgent)
	
	/* append host ip address to hash of user-agent */
	hash := fmt.Sprintf("%x", h.Sum(nil))[:8] + "+" + remotehost	
	
	/* lookup a session in the cache */
	if x,found := host.Cache.Get(hash); found {
		
		w.Header().Set("X-Session","continue")

		/* continue the session with the last backend */
		if session,ok := x.(*Session); ok {
			
			session.Set(w)
			be,err := host.SelectBackEnd()
			if err != nil {
				http.Error(w,"Service Not Available",503)
				return
			}
			
			w.Header().Set("X-Served-By",be.RemoteURL.String())
			be.ServeHTTP(w,req)
			return
		}	
	} /* else lookup a back end to try : */
	
	log.Printf("looking up - %s\n",req.Host)
	
	session := NewSession()

	for _,xhead := range host.XHeaders {
		
		session.Add(xhead.Key(),xhead.Set(w))
	}

	host.Cache.Set(hash,session,0)

	be,err := host.SelectBackEnd()
	if err != nil {
		http.Error(w,"Service Not Available",503)
		return
	}

	w.Header().Set("X-Served-By",be.RemoteURL.String())
	be.ServeHTTP(w,req)
}

func (host *Host) AddXHeader(name,value string) error {

	xhead,err := ParseXHeader(name,value)
	if err != nil {
		return err
	}

	host.XHeaders = append(host.XHeaders,xhead)
	return nil
}

func (host *Host) AddBackEnd(url string) (*BackEnd,error) {

	if be,exists := host.BackEnds[url]; exists {

		return be,BackEndAlreadyPresent
	}

	be,err := NewBackEnd(url)
	if err != nil {
		return nil,err
	}

	host.BackEnds[url] = be
	return be,nil
}

func NewHost(cname string) *Host {

	h := new(Host)
	h.CName = cname
	h.BackEnds = make(map[string]*BackEnd,0)
	h.XHeaders = make([]XHeader,0)

	h.Cache = cache.New(5 * time.Minute,30 * time.Second) /* TODO, add to create */

	return h
}
