---
comments: true
---

# Use of Inline Action in Smart Contracts

An action can also be initiated in a smart contract, such an action is called an inline action. It should be noted that actions are asynchronous, that is, the contract code corresponding to the inline action will only be called after the entire code has been executed. If the called contract does not define the relevant action or there is no deployed contract in the account, the call will have no effect, but no exceptions will be thrown. Empty inline actions like these are not without any effect, for example, they can be used as logs on the chain for applications to query.

The following is the main code of the Action structure in `action.go`:

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

Example:

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

Explanation:

The parameters of NewAction are explained as follows:

- `perm *PermissionLevel` Action permission, if you need multiple permissions, please use the NewActionEx function
- `account Name` Target contract account name
- `name Name` Action name
- `args ...interface{}` Variable parameters, used to pass in the original parameters of the Action, all parameters will be `pack`ed into `[]byte` type

This example actually demonstrates how to transfer EOS through an inline Action. Here, two ways of sending Actions are demonstrated. If you know the parameters of the Action, you can directly pass them into the `NewAction`

function. Another way is to encapsulate the parameters into a structure and use it as a parameter for NewAction.

Test code:

```python
@chain_test
def test_action(tester):
    deploy_contract(tester, 'action_example')

    old_balance = tester.get_balance('hello')
    r = tester.push_action('hello', 'testaction1', b'', {'hello': 'active'})
    logger.info('++++++elapsed: %s', r['elapsed'])
    tester.produce_block()
    new_balance = tester.get_balance('hello')
    assert abs(old_balance - new_balance - 1.0000) < 1e-6
    logger.info("+++++++old_balance: %s, new_balance: %s", old_balance, new_balance)

    old_balance = tester.get_balance('hello')
    r = tester.push_action('hello', 'testaction2', b'', {'hello': 'active'})
    logger.info('++++++elapsed: %s', r['elapsed'])
    tester.produce_block()
    new_balance = tester.get_balance('hello')
    logger.info("+++++++old_balance: %s, new_balance: %s", old_balance, new_balance)
    assert abs(old_balance - new_balance - 1.0000) < 1e-6
```

Compilation:

```bash
cd examples/action_example
go-contract build .
```

Running the test:

```
ipyeos -m pytest -s -x test.py -k test_action
```

Output:

```
INFO     test:test.py:77 +++++++old_balance: 5000000.0, new_balance: 4999999.0
INFO     test:test.py:84 +++++++old_balance: 4999999.0, new_balance: 4999998.0
```

Note that in order to be able to call an inline action in a contract, you need to add the virtual permission `eosio.code` to the account's `active` permission. In the test code, the following function is used to add the virtual permission `eosio.code` to the `active` permission. In actual use, if the `active` permission is inherited from multiple accounts or public keys, you must ensure that the weight of `eosio.code` is greater than or equal to the threshold of the permission.

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

[Code Link](https://github.com/learnforpractice/gscdk-book/tree/master/examples/action_example)
