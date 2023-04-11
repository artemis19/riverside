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