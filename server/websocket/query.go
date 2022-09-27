package websocket

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/Jeffail/gabs/v2"
	"github.com/artemis19/viz/server/database"
)

// Parent function to determine what type of query the client is requesting from the server
func Query(message *WS_message) ([]byte, error, string) {
	if message.To != "server" {
		return nil, nil, ""
	}

	var response []byte
	var err error
	var dataType string

	switch message.Action {
	case "select":
		response, err, dataType = Select(message)
	case "update":
	case "delete":
	case "create":
	}
	return response, err, dataType
}

// Takes in SELECT query and gets specified data and fields from the table
func Select(message *WS_message) ([]byte, error, string) {
	// Prepare list of hosts to retrieve from database
	var data []byte
	var err error
	var dataType string

	switch message.Table {
	case "hosts":
		var hosts []database.Host
		database.DB.Where(message.Fields).Find(&hosts)
		// Serialize data to send as JSON object
		data, err = AddInterfacesToHostData(hosts)

		dataType = "hosts"

	case "network_interfaces":
		var network_interfaces []database.NetworkInterface
		// TODO: Need to add error checking if no hosts
		database.DB.Where(message.Fields).Find(&network_interfaces)
		// Serialize data to send as JSON object
		data, err = json.Marshal(network_interfaces)
		dataType = "network_interfaces"

	case "net_flow":
		var netflow []database.NetFlow
		database.DB.Where(message.Fields).Find(&netflow)
		data, err = json.Marshal(netflow)
		dataType = "netflow"
	}

	if err != nil {
		log.Printf("Failed to complete SELECT query: %v", err)
		return nil, err, ""
	}

	return data, nil, dataType
}

// Function to handle server's response to the client's query
// Allows you to dynamically manipulate message back to the client if necessary
func Response(message *WS_message, response []byte, dataType string) ([]byte, error) {
	// Generate new JSON message
	jsonObj := gabs.New()
	jsonObj.Set("server", "from")
	jsonObj.Set(message.From, "to")

	jsonParsed, err := gabs.ParseJSON([]byte(response))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Response Query - Failed to parse JSON reply: %v", err))
	}

	jsonObj.Set(jsonParsed, "data")
	jsonObj.Set(dataType, "type")

	return []byte(jsonObj.String()), nil
}

// Function that allows server to send new data to all the websocket clients
func SendAll(data []byte, dataType string) ([]byte, error) {
	jsonObj := gabs.New()
	jsonObj.Set("server", "from")
	jsonObj.Set("", "to")

	jsonParsed, err := gabs.ParseJSON(data)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("SendAll Query - Failed to parse JSON reply: %v", err))
	}

	jsonObj.Set(dataType, "type")
	jsonObj.Set(jsonParsed, "data")

	return []byte(jsonObj.String()), nil
}

// Function to add network interfaces to marshalled host data
func AddInterfacesToHostData(hosts []database.Host) ([]byte, error) {
	var data []byte
	var err error
	data, err = json.Marshal(hosts)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to marshal JSON hosts data: %v", err))
	}

	for index, eachHost := range hosts {
		var interfaces []database.NetworkInterface
		database.DB.Where("host_id = ?", eachHost.MachineID).Find(&interfaces)

		finalJSON := gabs.New()
		hostJSON, err := gabs.ParseJSON(data)

		finalJSON.Set(hostJSON.Data(), "host_data")
		niData, err := json.Marshal(interfaces)
		if err != nil {
			log.Printf("SELECT: Failed to parse interfaces json: %s", err)
		}
		niJSON, err := gabs.ParseJSON(niData)
		finalJSON.SetP(niJSON.Data(), fmt.Sprintf("host_data.%d.network_interfaces", index))

		data = []byte(hostJSON.String())
	}
	return data, err
}
