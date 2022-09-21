# Riverside

*What is this tool meant to do?*

Riverside provides a web-based, dynamic network security visualization of your network through real-time network flow data. Users can traverse time to watch how their network topology changes over a given time period, while being provided high level insights of their network's security posture.

## Running Riverside

*__Important Note: Riverside's frontend visualization is still being developed and tested. I will be updating this repository in January of 2023.__*

This tool uses gRPC functionality to communicate between a server and agent-installed hosts. Protocol buffers are used to structure and serialize the data from the agents to the server. All agent and server binaries are written in Golang and have been tested successfully on Windows and various Linux environments. The GORM library was used to handle all database functionality. The database used for testing and initial deployments was a Postgres docker container. The client, or frontend in this case, communicates with the server via the Gorilla Websocket library to display batched network in a web-based network visualization.

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

