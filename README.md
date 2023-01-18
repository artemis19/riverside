# Riverside

*What is this tool meant to do?*

Riverside provides a web-based, dynamic network security visualization of real-time network flow data. Users can traverse time to watch how their network topology changes over a given time period, while being provided high level insights of their network's security posture.

![teaser](https://user-images.githubusercontent.com/21197485/191615126-ad53f8a4-55cf-491b-b991-f85e3488a318.png)

## Running Riverside

*__Important Note: Riverside is still being actively developed and tested in various environments. Session management has not yet been completely implemented or tested.__*

This tool uses gRPC functionality to communicate between a server and agent-installed hosts. Protocol buffers are used to structure and serialize data for agent and server communications. The agent and server binaries are written in Golang and have been tested successfully on Windows and various Linux architectures. The GORM library was used to handle all database functionality, and the supported database type is Postgres but can be easily changed with one line in the server source code. The client, or frontend, communicates with the server via the Gorilla WebSockets library to display batched data in a web-based network visualization.


## Frontend Visualization

The frontend layout was built using the [vis.js](https://visjs.org/) library along with custom Javascript code all contained within the `frontend` folder in this project. Agents, or internal hosts, will show up as cyan-colored nodes by default with remote, or external hosts, being grey. The timeline on the bottom controls what moment in time the visualization is displaying. Nodes will be connected with links that show whether traffic is to or from an agent node or bidirectional between two hosts. Communication is colored based on the protocol and edge thickness is correlated to the amount of traffic.

### Riverside Deployment

To run the frontend locally, you will just need to run a local web server, like the below, in the `frontend` folder:

```python
python3 -m http.server 9000
```

If you wish to host this live on a publicly-accessible web server, please follow the below instructions. This was done using a Let's Encrypt certificate on a Digital Ocean droplet. _Some changes may be needed if using a different CA or deployment method._

You will need to change `ListenAndServe` to `ListenAndServeTLS` in `server/websocket.go` to the following line:

```go
// log.Fatal(http.ListenAndServe(bind, nil))
log.Fatal(http.ListenAndServeTLS(bind, "/etc/letsencrypt/live/www.riverside-vis.dev/fullchain.pem", "/etc/letsencrypt/live/www.riverside-vis.dev/privkey.pem", nil))
```

I then had to change this line in `frontend/static/js/websocket_comms.js` to the following on my public droplet to enable TLS through the Websocket protocol:

```js
// let ws = new WebSocket("ws://localhost:8000/ws")
let ws = new WebSocket("wss://riverside-vis.dev:8000/ws")
```

## Riverside Binaries

You can download the latest pre-compiled binaries in the Releases section of this Github repository.

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

## Database

I use an ORM to handle database operations with the [GORM libray](https://gorm.io/). The following go dependencies may need to be installed before the server binary will compile correctly.

```go
go get -u gorm.io/gorm
go get -u gorm.io/driver/postgres
```

The driver for the database type can be changed depending on your needs or preferences in the `server/database` folder.

For local testing, you can set up a `postgres` Docker instance:

```sh
docker run -p 5432:5432 -e POSTGRES_USER=viz -e POSTGRES_PASSWORD=password postgres
```

## Server

The server reads from `server.yml`, and if one does not already exist, it will be created when the server is first run. An example is included in this repository as `server.yml.example`. All server logging will be saved in an `server.log` file in the `riverside/server` folder.

### Building & Running

Navigate to `riverside/server` and run the following:

```sh
go build
```

The `server` has the following commands with accompanying options when running:

```sh
./server -h
Server to store agent traffic

Usage:
  server [flags]
  server [command]

Available Commands:
  config      Show configuration settings
  help        Help about any command

Flags:
  -c, --configFile string      Location of config file to read from (default "server.yml")
  -a, --dbAddress string       Postgres database address (default "localhost")
  -p, --dbPassword string      Postgres database password (default "mysecretpassword")
  -b, --dbPort string          Postgres database port (default "5432")
  -s, --dbSSL                  Turn Postgres SSL mode on
  -u, --dbUsername string      Postgres database user (default "viz")
  -d, --debug                  Turn debug mode on
  -h, --help                   help for server
  -o, --outputFile string      Location of log file output (default "server.log")
  -l, --port string            Port for server to listen on (default "1533")
  -v, --version                version for server
  -w, --websocketPort string   Websocker server listening port (default "8000")

Use "server [command] --help" for more information about a command.
```

The server binary can be run with the following line as long as a database is currently running and correctly specified in the configuration file.

```sh
./server
```