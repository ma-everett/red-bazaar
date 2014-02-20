
/* red-bazaar/cache
 * helpers.go
 * mae 02014-02
 */

package cache

import (
	"time"
	"bytes"
	"encoding/json"
)

type addT struct {

	Key string
	Content []byte
	TTL time.Duration
}

type addR struct {

	Key string
	Error error
}

type setT struct {

	Key string
	Content []byte
	TTL time.Duration
}

type setR struct {

	Key string
	Error error
}

type getT struct {

	Key string
}

type getR struct {

	Key string
	Content []byte
	Error error
}

type removeT struct {

	Key string
}

type removeR struct {

	Key string
	Error error
}

/* CheckAdd : check add operation */
func checkAdd(packet []byte) (*addT,bool,error) {

	if packet[0] != TAddToCache {

		return nil,false,OperationInvalid
	}

	var op addT

	err := json.Unmarshal(packet[1:],&op)
	if err != nil {

		return nil,false,DataInvalid
	}

	return &op,true,nil
}

func replyAdd(key string,err error) ([]byte,error) {

	buf := bytes.NewBuffer(nil)
	buf.WriteByte(RAddToCache)
		
	d,err := json.Marshal(addR{key,err})
	if err != nil {
		return nil,err
	}

	buf.Write(d)
	return buf.Bytes(),nil
}

func checkSet(packet []byte) (*setT,bool,error) {

	if packet[0] != TSetToCache {

		return nil,false,OperationInvalid
	}

	var op setT

	err := json.Unmarshal(packet[1:],&op)
	if err != nil {

		return nil,false,DataInvalid
	}

	return &op,true,nil
}

func replySet(key string,err error) ([]byte,error) {

	buf := bytes.NewBuffer(nil)
	buf.WriteByte(RSetToCache)

	d,err := json.Marshal(setR{key,err})
	if err != nil {
		return nil,err
	}

	buf.Write(d)
	return buf.Bytes(),nil
}

func checkGet(packet []byte) (*getT,bool,error) {

	if packet[0] != TGetFromCache {

		return nil,false,OperationInvalid
	}

	var op getT

	err := json.Unmarshal(packet[1:],&op)
	if err != nil {
		
		return nil,false,DataInvalid
	}

	return &op,true,nil
}

func replyGet(key string,content []byte,err error) ([]byte,error) {

	buf := bytes.NewBuffer(nil)
	buf.WriteByte(RGetFromCache)

	d,err := json.Marshal(getR{key,content,err})
	if err != nil {
		return nil,err
	}
	
	buf.Write(d)
	return buf.Bytes(),nil
}

func checkRemove(packet []byte) (*removeT,bool,error) {

	if packet[0] != TRemoveFromCache {

		return nil,false,OperationInvalid
	}

	var op removeT

	err := json.Unmarshal(packet[1:],&op)
	if err != nil {

		return nil,false,DataInvalid
	}

	return &op,true,nil
}

func replyRemove(key string,err error) ([]byte,error) {

	buf := bytes.NewBuffer(nil)
	buf.WriteByte(RRemoveFromCache)

	d,err := json.Marshal(removeR{key,err})
	if err != nil {
		return nil,err
	}

	buf.Write(d)
	return buf.Bytes(),nil
}
