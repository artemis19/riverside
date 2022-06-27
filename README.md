# Riverside - Dynamic Network Map Visualization

Riverside is an open-source network visualization tool operating from inside the network, showcasing live traffic between internal hosts and external remote hosts in a real-time network graph. While capturing netflow and packet information inside of a database, users can traverse backwards in time to analyze previous network activity for enriched situational awareness and a thorough understanding of their network security posture. This utility supplements existing tooling to provide more insight for use cases such as incident response, analysis and investigation, and identification of true assets used within a network environment.

This tool uses gRPC functionality to communicate between a server and hosts with installed agents. The data being transferred is handled using protocol buffers, and agents can function within Windows and Linux environments. MacOS has not been developed at this point. Additionally, GORM was used to handle all database functionality, which is a Golang ORM library. The database that was used for testing and initial deployments was Postgres from a docker container, but database arguments can be provided to the server configuration file if an alternative storage backend is preferred. The client, or `frontend` in this case, communicates with the server via the Gorilla Websocket library. The frontend was created primarily using the vis.js, community-supported library.

## Key Features

*To be added*

## Visualization Frontend

* [vis.js](https://visjs.org/)

This is contained within the `frontend` folder in this project.

## Websocket Communication

* [Gorilla Websockets](https://github.com/gorilla/websocket)

This is contained within the `server/websocket` folder in this project.

## Infrastructure

The infrastructure for this tool is split into three parts: a database, server, and agents. The database stores the netflow information being received by the server. Agents are installed on hosts within the network and send their network traffic to the server.

### Database

I am using an ORM library called [GORM](https://gorm.io/) to manipulate the data being stored in my database.

```go
go get -u gorm.io/gorm
go get -u gorm.io/driver/postgres
```

The driver for the database type can be changed depending on your needs or preferences. We are using a temporary Postgres database run inside a docker container and creating the user `viz` and subsequent database for that user:

```sh
docker run -p 5432:5432 -e POSTGRES_USER=viz -e POSTGRES_PASSWORD=password postgres
```

Hosts have a `has many` relationship with network interfaces.

### Postgres Database

I am currently using a docker instance of a local postgres database. To interact with my local database:

```sh
PGPASSWORD=password psql -h localhost -U viz
```

### Agent

_Insert blurb about what this is at some point_

```go
go mod init

go mod tidy
```

#### Dependencies

```bash
sudo apt update
sudo apt install build-essential libpcap-dev
```

#### Installing

For Linux:

```shell
go build
```

For Windows:

```shell
GOOS=windows GOARCH=amd64 go build
```

#### Configuration

The agent reads from a `agent.yml` if present. An example is included in this repository as `agent.yml.example`:

```yaml
configfilename: agent.yml
debug: false
filter: ""
interface: eth0
length: "262144"
outputfile: agent.log
serveraddress: ""
```

#### Usage

```shell
./agent

Local traffic capture agent

Usage:
  agent [command]

Available Commands:
  config      Show configuration settings
  help        Help about any command
  listen      Listen for traffic
  ping        Test connectivity to collection server

Flags:
  -c, --configFilename string   Location of config file to read from (default "./agent.yml")
  -d, --debug                   Toggle debug mode on
  -f, --filter string           filter for packet capture, defaults to nothing
  -h, --help                    help for agent
  -i, --interface string        Network interface to listen on
  -l, --length string           snap length for packet size, default is 262144 (default "262144")
  -o, --outputFile string       Location of output log file to write to (default "agent.log")
  -s, --serverAddress string    Server to connect to (host:port format)
  -v, --version                 version for agent

Use "agent [command] --help" for more information about a command.
```

#### Debug Mode

In debug mode, all logging will be sent to stdout as well as the configured log file (`agent.log` by default)

```sh
./agent -d
```

#### Capturing Traffic

**There must be a running server for the agent to listen and you must have administrator or root privileges to capture traffic.**

To listen on a specific interface:

```bash
./agent listen -i eth0
```

#### Testing Dummy Interfaces

To test new interfaces being aded, create a dummy interface with an IP address like the example below:

```sh
sudo nmcli connection add type dummy ifname dummy0 ipv4.method manual ipv4.addresses 192.0.1.1/24

sudo ip link del dummy0
```

### Server

_Insert blurb about what this is at some point_


#### Dependencies

```bash
sudo apt install build-essentials libpcap-dev
```

#### Installing

For Linux:

```shell
go build
```

For Windows:

```shell
GOOS=windows GOARCH=amd64 go build
```

#### Configuration

The agent reads from a `server.yml` if present. An example YAML file is below:

```yaml
configfilename: server.yml
debug: false
outputfile: server.log
port: "1533"
websocketport: "8000"
```

## Setting up Protocol Buffers & gRPC

The backbone of the communication for this tool is located in the `rpc` and `pb` folders.

### Linux

```sh
sudo apt install protobuf-compiler
```

Install the following two `go` packages and export your `GOPATH`:

```go
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
export PATH="$PATH:$(go env GOPATH)/bin"
```

Add this line to your `bashrc` file:

```sh
export PATH="$PATH:$(go env GOPATH)/bin"
```

### Windows

For installing `protoc` command on Windows, download the Windows binary [here](https://github.com/protocolbuffers/protobuf/releases/tag/v3.20.1). Then add the binary to your PATH environment variable.

Compiling the protocol buffer & gRPC:

```shell
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative viz.proto
```

## Current TO-DO:

* Go binary static linking due to `gopacket` requiring `CGO`
* Fix the interfaceName that's displayed on the agent side (specifically for Windows)
  * Shows `\Device\NPF_{60E38DFC-F320-4F5E-98C5-935AB32DD21D}` as an example since it's Windows
* Filter out database & server traffic when database is on a host separate from the server
* Build out `SerializePacket()` functionality