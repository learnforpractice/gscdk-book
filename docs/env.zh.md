---
comments: true
---

# 开发环境搭建

## 安装开发工具包

安装gscdk包，这个包用于编译Python智能合约

```bash
python3 -m pip install gscdk
```

安装ipyeos，这个包用于测试智能合约或者运行一个节点

```bash
python3 -m pip install ipyeos
```

安装pyeoskit，这个工具用于和EOS节点进行交互，如发布智能合约等等：

```bash
python3 -m pip install pyeoskit
```

在Docker中运行

目前，该开发工具包不支持Window和Macbook M1/M2，开发工具在这些平台上需要利用docker来运行。

macOS平台推荐使用[OrbStack](https://orbstack.dev/download)软件来安装docker和运行docker。其它平台可以使用[Docker Desktop](https://www.docker.com/products/docker-desktop)。


下载Docker镜像

```bash
docker pull ghcr.io/uuosio/scdk:latest
```

运行container:

```bash
docker run --entrypoint bash -it --rm -v "$(pwd)":/work -w /work -t ghcr.io/uuosio/scdk
```

## 测试安装环境是否安装成功：

新建一个测试项目：

```bash
go-contract init mytest
cd mytest
```

编译合约代码：
```bash
./build.sh
```


不出异常会生成`mytest.wasm`这个WebAssembly的二进制文件

测试：

```bash
./test.sh
```

正常会看到输出：

```
count:  1
count:  2
```
