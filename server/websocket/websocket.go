package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"time"

	"github.com/artemis19/viz/server/database"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
)

// Channels for sending websocket data
var WSNetflowWriter = make(chan database.NetFlow)
var WSHostWriter = make(chan database.Host)

// Golang "stopwatch" for batching
var Batch time.Ticker

// Create struct for websocket type
type WebSocketClient struct {
	writer chan []byte
	conn   *websocket.Conn
}

// List of clients
var WSClients []*WebSocketClient

// Struct for websocket data
var upgrader = websocket.Upgrader{}

func setupRoutes() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", wsEndpoint)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Viz server is online.")
}

func WSWriteChannel() {
	// Create variables to send batched messages to frontend
	BatchTicker := time.NewTicker(1 * time.Second)
	hostBatch := make([]database.Host, 0)
	netflowBatch := make([]database.NetFlow, 0)

	for {
		select {
		case hostData := <-WSHostWriter:
			hostBatch = append(hostBatch, hostData)
		case netflowData := <-WSNetflowWriter:
			netflowBatch = append(netflowBatch, netflowData)
		case <-BatchTicker.C:
			// log.Printf("Ticker: %v, clearing batch", t)
			for _, eachClient := range WSClients {
				// log.Printf("Trying to write to client: %v", eachClient)
				if len(hostBatch) != 0 {
					allHostData, err := AddInterfacesToHostData(hostBatch)
					if err != nil {
						log.Printf("Failed to marshal host batch JSON data: %v", err)
					}
					sendHostData, err := SendAll(allHostData, "hosts")
					if err != nil {
						log.Printf("Failed to send host batch JSON data: %v", err)
					}
					if err := eachClient.conn.WriteMessage(websocket.TextMessage, sendHostData); err != nil {
						log.Printf("Websocket write error: %v", err)
					} else {
						log.Printf("Websocket wrote: %v", string(sendHostData))
					}
				}
				if len(netflowBatch) != 0 {
					allNetflowData, err := json.Marshal(netflowBatch)
					if err != nil {
						log.Printf("Failed to marshal host batch JSON data: %v", err)
					}
					sendNetflowData, err := SendAll(allNetflowData, "netflow")
					if err != nil {
						log.Printf("Failed to send host batch JSON data: %v", err)
					}
					if err := eachClient.conn.WriteMessage(websocket.TextMessage, sendNetflowData); err != nil {
						log.Printf("Websocket write error: %v", err)
					} else {
						log.Printf("Websocket wrote: %v", string(sendNetflowData))
					}
				}
			}
			// Clear batches
			hostBatch = hostBatch[:0]
			netflowBatch = netflowBatch[:0]
		}
	}
}

// Remove websocket client if no longer there
func removeWSClient(client *WebSocketClient) {
	var newWSClients []*WebSocketClient

	for _, eachClient := range WSClients {
		if eachClient == client {
			continue
		}
		newWSClients = append(newWSClients, eachClient)
	}

	// Update global list of websocket clients
	WSClients = newWSClients

}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to connect %v", err)
	}

	defer conn.Close()
	log.Printf("Websocket client successfully connected.")

	client := &WebSocketClient{conn: conn}
	WSClients = append(WSClients, client)
	defer removeWSClient(client)

	for {
		messageType, message, err := client.conn.ReadMessage()
		if err != nil {
			log.Printf("Could not read data: %v", err)
			return
		}

		log.Printf("Websocket received message %v with type %v", string(message), messageType)

		msg, err := NewMessage(message)
		if err != nil {
			log.Printf("Failed to validate websocket message: %v", err)
			continue
		}

		log.Printf("Serialized JSON as struct: %v", msg)
		reply, err, dataType := Query(msg)
		if err != nil {
			log.Printf("Failed to query DB from websocket request: %v", err)
		}
		response, err := Response(msg, reply, dataType)
		if err != nil {
			log.Printf("Failed to create JSON response message: %v", err)
		}
		log.Printf("Responding with: %s", response)
		// Sends message through websocket with response data to client
		client.conn.WriteMessage(websocket.TextMessage, response)
	}
}

func Serve() {
	bind := fmt.Sprintf(":%v", viper.Get("websocketPort").(string))
	log.Printf("Web socket is listening on %v...\n", bind)
	setupRoutes()

	// Websocket's write channel
	go WSWriteChannel()

	log.Fatal(http.ListenAndServe(bind, nil))
	// Use http.ListenAndServeTLS with a valid cert to have websockets function when not working locally
	// log.Fatal(http.ListenAndServeTLS(bind, "/etc/letsencrypt/live/www.riverside-vis.dev/fullchain.pem", "/etc/letsencrypt/live/www.riverside-vis.dev/privkey.pem", nil))
}
