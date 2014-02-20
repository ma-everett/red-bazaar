
package cache

import (
	"testing"
)

const (
	url = "loopback:7070"
)
	

func TestAdd(t *testing.T) {

	client := NewClient()
	err := client.Dial(url)
	if err != nil {
		t.Fatal(err)
	}

	defer client.Close()

	err = client.Add("hello",[]byte("there"),0)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSet(t *testing.T) {

	client := NewClient()
	err := client.Dial(url)
	if err != nil {
		t.Fatal(err)
	}

	defer client.Close()

	err = client.Set("hello",[]byte("there,yo yo"),0)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGet(t *testing.T) {

	client := NewClient()
	err := client.Dial(url)
	if err != nil {
		t.Fatal(err)
	}

	defer client.Close()

	_,_,err = client.Get("hello")
	if err != nil {
		t.Fatal(err)
	}
}

func TestRemove(t *testing.T) {

	client := NewClient()
	err := client.Dial(url)
	if err != nil {
		t.Fatal(err)
	}

	defer client.Close()

	err = client.Remove("hello")
	if err != nil {
		t.Fatal(err)
	}
}


