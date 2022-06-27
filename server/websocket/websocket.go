package websocket

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
)

var WSWriter chan []byte

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
	for {
		writeMsg := <-WSWriter
		for _, eachClient := range WSClients {
			log.Printf("Trying to write to client %v", eachClient)
			if err := eachClient.conn.WriteMessage(websocket.TextMessage, writeMsg); err != nil {
				log.Printf("Websocket write error: %v", err)
			} else {
				log.Printf("Websocket wrote: %v", string(writeMsg))
			}
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

	// Initialize websocket's write channel
	WSWriter = make(chan []byte)
	go WSWriteChannel()

	log.Fatal(http.ListenAndServe(bind, nil))
}
