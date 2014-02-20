/* red-bazaar/http-proxy
 * headers.go
 * mae 02014-02
 */
package main

import (
	"math/rand"
	"errors"
	"strings"
	"net/http"
)

var (	
	UnknownXHeaderType = errors.New("Unknown X-Header Type")
)

type XHeader interface {

	Set(w http.ResponseWriter) string
	Key() string
}

type XHeaderStatic struct {

	name string
	value string
}

func (r XHeaderStatic) Set(w http.ResponseWriter) string {

	w.Header().Set(r.name,r.value)
	return r.value
}

func (r XHeaderStatic) Key() string {
	return r.name
}

type XHeaderRandom struct {

	name string
	choices []string
}

func (r XHeaderRandom) Set(w http.ResponseWriter) string {

	n := rand.Intn(len(r.choices) - 1)
	w.Header().Set(r.name,r.choices[n])
	return r.choices[n]
}

func (r XHeaderRandom) Key() string {
	return r.name
}

	

func ParseXHeader(name,value string) (XHeader,error) {

	switch value[0] {
	case '!': /* then static */
		
		return XHeaderStatic{name,value[1:]},nil
		break
	case '?': /* then is a random XHeader */
		
		/* should be a comma seperated list of options */
		options := strings.Split(value[1:],",")
		if len(options) == 1 { /* not very random, return a static XHeader instead */

			return XHeaderStatic{name,options[0]},nil
		}

		return XHeaderRandom{name,options},nil
		break	
	default:

		break
	}

	return nil,UnknownXHeaderType
}
