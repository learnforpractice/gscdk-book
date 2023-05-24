---
comments: true
---

# HelloWorld

## First Smart Contract

The following shows a simplest smart contract code and test code:

```go
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
```

Test code:

```python
# test.py

import os
import sys
import json
import struct
import pytest

test_dir = os.path.dirname(__file__)
sys.path.append(os.path.join(test_dir, '..'))

from ipyeos import log
from ipyeos import chaintester
chaintester.chain_config['contracts_console'] = True

logger = log.get_logger(__name__)

def init_tester():
    chain = chaintester.ChainTester()
    return chain

def chain_test(fn):
    def call():
        chain = init_tester()
        ret = fn(chain)
        chain.free()
        return ret
    return call

class NewChainTester():
    def __init__(self):
        self.tester = None

    def __enter__(self):
        self.tester = init_tester()
        return self.tester

    def __exit__(self, type, value, traceback):
        self.tester.free()

test_dir = os.path.dirname(__file__)
def deploy_contract(tester, package_name):
    with open(f'{test_dir}/{package_name}.wasm', 'rb') as f:
        code = f.read()
    with open(f'{test_dir}/{package_name}.abi', 'rb') as f:
        abi = f.read()
    tester.deploy_contract('hello', code, abi)

@chain_test
def test_action(tester):
    deploy_contract(tester, 'helloworld')

    r = tester.push_action('hello', 'sayhello', {}, {'hello': 'active'})
    logger.info('++++++elapsed: %s', r['elapsed'])
    tester.produce_block()
```


Compile:

```bash
cd examples/helloworld
go-contract build .
```


Run the test code:
```bash
ipyeos -m pytest -s -x test.py -k test_helloworld
```

Output:

```
Hello, World!
```

[Complete code](https://github.com/learnforpractice/gscdk-book/tree/master/examples/helloworld)

## Creating an Initial Project

You can use the `go-contract init` command to create an initial project. For example, the following code creates an initial project named `mycontract`:

```bash
go-contract init mycontract
```

After creation, you can use the following command to compile the contract:

```
cd mycontract
./build.sh
```

After successful execution, two files `mycontract.wasm` and `mycontract.abi` will be generated.

You can run the following command to test:

```
./test.sh
```

The following text information will be output in green:

(hello,inc)->hello]: CONSOLE OUTPUT BEGIN =====================
count:  1

[(hello,inc)->hello]: CONSOLE OUTPUT END   =====================
INFO     test:test.py:76 ++++++elapsed: 374
debug 2023-05-24T01:51:49.482 thread-0  apply_context.cpp:40          print_debug          ] 
[(hello,inc)->hello]: CONSOLE OUTPUT BEGIN =====================
count:  2

[(hello,inc)->hello]: CONSOLE OUTPUT END   =====================
```

It's worth noting that the above output is the debug formation. If running on the main network, the content output by the `chain.Println` function cannot be seen. If running on a test network, you need to add the `--contracts-console` parameter when running the nodeos command to see the debug output in the return.

In the above test code, the debug information is output directly through the following line of code:

```python
chaintester.chain_config['contracts_console'] = True
```

Furthermore, in the production version of the code, to improve the performance of the program, you should not include the debug output code.
