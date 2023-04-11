## Agent

Riverside uses Golang agents deployed on internal network hosts to collect traffic. All agent logging will be saved in an `agent.log` file in the `riverside/agent` folder.

### Configuration

The agent reads from `agent.yml`, and if one does not already exist, it will be created when the agent is first run. An example is included in this repository as `agent.yml.example`.

If you would like to compile the binaries natively, please follow the below instructions depending on the OS.

The `agent` has the following commands with accompanying options when running:

```sh
./agent -h
Local traffic capture agent

Usage:
  agent [command]

Available Commands:
  config      Show configuration settings
  help        Help about any command
  listen      Listen for traffic
  ping        Test connectivity to collection server

Flags:
  -d, --debug                    Turn debug mode on
  -h, --help                     help for agent
  -p, --operatingSystem string   Force operating system for debugging purposes
  -s, --serverAddress string     Server to connect to
  -v, --version                  version for agent

Use "agent [command] --help" for more information about a command.
```

As a note, I listen listen on all interfaces by default, but if you uncomment this line, you will only listen on the primary interfaces.

```go
// This was removes virtual, docker, and loopback interfaces.
if strings.HasPrefix(interfaceName, "veth") || strings.HasPrefix(interfaceName, "docker") || strings.HasPrefix(interfaceName, "lo") || strings.HasPrefix(interfaceName, "vmnet") || strings.HasPrefix(interfaceName, "br-") {
 // Ignore virtual, docker, and loopback interfaces
 continue
}
```

### Linux

There are some dependencies that will be required to compile the agent natively on your host.

#### Dependencies

```sh
sudo apt update
sudo apt install build-essential libpcap-dev
```

Install the following two `go` packages and export your `GOPATH`:

```sh
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
export PATH="$PATH:$(go env GOPATH)/bin"
```

Add this line to your `bashrc` file:

```sh
export PATH="$PATH:$(go env GOPATH)/bin"
```

#### DPI Dependencies

I use the [go-dpi library](https://github.com/mushorg/go-dpi/wiki/Installation-guide) for this project on top of [gopacket](https://github.com/google/gopacket). Because of this, there are some extra dependencies if you wish to compile this on your own versus using the pre-compiled binaries provided. Currently, DPI is only supported for Linux agents.

```sh
curl -1sLf 'https://dl.cloudsmith.io/public/wand/libwandio/cfg/setup/bash.deb.sh' | bash
curl -1sLf 'https://dl.cloudsmith.io/public/wand/libwandder/cfg/setup/bash.deb.sh' | bash
curl -1sLf 'https://dl.cloudsmith.io/public/wand/libtrace/cfg/setup/bash.deb.sh' | bash
curl -1sLf 'https://dl.cloudsmith.io/public/wand/libflowmanager/cfg/setup/bash.deb.sh' | bash
curl -1sLf 'https://dl.cloudsmith.io/public/wand/libprotoident/cfg/setup/bash.deb.sh' | bash

sudo apt update
sudo apt install -y liblinear4 liblinear-dev libtrace4-dev libtrace4 autoconf libtool git make build-essential
git clone --branch 3.2-stable https://github.com/ntop/nDPI/ /tmp/nDPI
cd /tmp/nDPI && ./autogen.sh && ./configure && make && make install && cd -

git clone https://github.com/LibtraceTeam/libflowmanager /opt/libflowmanager
cd /opt/libflowmanager
./bootstrap.sh
./configure
make
make install

git clone https://github.com/LibtraceTeam/libprotoident /opt/libprotoident
cd /opt/libprotoident
./bootstrap.sh
./configure
make
make install
```

When building, you will need to include the specific header file for the `liblinear` dependency:

```go
CGO_CFLAGS="-I/usr/include/liblinear" go build
```

#### Building & Running

To build the agent on a Linux host, navigate to `riverside/agent` and run the following line:

```sh
go build
```

Use the below syntax to run the agent and listen on a Linux host's network interfaces, where `-s` specifies the server IP and port.

```sh
./agent listen -s localhost:1533
```

You will need root privileges to run the agent as it needs permission to capture traffic on the host interfaces.

### Windows

The Windows agent does not support DPI at this moment in time but does require some additional dependencies before it will compile.

#### Dependencies

For installing `protoc` command on Windows, download the Windows binary [here](https://github.com/protocolbuffers/protobuf/releases/tag/v3.20.1). Then add the binary to your PATH environment variable.

To compile the protocol buffer file using `protoc`: 

```sh
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative viz.proto
```

#### Building & Running

To build the agent on a Windows host, navigate to `riverside\agent` and run the following line:

```sh
GOOS=windows GOARCH=amd64 go build
```

Use the below syntax to run the agent and listen on a Windows host's network interfaces, where `-s` specifies the server IP and port.

```sh
./agent listen -s localhost:1533
```

To run the agent on Windows, you will need administrator privileges to capture traffic on the host interfaces.