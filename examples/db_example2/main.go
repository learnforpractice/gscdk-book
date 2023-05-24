package main

import (
	"github.com/uuosio/chain"
)

// table mytable
type A struct {
	a uint64        //primary
	b uint64        //secondary
	c chain.Uint128 //secondary
	d string
}

// contract db_example2
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
	data := &A{1, 2, chain.NewUint128(3, 0), "1"}
	mytable.Store(data, payer)
	chain.Println("++++++++teststore done!")
}

// action testupdate
func (c *MyContract) TestUpdate() {
	code := c.Receiver
	payer := c.Receiver
	mytable := NewATable(code)
	idxb := mytable.GetIdxTableByb()
	secondaryIt := idxb.Find(2)
	chain.Check(secondaryIt.IsOk(), "secondary index 2 not found")
	mytable.Updateb(secondaryIt, 3, payer)

	secondaryIt = idxb.Find(3)
	chain.Check(secondaryIt.IsOk() && secondaryIt.Primary == 1, "secondary index 3 not found")
	chain.Println("+++++++test update done!")
}

// action testremove
func (c *MyContract) TestRemove() {
	code := c.Receiver
	mytable := NewATable(code)

	idxb := mytable.GetIdxTableByb()
	secondaryIt := idxb.Find(2)

	it := mytable.Find(secondaryIt.Primary)
	chain.Check(it.IsOk(), "key does not exists!")

	mytable.Remove(it)

	secondaryIt = idxb.Find(2)
	chain.Check(!secondaryIt.IsOk(), "something went wrong")
	chain.Println("+++++testremove done!")
}

// action testbound
func (c *MyContract) TestBound() {
	code := c.Receiver
	payer := c.Receiver

	mytable := NewATable(code)
	data := &A{1, 2, chain.NewUint128(3, 0), "1"}
	mytable.Store(data, payer)
	data = &A{11, 22, chain.NewUint128(33, 0), "11"}
	mytable.Store(data, payer)
	data = &A{111, 222, chain.NewUint128(333, 0), "111"}
	mytable.Store(data, payer)

	{
		idxb := mytable.GetIdxTableByb()
		secondaryIt, secondaryValue := idxb.Lowerbound(2)
		chain.Check(secondaryIt.IsOk() && secondaryIt.Primary == 1 && secondaryValue == 2, "bad Lowerbound value")

		secondaryIt, secondaryValue = idxb.Upperbound(22)
		chain.Check(secondaryIt.IsOk() && secondaryIt.Primary == 111 && secondaryValue == 222, "bad Upperbound value")
	}

	{
		idxc := mytable.GetIdxTableByc()
		secondaryIt, secondaryValue := idxc.Lowerbound(chain.NewUint128(3, 0))
		chain.Check(secondaryIt.IsOk() && secondaryIt.Primary == 1 && secondaryValue == chain.NewUint128(3, 0), "bad Lowerbound value")

		secondaryIt, secondaryValue = idxc.Upperbound(chain.NewUint128(33, 0))
		chain.Check(secondaryIt.IsOk() && secondaryIt.Primary == 111 && secondaryValue == chain.NewUint128(333, 0), "bad Upperbound value")
		chain.Println("++++testbound done!")
	}

}
