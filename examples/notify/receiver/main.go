package main

import (
	"github.com/uuosio/chain"
)

// contract receiver
type Contract struct {
	receiver      chain.Name
	firstReceiver chain.Name
	action        chain.Name
}

func NewContract(receiver, firstReceiver, action chain.Name) *Contract {
	return &Contract{
		receiver,
		firstReceiver,
		action,
	}
}

// notify sayhello
func (c *Contract) SayHello() {
	chain.Println("hello world from alice!")
}
