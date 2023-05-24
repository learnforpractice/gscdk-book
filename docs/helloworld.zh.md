---
comments: true
---

# HelloWorld

## 第一个智能合约

以下展示了一个最简单的智能合约代码和测试代码

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

测试代码：

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


编译：

```bash
cd examples/helloworld
go-contract build .
```


运行测试代码：
```bash
ipyeos -m pytest -s -x test.py -k test_helloworld
```

输出：

```
Hello, World!
```

[完整代码](https://github.com/learnforpractice/gscdk-book/tree/master/examples/helloworld)

## 创建一个初始项目

可以用`go-contract init`命令来创建一个初始项目，例如下面的代码创建了一个`mycontract`的初始项目：

```bash
go-contract init mycontract
```

创建完后可以用下面的命令编译合约：

```
cd mycontract
./build.sh
```

执行成功后会生成`mycontract.wasm`和`mycontract.abi`这两个文件

可以运行下面的命令进行测试：

```
./test.sh
```

会以绿色字体输出以下的的文字信息：

```
(hello,inc)->hello]: CONSOLE OUTPUT BEGIN =====================
count:  1

[(hello,inc)->hello]: CONSOLE OUTPUT END   =====================
INFO     test:test.py:76 ++++++elapsed: 374
debug 2023-05-24T01:51:49.481 thread-0  controller.cpp:2499           clear_expired_input_ ] removed 0 expired transactions of the 50 input dedup list, pending block time 2018-06-01T12:00:03.500
debug 2023-05-24T01:51:49.482 thread-0  apply_context.cpp:40          print_debug          ] 
[(hello,inc)->hello]: CONSOLE OUTPUT BEGIN =====================
count:  2

[(hello,inc)->hello]: CONSOLE OUTPUT END   =====================
```

需要注意的是上面的输出是调用信息，如果是在主网上运行,`chain.Println`函数输出的内容是看不到的，如果是运行在测试网，则在运行nodeos命令的时候要加上参数`--contracts-console`才能在返回中看调试输出。

在上面测试代码中，则是直接通过下面的这行代码来输出调试信息：

```python
chaintester.chain_config['contracts_console'] = True
```

另外，在发布版本的代码中，为了提高程序运行的性能，也不应该包含调试输出的代码。
