let socket = new WebSocket("ws://localhost:8000/ws")
console.log("Attempting web socket connection...");

// Sent upon initial websocket connection
socket.onopen = () => {
    console.log("Successfully connected to websocket.");
    // Send client info through websocket
    viz_select(socket, table = "hosts")
}

// Sent when client receives message from server
socket.onmessage = (event) => {
    var msg = JSON.parse(event.data);
    switch (msg.type) {
        case "hosts":
            for (index in msg.data) {
                host = msg.data[index]
                addAgentNode(host)
            }
            break
        case "netflow":
            netflowData = msg.data
            handleNetflow(netflowData)
            break
        default:
            console.log("Received WS message with unknown type: " + msg.type)
    }
}

// Sent when client closes connection
socket.onclose = (event) => {
    console.log("Socket closed connection: ", event);
}

// Sent when error occurs in connection
socket.onerror = (error) => {
    console.log("Socket error: ", error);
}