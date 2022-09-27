# Riverside

*What is this tool meant to do?*

Riverside provides a web-based, dynamic network security visualization of real-time network flow data. Users can traverse time to watch how their network topology changes over a given time period, while being provided high level insights of their network's security posture.

![teaser](https://user-images.githubusercontent.com/21197485/191615126-ad53f8a4-55cf-491b-b991-f85e3488a318.png)

## Running Riverside

*__Important Note: Riverside's frontend visualization is still being developed and tested. I will be updating this repository in January of 2023.__*

This tool uses gRPC functionality to communicate between a server and agent-installed hosts. Protocol buffers are used to structure and serialize data for agent and server communications. The agent and server binaries are written in Golang and have been tested successfully on Windows and various Linux architectures. The GORM library was used to handle all database functionality, and the supported database type is Postgres but can be easily changed with one line in the server soruce code. The client, or frontend, communicates with the server via the Gorilla WebSockets library to display batched data in a web-based network visualization.

## Binaries

You can download the latest pre-compiled binaries in the Releases section of this Github repository.

## Agent

Riverside uses Golang agents deployed on internal network hosts to collect traffic. By default, the agent is set to listen on all interfaces unless otherwise specified. All agent logging will be saved in an `agent.log` file in the `riverside/agent` folder.

### Configuration

The agent reads from `agent.yml`, and if one does not already exist, it will be created when the agent is first run. An example is included in this repository as `agent.yml.example`.

If you would like to compile the binaries natively, please follow the below instructions depending on the OS.

### Linux

There are some dependencies that will be required to compile the agent natively on your host.

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

I use the [go-dpi library](https://github.com/mushorg/go-dpi/wiki/Installation-guide) for this project on top of [gopacket](https://github.com/google/gopacket). Because of this, there are some extra dependencies if you wish to compile this on your own versus using the pre-compiled binaries provided. DPI is only supported for Linux agents at the moment.

```sh
sudo bash -c 'echo "deb http://packages.wand.net.nz trusty main" | tee -a /etc/apt/sources.list'
sudo bash -c 'curl https://packages.wand.net.nz/keyring.gpg -o /etc/apt/trusted.gpg.d/wand.gpg'

sudo apt update
sudo apt -y install libflowmanager
sudo apt -y install libtrace4 libtrace4-dev
sudo apt -y install liblinear4 liblinear-dev libprotoident libprotoident-dev autoconf libtool

git clone --branch 3.2-stable https://github.com/ntop/nDPI/ /tmp/nDPI
cd /tmp/nDPI && ./autogen.sh && ./configure && make && sudo make install && cd -
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

You will likely need root privileges to run the agent as it needs permission to capture traffic on the host interfaces.

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

To run the agent on Windows, you will likely need administrator privileges to capture traffic on the host interfaces.

### Database

I use an ORM to handle database operations with the [GORM libray](https://gorm.io/). The following go dependencies may need to be installed before the server binary will compile correctly.

```go
go get -u gorm.io/gorm
go get -u gorm.io/driver/postgres
```

The driver for the database type can be changed depending on your needs or preferences in the `server/database` folder.

### Server

The server reads from `server.yml`, and if one does not already exist, it will be created when the server is first run. An example is included in this repository as `server.yml.example`. All server logging will be saved in an `server.log` file in the `riverside/server` folder.

### Building & Running

Navigate to `riverside/server` and run the following:

```sh
go build
```

The server binary can be run with the following line as long as a database is currently running and specified in the configuration file.

```sh
./server
```