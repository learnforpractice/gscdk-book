---
comments: true
---

# Explanation of ABI types

## Built-in ABI types

Here are the built-in ABI types, a total of 31:

- Basic types: name bytes string
- Numerical types: bool int8 uint8 int16 uint16 int32 uint32 int64 uint64 int128 uint128 varint32 varuint32 float32 float64 float128
- Time-related: time_point time_point_sec block_timestamp_type
- Password-related functions: checksum160 checksum256 checksum512 public_key signature
- Token-related: symbol symbol_code asset extended_asset

The following are the most commonly used:

```
name bytes string bool uint64 checksum256
public_key signature symbol asset extended_asset
```

## Correspondence table of built-in data types in ABI and data types in Python

The table below shows the correspondence between the built-in types in ABI and the types in Python.

|         ABI Type     |   Go Type       |      Module    |
|:--------------------:|:------------------:|:------------------:|
|         bool         |        bool        |   Built-in    |
|         int8         |         int8         |   Built-in    |
|         uint8        |         uint8         |   Built-in    |
|         int16        |         int16        |   Built-in    |
|         int32        |         int32        |   Built-in    |
|        uint32        |         uint32        |   Built-in    |
|         int64        |         int64        |   Built-in    |
|        uint64        |         uint64        |   Built-in    |
|        int128        |        chain.Int128        |   Built-in    |
|        uint128       |        chain.Uint128        |   Built-in    |
|       varint32       |      chain.VarInt32      |   chain |
|       varuint32      |      chain.VarUint32     |   chain |
|        float32       |     float32        |  Built-in     |
|        float64       |       float64        |  Built-in     |
|       float128       |      chain.Float128      |  chain  |
|      time_point      |      chain.TimePoint     |  chain  |
|    time_point_sec    |    chain.TimePointSec    |  chain  |
| block_timestamp_type | chain.BlockTimestampType |  chain  |
|         name         |        chain.Name        |  name  |
|         bytes        |        []byte       |  Built-in  |
|        string        |        string         |  Built-in  |
|      checksum160     |     chain.Checksum160    |  chain  |
|      checksum256     |   chain.Checksum256 |  chain  |
|      checksum512     |     chain.Checksum512    |  chain  |
|      public_key      |      chain.PublicKey     |  chain  |
|       signature      |      chain.Signature     |  chain  |
|        symbol        |       chain.Symbol       | asset   |
|      symbol_code     |     chain.SymbolCode     | asset   |
|         asset        |        chain.Asset       | asset   |
|    extended_asset    |    chain.ExtendedAsset   | asset   |

## Special ABI types

### Optional

### Variant

### Binaryextension
