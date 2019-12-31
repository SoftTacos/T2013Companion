package main

import (
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
			response := EncodeResponse(req.RequestType, req.Message, websocket.TextMessage)
			fmt.Println("SENDING RESPONSE: ", response[0], string(response[1:]))
			if err := req.Client.conn.WriteMessage(2, response); err != nil {
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

		requestType, data := DecodeRequest(messageType, rawMessage) //websocket.BinaryMessage, websocket.TextMessage
		fmt.Println("REQUEST RECEIVED: ", requestType, string(data))
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

func DecodeRequest(msgType int, rawMsg []byte) (uint8, []byte) {
	/*fmt.Println(msgType, string(rawMsg))
	if len(rawMsg) < 3 {
		fmt.Println("Message too small, skipping")
		return 0, []byte("")
	}
	rrt := rawMsg[0:2]
	fmt.Println("RAW: ", rrt)
	requestType := binary.LittleEndian.Uint16(rrt)
	fmt.Println("TYPE: ", requestType)
	*/

	return rawMsg[0], rawMsg[1:]
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
