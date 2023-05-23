---
comments: true
---

# ABI类型详解

ABI是Application Binary Interface的缩写，ABI文件为json格式，描述了action相关的参数和表结构，方便开发者和链上的智能合约和数据进行交互。

## 内置的ABI类型

以下是内置的ABI类型，一共31个

- 基本类型：name bytes string
- 数值类型： bool int8 uint8 int16 uint16 int32 uint32 int64 uint64 int128 uint128 varint32 varuint32 float32 float64 float128
- 时间相关：time_point time_point_sec block_timestamp_type
- 密码函数相关：checksum160 checksum256 checksum512 public_key signature
- Token相关：symbol symbol_code asset extended_asset

比较常用的有下面这些：

```
name bytes string bool uint64 checksum256
public_key signature symbol asset extended_asset
```
                                                                                                    
## ABI中的内置数据类型和Python中的数据类型的对应关系表

关系表：

|         ABI 类型     |   Go 类型       |      所属模块    |
|:--------------------:|:------------------:|:------------------:|
|         bool         |        bool        |   内置    |
|         int8         |         int8         |   内置    |
|         uint8        |         uint8         |   内置    |
|         int16        |         int16        |   内置    |
|         int32        |         int32        |   内置    |
|        uint32        |         uint32        |   内置    |
|         int64        |         int64        |   内置    |
|        uint64        |         uint64        |   内置    |
|        int128        |        chain.Int128        |   内置    |
|        uint128       |        chain.Uint128        |   内置    |
|       varint32       |      chain.VarInt32      |   chain |
|       varuint32      |      chain.VarUint32     |   chain |
|        float32       |     float32        |  内置     |
|        float64       |       float64        |  内置     |
|       float128       |      chain.Float128      |  chain  |
|      time_point      |      chain.TimePoint     |  chain  |
|    time_point_sec    |    chain.TimePointSec    |  chain  |
| block_timestamp_type | chain.BlockTimestampType |  chain  |
|         name         |        chain.Name        |  name  |
|         bytes        |        []byte       |  内置  |
|        string        |        string         |  内置  |
|      checksum160     |     chain.Checksum160    |  chain  |
|      checksum256     |   chain.Checksum256 |  chain  |
|      checksum512     |     chain.Checksum512    |  chain  |
|      public_key      |      chain.PublicKey     |  chain  |
|       signature      |      chain.Signature     |  chain  |
|        symbol        |       chain.Symbol       | asset   |
|      symbol_code     |     chain.SymbolCode     | asset   |
|         asset        |        chain.Asset       | asset   |
|    extended_asset    |    chain.ExtendedAsset   | asset   |
                                                                                                    
## 特别的ABI类型

### optional

### variant

### binaryextension
