<div align="center">
  <img src="resources/OdysseyLogoRed.png?raw=true">
</div>

---

Node implementation for the [Odyssey](https://dione.network) network -
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

### Binary Repository

Install OdysseyGo using an `apt` repository.

#### Adding the APT Repository

If you already have the APT repository added, you do not need to add it again.

To add the repository on Ubuntu, run:

```sh
sudo su -
wget -qO - https://downloads.dione.network/odysseygo.gpg.key | tee /etc/apt/trusted.gpg.d/odysseygo.asc
source /etc/os-release && echo "deb https://downloads.dione.network/apt $UBUNTU_CODENAME main" > /etc/apt/sources.list.d/odyssey.list
exit
```

#### Installing the Latest Version

After adding the APT repository, install `odysseygo` by running:

```sh
sudo apt update
sudo apt install odysseygo
```

### Binary Install

Download the [latest build](https://github.com/DioneProtocol/odysseygo/releases/latest) for your operating system and architecture.

The Odyssey binary to be executed is named `odysseygo`.

### Docker Install

Make sure Docker is installed on the machine - so commands like `docker run` etc. are available.

Building the Docker image of latest `odysseygo` branch can be done by running:

```sh
./scripts/build_image.sh
```

To check the built image, run:

```sh
docker image ls
```

The image should be tagged as `avaplatform/odysseygo:xxxxxxxx`, where `xxxxxxxx` is the shortened commit of the Odyssey source it was built from. To run the Odyssey node, run:

```sh
docker run -ti -p 9650:9650 -p 9651:9651 avaplatform/odysseygo:xxxxxxxx /odysseygo/build/odysseygo
```

## Running Odyssey

### Connecting to Mainnet

To connect to the Odyssey Mainnet, run:

```sh
./build/odysseygo
```

You should see some pretty ASCII art and log messages.

You can use `Ctrl+C` to kill the node.

### Connecting to Fuji

To connect to the Fuji Testnet, run:

```sh
./build/odysseygo --network-id=fuji
```

### Creating a Local Testnet

See [this tutorial.](https://docs.dione.network/build/tutorials/platform/create-a-local-test-network/)

## Bootstrapping

A node needs to catch up to the latest network state before it can participate in consensus and serve API calls. This process (called bootstrapping) currently takes several days for a new node connected to Mainnet.

A node will not [report healthy](https://docs.dione.network/build/odysseygo-apis/health) until it is done bootstrapping.

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

## Security Bugs

**We and our community welcome responsible disclosures.**

Please refer to our [Security Policy](SECURITY.md) and [Security Advisories](https://github.com/DioneProtocol/odysseygo/security/advisories).
