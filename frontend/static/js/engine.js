// Create network container
var container = document.getElementById('viz');

// Create DataSet objects for individual entities
var agentNodes = new vis.DataSet([])
var remoteNodes = new vis.DataSet([])
// Create DataSet objects for groups of entities
var allNodes = new vis.DataSet([])
var allEdges = new vis.DataSet([])

// Create DataPipe objects to funnel data
var agentPipe = new vis.createNewDataPipeFrom(agentNodes).to(allNodes)
var remotePipe = new vis.createNewDataPipeFrom(remoteNodes).to(allNodes)

// Initiate pipes
agentPipe.all().start()
remotePipe.all().start()

// Time we will be using to update visualization based on cursor drag
displayTime = new Date()

// Variable for determining whether to alter display time
previousTimeTick = new Date((parseInt((displayTime).getTime() / 1000)) * 1000)

//-------------------------------
// DataView filters
function nodeTimeFilter(node) {

    node.nodeCurrentTime = new Date((parseInt((displayTime).getTime()) / 1000) * 1000)
    if (node.nodeCurrentTime != node.nodePreviousTime) {
        node.nodePreviousTime = node.nodeCurrentTime
        withinTimeWindow = false

        for ([nodeStartTime, nodeEndTime] of node.timeWindows) {
            endTimeClone = new Date(nodeEndTime.getTime())
            // Look one second "into the future"
            endTimeWindow = new Date(endTimeClone.setSeconds(endTimeClone.getSeconds() + 2))
            // Look through each time window to see if we are in one
            if (displayTime > nodeStartTime && displayTime < endTimeWindow) {
                withinTimeWindow = true
                break
            }
        }
        node.withinTimeWindow = withinTimeWindow
        return withinTimeWindow
    }
    return node.withinTimeWindow
}

const edgeWidthCap = 15
const nodeOpacityMin = 0.3

function edgeTimeFilter(edge) {
    // Set each edge's current time to the current display time to the latest second
    edge.edgeCurrentTime = new Date((parseInt((new Date()).getTime() / 1000)) * 1000)
    if (edge.edgeCurrentTime != edge.edgePreviousTime) {
        edge.edgePreviousTime = edge.edgeCurrentTime

        node = remoteNodes.get(edge.to)
        isRemoteNode = true
        isAgentNode = false

        if (node == undefined) {
            node = agentNodes.get(edge.to)
            isRemoteNode = false
            isAgentNode = true
            focusObject = edge
        }

        // Force this so that we are correctly gathering edges in time
        focusObject = edge

        withinTimeWindow = false

        // Retrieve previous color
        color = parseRGBA(node.color)
        r = color[0]
        g = color[1]
        b = color[2]
        a = color[3]
        opacityChange = 0

        for ([nodeStartTime, nodeEndTime] of focusObject.timeWindows) {
            if (nodeStartTime > displayTime) {
                withinTimeWindow = false
                break
            }
            endTimeClone = new Date(nodeEndTime.getTime())
            // Look one second "into the future"
            endTimeWindow = new Date(endTimeClone.setSeconds(endTimeClone.getSeconds() + 2))
            // Look through each time window to see if we are in one
            if (displayTime > nodeStartTime && displayTime < endTimeWindow) {
                withinTimeWindow = true
                before = nodeStartTime
                after = displayTime
                secondsChange = Math.floor((after.getTime() - before.getTime()) / 1000)
                // Throughput traffic needs to be doubled for bidirectional
                if (edge.direction == "bidirectional") {
                    edge.throughput = secondsChange * 2
                } else {
                    edge.throughput = secondsChange
                }
                // Update title with new throughput value
                edge.startTime = nodeStartTime
                // Direction may change depending on the conversation, 
                // and who "started" the conversation.
                // These need to be set in a "Map" based off our time windows
                // when creating the edge.
                edge.originAddress = edge.originWindows.get(before.toString())
                edge.direction = edge.directionWindows.get(before.toString())
                
                if (edge.direction == "to") {
                    edge.arrows = { to: true }
                }
                if (edge.direction == "from") {
                    edge.arrows = { from: true }
                }
                if (edge.direction == "bidirectional") {
                    edge.arrows = { to: true, from: true }
                }
                
                edge.endTime = nodeEndTime
                edge.title = edgeHover(edge)
                
                if (isRemoteNode) {
                    // Update node transparency through time
                    finalOpacity = opacityChange + secondsChange
                    newAlpha = 1 - (opacityChange / 100)
                    if (newAlpha <= nodeOpacityMin) {
                        newAlpha = nodeOpacityMin
                    }
                    if (newAlpha != a) {
                        node.color = "rgba(" + r + "," + g + "," + b + "," + newAlpha + ")"
                        remoteNodes.update(node)
                        opacityChange = finalOpacity
                    }
                }
                if (secondsChange >= edgeWidthCap) {
                    secondsChange = edgeWidthCap
                }
                // Fake throughput in past time
                if (edge.width != secondsChange) {
                    edge.width = secondsChange
                    allEdges.update(edge)
                }
                break
            } else {
                // Look over previous time to add to the opacity delta
                opacityChange += Math.floor((nodeEndTime.getTime() - nodeStartTime.getTime()) / 1000)
            }
        }
        edge.withinTimeWindow = withinTimeWindow
        return withinTimeWindow
    }
    return edge.withinTimeWindow
}

//--------------------------------------
// Create DataView to populate the visualization based on time
nodeTimeView = new vis.DataView(allNodes, { filter: nodeTimeFilter })
edgeTimeView = new vis.DataView(allEdges, { filter: edgeTimeFilter })

var networkOptions = {
    width: "100%",
    height: "100%",
    clickToUse: false,
    physics: {
        enabled: true,
        barnesHut: {
            theta: 0,
            gravitationalConstant: 200,
            centralGravity: 0,
            springLength: 400,
            springConstant: 1,
            damping: 1,
            avoidOverlap: 1
        },
        forceAtlas2Based: {
            theta: 0,
            gravitationalConstant: 0,
            centralGravity: 0,
            springLength: 200,
            springConstant: 1,
            damping: 1,
            avoidOverlap: 1
        },
        repulsion: {
            centralGravity: 0,
            springLength: 0,
            springConstant: 1,
            nodeDistance: 100,
            damping: 1
        },
        hierarchicalRepulsion: {
            centralGravity: 0.0,
            springLength: 200,
            springConstant: 1,
            nodeDistance: 0,
            damping: 1,
            avoidOverlap: 1
        },
        maxVelocity: 100,
        minVelocity: 0.1,
        solver: 'barnesHut',
        stabilization: {
            enabled: true,
            iterations: 1000,
            updateInterval: 100,
            onlyDynamicEdges: false,
            fit: true
        },
        timestep: 0.7,
        adaptiveTimestep: true,
    }
};

// Initialize network
var network = new vis.Network(container, { nodes: nodeTimeView, edges: edgeTimeView },
    networkOptions);

// Initial zoom
network.moveTo({ scale: 0.5 })

//-------------------------------------------

// Create window of time (1 minute)
startTime = new Date();
startWindow = new Date();
startWindow = new Date(startWindow.setSeconds(startTime.getSeconds() + 27));
timelineOptions = {
    height: "110px",
    start: startTime,
    end: startWindow,
    loadingScreenTemplate: function() {
        return '<p>Loading timeline...</p>';
    }
}

var items = new vis.DataSet([]);

var timelineContainer = document.getElementById("timeline")
var timeline = new vis.Timeline(timelineContainer, items, timelineOptions)

cursorID = "cursor"
timeline.addCustomTime(timeline.getCurrentTime(), cursorID)
timeline.setCustomTimeMarker('', cursorID, true)
var withTime = true
var playForward = false
var playDistance = 0

// Keeps timeline bar up to current time when "playing"
timeline.on("currentTimeTick", function() {
    if (withTime) {
        timeline.setCustomTime(timeline.getCurrentTime().getTime(), cursorID)
        handleTime()
        // Move timeline with real-time when "playing"
        if (timeline.getCurrentTime().getTime() > timeline.range.end - 1000) {
            distance = timeline.range.end - timeline.range.start
            timeline.setWindow(timeline.range.start + 10000, timeline.getCurrentTime().getTime() + 10000), { animation: "linear" }
        }
    } else {
        if (playForward) {
            if (playDistance == 0) {
                playDistance = timeline.getCurrentTime().getTime() - timeline.getCustomTime(cursorID).getTime()
            }
            diff = timeline.getCurrentTime().getTime() - playDistance
            timeline.setCustomTime(diff, cursorID)
            handleTime()
        }
    }
})

// If we drag our timeline cursor, trigger an event
timeline.on("timechange", function(data) {
    //console.log("time change called")
    id = data.id
    time = data.time
    icon = $('#play-pause').children()
    if (id == cursorID) {
        if (timelineIsPlaying()) {
            pauseTimeline(icon)
        }
    }

    if (id == cursorID) {
        //console.log("time changed")
        if (time > timeline.getCurrentTime().getTime()) {
            timeline.setCustomTime(timeline.getCurrentTime().getTime(), cursorID)
        } else {
            handleTime()
        }
    }
})

//-------------------------------------------
// Convenience functions

// Get time that our cursor is at
// This is run frequently due to nanoseconds/milliseconds precision and we only care about seconds
// Could validate and call less often
function handleTime() {
    newTime = timeline.getCustomTime(cursorID)

    // displayTime = newTime
    // d = displayTime

    // Cutting new time for our timeline cursor to only seconds precision
    displayTime = new Date((parseInt((newTime).getTime() / 1000)) * 1000)
    if ((previousTimeTick.getTime() / 1000) != (displayTime.getTime() / 1000)) {

        previousTimeTick = displayTime
        d = displayTime

        $('input[type=datetime-local]').val(
            d.getFullYear() + "-" +
            zeroPadded(d.getMonth() + 1) + "-" +
            zeroPadded(d.getDate()) + "T" +
            zeroPadded(d.getHours()) + ":" +
            zeroPadded(d.getMinutes()) + ":" +
            zeroPadded(d.getSeconds())
        );

        // Look back 5 minutes
        // currentTime = new Date(displayTime.getTime())
        // pastTimeWindow = new Date(currentTime.setMinutes(currentTime.getMinutes() - 5))
        // // Look forward 5 minutes, if there is data that far ahead
        // currentTime = new Date(displayTime.getTime())
        // futureTimeWindow = new Date(currentTime.setMinutes(currentTime.getMinutes() + 5))

        nodeTimeView.refresh()
        edgeTimeView.refresh()
    }
}

// Places our nodes in random positions on the visualization
// so you can see everything and no overlap occurs
function random(min, max) {
    min = Math.ceil(min);
    max = Math.floor(max);
    return Math.floor(Math.random() * (max - min) + min);
}

// Choose random value in a list
function choose(choices) {
    var index = Math.floor(Math.random() * choices.length);
    return choices[index];
}

// Get the alpha value of the node each time
function parseRGBA(input) {
    return input.split("(")[1].split(")")[0].split(",");
}

// Color picker for changing colors needs to be converted
function hexToRGBA(hex, alpha) {
    var r = parseInt(hex.slice(1, 3), 16),
        g = parseInt(hex.slice(3, 5), 16),
        b = parseInt(hex.slice(5, 7), 16);

    if (alpha) {
        return "rgba(" + r + "," + g + "," + b + "," + alpha + ")"
    } else {
        return "rgba(" + r + "," + g + "," + b + ")";
    }
}

//-------------------------------------------

function addAgentNode(hostObj) {
    // Create maps of agent nodes to keep track of previous traffic w/same remote node
    pastAgentNode = agentNodes.get(hostObj.machine_id)

    agentNode = {
        id: hostObj.machine_id,
        label: hostObj.hostname,
        shape: 'circle',
        size: 50,
        margin: 10,
        font: {
            size: 20,
        },
        // Custom properties added
        timeWindows: new Map(),
        startTime: new Date(hostObj.created_at),
        endTime: "",
        createdAt: hostObj.created_at,
        updatedAt: hostObj.updated_at,
        hostname: hostObj.hostname,
        os: hostObj.os,
        arch: hostObj.arch,
        networkInterfaces: hostObj.network_interfaces,
        nodeCurrentTime: new Date((parseInt((new Date()).getTime() / 1000)) * 1000),
        nodePreviousTime: new Date((parseInt((new Date()).getTime() / 1000)) * 1000),
        withinTimeWindow: true
    }

    // Setting the title enables hover by default
    agentNode.title = agentNodeHover(agentNode)

    // Handle agent end time
    if (hostObj.deleted_at != undefined) {
        agentNode.endTime = hostObj.deleted_at
    } else {
        endDate = new Date()
        endDate = new Date(endDate.setFullYear(startTime.getFullYear() + 100));
        agentNode.endTime = endDate
    }

    // If we have seen this node before, keep track of its past time windows
    if (pastAgentNode != null) {
        agentNode.timeWindows = pastAgentNode.timeWindows
    }

    // Update time windows for the current node
    agentNode.timeWindows.set(agentNode.startTime, agentNode.endTime)

    // Make agent nodes appear away from other agents
    otherAgentNodes = agentNodes.getIds()
    if (otherAgentNodes.length != 0) {
        otherNodeID = otherAgentNodes[0]
        otherNodeID = otherAgentNodes[otherAgentNodes.length - 1]
        agentNode.x = network.getPosition(otherNodeID).x + 200
        if (agentNodes.length % 2) {
            agentNode.y = network.getPosition(otherNodeID).y + random(150, 400)
        } else {
            agentNode.y = network.getPosition(otherNodeID).y - random(150, 400)
        }
        // agentNode.x = network.getPosition(otherNodeID).x + (choose([-1, 1]) * random(400, 600))
        // agentNode.y = network.getPosition(otherNodeID).y + (choose([-1, 1]) * random(400, 600))
    }

    // Adds new node to agentNodes DataSet
    previous = agentNodes.get(hostObj.machine_id)
    if (previous == null) {
        agentNode.color = "rgba(100,250,255,1)"
        agentNodes.add(agentNode)
    }
}

//-------------------------------------------

function handleNetflow(netflowObj) {
    // Create maps of remote nodes to keep track of previous traffic w/same remote node
    pastRemoteNode = remoteNodes.get(netflowObj.remote_host_id)

    remoteNode = {
        id: netflowObj.remote_host_id,
        size: 50,
        // Custom properties added
        timeWindows: new Map(),
        startTime: new Date(netflowObj.start_time),
        endTime: new Date(netflowObj.end_time),
        x: network.getPosition(netflowObj.host_id).x + (choose([-1, 1]) * random(100, 400)),
        y: network.getPosition(netflowObj.host_id).y + (choose([-1, 1]) * random(100, 400)),
        createdAt: netflowObj.created_at,
        updatedAt: netflowObj.updated_at,
        ipAddress: netflowObj.dst_address,
        nodeCurrentTime: new Date((parseInt((new Date()).getTime() / 1000)) * 1000),
        nodePreviousTime: new Date((parseInt((new Date()).getTime() / 1000)) * 1000),
        withinTimeWindow: true
    }

    // Setting the title enables hover by default
    // remoteNode.title = remoteNodeHover(remoteNode)

    // If we have seen this node before, keep track of its past time windows
    if (pastRemoteNode != null) {
        remoteNode.timeWindows = pastRemoteNode.timeWindows
        remoteNode.x = pastRemoteNode.x
        remoteNode.y = pastRemoteNode.y
    }

    // Update time windows for the current node
    remoteNode.timeWindows.set(remoteNode.startTime, remoteNode.endTime)

    newEdge = {
        id: netflowObj.host_id + "->" + netflowObj.remote_host_id,
        to: netflowObj.remote_host_id,
        from: netflowObj.host_id,
        arrows: {
            to: true,
        },
        physics: false,
        timeWindows: new Map(),
        // Hover fields
        startTime: new Date(netflowObj.start_time),
        endTime: new Date(netflowObj.end_time),
        originWindows: new Map(),
        directionWindows: new Map(),
        originAddress: netflowObj.origin_address,
        srcAddress: netflowObj.src_address,
        dstAddress: netflowObj.dst_address,
        srcPort: netflowObj.src_port,
        dstPort: netflowObj.dst_port,
        direction: netflowObj.direction,
        protocol: netflowObj.protocol,
        edgeCurrentTime: new Date((parseInt((new Date()).getTime() / 1000)) * 1000),
        edgePreviousTime: new Date((parseInt((new Date()).getTime() / 1000)) * 1000),
        withinTimeWindow: true
    }

    // Setting the title enables hover by default
    newEdge.title = edgeHover(newEdge);

    // If we have seen this edge before, keep track of time windows
    pastEdge = allEdges.get(newEdge.id)
    if (pastEdge != null) {
        newEdge.timeWindows = pastEdge.timeWindows
        newEdge.originWindows = pastEdge.originWindows
        newEdge.directionWindows = pastEdge.directionWindows
    }

    // Update time windows for this edge
    newEdge.timeWindows.set(newEdge.startTime, newEdge.endTime)
    newEdge.originWindows.set(newEdge.startTime.toString(), netflowObj.origin_address)
    newEdge.directionWindows.set(newEdge.startTime.toString(), netflowObj.direction)

    // Add arrows in correct directions
    if (netflowObj.direction == "to") {
        newEdge.arrows = { to: true }
    }
    if (netflowObj.direction == "from") {
        newEdge.arrows = { from: true }
    }
    if (netflowObj.direction == "bidirectional") {
        newEdge.arrows = { to: true, from: true }
    }

    // Color edges by protocol
    if (newEdge.protocol == "SSH") {
        newEdge.color = "salmon"
    }
    if (newEdge.protocol == "FTP") {
        newEdge.color = "orange"
    }
    if (newEdge.protocol == "HTTP") {
        newEdge.color = "green"
    }
    if (newEdge.protocol == "SSL") {
        newEdge.color = "darkgreen"
    }
    if (newEdge.protocol == "RDP") {
        newEdge.color = "blue"
    }
    if (newEdge.protocol == "SMB") {
        newEdge.color = "cyan"
    }
    if (newEdge.protocol == "SMTP") {
        newEdge.color = "lightsteelblue"
    }
    if (newEdge.protocol == "TCP") {
        newEdge.color = "silver"
    }
    if (newEdge.protocol == "UDP") {
        newEdge.color = "maroon"
    }
    if (newEdge.protocol == "ICMP") {
        newEdge.color = "gold"
    }
    if (newEdge.protocol == "N/A") {
        newEdge.color = "lavender"
    }
    if (newEdge.protocol == "DNS") {
        newEdge.color = "slateblue"
    }

    // Add remote nodes to DataSet if it doesn't already exist
    isAgent = agentNodes.get(remoteNode.id)
    if (!isAgent) {
        previous = remoteNodes.get(remoteNode.id)
        if (previous == null) {
            remoteNode.color = "rgba(200,200,200,1)"
            remoteNode.shape = "circle"
            remoteNode.label = netflowObj.dst_address
            remoteNode.title = remoteNodeHover(remoteNode)
            remoteNodes.add(remoteNode)
            newEdge.width = 1
            allEdges.add(newEdge)
        } else {
            remoteNode.label = previous.label
            remoteNode.title = remoteNodeHover(remoteNode)
            remoteNodes.update(remoteNode)
            allEdges.update(newEdge)
            // In the event an edge's metadata has changed
            newEdge.title = edgeHover(newEdge);
        }
    } else {
        // If it is an agent, still need to add the edge
        allEdges.update(newEdge)
        // In case edge changed, update title for hover data
        newEdge.title = edgeHover(newEdge)
    }
}

//-------------------------------------------
// Network events

// Change location of node in DataSet if user drags or moves it in the visualization
network.on("dragging", function(event) {
    if (event.nodes.length != 0) {
        for (index in event.nodes) {
            nodeID = event.nodes[index]
            // Handle remote node to keep it from moving during drag since we are still collecing real-time traffic
            node = remoteNodes.get(nodeID)
            if (node != null) {
                node.x = event.pointer.canvas.x
                node.y = event.pointer.canvas.y
                // Update DataSet
                remoteNodes.update(node)
            }
        }
    }
})


// Context menu to change properties of nodes when right-clicking
$(".vis-network canvas").contextMenu(function(event) {
    event.preventDefault();
})

trigger = $(".vis-network canvas");
trigger.contextMenu(false);

network.on("oncontext", function(event) {
    // Check if there is a node at the current position
    nodeID = network.getNodeAt({
        x: event.pointer.DOM.x,
        y: event.pointer.DOM.y
    })
    // If there is a node, select it
    if (nodeID != undefined) {
        node = agentNodes.get(nodeID)
        nodesSelected = network.getSelectedNodes().concat([nodeID])
        network.selectNodes(nodesSelected)
        // Turn on the context menu so it's displayed right away
        trigger = $('.vis-network canvas');
        if (trigger.hasClass('context-menu-disabled')) {
            trigger.contextMenu(true);
        }
    } else {
        // If no node, unselect previous nodes and turn off context menu
        trigger.contextMenu(false);
        network.unselectAll();
    }
});