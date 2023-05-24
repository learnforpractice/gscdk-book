---
comments: true
---

# 密码学相关函数

密码学相关的函数在`chain`模块的[crypto.go](https://github.com/uuosio/chain/blob/master/crypto.go)中定义.

定义了如下的hash相关的函数：

```
Sha256
Sha512
Ripemd160
Sha1
```

类似的函数声明如下：

```go
func Sha256(data []byte) Checksum256
```

以及相关的验证函数：

```go
func AssertSha256(data []byte, hash Checksum256)
```

另外，还提供了下面这两个函数，
```go
func RecoverKey(digest Checksum256, sig *Signature) *PublicKey
func AssertRecoverKey(digest Checksum256, sig Signature, pub PublicKey)
```

用于从digest和signture中恢复出公钥，可以用于在合约中进行签名的效验


## 示例：

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

测试代码：

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

编译：

```
cd examples/crypto_example
go-contract build .
```

测试：

```bash
ipyeos -m pytest -s -x test.py -k test_crypto
```

在这个示例代码中，分别演示了常用的hash函数的用法以及`RecoverKey`的用法。hash函数的用法比较简单，这里解释一下RecoverKey的测试代码：
RecoverKey接受二个参数，分别是`digest`和`signature`，digest是对一个二进制数据进行sha256运行的结果。在上面的代码中是通过对`hello,world`进行sha256算法的hash计算。


在实际的智能合约的应用中，如果要在智能合约里判断某段二进制数据是否是用特定的私钥进行的签名也可以用上面的办法。过程如下：

- 合约中保存用户一个私钥对应的公钥
- 用户用自己的私钥对数据进行签名
- 用户将数据，以及对应的签名传给智能合约
- 智能合约可以调用`RecoverKey`从用户数据，以及对数据的签名中还原出公钥
- 智能合约读取保存在链上的用户公钥，与通过调用`RecoverKey`还原出的公钥进行比较，相同即可以确定数据是对应的用户签的名
