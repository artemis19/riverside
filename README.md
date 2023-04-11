# Riverside

*What is this tool meant to do?*

Riverside provides a web-based, dynamic network security visualization of real-time network flow data. Users can traverse time to watch how their network topology changes over a given time period, while being provided high level insights of their network's security posture.

![teaser](https://user-images.githubusercontent.com/21197485/191615126-ad53f8a4-55cf-491b-b991-f85e3488a318.png)

## Running Riverside

*__Important Note: Riverside is still being actively developed and tested in various environments. Session management has not yet been completely implemented or tested.__*

This tool uses gRPC functionality to communicate between a server and agent-installed hosts. Protocol buffers are used to structure and serialize data for agent and server communications. The agent and server binaries are written in Golang and have been tested successfully on Windows and various Linux architectures. The GORM library was used to handle all database functionality, and the supported database type is Postgres but can be easily changed with one line in the server source code. The client, or frontend, communicates with the server via the Gorilla WebSockets library to display batched data in a web-based network visualization.


## Frontend Visualization

The frontend layout was built using the [vis.js](https://visjs.org/) library along with custom Javascript code all contained within the `frontend` folder in this project. Agents, or internal hosts, will show up as cyan-colored nodes by default with remote, or external hosts, being grey. The timeline on the bottom controls what moment in time the visualization is displaying. Nodes will be connected with links that show whether traffic is to or from an agent node or bidirectional between two hosts. Communication is colored based on the protocol and edge thickness is correlated to the amount of traffic.

*__An initial pop-up will appear on the frontend for logging in to a server or registering a user, but this has not been fully implemented. If you wish to remove it, you can comment out lines 44-123 in `frontend/index.html` or simply click out of it upon loading the page.__*

## Riverside Setup

All of the instructions for Riverside's setup and deployment are located in the [Riverside Wiki](https://github.com/artemis19/riverside/wiki).