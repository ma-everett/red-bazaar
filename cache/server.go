
/* red-bazaar/cache
 * server.go
 * mae 02014-02
 */

package cache

import (
	"net"
	"log"
	"io"
)

type Server struct {

	local CacheServer
	conn *net.TCPConn
}

func (s *Server) add(buf []byte) ([]byte,error) {

	op,ok,err := checkAdd(buf)
	if err != nil {
		log.Printf("error adding - %v\n",err)
		
		return replyAdd("__unknown__",err)		
	}
	
	if !ok {
		return replyAdd("__unknown__",DataInvalid)
	}			
	
	err = s.local.Add(op.Key,op.Content,op.TTL) 
	
	if err != nil {
		
		log.Printf("error adding - %v\n",err)
	}

	return replyAdd(op.Key,err)
}

func (s *Server) set(buf []byte) ([]byte,error) {

	op,ok,err := checkSet(buf)
	if err != nil {
		log.Printf("error setting - %v\n",err)
		
		return replySet("__unknown__",err)		
	}
	
	if !ok {
		return replySet("__unknown__",DataInvalid)
	}			
	
	err = s.local.Set(op.Key,op.Content,op.TTL) 
	
	if err != nil {
		
		log.Printf("error adding - %v\n",err)
	}

	return replySet(op.Key,err)
}

func (s *Server) get(buf []byte) ([]byte,error) {

	op,ok,err := checkGet(buf)
	if err != nil {
		log.Printf("error setting - %v\n",err)
		
		return replyGet("__unknown__",nil,err)		
	}
	
	if !ok {
		return replyGet("__unknown__",nil,DataInvalid)
	}			
	
	content,ok := s.local.Get(op.Key) 
	if !ok {
		return replyGet(op.Key,nil,CacheMiss)
	}

	return replyGet(op.Key,content,nil)
}

func (s *Server) remove(buf []byte) ([]byte,error) {

	op,ok,err := checkRemove(buf)
	if err != nil {
		log.Printf("error setting - %v\n",err)
		
		return replyRemove("__unknown__",err)		
	}
	
	if !ok {
		return replyRemove("__unknown__",DataInvalid)
	}			
	
	err = s.local.Remove(op.Key) 

	return replyRemove(op.Key,err)
}



func (s *Server) Handle() error {
	
	/* process the connection */
	defer s.conn.Close()
	
	buf := make([]byte,1024 * 5)

	for {	
		n,err := s.conn.Read(buf)
		if err != nil {
			
			if err == io.EOF {
				
				//log.Printf("client dropped - %v\n",err)
				return ClientDropped
			}
			
			//log.Printf("error handling connection on read - %v\n",err)
			continue
		}
		
		if n == 0 {
			
			//log.Printf("zero size packet\n")
			continue
		}

		/* read the operation byte */
		switch buf[0] {
		case TAddToCache:
			
			reply,err := s.add(buf[:n])
			if err != nil {
				log.Printf("write reply error - %v\n",err)
				continue
			}

			_,err = s.conn.Write(reply)
			if err != nil {
				log.Printf("write error - %v\n",err)
			}
			
			break
		case TSetToCache:
			
			reply,err := s.set(buf[:n])
			if err != nil {
				log.Printf("write reply error - %v\n",err)
				continue
			}

			_,err = s.conn.Write(reply)
			if err != nil {
				log.Printf("write error - %v\n",err)
			}

			break
		case TGetFromCache:
			
			reply,err := s.get(buf[:n])
			if err != nil {
				log.Printf("write reply error - %v\n",err)
				continue
			}
			
			_,err = s.conn.Write(reply)
			if err != nil {
				log.Printf("write error - %v\n",err)
			}			
			
			break
		case TRemoveFromCache:

			reply,err := s.remove(buf[:n])
			if err != nil {
				log.Printf("write reply error - %v\n",err)
				continue
			}

			_,err = s.conn.Write(reply)
			if err != nil {
				log.Printf("write error - %v\n",err)
			}
			
			break
		}

	
	} /* end forloop */		

	return nil
}

func NewServer(local CacheServer,conn *net.TCPConn) (*Server,error) {

	if conn == nil || local == nil {

		return nil,NoConnection
	}

	s := new(Server)
	s.local = local
	s.conn = conn

	return s,nil
}

