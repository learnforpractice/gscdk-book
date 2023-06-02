---
comments: true
---

# Database Operations

Storing and retrieving data on the chain is a critical feature of smart contracts. The EOS chain has implemented a memory database that supports data storage in the form of tables. Each entry in each table has a unique primary index, known as the primary key, which is a uint64 type. The raw data stored in the table is binary data of any length. When a smart contract calls the data storage function, it serializes the class data and stores it in the table. When reading, it deserializes the raw data into class objects. The system also supports secondary index tables with uint64, Uint128, Uint256, Float64, Float128 types. Secondary index tables can be seen as special tables with fixed data length. Primary index tables and secondary index tables can be used together to implement multi-index functions. There can be multiple secondary index tables. The values in the secondary index tables can be repeated, but the primary keys in the primary index tables must be unique.

Below, the use of EOS's on-chain memory database is explained with examples.

## Store

The storage function is the simplest function of the database. The following code demonstrates this function.


[db_example1](https://github.com/learnforpractice/gscdk-book/tree/master/examples/db_example1)

```go
package main

import (
	"github.com/uuosio/chain"
)

// table mytable
type A struct {
	a uint64 //primary
	b string
}

// contract test
type MyContract struct {
	Receiver      chain.Name
	FirstReceiver chain.Name
	Action        chain.Name
}

func NewContract(receiver, firstReceiver, action chain.Name) *MyContract {
	return &MyContract{receiver, firstReceiver, action}
}

// action teststore
func (c *MyContract) TestStore(name string) {
	code := c.Receiver
	payer := c.Receiver
	mytable := NewATable(code)
	data := &A{123, "hello, world"}
	mytable.Store(data, payer)
}
```

Let's explain the above code:

- The comment `// table mytable` guides the compiler to generate code related to the table, such as NewATable, which is generated code saved in the `generated.go` file.
- The comment `// contract test` indicates that `MyContract` is a smart contract class, and it will also guide the compiler to generate additional code.
- `// action teststore` means that the `TestStore` method is an `action`, and it will be triggered by the Action structure included in the Transaction.
- `NewATable(code)` specifies the creation of a table, which is stored in the account specified by `code`. In this test case, it's the `hello` account.
- The line of code `mytable.Store(data, payer)` stores the data into the blockchain database. The `payer` is used to specify which account will pay for RAM resources and needs to have already signed with the account's `active` permission in the Transaction.

Compilation:

```bash
cd examples/db_example1
go-contract build .
```

Test:

```bash
ipyeos -m pytest -s -x test.py -k test_store
```

The test code run is as follows:

```python
def test_store():
    t = init_db_test('db_example1')
    ret = t.push_action('hello', 'teststore', "", {'hello': 'active'})
    t.produce_block()
    logger.info("++++++++++%s\n", ret['elapsed'])
```

Note in this example, if there is already data in the table with a primary index of `123` of type `uint64`, then the function will throw an exception.

If you modify the above test case to the following code:

```python
def test_example1():
    t = init_db_test('db_example1')
    ret = t.push_action('hello', 'teststore', "", {'hello': 'active'})
    t.produce_block()
    logger.info("++++++++++%s\n", ret['elapsed'])

    # will raise exception
    ret = t.push_action('hello', 'teststore', "", {'hello': 'active'})
    t.produce_block()
```

Using the same command to run the test, the function will throw an exception like the following when calling `push_action` for the second time:

```
could not insert object, most likely a uniqueness constraint was violated
```

In order to avoid exceptions, when updating data in the table, you need to use the `Update` method.
Before calling `Store`, you need to check whether the primary index exists in the table. If it already exists, you cannot call the `Store` method, but must call the `Update` method.
The following example shows the usage:

## Find/Update

This section demonstrates the database's search and update functions.

```go
// db_example1

// action testupdate
func (c *MyContract) TestUpdate() {
	code := c.Receiver
	payer := c.Receiver
	mytable := NewATable(code)
	it, data := mytable.GetByKey(123)
	chain.Check(it.IsOk(), "bad key")
	chain.Println("+++++++old value:", data.b)
	data.b = "goodbye world"
	mytable.Update(it, data, payer)
	chain.Println("done!")
}

```

Testing code:

```python
@chain_test
def test_update(tester):
    deploy_contract(tester, 'db_example1')

    r = tester.push_action('hello', 'teststore', b'', {'hello': 'active'})
    tester.produce_block()

    r = tester.push_action('hello', 'testupdate', b'', {'hello': 'active'})
    logger.info('++++++elapsed: %s', r['elapsed'])
    tester.produce_block()
```

Compilation:

```bash
cd examples/db_example1
go-contract build .
```

Testing:

```bash
ipyeos -m pytest -s -x test.py -k test_update
```

Output:

```
+++++++old value: hello, world
```

As you can see, the above code is a bit complex. You first need to call `GetByKey` to get the `Iterator` and the stored value, then use `it.IsOk()` to determine whether the value corresponding to the primary index exists or not, and finally call `Update` to update the data. The `payer` is used to specify which account will pay for RAM resources and needs to have already signed with the account's `active` permission in the Transaction. It's important to note that during the update process, **the value of the primary index cannot be changed**; otherwise, an exception will be thrown.

You can try changing the update code to:

```go
data.a = 1
data.b = "goodbye world"
```

You will see an exception thrown in the smart contract, indicating:

```
mi.Update: Can not change primary key during update
```
                                                                                                    
## Remove

The following code shows how to delete an item from the database.

```go
// db_example1
// action testremove
func (c *MyContract) TestRemove() {
	code := c.Receiver
	mytable := NewATable(code)
	it := mytable.Find(123)
	chain.Check(it.IsOk(), "key 123 does not exists!")

	mytable.Remove(it)

	it = mytable.Find(123)
	chain.Check(!it.IsOk(), "something went wrong")
	chain.Println("+++++done!")
}
```

The above code first calls the `mytable.Find(123)` method to find the specified data, then calls `Remove` to delete it, and uses `it.IsOk()` to check whether the data at the specified index exists or not.

**Note:**

The `Remove` here does not need to call the `payer` account's permission specified in `Store` or `Update` to delete the data. So, in actual applications, you need to call `chain.RequireAuth` to ensure the specified account's permission can delete the data, for example:
```go
	chain.RequireAuth(chain.NewName("hello"))
```

Test code:

```python
@chain_test
def test_remove(tester):
    deploy_contract(tester, 'db_example1')

    r = tester.push_action('hello', 'teststore', b'', {'hello': 'active'})
    tester.produce_block()

    r = tester.push_action('hello', 'testremove', b'', {'hello': 'active'})
    logger.info('++++++elapsed: %s', r['elapsed'])
    tester.produce_block()
```

Compilation:

```bash
cd examples/db_example1
go-contract build .
```

Test:

```bash
ipyeos -m pytest -s -x test.py -k test_remove
```

## Lowerbound/Upperbound

These two methods are also used to find elements in the table. Unlike the `find` method, these two functions are used for fuzzy searching. The `lowerbound` method returns an `Iterator` that is `>=` to the specified `id`, while the `upperbound` method returns an `Iterator` that is `>` than the specified `id`. Let's see how to use them:

```go
// examples/db_example1

// action testbound
func (c *MyContract) TestBound() {
	code := c.Receiver
	payer := c.Receiver

	mytable := NewATable(code)
	mytable.Store(&A{1, "1"}, payer)
	mytable.Store(&A{2, "2"}, payer)
	mytable.Store(&A{5, "3"}, payer)

	it := mytable.Lowerbound(1)
	chain.Check(it.IsOk() && it.GetPrimary() == 1, "bad Lowerbound value")

	it = mytable.Upperbound(2)
	chain.Check(it.IsOk() && it.GetPrimary() == 5, "bad Upperbound value")
}
```

Test code:

```python
@chain_test
def test_bound(tester):
    deploy_contract(tester, 'db_example1')

    r = tester.push_action('hello', 'testbound', b'', {'hello': 'active'})
    tester.produce_block()
```

Compilation:

```bash
cd examples/db_example1
go-contract build .
```

Run the test:

```bash
ipyeos -m pytest -s -x test.py -k test_bound
```

Output:

```
++++testbound done!
```
                                                                                                    
## Using API to Query the Table

The above examples discuss how to operate the table in the blockchain database through smart contracts. In fact, the `get_table_rows` API interface provided by EOS can also be used to query the table on the chain. Both the `ChainTester` class in `ipyeos` and the `ChainApiAsync` and `ChainApi` classes in `pyeoskit` provide the `get_table_rows` interface to facilitate table query operations.

In Python code, the definition of `get_table_rows` is as follows:

```python
def get_table_rows(self, _json, code, scope, table,
                                lower_bound, upper_bound,
                                limit,
                                key_type='',
                                index_position='', 
                                encode_type='',
                                reverse = False,
                                show_payer = False):
    """ Fetch smart contract data from an account. 
    key_type: "i64"|"i128"|"i256"|"float64"|"float128"|"sha256"|"ripemd160"
    index_position: "2"|"3"|"4"|"5"|"6"|"7"|"8"|"9"|"10"
    encode_type: "dec" or "hex", default to "dec"
    """
```

Let's break down the parameters for this function:

- `_json`: If True, it returns JSON-formatted database records, if False, it returns the raw data in hexadecimal format.
- `code`: Represents the account where the table is located.
- `scope`: Generally set as an empty string. When `code` and `table` are the same, different `scope` can be used to differentiate different tables.
- `table`: The name of the table to be queried.
- `lower_bound`: The starting value of the primary or secondary index, specified by `key_type`. It can be numeric, a string of numbers, or a hexadecimal string.
- `upper_bound`: The ending value of the primary or secondary index, specified by `key_type`. It can be numeric, a string of numbers, or a hexadecimal string. If empty, it means no upper limit has been set. If a non-empty value is set, the results will return all values that are `>=lower_bound` and `<=upper_bound`.
- `limit`: Limits the number of returned values. If the queried records exceed the limit, `more` will be set to `true` in the returned values, and `next_key` will point to the next valid index.
- `key_type`: The values can be: `"name"`, `"i64"`, `"i128"`, `"i256"`, `"float64"`, `"float128"`, `"sha256"`, `"ripemd160"`. For the primary index (i.e., `index_position` is `1`), the value can only be `"name"` or `"i64"`. For secondary index (i.e., `index_position >= 2`), the value could be any of the listed types. The encoding method of `lower_bound` and `upper_bound` under each value will be explained separately below.
- `index_position`: Specifies the relative position of the index. If it's empty or `1`, it denotes the primary index. Any number above `2` denotes the position of the secondary index.
- `encode_type`: It's either `"dec"` or `"hex"`, defaulting to `"dec"`. It specifies the encoding format of `lower_bound`, `upper_bound`, and the return value `next_key`.
- `reverse`: Specifies whether the data should be returned in reverse order.
- `show_payer`: Specifies whether to display the RAM resource paying account.

Detailed explanation for `key_type`:

- "name" is a `name` type string.
- "i64" can be a numeric type or a string of numbers, such as 123 or "123".
- "i128" can be a numeric type, a string of numbers, or a hexadecimal string, such as: 123, "123", "0xaabb", "aabb".
- "i256" when the value of `encode_type` is `"dec"` or an empty string `""`, the encoding format is: a hexadecimal string, represented in **little-endian mode**, 64 characters in length. For example: `fb54b91bfed2fe7fe39a92d999d002c550f0fa8360ec998f4bb65b00c86282f5` will be converted into two `uint128_t` type values in little-endian mode: `50f0fa8360ec998f4bb65b00c86282f5` and `fb54b91bfed2fe7fe39a92d999d002c5`. When the value of `encode_type` is `"hex"`, it uses the same encoding method as the `"sha256"` type, which is big-endian mode.
- "float64": The value is a floating-point string, like `"123.456"`.
- "float128": When the value of `encode_type` is `"dec"` or an empty string `""`, the value is a floating-point string, like `"123.456"`, and the range it represents can only be within the range allowed by `float64`. When the value of `encode_type` is `"hex"`, `encode_type` represents the data as a hexadecimal string in little-endian mode.
- "sha256": A hexadecimal string represented in **big-endian mode**, 64 characters long, will be converted into two `uint128_t` type values in little-endian mode: such as `f58262c8005bb64b8f99ec6083faf050c502d099d9929ae37ffed2fe1bb954fb` will be converted into `50f0fa8360ec998f4bb65b00c86282f5` and `fb54b91bfed2fe7fe39a92d999d002c5`. Refer to the [keytype_converter](https://github.com/AntelopeIO/leap/blob/db132c5fd44e0b1c492e46e3f51e185cd5c59ed0/plugins/chain_plugin/include/eosio/chain_plugin/chain_plugin.hpp#L900) structure's code for more details.
- "ripemd160": A hexadecimal string, 64 characters long, big-endian mode, will be converted into two `uint128_t` type values in little-endian mode: such as `83a83a3876c64c33f66f33c54f1869edef5b5d4a000000000000000000000000` will be converted into `ed69184fc5336ff6334cc676383aa883` and `0000000000000000000000004a5d5bef`. Refer to the [keytype_converter](https://github.com/AntelopeIO/leap/blob/db132c5fd44e0b1c492e46e3f51e185cd5c59ed0/plugins/chain_plugin/include/eosio/chain_plugin/chain_plugin.hpp#L918) structure's code for more details.

The `get_table_rows` function's parameters are quite complex, here's a summary:

- If `lower_bound` and `upper_bound` are empty, it means the query has no range limit.
- When the value of `key_type` is `"i256"` and `"float128"`, the encoding method of `lower_bound` and `upper_bound` is also affected by `encode_type`.


```python
@chain_test
def test_offchain_find(tester):
    deploy_contract(tester, 'db_example1')

    r = tester.push_action('hello', 'testbound', b'', {'hello': 'active'})
    tester.produce_block()

    r = tester.get_table_rows(False, 'hello', '', 'mytable', '', '', 10)
    logger.info("+++++++rows: %s", r)

    r = tester.get_table_rows(True, 'hello', '', 'mytable', '', '', 10)
    logger.info("+++++++rows: %s", r)

    r = tester.get_table_rows(True, 'hello', '', 'mytable', '1', '2', 10)
    logger.info("+++++++rows: %s", r)
```

Output:

```
+++++++rows: {'rows': ['01000000000000000131', '02000000000000000132', '05000000000000000133'], 'more': False, 'next_key': ''}
+++++++rows: {'rows': [{'a': 1, 'b': '1'}, {'a': 2, 'b': '2'}, {'a': 5, 'b': '3'}], 'more': False, 'next_key': ''}
+++++++rows: {'rows': [{'a': 1, 'b': '1'}, {'a': 2, 'b': '2'}], 'more': False, 'next_key': ''}
```

## Secondary Index Operations

First, please look at the following example:

[db_example2](https://github.com/learnforpractice/gscdk-book/tree/master/examples/db_example2)

```go
// db_example2
package main

import (
	"github.com/uuosio/chain"
)

// table mytable
type A struct {
	a uint64        //primary
	b uint64        //secondary
	c chain.Uint128 //secondary
	d string
}

// contract db_example2
type MyContract struct {
	Receiver      chain.Name
	FirstReceiver chain.Name
	Action        chain.Name
}

func NewContract(receiver, firstReceiver, action chain.Name) *MyContract {
	return &MyContract{receiver, firstReceiver, action}
}

// action teststore
func (c *MyContract) TestStore() {
	code := c.Receiver
	payer := c.Receiver
	mytable := NewATable(code)
	data := &A{1, 2, chain.NewUint128(3, 0), "1"}
	mytable.Store(data, payer)
	chain.Println("++++++++teststore done!")
}
```

In this example, two secondary indexes are defined:

```go
b uint64        //secondary
c chain.Uint128 //secondary
```

Test code:

```python
# test.py
@chain_test
def test_store(tester):
    deploy_contract(tester, 'db_example2')
    r = tester.push_action('hello', 'teststore', b'', {'hello': 'active'})
    logger.info('++++++elapsed: %s', r['elapsed'])
    tester.produce_block()
```

Compilation:

```bash
cd examples/db_example2
go-contract build .
```

Run the test:

```bash
ipyeos -m pytest -s -x test.py -k test_store
```

Summary: Compared with the primary index examples, if a table contains secondary indexes, the method called for storage is the same, both call the `Store` method.

                                                                                                    
## Updating Secondary Index

In practical applications, sometimes we need to update secondary indexes. First, look at the following code:

```go
// db_example2

// action testupdate
func (c *MyContract) TestUpdate() {
	code := c.Receiver
	payer := c.Receiver
	mytable := NewATable(code)

    idxb := mytable.GetIdxTableByb()
	secondaryIt := idxb.Find(2)
	chain.Check(secondaryIt.IsOk(), "secondary index 2 not found")
	mytable.Updateb(secondaryIt, 3, payer)

	secondaryIt = idxb.Find(3)
	chain.Check(secondaryIt.IsOk() && secondaryIt.Primary == 1, "secondary index 3 not found")
	chain.Println("+++++++test update done!")
}
```

Pay attention to the code above:

```go
idxb := mytable.GetIdxTableByb()
secondaryIt := idxb.Find(2)
chain.Check(secondaryIt.IsOk(), "secondary index 2 not found")
mytable.Updateb(secondaryIt, 3, payer)

secondaryIt = idxb.Find(3)
chain.Check(secondaryIt.IsOk() && secondaryIt.Primary == 1, "secondary index 3 not found")
chain.Println("+++++++test update done!")
```

Here is a brief description of the process:

- `idxb := mytable.GetIdxTableByb()` fetches the secondary index of `b`, `GetIdxTableByb` is an auto-generated function, you can find the code in `generated.go`.
- `secondaryIt := idxb.Find(2)` looks for the secondary index with type `uint64` value `2`. The returned value `secondaryIt` is of type `SecondaryIterator`.
- **`mytable.Updateb(secondaryIt, uint64(3), payer)`** This line of code implements the update functionality, updating the value of `b` to `3`. `Updateb` is an auto-generated function defined in `generated.go`.
- `secondaryIt = idxb.Find(3)` looks for the new secondary index.
- `chain.Check(secondaryIt.IsOk() && secondaryIt.Primary == 1, "secondary index 3 not found")` is used to confirm whether the secondary index update was successful. Note that here we are also checking if the primary index is `1`.

## Querying a Secondary Index

Secondary indices also support query methods such as `Find`, `Lowerbound`, and `Upperbound`. The following example demonstrates how to query the values in the `b` and `c` secondary indices:

```go
// action testbound
func (c *MyContract) TestBound() {
	...
}
```

You can run the test case in the example directory using the following command:

```bash
ipyeos -m pytest -s -x test.py -k test_bound
```

## Removing a Secondary Index

```go
// action testremove
func (c *MyContract) TestRemove() {
	...
}
```

Here is a brief explanation of the above code:

- `secondaryIt := idxb.Find(2)` looks for the secondary index.
- `it := mytable.Find(secondaryIt.Primary)` uses `SecondaryIterator` to get the primary index, and then returns the primary index's `Iterator`.
- `mytable.Remove(it)` removes elements from the table, including the primary index and all secondary indices.

From the above examples, we can see that the removal of a secondary index involves first locating the primary index through the secondary index, and then deleting via the primary index.

To compile:

```bash
cd examples/db_example2
go-contract build .
```

To run the test:

```bash
ipyeos -m pytest -s -x test.py -k test_remove
```

## Using APIs to Query a Table with Secondary Index

The `get_table_rows` API also supports finding corresponding values through secondary indices.

```python
@chain_test
def test_offchain_find(tester):
    ...
```

Run the test case:

```bash
ipyeos -m pytest -s -x test.py -k test_offchain_find
```

The results of running the above test code are as follows:

```
{'rows': [{'a': 1, 'b': 2, 'c': '3', 'd': '1'}, {'a': 11, 'b': 22, 'c': '33', 'd': '11'}, {'a': 111, 'b': 222, 'c': '333', 'd': '111'}], 'more': False, 'next_key': ''}
{'rows': [{'a': 1, 'b': 2, 'c': '3', 'd': '1'}, {'a': 11, 'b': 22, 'c': '33', 'd': '11'}, {'a': 111, 'b': 222, 'c': '333', 'd': '111'}], 'more': False, 'next_key': ''}
```

## Conclusion

EOS's data storage function is quite comprehensive, and the availability of secondary index tables makes data retrieval extremely flexible. This chapter provided a detailed explanation of the add, delete, update, and query operations for database tables. There's a lot of content in this chapter, so take your time to digest it all. Try to make some modifications to the examples and run them to enhance your understanding of the knowledge points in this chapter.
