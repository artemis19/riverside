function agentNodeHover(node) {
    const container = document.createElement("div");
    interfaces = "<ol>"
    for (index in node.networkInterfaces) {
        interface = node.networkInterfaces[index]
        interfaces += `<li>${interface.ip_address}</li>`
    }
    interfaces += "</ol>"

    container.innerHTML = `
	<h5>${node.label}</h5>
		<table>
			<tr>
			<th>Field</th>
			<th>Value</th>
			</tr>
			<tr>
			<td>Created At: </td>
			<td>${node.createdAt}</td>
			</tr>
			<tr>
			<td>Updated At: </td>
			<td>${node.updatedAt}</td>
			</tr>
			<tr>
    		<td>Hostname: </td>
    		<td>${node.hostname}</td>
  			</tr>
			<tr>
			<td>OS: </td>
			<td>${node.os}</td>
			</tr>
			<tr>
			<td>Architecture: </td>
			<td>${node.arch}</td>
		</table>
		<h6><b>Known Interfaces:</b></h6>
		${interfaces}`
    return container
}

function remoteNodeHover(node) {
    const container = document.createElement("div");
    container.innerHTML = `
	<h5>${node.label}</h5>
		<table>
			<tr>
			<th>Field</th>
			<th>Value</th>
			</tr>
			<tr>
			<td>IP Address: </td>
			<td>${node.ipAddress}</td>
			</tr>
			<tr>
			<td>Created At: </td>
			<td>${node.createdAt}</td>
			</tr>
			<tr>
			<td>Updated At: </td>
			<td>${node.updatedAt}</td>
			</tr>
		</table>`
    return container
}

function edgeHover(edge) {
    const container = document.createElement("div");
    // Label source/dest based on "originator"
    startAddress = edge.originAddress
    if (startAddress == edge.dstAddress) {
        startPort = edge.dstPort
        endAddress = edge.srcAddress
        endPort = edge.srcPort
    }
    if (startAddress == edge.srcAddress) {
        startPort = edge.srcPort
        endAddress = edge.dstAddress
        endPort = edge.dstPort
    }
    // Handle ports if ICMP or N/A protocol
    if (startPort == 99999) {
        startPort = ""
    } else {
        startPort = ":" + startPort
    }
    if (endPort == 99999) {
        endPort = ""
    } else {
        endPort = ":" + endPort
    }
    // Display arrows for directions
    // if (edge.direction == "to") {
    //     direction = "&xrarr;"
    // }
    // if (edge.direction == "from") {
    //     direction = "&xlarr;"
    // }
    direction = "&xrarr;"
    if (edge.direction == "bidirectional") {
        direction = "&xharr;"
    }
    container.innerHTML = `
	<h5>${edge.protocol}: ${startAddress}${startPort}${direction}${endAddress}${endPort}</h5>
		<table>
			<tr>
			<th>Field</th>
			<th>Value</th>
			</tr>
			<tr>
			<td>Start Time: </td>
			<td>${edge.startTime}: </td>
			</tr>
			<tr>
			<td>End Time: </td>
			<td>${edge.endTime}: </td>
			</tr>
			<td>Throughput: </td>
			<td>${edge.throughput}</td>
			</tr>
		</table>`
    return container
}