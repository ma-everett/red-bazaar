
/* red-bazaar/http-proxy
 * context.go
 * mae 02014-02
 */
package main

import (
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
	
