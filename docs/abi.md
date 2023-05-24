---
comments: true
---

# Detailed Explanation of ABI Types

ABI stands for Application Binary Interface. The ABI file is in JSON format and describes the parameters related to the action and the table structure, which makes it easier for developers to interact with smart contracts and data on the chain.

## Built-in ABI types

The following are the built-in ABI types, a total of 31

- Basic types: name bytes string
- Numeric types: bool int8 uint8 int16 uint16 int32 uint32 int64 uint64 int128 uint128 varint32 varuint32 float32 float64 float128
- Time-related: time_point time_point_sec block_timestamp_type
- Cryptography-related: checksum160 checksum256 checksum512 public_key signature
- Token-related: symbol symbol_code asset extended_asset

The more commonly used ones are the following:

```
name bytes string bool uint64 checksum256
public_key signature symbol asset extended_asset
```

## Table of Corresponding Built-in Data Types in ABI and Go

Relationship table:

|         ABI Type     |   Go Type       |      Module    |
|:--------------------:|:------------------:|:------------------:|
|         bool         |        bool        |   built-in    |
|         int8         |         int8         |   built-in    |
|         uint8        |         uint8         |   built-in    |
|         int16        |         int16        |   built-in    |
|         int32        |         int32        |   built-in    |
|        uint32        |         uint32        |   built-in    |
|         int64        |         int64        |   built-in    |
|        uint64        |         uint64        |   built-in    |
|        int128        |        chain.Int128        |   built-in    |
|        uint128       |        chain.Uint128        |   built-in    |
|       varint32       |      chain.VarInt32      |   chain |
|       varuint32      |      chain.VarUint32     |   chain |
|        float32       |     float32        |  built-in     |
|        float64       |       float64        |  built-in     |
|       float128       |      chain.Float128      |  chain  |
|      time_point      |      chain.TimePoint     |  chain  |
|    time_point_sec    |    chain.TimePointSec    |  chain  |
| block_timestamp_type | chain.BlockTimestampType |  chain  |
|         name         |        chain.Name        |  name  |
|         bytes        |        []byte       |  built-in  |
|        string        |        string         |  built-in  |
|      checksum160     |     chain.Checksum160    |  chain  |
|      checksum256     |   chain.Checksum256 |  chain  |
|      checksum512     |     chain.Checksum512    |  chain  |
|      public_key      |      chain.PublicKey     |  chain  |
|       signature      |      chain.Signature     |  chain  |
|        symbol        |       chain.Symbol       | asset   |
|      symbol_code     |     chain.SymbolCode     | asset   |
|         asset        |        chain.Asset       | asset   |
|    extended_asset    |    chain.ExtendedAsset   | asset   |

## Special ABI Types

### optional

### variant

### binaryextension