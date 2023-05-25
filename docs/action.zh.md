---
comments: true
---

# 内联Action在智能合约的使用

在智能合约中也可以发起一个action，这样的action称之为内联action(inline action)。需要注意的是，action是异步的，也就是说，只有在整个代码执行完后，内联action对应的合约代码才会被调用，如果被调用的合约没有定义相关的action或者账号中没有部属合约，那么调用将没有影响，但也不会有异常抛出。像这些空的内联action也不是没有任何作用，例如可以当作链上的日志，以供应用程序来查询。

以下是Action结构在`action.go`中的主要代码：

```go
type Action struct {
	Account       Name
	Name          Name
	Authorization []*PermissionLevel
	Data          []byte
}

func NewAction(perm *PermissionLevel, account Name, name Name, args ...interface{}) *Action {
...
}

func NewActionEx(perms []*PermissionLevel, account Name, name Name, args ...interface{}) *Action {
}

func (a *Action) Send() {
	data := EncoderPack(a)
	SendInline(data)
}
```

示例：

```go
// example/action_example

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
```

解释：

NewAction的参数解释如下：

- `perm *PermissionLevel` Action权限，如需要多个权限，请使用NewActionEx这个函数
- `account Name` 目标合约账号名
- `name Name` Action名称
- `args ...interface{}` 可变参数, 用于传入Action的原始参数，所有参数将被`pack`成`[]byte`类型

这个例子实际展示了如何通过内联Action来转账EOS。这里演示了两种发送Action的办法，如果你知道Action的参数，可以在`NewAction`函数中直接传入，还有一种是直接把参数封装成一个结构体，然后当作NewAction的参数。

测试代码：

```python
@chain_test
def test_action(tester):
    deploy_contract(tester, 'action_example')

    old_balance = tester.get_balance('hello')
    r = tester.push_action('hello', 'testaction1', b'', {'hello': 'active'})
    logger.info('++++++elapsed: %s', r['elapsed'])
    tester.produce_block()
    new_balance = tester.get_balance('hello')
    assert old_balance - new_balance == 10000
    logger.info("+++++++old_balance: %s, new_balance: %s", old_balance, new_balance)

    old_balance = tester.get_balance('hello')
    r = tester.push_action('hello', 'testaction2', b'', {'hello': 'active'})
    logger.info('++++++elapsed: %s', r['elapsed'])
    tester.produce_block()
    new_balance = tester.get_balance('hello')
    logger.info("+++++++old_balance: %s, new_balance: %s", old_balance, new_balance)
    assert old_balance - new_balance == 10000
```

编译：

```bash
cd examples/action_example
go-contract build .
```

运行测试：

```
ipyeos -m pytest -s -x test.py -k test_action
```

输出：

```
INFO     test:test.py:77 +++++++old_balance: 5000000.0, new_balance: 4999999.0
INFO     test:test.py:84 +++++++old_balance: 4999999.0, new_balance: 4999998.0
```

需要注意的是，为了在合约中能够调用inline action，需要在账号的`active`权限中添加`eosio.code`这个虚拟权限,在测试代码中，通过下面的函数来将`eosio.code`这个虚拟权限添加到`active`权限中。在实际的使用中，如果`active`权限继承自多个账号或者公钥，必须确保`eosio.code`的weight大于等于权限的`threshold`。

```python
def update_auth(chain, account):
    a = {
        "account": account,
        "permission": "active",
        "parent": "owner",
        "auth": {
            "threshold": 1,
            "keys": [
                {
                    "key": 'EOS6AjF6hvF7GSuSd4sCgfPKq5uWaXvGM2aQtEUCwmEHygQaqxBSV',
                    "weight": 1
                }
            ],
            "accounts": [{"permission":{"actor":account,"permission": 'eosio.code'}, "weight":1}],
            "waits": []
        }
    }
    chain.push_action('eosio', 'updateauth', a, {account:'active'})
```

[完整代码链接](https://github.com/learnforpractice/gscdk-book/tree/master/examples/action_example)
