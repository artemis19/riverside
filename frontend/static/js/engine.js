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

var globalOptions = {
    width: "100%",
    height: "100%",
    clickToUse: true,
    physics: {
        enabled: false,
    }
};

// Initialize network
var network = new vis.Network(container, { nodes: allNodes, edges: allEdges },
    globalOptions);

//-------------------------------------------

timelineOptions = {
    height: "180px",
}

var items = new vis.DataSet([]);

var timelineContainer = document.getElementById("timeline")
var timeline = new vis.Timeline(timelineContainer, items, timelineOptions)

//-------------------------------------------

function addAgentNode(hostObj) {
    agentNode = {
        id: hostObj.machine_id,
        label: hostObj.hostname,
        shape: 'circle',
        size: 50,
    }
    // Adds new node to agentNodes DataSet
    agentNodes.add(agentNode)
    console.log("agent node: ", agentNode.id)
}

function handleNetflow(netflowObj) {
    remoteNode = {
        id: netflowObj.remote_host_id,
        label: netflowObj.dst_address,
        shape: 'circle',
        size: 50,
        color: {
            background: 'silver',
            border: 'gray',
            highlight: {
                background: '#eaeaea',
                border: 'silver',
            }
        }
    }

    newEdge = {
        id: netflowObj.host_id + "->" + netflowObj.remote_host_id,
        to: netflowObj.remote_host_id,
        from: netflowObj.host_id,
        arrows: {
            to: true,
        },
        length: 10,
    }

    // Add remote nodes to DataSet if it doesn't already exist
    previous = remoteNodes.get(remoteNode.id)
    if (previous == null) {
        remoteNodes.add(remoteNode)
        allEdges.add(newEdge)
    }
    console.log(allEdges)

}