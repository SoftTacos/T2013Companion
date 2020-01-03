package main

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID         uint
	conn       *websocket.Conn
	character  *Character
	requests   chan *GameRequest
	gameServer *GameServer
	//responses chan []byte
}

func (c *Client) Start() {
	go c.listener()
	//go pc.writer()
}

//string request format: <uint16 requestType>:<string data>
func (c *Client) listener() {
	for {
		messageType, rawMessage, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Printf("error: %v", err)
			}
			break
		}

		requestType, data := DecodeRequest(messageType, rawMessage) //websocket.BinaryMessage, websocket.TextMessage
		//fmt.Println("REQUEST RECEIVED: ", requestType, string(data))
		c.requests <- &GameRequest{
			MessageType: messageType,
			RequestType: requestType,
			Message:     data,
			Client:      c,
		}

		//testBytes := make([]byte, 16)
		//binary.LittleEndian.PutUint16(testBytes[0:], uint16(12))
		//testBytes[2:] = [12]byte{`0`, `1`}
		//testreq, testdata := ParseRequest(websocket.TextMessage, testBytes)
		//fmt.Println("TEST: ", testreq, testdata)
	}
	//client has closed the connection
	fmt.Println("Client has closed the connection: ", c.ID)
	delete(c.gameServer.Clients, c.ID)
}

func DecodeRequest(msgType int, rawMsg []byte) (uint8, []byte) {
	return rawMsg[0], rawMsg[1:]
} /*fmt.Println(msgType, string(rawMsg))
if len(rawMsg) < 3 {
	fmt.Println("Message too small, skipping")
	return 0, []byte("")
}
rrt := rawMsg[0:2]
fmt.Println("RAW: ", rrt)
requestType := binary.LittleEndian.Uint16(rrt)
fmt.Println("TYPE: ", requestType)
*/
