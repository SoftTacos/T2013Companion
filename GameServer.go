package main

import (
	"container/list"
	"fmt"

	"github.com/gorilla/websocket"
)

type GameRequest struct {
	RequestType uint8
	MessageType int
	Message     []byte
	Client      *Client
}

type GameServer struct {
	Clients      map[uint]*Client
	NextClientID uint
	DM           *Client
	Requests     chan *GameRequest
	ChatMessages *list.List
}

func (gs *GameServer) AddClient(char *Character, ws *websocket.Conn) {
	newClient := &Client{
		ID:         gs.NextClientID,
		conn:       ws,
		character:  char,
		requests:   gs.Requests,
		gameServer: gs,
	}
	if char == nil {
		gs.Clients[0] = newClient
	} else {
		gs.Clients[gs.NextClientID] = newClient
		gs.NextClientID++
	}
	newClient.Start()
}

func (gs *GameServer) Handle() {
	for {
		select {
		case req := <-gs.Requests:
			fmt.Println("Request: ", req.RequestType, req.Client.conn.RemoteAddr(), string(req.Message))
			/*
				response := EncodeResponse(req.RequestType, req.Message, websocket.TextMessage)
				fmt.Println("SENDING RESPONSE: ", response[0], string(response[1:]))
				if err := req.Client.conn.WriteMessage(2, response); err != nil {
					fmt.Println(err)
					continue
				}
			*/
			funk, ok := RequestRoutingMap[req.RequestType]
			if !ok {
				fmt.Println("ERROR: Request type not valid", req.RequestType)
				response := EncodeResponse(req.RequestType, []byte("ERR"), websocket.BinaryMessage)
				if err := req.Client.conn.WriteMessage(2, response); err != nil {
					fmt.Println(err)
					continue
				}
				continue
			}
			funk(req, gs) //maybe make async?
		}
	}
}

func (gs *GameServer) RemoveClient(client *Client) {
	if gs.DM == client {

	}
}

var RequestRoutingMap = map[uint8]func(*GameRequest, *GameServer){
	0: skillCheck,
	1: getAllChatMessages,
	2: sendChatMessage,
}

func skillCheck(gr *GameRequest, gs *GameServer) {

}

func getAllChatMessages(gr *GameRequest, gs *GameServer) {
	if gs.ChatMessages.Len() == 0 {
		return
	}
	messages := []byte{}
	for e := gs.ChatMessages.Front(); e != nil; e = e.Next() {
		message := append(e.Value.([]byte), []byte(",")...)
		fmt.Println(message)
		messages = append(messages, message...)
	}
	sendResponse(EncodeResponse(1, messages[0:len(messages)-1], gr.MessageType), gr.Client)
}

func sendChatMessage(gr *GameRequest, gs *GameServer) {
	//add chat message to chat messages
	fmt.Println(gr.Message)
	gs.ChatMessages.PushBack(gr.Message)
	//send message to everyone
	for _, client := range gs.Clients {
		sendResponse(EncodeResponse(2, gr.Message, gr.MessageType), client)
	}
}

func EncodeResponse(responseType uint8, message []byte, messageType int) []byte {
	/*response := make([]byte, len(message)+2)
	binary.LittleEndian.PutUint16(response[0:], responseType)
	copy(response[2:], message)
	fmt.Println(string(response), response)*/
	response := make([]byte, 1+len(message))
	response[0] = responseType
	copy(response[1:], message)
	return response //[]byte("ENCODED RESPONSE")
}

func sendResponse(response []byte, client *Client) {
	//response := EncodeResponse(req.RequestType, req.Message, websocket.TextMessage)
	if err := client.conn.WriteMessage(2, response); err != nil {
		fmt.Println("Error?", err)
		return
	}
	//fmt.Println("SENT RESPONSE: ", response[0], string(response[1:]))

}

/*
//this might not be needed after all
func (pc *Client) writer() {
	for {
		select {
		case response := <-pc.responses:
			if err := pc.conn.WriteMessage(1, response); err != nil {
				fmt.Println(err)
				continue
			}
		}
	}
}
*/
