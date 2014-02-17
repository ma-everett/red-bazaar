
/* red-bazaar/http-proxy
 * session.go
 * mae 02014-02
 */

package main

import (
	"net/http"
)

type Session struct {

	XHeaders map[string]string
}

func (s *Session) Add(key,value string) {

	s.XHeaders[key] = value
}

func (s *Session) Set(w http.ResponseWriter) {

	for key,value := range s.XHeaders {

		w.Header().Set(key,value)
	}
}

func NewSession() *Session {

	s := new(Session)
	s.XHeaders = make(map[string]string,0)
	return s
}
