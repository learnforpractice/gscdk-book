---
comments: true
---

# Cryptography Related Functions

Functions related to cryptography are defined in the `chain` module's [crypto.go](https://github.com/uuosio/chain/blob/master/crypto.go).

The following hash-related functions are defined:

```
Sha256
Sha512
Ripemd160
Sha1
```

Similar function declarations are as follows:

```go
func Sha256(data []byte) Checksum256
```

As well as the related verification functions:

```go
func AssertSha256(data []byte, hash Checksum256)
```

In addition, the following two functions are provided,
```go
func RecoverKey(digest Checksum256, sig *Signature) *PublicKey
func AssertRecoverKey(digest Checksum256, sig Signature, pub PublicKey)
```

They are used to recover the public key from the digest and signature, which can be used for signature verification within contracts.


## Example:

```go
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
```

Test code:

```python
@chain_test
def test_crypto(tester):
    deploy_contract(tester, 'crypto_example')

    r = tester.push_action('hello', 'test', {})
    tester.produce_block()

    sig = 'SIG_K1_KiXXExwMGG5NvAngS3X58fXVVcnmPc7fxgwLQAbbkSDj9gwcxWHxHwgpUegSCfgp4nFMMgjLDAKSQWZ2NLEmcJJn1m2UUg'
    pub = 'EOS7wy4M8ZTYqtoghhDRtE37yRoSNGc6zC2zFgdVmaQnKV5ZXe4kV'
    data = b'hello,world'
    args = {
        'data': data.hex(),
        'sig': sig,
        'pub': pub,
    }
    r = tester.push_action('hello', 'testrecover', args)
    logger.info('++++++elapsed: %s', r['elapsed'])
    tester.produce_block()
```

Compilation:

```
cd examples/crypto_example
go-contract build .
```

Test:

```bash
ipyeos -m pytest -s -x test.py -k test_crypto
```

In this example code, the common uses of hash functions and the `RecoverKey` function are demonstrated. The use of hash functions is quite straightforward, but let's explain the `RecoverKey` test code:

`RecoverKey` accepts two arguments, `digest` and `signature`. The `digest` is the result of running sha256 on some binary data. In the code above, the sha256 hash is computed by running the algorithm on the string `hello,world`.

In actual smart contract applications, if you want to determine whether a piece of binary data was signed with a specific private key within the smart contract, you can use the method above. The process is as follows:

- The contract stores the public key corresponding to a user's private key
- The user signs data with their private key
- The user sends the data and its corresponding signature to the smart contract
- The smart contract can call `RecoverKey` to recover the public key from the user data and the signature on the data
- The smart contract reads the user's public key stored on the chain and compares it to the public key recovered from calling `RecoverKey`. If they are the same, the contract can be certain that the data was signed by the corresponding user.