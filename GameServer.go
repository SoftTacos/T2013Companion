package main

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type GameRequest struct {
	Message []byte
	Client  *Client
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
	responses chan []byte
}

func (pc *Client) Start() {
	go pc.listener()
	//go pc.writer()
}

func (pc *Client) listener() {
	for {
		_, message, err := pc.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Printf("error: %v", err)
			}
			break
		}
		//message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		pc.requests <- &GameRequest{
			Message: message,
			Client:  pc,
		}
	}
}

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
