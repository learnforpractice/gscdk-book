package main

import (
	"github.com/uuosio/chain"
)

// table mytable
type A struct {
	a uint64 //primary
	b string
}

// contract db_example1
type MyContract struct {
	Receiver      chain.Name
	FirstReceiver chain.Name
	Action        chain.Name
}

func NewContract(receiver, firstReceiver, action chain.Name) *MyContract {
	return &MyContract{receiver, firstReceiver, action}
}

// action teststore
func (c *MyContract) TestStore() {
	code := c.Receiver
	payer := c.Receiver
	mydb := NewATable(code)
	data := &A{123, "hello, world"}
	mydb.Store(data, payer)
	chain.Println("done!")
}

// action testupdate
func (c *MyContract) TestUpdate() {
	code := c.Receiver
	payer := c.Receiver
	mydb := NewATable(code)
	it, data := mydb.GetByKey(123)
	chain.Check(it.IsOk(), "bad key")
	chain.Println("+++++++old value:", data.b)
	data.b = "goodbye world"
	mydb.Update(it, data, payer)
	chain.Println("done!")
}

// action testremove
func (c *MyContract) TestRemove() {
	code := c.Receiver
	mydb := NewATable(code)
	it := mydb.Find(123)
	chain.Check(it.IsOk(), "key 123 does not exists!")

	mydb.Remove(it)

	it = mydb.Find(123)
	chain.Check(!it.IsOk(), "something went wrong")
	chain.Println("+++++done!")
}
