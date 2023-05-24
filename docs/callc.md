---
comments: true
---

# Calling C/C++ Code in Go

The following uses the compilation of the `say_hello` function as an example to demonstrate how to compile code:

say_hello.cpp

```cpp
#include <eosio/print.hpp>

extern "C" void say_hello(const char *s, uint32_t len) {
    std::string _s(s, len);
    eosio::print(_s);
}
```

Next, let's look at how to use the `say_hello` function in Go:

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

Here,

```go
/*
#include <stdint.h>
void say_hello(const char *s, uint32_t len);
*/
import "C"
```

tells the Go compiler to link the `say_hello` C function.

The following code shows how to call this function:

```go
str := "hello, world\n"
C.say_hello(GetStringPtr(str), C.uint32_t(len(str)))
```

Here, `GetStringPtr` is used to get the raw pointer of the `string` type. It should be noted that the raw string in Go does not end with `\0`, so here you need to specify a length parameter, otherwise, if it is used directly in C++ code, it may cause unpredictable results.

Next, use the following code to test:

test.py

```python
@chain_test
def test_say_hello(tester):
    deploy_contract(tester, 'callcpp')
    r = tester.push_action('hello', 'sayhello', b'', {'hello': 'active'})
    logger.info('++++++elapsed: %s', r['elapsed'])
    tester.produce_block()
```

Run the test:

```bash
ipyeos -m pytest -s -x test.py -k test_say_hello
```

The output will be:

```
hello, world
```

[Complete Code Link](https://github.com/learnforpractice/gscdk-book/tree/master/examples/callcpp)
