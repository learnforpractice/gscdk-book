---
comments: true
---

# Common Smart Contract Functions

## IsAccount

Declaration:

```go
func IsAccount(name Name) bool
```

Description:

Used to determine whether an account exists or not.

## HasAuth

Declaration:

```go
func HasAuth(name Name) bool
```

Description:

Used to determine whether there is `active` authority for a specified account, that is, whether the Transaction has been signed with the private key corresponding to the `active` authority of the specified account. There must be at least one corresponding private key, possibly more.

## RequireAuth/RequireAuth2

Declaration:

```go
func RequireAuth(name Name)
func RequireAuth2(name Name, permission Name)
```

Description:

These two functions will throw exceptions if the account does not exist or if the specified account's permission is not detected. The difference is that `RequireAuth` checks for the existence of `active` permission, while `RequireAuth2` can check for specified permissions.

## CurrentTime

```go
func CurrentTime() TimePoint
```

Used to get the time of the block where the Transaction is located.

## Check

Declaration:

```go
func Check(test bool, msg string)
```

Description:

If test is false, it will throw an exception containing the `msg` message.
