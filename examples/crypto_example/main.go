package main

import (
	"github.com/uuosio/chain"
)

// contract crypto_example
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

// action test
func (c *Contract) Test() {
	message := []byte("hello, world")
	{
		h := chain.Sha256(message)
		chain.AssertSha256(message, h)
	}

	{
		h := chain.Sha512(message)
		chain.AssertSha512(message, h)
	}

	{
		h := chain.Ripemd160(message)
		chain.AssertRipemd160(message, h)
	}

	{
		h := chain.Sha1(message)
		chain.AssertSha1(message, h)
	}

	chain.Println("Test Done!")
}

// action testrecover
func (c *Contract) TestRecover(data []byte, sig *chain.Signature, pub *chain.PublicKey) {
	hash := chain.Sha256([]byte("hello,world"))
	pub2 := chain.RecoverKey(hash, sig)
	chain.Check(*pub == *pub2, "bad recovery")
	chain.Println("TestRecover Done!")
}
