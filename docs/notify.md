# RequireRecipient Function

The function is declared in `action.go` as follows:

```go
func RequireRecipient(name Name)
```

The `RequireRecipient` function is used to notify other contracts. If an account contract has the same action, then this action will be called.

The following `sender` and `receiver` code demonstrates how to send a notification from one contract to another contract.

[Complete code link](https://github.com/learnforpractice/gscdk-book/tree/master/examples/notify)

```go
// sender
package main

import (
	"github.com/uuosio/chain"
)

// contract sender
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
	chain.RequireRecipient(chain.NewName("alice"))
}
```

```go
// receiver
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
```

Let's explain the code:

The sender calls the following code to initiate the notification:

```go
chain.RequireRecipient(chain.NewName("alice"))
```

The `SayHello` function in `receiver` is different from the `SayHello` function in `sender`, the comment added to `SayHello` in `receiver` is `notify sayhello`, indicating that this action is used to receive notifications and can only be triggered by calling `RequireRecipient`.

The following is the test code:

```python
test_dir = os.path.dirname(__file__)
def deploy_contract(tester):
    with open(f'{test_dir}/sender/sender.wasm', 'rb') as f:
        code = f.read()
    with open(f'{test_dir}/sender/sender.abi', 'rb') as f:
        abi = f.read()
    tester.deploy_contract('hello', code, abi)

    with open(f'{test_dir}/receiver/receiver.wasm', 'rb') as f:
        code = f.read()
    with open(f'{test_dir}/receiver/receiver.abi', 'rb') as f:
        abi = f.read()
    tester.deploy_contract('alice', code, abi)

@chain_test
def test_notify(tester):
    deploy_contract(tester)

    r = tester.push_action('hello', 'sayhello', {}, {'hello': 'active'})
    logger.info('++++++elapsed: %s', r['elapsed'])
    tester.produce_block()
```

Compile:

```bash
cd examples/notify
./build.sh
```

Test:

```bash
ipyeos -m pytest -s -x test.py -k test_notify
```

Output:

```
[(hello,sayhello)->hello]: CONSOLE OUTPUT BEGIN =====================
Hello, World!

[(hello,sayhello)->hello]: CONSOLE OUTPUT END   =====================
debug 2023-05-24T03:26:47.224 thread-0  apply_context.cpp:40          print_debug          ] 
[(hello,sayhello)->alice]: CONSOLE OUTPUT BEGIN =====================
hello world from alice!

[(hello,sayhello)->alice]: CONSOLE OUTPUT END   =====================
INFO     test:test.py:58 ++++++elapsed: 409
```

You can see that the output includes both the message from the sender ("Hello, World!") and the receiver ("hello world from alice!"). This is because the `RequireRecipient` function in the sender contract triggered the `SayHello` function in the receiver contract.
