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
