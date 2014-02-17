
/* red-bazaar/http-proxy
 * context.go
 * mae 02014-02
 */
package main

import (
	//"io"
	//"log"
	//"fmt"
	//"strings"
	//"time"
	//"crypto/md5"
	"net/http"	


)


type Context struct {
	
	Hosts map[string]*Host

}

func (ctx *Context) Handler() func(http.ResponseWriter,*http.Request) {

	return func(w http.ResponseWriter, req *http.Request) {

		w.Header().Set("X-Proxy","http-proxy")

		if host,exists := ctx.Hosts[req.Host]; exists {
			
			host.ServeHTTP(w,req)
			return
		}

		/*
		userAgent := req.Header.Get("User-Agent")
		
	
		host := ""

	
		parts := strings.Split(req.Host,":")
		if len(parts) != 2 {

			host = req.Host
		} else {

			host = parts[0]
		}
	
		h := md5.New()
		io.WriteString(h, userAgent)
		
	
		hash := fmt.Sprintf("%x", h.Sum(nil))[:8] + "+" + host	
		
		
		if x,found := ctx.Cache.Get(hash); found {

			
			if be,ok := x.(*BackEnd); ok {

			
				if be.isOnline() {

					be.ServeHTTP(w,req)
					return

				}  else {

					ctx.Cache.Delete(hash) 
				}
			}	
		} 

		log.Printf("looking up - %s\n",req.Host)
		
		if h,exists := ctx.Hosts[req.Host]; exists {

			h.ServeHTTP(w,req)
			return
		}
                */

		http.Error(w,"Service Not Found",404)
	}
}

func (ctx *Context) NewHost(cname string) (*Host,error) {

	if host,exists := ctx.Hosts[cname]; exists {

		return host,HostAlreadyPresent
	}

	h := NewHost(cname)
	ctx.Hosts[cname] = h

	return h,nil
}

func NewContext() *Context {

	ctx := new(Context)

	ctx.Hosts = make(map[string]*Host,0)
	

	return ctx
}
	
