---
comments: true
---

# Go代码里调用C/C++代码

下面以编译`say_hello`函数为例，演示如何编译代码：

say_hello.cpp

```cpp
#include <eosio/print.hpp>

extern "C" void say_hello(const char *s, uint32_t len) {
    std::string _s(s, len);
    eosio::print(_s);
}
```

接下来看下如何在Go中使用`say_hello`这个函数：

```go
package main

/*
#include <stdint.h>
void say_hello(const char *s, uint32_t len);
*/
import "C"

import (
	"unsafe"

	"github.com/uuosio/chain"
)

// table mytable
type A struct {
	a uint64        //primary
	b uint64        //secondary
	c chain.Uint128 //secondary
	d string
}

// contract callcpp
type MyContract struct {
	Receiver      chain.Name
	FirstReceiver chain.Name
	Action        chain.Name
}

func NewContract(receiver, firstReceiver, action chain.Name) *MyContract {
	return &MyContract{receiver, firstReceiver, action}
}

type stringHeader struct {
	data unsafe.Pointer
	len  uintptr
}

func GetStringPtr(str string) *C.char {
	if len(str) != 0 {
		_str := (*stringHeader)(unsafe.Pointer(&str))
		return (*C.char)(_str.data)
	}
	return (*C.char)(unsafe.Pointer(uintptr(0)))
}

// action sayhello
func (c *MyContract) say_hello() {
	str := "hello, world\n"
	C.say_hello(GetStringPtr(str), C.uint32_t(len(str)))
	chain.Println("++++++++test done!\n")
}
```

这里的

```go
/*
#include <stdint.h>
void say_hello(const char *s, uint32_t len);
*/
import "C"
```

即告诉go编译器要链接`say_hello`这个c函数。

下面的代码展示了如何调用这个函数：
```go
str := "hello, world\n"
C.say_hello(GetStringPtr(str), C.uint32_t(len(str)))
```

这里用到了`GetStringPtr`用于获取`string`类型的原始指针。需要注意的是， Go里的原始字符串并不以`\0`结尾，所以这里需要指定一下长度参数，否则如果直接在C++代码中使用可能会引起不可预知的结果。

接下来用下面的代码来测试：

test.py

```python
@chain_test
def test_say_hello(tester):
    deploy_contract(tester, 'callcpp')
    r = tester.push_action('hello', 'sayhello', b'', {'hello': 'active'})
    logger.info('++++++elapsed: %s', r['elapsed'])
    tester.produce_block()
```

运行测试：

```bash
ipyeos -m pytest -s -x test.py -k test_say_hello
```

会有下面的输出：

```
hello, world
```

[完整代码链接](https://github.com/learnforpractice/gscdk-book/tree/master/examples/callcpp)
