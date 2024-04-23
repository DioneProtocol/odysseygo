Node implementation for the Odyssey network -
a blockchains platform with high throughput, and blazing fast transactions.

## Installation

Odyssey is an incredibly lightweight protocol, so the minimum computer requirements are quite modest.
Note that as network usage increases, hardware requirements may change.

The minimum recommended hardware specification for nodes connected to Mainnet is:

- CPU: Equivalent of 8 AWS vCPU
- RAM: 16 GiB
- Storage: 1 TiB
- OS: Ubuntu 20.04/22.04 or macOS >= 12
- Network: Reliable IPv4 or IPv6 network connection, with an open public port.

If you plan to build OdysseyGo from source, you will also need the following software:

- [Go](https://golang.org/doc/install) version >= 1.20.8
- [gcc](https://gcc.gnu.org/)
- g++

### Building From Source

#### Clone The Repository

Clone the OdysseyGo repository:

```sh
git clone git@github.com:DioneProtocol/odysseygo.git
cd odysseygo
```

This will clone and checkout the `master` branch.

#### Building OdysseyGo

Build OdysseyGo by running the build script:

```sh
./scripts/build.sh
```

The `odysseygo` binary is now in the `build` directory. To run:

```sh
./build/odysseygo
```

## Running Odyssey

### Connecting to Mainnet

To connect to the Odyssey Mainnet, run:

```sh
./build/odysseygo
```

You should see some pretty ASCII art and log messages.

You can use `Ctrl+C` to kill the node.

### Connecting to Testnet

To connect to the Testnet, run:

```sh
./build/odysseygo --network-id=testnet
```

### Run local testnet

To run the local testnet, run:
```sh
./scripts/run_blockchain.sh
```

It is possible to configure the launch parameters of each individual node in the run_config.json file in the `bootstrapNodes` block. In the `bootstrapNodes` block the parameters of validator nodes are set. In the `nodes` block the parameters of ordinary nodes are set, which can be launched for a local testnet. You can set the following options:
- id - Identifier that was generated by the node using the tls private key
- networkId - Network identifier with the parameters of which the local testnet will be launched
- ip - Public IP of this node for P2P communication
- hostIp - Address of the HTTP server
- httpPort - Port of the HTTP server
- stakingPort - Port of the consensus server
- dbDir - The directory in which node data and logs will be saved
- tlsKeyFilePath -TLS private key file for the HTTPs server
- tlsCertFilePath - TLS certificate file for the HTTPs server

#### Run node for local testnet

To run the ordinary node for local testnet, run:
```sh
./scripts/run_node.sh -n node5
```

It is possible to configure the launch parameters of each individual node in the run_config.json file in the `nodes` block. You can set the following options:
- networkId - Network identifier with the parameters of which the local testnet will be launched
- ip - Public IP of this node for P2P communication
- hostIp - Address of the HTTP server
- httpPort - Port of the HTTP server
- stakingPort - Port of the consensus server
- dbDir - The directory in which node data and logs will be saved

## Bootstrapping

A node needs to catch up to the latest network state before it can participate in consensus and serve API calls. This process (called bootstrapping) currently takes several days for a new node connected to Mainnet.

A node will not report healthy until it is done bootstrapping.

Improvements that reduce the amount of time it takes to bootstrap are under development.

The bottleneck during bootstrapping is typically database IO. Using a more powerful CPU or increasing the database IOPS on the computer running a node will decrease the amount of time bootstrapping takes.

## Generating Code

OdysseyGo uses multiple tools to generate efficient and boilerplate code.

### Running protobuf codegen

To regenerate the protobuf go code, run `scripts/protobuf_codegen.sh` from the root of the repo.

This should only be necessary when upgrading protobuf versions or modifying .proto definition files.

To use this script, you must have [buf](https://docs.buf.build/installation) (v1.26.1), protoc-gen-go (v1.30.0) and protoc-gen-go-grpc (v1.3.0) installed.

To install the buf dependencies:

```sh
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.30.0
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0
```

If you have not already, you may need to add `$GOPATH/bin` to your `$PATH`:

```sh
export PATH="$PATH:$(go env GOPATH)/bin"
```

If you extract buf to ~/software/buf/bin, the following should work:

```sh
export PATH=$PATH:~/software/buf/bin/:~/go/bin
go get google.golang.org/protobuf/cmd/protoc-gen-go
go get google.golang.org/protobuf/cmd/protoc-gen-go-grpc
scripts/protobuf_codegen.sh
```

For more information, refer to the [GRPC Golang Quick Start Guide](https://grpc.io/docs/languages/go/quickstart/).

### Running protobuf codegen from docker

```sh
docker build -t odyssey:protobuf_codegen -f api/Dockerfile.buf .
docker run -t -i -v $(pwd):/opt/odyssey -w/opt/odyssey odyssey:protobuf_codegen bash -c "scripts/protobuf_codegen.sh"
```

### Running mock codegen

To regenerate the [gomock](https://github.com/uber-go/mock) code, run `scripts/mock.gen.sh` from the root of the repo.

This should only be necessary when modifying exported interfaces or after modifying `scripts/mock.mockgen.txt`.

## Versioning

### Version Semantics

OdysseyGo is first and foremost a client for the Odyssey network. The versioning of OdysseyGo follows that of the Odyssey network.

- `v0.x.x` indicates a development network version.
- `v1.x.x` indicates a production network version.
- `vx.[Upgrade].x` indicates the number of network upgrades that have occurred.
- `vx.x.[Patch]` indicates the number of client upgrades that have occurred since the last network upgrade.

### Library Compatibility Guarantees

Because OdysseyGo's version denotes the network version, it is expected that interfaces exported by OdysseyGo's packages may change in `Patch` version updates.

### API Compatibility Guarantees

APIs exposed when running OdysseyGo will maintain backwards compatibility, unless the functionality is explicitly deprecated and announced when removed.

## Supported Platforms

OdysseyGo can run on different platforms, with different support tiers:

- **Tier 1**: Fully supported by the maintainers, guaranteed to pass all tests including e2e and stress tests.
- **Tier 2**: Passes all unit and integration tests but not necessarily e2e tests.
- **Tier 3**: Builds but lightly tested (or not), considered _experimental_.
- **Not supported**: May not build and not tested, considered _unsafe_. To be supported in the future.

The following table lists currently supported platforms and their corresponding
OdysseyGo support tiers:

| Architecture | Operating system | Support tier  |
| :----------: | :--------------: | :-----------: |
|    amd64     |      Linux       |       1       |
|    arm64     |      Linux       |       2       |
|    amd64     |      Darwin      |       2       |
|    amd64     |     Windows      |       3       |
|     arm      |      Linux       | Not supported |
|     i386     |      Linux       | Not supported |
|    arm64     |      Darwin      | Not supported |

To officially support a new platform, one must satisfy the following requirements:

| OdysseyGo continuous integration | Tier 1  | Tier 2  | Tier 3  |
| ---------------------------------- | :-----: | :-----: | :-----: |
| Build passes                       | &check; | &check; | &check; |
| Unit and integration tests pass    | &check; | &check; |         |
| End-to-end and stress tests pass   | &check; |         |         |
