package main

import (
	"encoding/binary"
	"fmt"

	"github.com/gorilla/websocket"
)

type GameRequest struct {
	RequestType uint16
	MessageType int
	Message     []byte
	Client      *Client
}

type GameServer struct {
	Players  []*Client
	DM       *Client
	Requests chan *GameRequest
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
			if err := req.Client.conn.WriteMessage(1, []byte("[response!]")); err != nil {
				fmt.Println(err)
				continue
			}
		}
	}
}

type Client struct {
	conn      *websocket.Conn
	character *Character
	requests  chan *GameRequest
	//responses chan []byte
}

func (pc *Client) Start() {
	go pc.listener()
	//go pc.writer()
}

//string request format: <uint16 requestType>:<string data>
func (pc *Client) listener() {
	for {
		messageType, rawMessage, err := pc.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Printf("error: %v", err)
			}
			break
		}

		requestType, data := ParseRequest(messageType, rawMessage) //websocket.BinaryMessage, websocket.TextMessage
		pc.requests <- &GameRequest{
			MessageType: messageType,
			RequestType: requestType,
			Message:     data,
			Client:      pc,
		}

		//testBytes := make([]byte, 16)
		//binary.LittleEndian.PutUint16(testBytes[0:], uint16(12))
		//testBytes[2:] = [12]byte{`0`, `1`}
		//testreq, testdata := ParseRequest(websocket.TextMessage, testBytes)
		//fmt.Println("TEST: ", testreq, testdata)
	}
}

func ParseRequest(msgType int, rawMsg []byte) (uint16, []byte) {
	fmt.Println(msgType, string(rawMsg))
	if len(rawMsg) < 3 {
		fmt.Println("Message too small, skipping")
		return 0, []byte("")
	}
	rrt := rawMsg[0:2]
	fmt.Println("RAW: ", rrt)
	requestType := binary.LittleEndian.Uint16(rrt)
	fmt.Println("TYPE: ", requestType)

	return requestType, rawMsg[2:]
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
