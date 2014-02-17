
/* red-bazaar/http-proxy 
 * reverse HTTP proxy
 * mae 02014-02
 */

package main

import (
	"log"
	"time"
	"net/http"
)

func main() {

	//err = syscall.Setgid(gid); 
        //err = syscall.Setuid(uid);

	/* create a context */
	ctx := NewContext() 


	/* test configuration */
	host,err := ctx.NewHost("test.intranet:8080") /* notice the port */
	if err != nil {
		
		log.Fatalf("error adding host test.intranet - %v\n",err)
	}

	host.AddXHeader("X-name","!foo")
	host.AddXHeader("X-test","?correct")
	host.AddXHeader("X-css","?set1,set2,set3")

	be,err := host.AddBackEnd("http://loopback:6060")
	if err != nil {

		log.Fatalf("unable to add backend - %v\n",err)
	}

	be.Watch(1 * time.Second)
	

	be,err = host.AddBackEnd("http://loopback:6061")
	if err != nil {

		log.Fatalf("unable to add backend - %v\n",err)
	}

	be.Watch(1 * time.Second)


	/* setup http reverse handler */
	http.HandleFunc("/",ctx.Handler())
	if err = http.ListenAndServe(":8080",nil); err != nil {
	
		log.Fatalf(err.Error())
	}
}
