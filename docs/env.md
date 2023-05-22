Here is the translated version of your markdown text:

---
comments: true
---

# Setting Up the Development Environment

## Direct Installation of Development Tool Package

Install the `gscdk` package, which is used for compiling Python smart contracts:

```bash
python3 -m pip install gscdk
```

Install `ipyeos`, this package is used for testing smart contracts or running a eos node:

```bash
python3 -m pip install ipyeos
```

Install `pyeoskit`, this tool is used for interacting with eos nodes, such as deploying smart contracts, etc.:

```bash
python3 -m pip install pyeoskit
```

### Running in Docker

Currently, this development tool package does not support Window and MacBook M1/M2, the development tools need to be run in docker on these platforms.

For macOS platform, it is recommended to use [OrbStack](https://orbstack.dev/download) software to install docker and run docker. Other platforms can use [Docker Desktop](https://www.docker.com/products/docker-desktop).

Download Docker image:

```bash
docker pull ghcr.io/uuosio/scdk:latest
```

Run the container:

```bash
docker run --entrypoint bash -it --rm -v "$(pwd)":/work -w /work -t ghcr.io/uuosio/scdk
```

## Test Whether the Installation Environment is Successful:

Create a new test project:

```bash
go-contract init mytest
cd mytest
```

Compile the contract code:

```bash
./build.sh
```

If there are no exceptions, the WebAssembly binary file `mytest.wasm` will be generated.

Test:

```bash
./test.sh
```

Normally, you will see the output:

```
count:  1
count:  2
```