package websocket

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Jeffail/gabs/v2"
	"github.com/artemis19/viz/server/database"
	"log"
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
		// TODO: Need to add error checking if no hosts
		database.DB.Where(message.Fields).Find(&hosts)
		// Serialize data to send as JSON object
		data, err = json.Marshal(hosts)
		dataType = "hosts"
	case "network_interfaces":
		var network_interfaces []database.NetworkInterface
		// TODO: Need to add error checking if no hosts
		database.DB.Where(message.Fields).Find(&network_interfaces)
		// Serialize data to send as JSON object
		data, err = json.Marshal(network_interfaces)
		dataType = "network_interfaces"
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
