# Riverside

*What is this tool meant to do?*

Riverside provides a web-based, dynamic network security visualization of real-time network flow data. Users can traverse time to watch how their network topology changes over a given time period, while being provided high level insights of their network's security posture.

![teaser](https://user-images.githubusercontent.com/21197485/191615126-ad53f8a4-55cf-491b-b991-f85e3488a318.png)

## Running Riverside

*__Important Note: Riverside's frontend visualization is still being developed and tested. I will be updating this repository in January of 2023.__*

This tool uses gRPC functionality to communicate between a server and agent-installed hosts. Protocol buffers are used to structure and serialize data for agent and server communications. The agent and server binaries are written in Golang and have been tested successfully on Windows and various Linux architectures. The GORM library was used to handle all database functionality, and the supported database type is Postgre but can be easily changed with one line in the server soruce code. The client, or frontend, communicates with the server via the Gorilla WebSockets library to display batched data in a web-based network visualization.

## Binaries

You can download the latest pre-compiled binaries here...

### Agent

Need to update.

#### Linux

Need to update.

#### Windows

Need to update.

### Server

Need to update.

