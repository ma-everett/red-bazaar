/* red-bazaar/cache
 * client.go
 * mae 02014-02
 */

package cache

import (
	"net"
	"time"
	"bytes"
	"encoding/json"
)

type Client struct {

	buffer []byte
	conn *net.TCPConn
}

/* Add : add content under key with a time-to-live in a local cache */
func (c *Client) Add(key string, content []byte, ttl time.Duration) (error) {	

	if c.conn == nil {
		return NoConnection
	}

	/* encode data */
	op := addT{key,content,ttl}
	
	data,err := json.Marshal(op)
	if err != nil {

		return DataInvalid
	}

	/* send data */
	buf := bytes.NewBuffer(nil)
	buf.WriteByte(TAddToCache)
	buf.Write(data)
	
	_,err = c.conn.Write(buf.Bytes())
	if err != nil {
		return err
	}

	n,err := c.conn.Read(c.buffer)
	if err != nil {
		return err
	}
	buf = bytes.NewBuffer(c.buffer[:n])

	/* read the command */
	if r,_ := buf.ReadByte(); r != RAddToCache {

		return DataInvalid
	}

	var reply addR
	err = json.Unmarshal(buf.Bytes(),&reply)

	if reply.Key != key {

		return DataInvalid
	}
	
	return reply.Error
}

/* Set : set the content under key with a time-to-live in a local cache */
func (c *Client) Set(key string, content []byte, ttl time.Duration) (error) {

	op := setT{key,content,ttl}

	data,err := json.Marshal(op)
	if err != nil {

		return DataInvalid
	}

	buf := bytes.NewBuffer(nil)
	buf.WriteByte(TSetToCache)
	buf.Write(data)

	
	_,err = c.conn.Write(buf.Bytes())
	if err != nil {
		return err
	}

	n,err := c.conn.Read(c.buffer)
	if err != nil {
		return err
	}

	buf = bytes.NewBuffer(c.buffer[:n])

	if r,_ := buf.ReadByte(); r != RSetToCache {

		return DataInvalid
	}
	
	var reply setR
	err = json.Unmarshal(buf.Bytes(),&reply)
	
	if reply.Key != key {
		return DataInvalid
	}

	return reply.Error
}

/* Get : get content under a key from a local cache */
func (c *Client) Get(key string) ([]byte,bool,error) {

	op := getT{key}

	data,err := json.Marshal(op)
	if err != nil {

		return nil,false,DataInvalid
	}

	buf := bytes.NewBuffer(nil)
	buf.WriteByte(TGetFromCache)
	buf.Write(data)

	_,err = c.conn.Write(buf.Bytes())
	if err != nil {
		return nil,false,err
	}

	n,err := c.conn.Read(c.buffer)
	
	buf = bytes.NewBuffer(c.buffer[:n])

	if r,_ := buf.ReadByte(); r != RGetFromCache {

		return nil,false,DataInvalid
	}

	var reply getR

	err = json.Unmarshal(buf.Bytes(),&reply)
	if err != nil {
		return nil,false,err
	}

	if reply.Key != key {
		return nil,false,DataInvalid
	}

	if reply.Content == nil {
		return nil,false,nil
	}

	d := make([]byte,len(reply.Content))
	copy(d,reply.Content)

	return d,true,nil
}

/* Remove : remove content under a key from a local cache */
func (c *Client) Remove(key string) (error) {

	op := removeT{key}

	data,err := json.Marshal(op)
	if err != nil {

		return DataInvalid
	}

	buf := bytes.NewBuffer(nil)
	buf.WriteByte(TRemoveFromCache)
	buf.Write(data)

	_,err = c.conn.Write(buf.Bytes())
	if err != nil {
		return err
	}

	n,err := c.conn.Read(c.buffer)
	buf = bytes.NewBuffer(c.buffer[:n])
	
	if r,_ := buf.ReadByte(); r != RRemoveFromCache {

		return DataInvalid
	}

	var reply removeR

	err = json.Unmarshal(buf.Bytes(),&reply)
	if err != nil {
		return err
	}

	if reply.Key != key {
		return DataInvalid
	}

	return reply.Error
}

func (c *Client) Dial(url string) error {

	naddr,err := net.ResolveTCPAddr("tcp",url)
	if err != nil {

		return err
	}

	conn,err := net.DialTCP(naddr.Network(), nil, naddr)
	if err != nil {

		return err
	}	
	
	c.conn = conn

	return nil
}

func (c *Client) Close() error {

	if c.conn != nil {

		c.conn.Close()
		return nil
	}

	return NoConnection
}

func NewClient() *Client {

	c := new(Client)
	c.buffer = make([]byte,1024 * 5) /* five megs max ?? */
	return c
}
