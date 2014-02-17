
/* red-bazaar/http-proxy
 * backend.go
 * mae 02014-02
 */
package main

import (
	"fmt"
	"time"
	"net/url"
	"net/http"
	"net/http/httputil"
)

var (
	client = &http.Client{}
)


type BackEnd struct {

	Weight int

	RemoteURL *url.URL
	*httputil.ReverseProxy

	Running bool

	quitWatching chan bool
}

func (be *BackEnd) isOnline() bool {

	if be.ReverseProxy == nil {
		return false
	}

	return be.Running
}

func (be *BackEnd) Check() (bool,error) {

	req,err := http.NewRequest("HEAD",be.RemoteURL.String(),nil)
	if err != nil {
		return false,err
	}

	resp,err := client.Do(req)
	if err != nil {
		return false,err
	}
	resp.Body.Close()
	return true,nil
}

func (be *BackEnd) Watch(checkEvery time.Duration) {

	go func() {

		quit := be.quitWatching
		running := false

		for {
			if quit == nil {
				return
			}
			
			current := running
			
			running,err := be.Check()
			if err != nil {
				/* do nothing */
			}
			
			if running != current { /* something has changed */

				if !current { /* then backend has dropped out or is congested */
					
				
	
				} else { /* backend has come back online */

				}
			}

			be.Running = running

			select {
			case <- time.After(checkEvery):
				break
			case <- quit:
				return
			}
		}
	}()
}

func (be *BackEnd) Stop() {

	if be.quitWatching != nil {
		close(be.quitWatching)
	}
}

func (be *BackEnd) ServeHTTP(w http.ResponseWriter,req *http.Request) {

	if !be.isOnline() {
		
		http.Error(w,"Service Unavailable",503)
		return
	}

	be.Weight += 1	
	w.Header().Set("X-BackEnd-Weight",fmt.Sprintf("%d",be.Weight))
	be.ReverseProxy.ServeHTTP(w,req)
}


func NewBackEnd(remoteUrl string) (*BackEnd,error) {

	remote,err := url.Parse(remoteUrl)
	if err != nil {
		return nil,err
	}

	be := new(BackEnd)
	be.Weight = 5
	be.RemoteURL = remote
	be.ReverseProxy = httputil.NewSingleHostReverseProxy(remote)
	be.Running = false

	be.quitWatching = make(chan bool,1)

	return be,nil
}

type ByWeight []*BackEnd

func (be ByWeight) Len() int {
    return len(be)
}
func (be ByWeight) Less(i, j int) bool {
    return be[i].Weight < be[j].Weight
}
func (be ByWeight) Swap(i, j int) {
    be[i], be[j] = be[j], be[i]
}
