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
	Players      []*Client
	DM           *Client
	Requests     chan *GameRequest
	ChatMessages *list.List
}

func (gs *GameServer) AddClient(char *Character, ws *websocket.Conn) {
	newClient := &Client{
		conn:      ws,
		character: char,
		requests:  gs.Requests,
		//responses: make(chan []byte, 1),
	}
	if char == nil {
		gs.DM = newClient
	} else {
		gs.Players = append(gs.Players, newClient)
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

var RequestRoutingMap = map[uint8]func(*GameRequest, *GameServer){
	0: skillCheck,
	1: getAllChatMessages,
	2: sendChatMessage,
}

func skillCheck(gr *GameRequest, gs *GameServer) {

}

func getAllChatMessages(gr *GameRequest, gs *GameServer) {

}

func sendChatMessage(gr *GameRequest, gs *GameServer) {
	//add chat message to chat messages
	//send message to everyone
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
