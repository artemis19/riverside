let ws = new WebSocket("ws://localhost:8000/ws")
console.log("Attempting websocket connection...");

initialized = false

// Sent upon initial WebSocket connection
ws.onopen = () => {
    console.log("Successfully connected to WebSocket.");
    // Send client info through websockets
    viz_select(ws, table = "hosts")
    viz_select(ws, table = "net_flow")
}

ws_message_queue = []

// Sent when client receives message from server
ws.onmessage = (event) => {
    var msg = JSON.parse(event.data);
    if (msg["to"] == FRONTEND_ID) {
        resolvePromise = ws_message_queue.pop()
        resolvePromise(msg)
    }
    switch (msg.type) {
        case "hosts":
            for (index in msg.data) {
                host = msg.data[index]
                addAgentNode(host)
            }
            initialized = true
            break
        case "netflow":
            if (initialized) {
                for (index in msg.data) {
                    netflowData = msg.data[index]
                    handleNetflow(netflowData)
                }
            }
            break
        case "host":
            if (initialized) {
                hostObj = msg.data
                addAgentNode(hostObj)
            }
            break
        default:
            console.log("Received WS message with unknown type: " + msg.type)
    }
}

// Sent when client closes connection
ws.onclose = (event) => {
    console.log("ws closed connection: ", event);
}

// Sent when error occurs in connection
ws.onerror = (error) => {
    console.log("ws error: ", error);
}