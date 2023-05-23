package main

import (
	"github.com/uuosio/chain"
)

// contract action_example
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

// action testaction1
func (c *Contract) TestAction1() {
	a := chain.NewAction(
		&chain.PermissionLevel{
			Actor:      c.receiver,
			Permission: chain.ActiveName,
		},
		chain.NewName("eosio.token"),
		chain.NewName("transfer"),
		c.receiver,                                       //from
		chain.NewName("eosio"),                           //to
		chain.NewAsset(10000, chain.NewSymbol("EOS", 4)), //quantity 1.0000 EOS
		"hello,world",                                    //memo
	)
	a.Send()
}

// packer
type Transfer struct {
	from     chain.Name
	to       chain.Name
	quantity chain.Asset
	memo     string
}

// action testaction2
func (c *Contract) TestAction2() {
	transfer := Transfer{
		from:     c.receiver,
		to:       chain.NewName("eosio"),
		quantity: chain.Asset{Amount: 10000, Symbol: chain.NewSymbol("EOS", 4)},
		memo:     "hello, world",
	}

	a := chain.NewAction(
		&chain.PermissionLevel{
			Actor:      c.receiver,
			Permission: chain.ActiveName,
		},
		chain.NewName("eosio.token"),
		chain.NewName("transfer"),
		&transfer,
	)
	a.Send()
}
