---
comments: true
---

# 常用智能合约函数

## IsAccount

声明：

```go
func IsAccount(name Name) bool
```

说明：

用来判断账号存不存在

## HasAuth

声明：

```go
func HasAuth(name Name) bool
```

说明：

用来判断是否有指定账号的`active`权限，也就是Transaction是否有用指定账号的`active`权限所对应的私钥进行签名。对应的私钥最少有一个，也可能二个以上。

## RequireAuth/RequireAuth2

声明：

```go
func RequireAuth(name Name)
func RequireAuth2(name Name, permission Name)
```

说明：

这两个函数在账号不存在或者没有检测到有指定账号的权限时都会抛出异常，不同的是`RequireAuth`为检测是否存在`active`权限，而`RequireAuth2`可以检测指定的权限。

## CurrentTime

```go
func CurrentTime() TimePoint
```

用于获取Transaction所在的区块的时间

## Check

声明：

```go
func Check(test bool, msg string)
```

说明：

如果test为false，则会抛出包含`msg`信息的异常。
