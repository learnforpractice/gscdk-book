---
comments: true
---

# RequireRecipient函数

函数在`action.go`中的声明如下：

```go
func RequireRecipient(name Name)
```

`RequireRecipient`函数用来通知其它合约. 如果account合约有相同的action，那么这个action将被调用。

以下的`sender`, `receiver`的代码演示了如何从一个合约发送通知到另一个合约。

[完整代码链接](https://github.com/learnforpractice/gscdk-book/tree/master/examples/notify)

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

解释下代码：

sender中调用下面的代码来发起通知：

```go
chain.RequireRecipient(chain.NewName("alice"))
```

`receiver`中的`SayHello`函数和`sender`中的`SayHello`函数的定义有些不同，`receiver`中的`SayHello`的添加的注释是`notify sayhello`，用来指定这个action是一个用来接收通知的action，只能通过调用`RequireRecipient`来触发。


以下是测试代码：

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

编译：
```bash
cd examples/notify
./build.sh
```

测试：

```bash
ipyeos -m pytest -s -x test.py -k test_notify
```

输出：

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
