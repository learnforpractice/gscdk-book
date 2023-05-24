package main

import (
	"github.com/uuosio/chain"
)

// contract helloworld
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

// action sayhello
func (c *Contract) SayHello() {
	chain.Println("Hello, World!")
}
