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
	mytable := NewATable(code)
	data := &A{123, "hello, world"}
	mytable.Store(data, payer)
	chain.Println("++++++++teststore done!")
}

// action testupdate
func (c *MyContract) TestUpdate() {
	code := c.Receiver
	payer := c.Receiver
	mytable := NewATable(code)
	it, data := mytable.GetByKey(123)
	chain.Check(it.IsOk(), "bad key")
	chain.Println("+++++++old value:", data.b)
	data.b = "goodbye world"
	mytable.Update(it, data, payer)
	chain.Println("++++++++testupdate done!")
}

// action testremove
func (c *MyContract) TestRemove() {
	code := c.Receiver
	mytable := NewATable(code)
	it := mytable.Find(123)
	chain.Check(it.IsOk(), "key 123 does not exists!")

	mytable.Remove(it)

	it = mytable.Find(123)
	chain.Check(!it.IsOk(), "something went wrong")
	chain.Println("+++++testremove done!")
}

// action testbound
func (c *MyContract) TestBound() {
	code := c.Receiver
	payer := c.Receiver

	mytable := NewATable(code)
	mytable.Store(&A{1, "1"}, payer)
	mytable.Store(&A{2, "2"}, payer)
	mytable.Store(&A{5, "3"}, payer)

	it := mytable.Lowerbound(1)
	chain.Check(it.IsOk() && it.GetPrimary() == 1, "bad Lowerbound value")

	it = mytable.Upperbound(2)
	chain.Check(it.IsOk() && it.GetPrimary() == 5, "bad Upperbound value")

	chain.Println("++++testbound done!")
}
